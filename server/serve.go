package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"samwang0723/jarvis/db"
	"samwang0723/jarvis/db/dal"
	"samwang0723/jarvis/handlers"
	"samwang0723/jarvis/logger"
	structuredlog "samwang0723/jarvis/logger/structured"
	"samwang0723/jarvis/services"
)

type IServer interface {
	Name() string
	Logger() structuredlog.ILogger
	Run(context.Context) error
	Start() error
	Stop() error
}

type server struct {
	opts Options
}

func Serve() {
	//TODO: Load configuration
	config := &db.Config{
		User:     "jarvis",
		Password: "password",
		Host:     "tcp(localhost:3306)",
		Database: "jarvis",
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
			sqld, _ := db.DB()
			return sqld.Close()
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

func (s *server) Start() error {
	for _, fn := range s.opts.BeforeStart {
		if err := fn(); err != nil {
			return err
		}
	}

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

	if err := s.Start(); err != nil {
		return err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	select {
	case <-c:
	case <-childCtx.Done():
	}

	return s.Stop()
}

func (s *server) Logger() structuredlog.ILogger {
	return s.opts.Logger
}

func (s *server) Name() string {
	return s.opts.Name
}
