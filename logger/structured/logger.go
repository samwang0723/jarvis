package structuredlog

import (
	"log"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var (
	instance ILogger
)

type ILogger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(s string, args ...interface{})
	Infof(s string, args ...interface{})
	Warnf(s string, args ...interface{})
	Errorf(s string, args ...interface{})
	Fatalf(s string, args ...interface{})
	Flush()
}

type structuredLogger struct {
	logger *logrus.Logger
	level  logrus.Level
}

func initialize(l ILogger) {
	instance = l
	instance.Info("initialized logger")
}

func Logger() ILogger {
	if instance == nil {
		resp := &structuredLogger{
			logger: logrus.New(),
			level:  logrus.InfoLevel,
		}
		resp.logger.SetLevel(resp.level)
		initialize(resp)
		initSentry()
	}
	return instance
}

func initSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://f3fb4890176c442aafef411fcf812312@o1049557.ingest.sentry.io/6030819",
		Environment: "development",
		// Specify a fixed sample rate:
		TracesSampleRate: 1.0,
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
}

func (l *structuredLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *structuredLogger) Debugf(s string, args ...interface{}) {
	l.logger.Debugf(s, args...)
}

func (l *structuredLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *structuredLogger) Infof(s string, args ...interface{}) {
	l.logger.Infof(s, args...)
}

func (l *structuredLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *structuredLogger) Warnf(s string, args ...interface{}) {
	l.logger.Warnf(s, args...)
}

func (l *structuredLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *structuredLogger) Fatalf(s string, args ...interface{}) {
	l.logger.Fatalf(s, args...)
}

func (l *structuredLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *structuredLogger) Errorf(s string, args ...interface{}) {
	l.logger.Errorf(s, args...)
}

func (l *structuredLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *structuredLogger) Panicf(s string, args ...interface{}) {
	l.logger.Panicf(s, args...)
}

func (log *structuredLogger) Flush() {
	// Flush buffered events before the program terminates.
	sentry.Flush(2 * time.Second)
}
