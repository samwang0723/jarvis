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

package structuredlog

import (
	"log"
	"samwang0723/jarvis/config"
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
}

func initialize(l ILogger) {
	instance = l
	instance.Info("initialized logger")
}

func Logger(cfg *config.Config) ILogger {
	if instance == nil {
		var level logrus.Level
		switch cfg.Log.Level {
		case "FATAL":
			level = logrus.FatalLevel
		case "INFO":
			level = logrus.InfoLevel
		case "WARN":
			level = logrus.WarnLevel
		case "ERROR":
			level = logrus.ErrorLevel
		default:
			level = logrus.DebugLevel
		}
		slog := &structuredLogger{
			logger: logrus.New(),
		}
		slog.logger.SetLevel(level)
		initialize(slog)
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
