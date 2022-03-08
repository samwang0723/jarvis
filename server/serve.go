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
	"os"
	"os/signal"
	"samwang0723/jarvis/concurrent"
	"samwang0723/jarvis/config"
	"samwang0723/jarvis/cronjob"
	"samwang0723/jarvis/db"
	"samwang0723/jarvis/db/dal"
	"samwang0723/jarvis/handlers"
	log "samwang0723/jarvis/logger"
	structuredlog "samwang0723/jarvis/logger/structured"
	pb "samwang0723/jarvis/pb"

	"google.golang.org/grpc"

	"samwang0723/jarvis/services"
	"syscall"
	"time"
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
	gRPCServer := grpc.NewServer()

	s := newServer(
		Name(cfg.Server.Name),
		Config(cfg),
		Logger(logger),
		Handler(handler),
		Dispatcher(concurrent.NewDispatcher(cfg.WorkerPool.MaxPoolSize)),
		GRPCServer(gRPCServer),
		BeforeStart(func() error {
			// initialize global job queue
			concurrent.JobQueue = make(concurrent.JobChan, cfg.WorkerPool.MaxQueueSize)
			dataService.StartCron()
			return nil
		}),
		BeforeStop(func() error {
			dataService.StopCron()
			close(concurrent.JobQueue)

			sqlDB, _ := db.DB()
			return sqlDB.Close()
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

	// starting the workerpool
	s.Dispatcher().Run(ctx)

	//TODO: separate cron action to other place
	err := s.Handler().CronDownload(ctx, "00 17 * * 1-5", []int32{handlers.DailyClose, handlers.ThreePrimary})
	err = s.Handler().CronDownload(ctx, "00 19 * * 1-5", []int32{handlers.Concentration})

	// start gRPC server
	cfg := config.GetCurrentConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
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
