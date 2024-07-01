// Copyright 2021 Wei (Sam) Wang <sam.wang.0723@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/heptiolabs/healthcheck"
	grpc_sentry "github.com/johnbellone/grpc-middleware-sentry"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/handlers"
	"github.com/samwang0723/jarvis/internal/app/middleware"
	pb "github.com/samwang0723/jarvis/internal/app/pb"
	gatewaypb "github.com/samwang0723/jarvis/internal/app/pb/gateway"
	"github.com/samwang0723/jarvis/internal/app/services"
	"github.com/samwang0723/jarvis/internal/database"
	"github.com/samwang0723/jarvis/internal/database/dal"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gracefulShutdownPeriod = 5 * time.Second
	maxGoRoutines          = 10000
	dnsResolveTimeout      = 200 * time.Millisecond
	databasePingTimeout    = 1 * time.Second
	readTimeout            = 5 * time.Second
	writeTimeout           = 10 * time.Second
	defaultHTTPTimeout     = 10 * time.Second
	smartproxy             = "SMART_PROXY"
)

type IServer interface {
	Name() string
	Logger() *zerolog.Logger
	Handler() handlers.IHandler
	Config() *config.Config
	GRPCServer() *grpc.Server
	Run(context.Context) error
	Start(context.Context) error
	Stop() error
}

type server struct {
	opts Options
}

//nolint:gosec //skip tls verification
func Serve(cfg *config.Config, logger *zerolog.Logger) {
	// sequence: handler(dto) -> service(dto to dao) -> DAL(dao) -> dbPool
	// initialize DAL layer
	dbPool := database.GormFactory(cfg)
	dalService := dal.New(
		dal.WithDB(dbPool),
		dal.WithBalanceRepository(dal.NewBalanceRepository(dbPool)),
		dal.WithTransactionRepository(dal.NewTransactionRepository(dbPool)),
		dal.WithOrderRepository(dal.NewOrderRepository(dbPool)),
	)

	// Initialize a HTTP client with proxy
	smartProxy := ""
	if proxy := os.Getenv(smartproxy); proxy != "" {
		smartProxy = proxy
	}
	proxyURL, pErr := url.Parse(smartProxy)
	if pErr != nil {
		logger.Fatal().Err(pErr).Msg("failed to parse proxy url")
	}
	proxy := &http.Client{
		Timeout: defaultHTTPTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Proxy:           http.ProxyURL(proxyURL),
		},
	}

	// bind DAL layer with service
	dataService := services.New(
		services.WithDAL(dalService),
		services.WithKafka(services.KafkaConfig{
			GroupID: cfg.Kafka.GroupID,
			Brokers: cfg.Kafka.Brokers,
			Topics:  cfg.Kafka.Topics,
			Logger:  logger,
		}),
		services.WithRedis(services.RedisConfig{
			Master:        cfg.RedisCache.Master,
			SentinelAddrs: cfg.RedisCache.SentinelAddrs,
			Logger:        logger,
			Password:      cfg.RedisCache.Password,
		}),
		services.WithCronJob(services.CronjobConfig{
			Logger: logger,
		}),
		services.WithLogger(logger),
		services.WithProxy(proxy),
	)
	// associate service with handler
	handler := handlers.New(dataService, logger)

	gRPCServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			selector.StreamServerInterceptor(
				auth.StreamServerInterceptor(middleware.Authenticate),
				selector.MatchFunc(middleware.AuthRoutes),
			),
			grpc_middleware.ChainStreamServer(
				grpc_sentry.StreamServerInterceptor(),
			),
		),
		grpc.ChainUnaryInterceptor(
			selector.UnaryServerInterceptor(
				auth.UnaryServerInterceptor(middleware.Authenticate),
				selector.MatchFunc(middleware.AuthRoutes),
			),
			grpc_middleware.ChainUnaryServer(
				grpc_sentry.UnaryServerInterceptor(),
			),
		),
	)

	// health check
	health := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(maxGoRoutines))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-database-read-dns",
		healthcheck.DNSResolveCheck(cfg.Replica.Host, dnsResolveTimeout),
	)
	health.AddReadinessCheck(
		"upstream-database-write-dns",
		healthcheck.DNSResolveCheck(cfg.Database.Host, dnsResolveTimeout),
	)

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	genericDB, dbErr := dbPool.DB()
	if dbErr != nil {
		logger.Fatal().Err(dbErr).Msg("failed to get db")
	}

	health.AddReadinessCheck(
		"database",
		healthcheck.DatabasePingCheck(genericDB, databasePingTimeout),
	)

	s := newServer(
		Name(cfg.Server.Name),
		Config(cfg),
		Logger(logger),
		Handler(handler),
		GRPCServer(gRPCServer),
		HealthCheck(health),
		BeforeStart(func() error {
			dataService.StartCron()

			return nil
		}),
		BeforeStop(func() error {
			dataService.StopCron()
			err := dataService.StopRedis()
			if err != nil {
				return fmt.Errorf("stop_redis: failed, reason: %w", err)
			}

			err = dataService.StopKafka()
			if err != nil {
				logger.Error().Err(err).Msg("StopKafka error")
			}
			sqlDB, err := dbPool.DB()
			if err != nil {
				return err
			}
			defer sqlDB.Close()

			return nil
		}),
	)

	err := s.Run(context.Background())
	if err != nil {
		s.Logger().Error().Err(err).Msg("server run failed")
	}
}

func newServer(opts ...Option) IServer {
	o := Options{}
	for _, opt := range opts {
		opt(&o)
	}

	return &server{
		opts: o,
	}
}

func (s *server) Start(ctx context.Context) error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

	signature := `
      _                  _                   _
     | |                (_)                 (_)
     | | __ _ _ ____   ___ ___    __ _ _ __  _
 _   | |/ _' | '__\ \ / / / __|  / _' | '_ \| |
| |__| | (_| | |   \ V /| \__ \ | (_| | |_) | |
 \____/ \__,_|_|    \_/ |_|___/  \__,_| .__/|_|
                                      | |
                                      |_|       Version (%s)
High performance stock analysis tool
Environment (%s)
_______________________________________________
`
	signatureOut := fmt.Sprintf(signature, s.Config().Server.Version, s.Config().Environment)
	//nolint:nolintlint, forbidigo
	fmt.Println(signatureOut)

	// start gRPC server
	cfg := config.GetCurrentConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GrpcPort)
	// start reversed proxy http server
	go s.startGRPCGateway(ctx, addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.Logger().Info().Msgf("gRPC server listening on %s", addr)

	pb.RegisterJarvisV1Server(s.GRPCServer(), s)

	go func() {
		if err = s.GRPCServer().Serve(lis); err != nil {
			s.Logger().Fatal().Msgf("gRPC server serve failed: %s", err.Error())
		}
	}()

	// listening kafka
	s.Handler().ListeningKafkaInput(ctx)

	return err
}

func (s *server) Stop() error {
	var err error
	for _, fn := range s.opts.BeforeStop {
		if err = fn(); err != nil {
			break
		}
	}

	s.GRPCServer().GracefulStop()

	<-time.After(gracefulShutdownPeriod)
	s.Logger().Warn().Msg("server being gracefully shuted down")

	return err
}

// Run starts the server and shut down gracefully afterwards
func (s *server) Run(ctx context.Context) error {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if err := s.Start(childCtx); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	if s.Config().Environment != "local" {
		// schedule to preset the stocks met expectation condition yesterday for realtime tracking
		var waitGroup sync.WaitGroup
		waitGroup.Add(1)

		go func(ctx context.Context, svc *server) {
			defer waitGroup.Done()

			err := svc.Handler().CronjobPresetRealtimeMonitoringKeys(childCtx, "40 8 * * 1-5")
			if err != nil {
				svc.Logger().Error().Err(err).Msg("CronjobPresetRealtimeMonitoringKeys error")
			}

			err = svc.Handler().CrawlingRealTimePrice(childCtx, "*/3 9-13 * * 1-5")
			if err != nil {
				svc.Logger().Error().Err(err).Msg("RetrieveRealTimePrice error")
			}

			<-ctx.Done()
		}(childCtx, s)

		waitGroup.Wait()
	}

	select {
	case <-quit:
		s.Logger().Warn().Msg("signal interrupt")
		cancel()
	case <-childCtx.Done():
		s.Logger().Warn().Msg("main context being canceled")
	}

	return s.Stop()
}

func (s *server) Logger() *zerolog.Logger {
	return s.opts.Logger
}

func (s *server) Name() string {
	return s.opts.Name
}

func (s *server) Handler() handlers.IHandler {
	return s.opts.Handler
}

func (s *server) Config() *config.Config {
	return s.opts.Config
}

func (s *server) GRPCServer() *grpc.Server {
	return s.opts.GRPCServer
}

func (s *server) HealthCheck() healthcheck.Handler {
	return s.opts.HealthCheck
}

func (s *server) startGRPCGateway(ctx context.Context, addr string) {
	c, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()

	err := gatewaypb.RegisterJarvisV1HandlerFromEndpoint(
		c,
		mux,
		addr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		s.Logger().Error().Err(err).Msg("cannot register grpc gateway")

		return
	}

	// add healthcheck into gRPC gateway mux
	err = mux.HandlePath(
		"GET",
		"/live",
		func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			s.HealthCheck().LiveEndpoint(w, r)
		},
	)
	if err != nil {
		s.Logger().Error().Err(err).Msg("cannot handle /live path")

		return
	}

	err = mux.HandlePath(
		"GET",
		"/ready",
		func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
			s.HealthCheck().ReadyEndpoint(w, r)
		},
	)
	if err != nil {
		s.Logger().Error().Err(err).Msg("cannot handle /ready path")

		return
	}

	httpMux := http.NewServeMux()
	// merge grpc gateway endpoint handling
	httpMux.Handle("/", mux)

	cfg := config.GetCurrentConfig()
	host := fmt.Sprintf(":%d", cfg.Server.Port)

	srv := &http.Server{
		ReadHeaderTimeout: readTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		Addr:              host,
		Handler:           cors.AllowAll().Handler(httpMux),
	}

	s.Logger().Info().Msgf("http server start listening on %s", host)
	err = srv.ListenAndServe()
	if err != nil {
		s.Logger().Error().Err(err).Msgf("cannot listen and serve http server on %s", host)

		return
	}
}
