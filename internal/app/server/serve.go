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
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/heptiolabs/healthcheck"
	config "github.com/samwang0723/jarvis/configs"
	"github.com/samwang0723/jarvis/internal/app/dto"
	"github.com/samwang0723/jarvis/internal/app/handlers"
	pb "github.com/samwang0723/jarvis/internal/app/pb"
	gatewaypb "github.com/samwang0723/jarvis/internal/app/pb/gateway"
	"github.com/samwang0723/jarvis/internal/app/services"
	"github.com/samwang0723/jarvis/internal/concurrent"
	"github.com/samwang0723/jarvis/internal/cronjob"
	"github.com/samwang0723/jarvis/internal/db"
	"github.com/samwang0723/jarvis/internal/db/dal"
	"github.com/samwang0723/jarvis/internal/helper"
	log "github.com/samwang0723/jarvis/internal/logger"
	structuredlog "github.com/samwang0723/jarvis/internal/logger/structured"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	grpc_sentry "github.com/johnbellone/grpc-middleware-sentry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	gracefulShutdownPeriod = 5 * time.Second
)

type IServer interface {
	Name() string
	Logger() structuredlog.ILogger
	Handler() handlers.IHandler
	Config() *config.Config
	Dispatcher() *concurrent.Dispatcher
	GRPCServer() *grpc.Server
	Run(context.Context) error
	Start(context.Context) error
	Stop() error
}

type server struct {
	opts Options
}

func Serve() {
	config.Load()
	cfg := config.GetCurrentConfig()
	logger := structuredlog.Logger(cfg)
	// sequence: handler(dto) -> service(dto to dao) -> DAL(dao) -> database
	// initialize DAL layer
	db := db.GormFactory(cfg)
	dalService := dal.New(dal.WithDB(db))
	// bind DAL layer with service
	dataService := services.New(
		services.WithDAL(dalService),
		services.WithCronJob(cronjob.New(logger)),
	)
	// associate service with handler
	handler := handlers.New(dataService)
	gRPCServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_sentry.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_sentry.UnaryServerInterceptor(),
		)),
	)

	//health check
	health := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(10000))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	health.AddReadinessCheck(
		"upstream-database-read-dns",
		healthcheck.DNSResolveCheck(cfg.Replica.Host, 200*time.Millisecond))
	health.AddReadinessCheck(
		"upstream-database-write-dns",
		healthcheck.DNSResolveCheck(cfg.Database.Host, 200*time.Millisecond))
	health.AddReadinessCheck(
		"upstream-redis-dns",
		healthcheck.DNSResolveCheck(cfg.RedisCache.Host, 200*time.Millisecond))

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	genericDB, _ := db.DB()
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(genericDB, 1*time.Second))

	s := newServer(
		Name(cfg.Server.Name),
		Config(cfg),
		Logger(logger),
		Handler(handler),
		Dispatcher(concurrent.NewDispatcher(cfg.WorkerPool.MaxPoolSize)),
		GRPCServer(gRPCServer),
		HealthCheck(health),
		BeforeStart(func() error {
			// initialize global job queue
			concurrent.JobQueue = make(concurrent.JobChan, cfg.WorkerPool.MaxQueueSize)
			dataService.StartCron()
			return nil
		}),
		BeforeStop(func() error {
			dataService.StopCron()
			//no need to explictly close a channel, it will be garbage collected
			//close(concurrent.JobQueue)

			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			defer sqlDB.Close()

			return nil
		}),
	)

	log.Initialize(s.Logger())
	err := s.Run(context.Background())
	if err != nil && s.Logger() != nil {
		log.Errorf("error returned by service.Run(): %s\n", err.Error())
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
	signatureOut := fmt.Sprintf(signature, "v1.1.1a", helper.GetCurrentEnv())
	fmt.Println(signatureOut)

	// starting the workerpool
	s.Dispatcher().Run(ctx)

	// by default starting cronjob for regular daily updates pulling
	// cronjob using redis distrubted lock to prevent multiple instances
	// pulling same content
	s.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
		Schedule: "30 16 * * 1-5",
		Types:    []dto.DownloadType{dto.DailyClose, dto.ThreePrimary},
	})
	s.Handler().CronDownload(ctx, &dto.StartCronjobRequest{
		Schedule: "30 18 * * 1-5",
		Types:    []dto.DownloadType{dto.Concentration},
	})

	// start gRPC server
	cfg := config.GetCurrentConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.GrpcPort)
	// start revered proxy http server
	go s.startGRPCGateway(ctx, addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("gRPC server running.")

	pb.RegisterJarvisServer(s.GRPCServer(), s)
	go func() {
		if err = s.GRPCServer().Serve(lis); err != nil {
			log.Fatalf("gRPC server serve failed: %s", err)
		}
	}()

	return err
}

func (s *server) Stop() error {
	var err error
	for _, fn := range s.opts.BeforeStop {
		if err = fn(); err != nil {
			break
		}
	}

	// graceful shutdown workerpool
	s.Dispatcher().WaitGroup().Wait()
	s.GRPCServer().GracefulStop()

	<-time.After(gracefulShutdownPeriod)
	log.Warn("server being gracefully shuted down")

	return err
}

// Run starts the server and shut down gracefully afterwards
func (s *server) Run(ctx context.Context) error {
	childCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	if s.Logger() != nil {
		defer s.Logger().Flush()
	}

	if err := s.Start(childCtx); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-quit:
		log.Warn("singal interrupt")
		cancel()
	case <-childCtx.Done():
		log.Warn("main context being cancelled")
	}
	return s.Stop()
}

func (s *server) Logger() structuredlog.ILogger {
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

func (s *server) Dispatcher() *concurrent.Dispatcher {
	return s.opts.Dispatcher
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
	err := gatewaypb.RegisterJarvisHandlerFromEndpoint(
		c,
		mux,
		addr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatalf("cannot start grpc gateway: %v", err)
	}

	// add healthcheck into gRPC gateway mux
	mux.HandlePath("GET", "/live", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		s.HealthCheck().LiveEndpoint(w, r)
	})
	mux.HandlePath("GET", "/ready", func(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
		s.HealthCheck().ReadyEndpoint(w, r)
	})

	cfg := config.GetCurrentConfig()
	host := fmt.Sprintf(":%d", cfg.Server.Port)
	err = http.ListenAndServe(host, mux)
	if err != nil {
		log.Fatalf("cannot listen and server: %v", err)
	}
}
