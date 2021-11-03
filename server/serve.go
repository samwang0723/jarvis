package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"samwang0723/jarvis/db"
	"samwang0723/jarvis/db/dal"
	"samwang0723/jarvis/dto"
	"samwang0723/jarvis/handlers"
	"samwang0723/jarvis/logger"
	structuredlog "samwang0723/jarvis/logger/structured"
	"samwang0723/jarvis/services"
)

const (
	gracefulShutdownPeriod = 5 * time.Second
)

type IServer interface {
	Name() string
	Logger() structuredlog.ILogger
	Handler() handlers.IHandler
	Run(context.Context) error
	Start(context.Context) error
	Stop() error
}

type server struct {
	opts Options
}

func Serve() {
	//TODO: Load configuration
	config := &db.Config{
		User:         "jarvis",
		Password:     "password",
		Host:         "tcp(localhost:3306)",
		Database:     "jarvis",
		MaxLifetime:  10,
		MaxIdleConns: 20,
		MaxOpenConns: 800,
	}
	db := db.GormFactory(config)
	dalService := dal.New(dal.WithDB(db))
	dataService := services.New(services.WithDAL(dalService))
	handler := handlers.New(dataService)

	s := newServer(
		Name("jarvis"),
		Logger(structuredlog.Logger()),
		Handler(handler),
		BeforeStop(func() error {
			sqlDB, _ := db.DB()
			return sqlDB.Close()
		}),
	)
	logger.Initialize(s.Logger())
	err := s.Run(context.Background())
	if err != nil && s.Logger() != nil {
		s.Logger().Errorf("error returned by service.Run(): %s\n", err.Error())
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

	go func() {
		s.Handler().BatchingDownload(ctx, &dto.DownloadRequest{
			RewindLimit: 20,
			RateLimit:   5000,
		})
	}()

	return nil
}

func (s *server) Stop() error {
	for _, fn := range s.opts.BeforeStop {
		if err := fn(); err != nil {
			return err
		}
	}

	return nil
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
		s.Logger().Warn("singal interrupt")
		cancel()
	case <-childCtx.Done():
		s.Logger().Warn("main context being cancelled")
	}
	<-time.After(gracefulShutdownPeriod)
	s.Logger().Warn("server being gracefully shuted down")

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
