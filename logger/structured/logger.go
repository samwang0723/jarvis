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

type LogLevel int

const (
	Fatal LogLevel = iota
	Error
	Warn
	Debug
	Info
)

type structuredLogger struct {
	logger *logrus.Logger
	level  LogLevel
}

func initialize(l ILogger) {
	instance = l
	instance.Info("initialized logger")
}

func Logger(cfg *config.Config) ILogger {
	if instance == nil {
		var level LogLevel
		switch cfg.Log.Level {
		case "FATAL":
			level = Fatal
		case "DEBUG":
			level = Debug
		case "WARN":
			level = Warn
		case "ERROR":
			level = Error
		default:
			level = Info
		}
		slog := &structuredLogger{
			logger: logrus.New(),
			level:  level,
		}
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
	if l.level < Debug {
		return
	}
	l.logger.Debug(args...)
}

func (l *structuredLogger) Debugf(s string, args ...interface{}) {
	if l.level < Debug {
		return
	}

	l.logger.Debugf(s, args...)
}

func (l *structuredLogger) Info(args ...interface{}) {
	if l.level < Info {
		return
	}

	l.logger.Info(args...)
}

func (l *structuredLogger) Infof(s string, args ...interface{}) {
	if l.level < Info {
		return
	}

	l.logger.Infof(s, args...)
}

func (l *structuredLogger) Warn(args ...interface{}) {
	if l.level < Warn {
		return
	}

	l.logger.Warn(args...)
}

func (l *structuredLogger) Warnf(s string, args ...interface{}) {
	if l.level < Warn {
		return
	}

	l.logger.Warnf(s, args...)
}

func (l *structuredLogger) Fatal(args ...interface{}) {
	if l.level < Fatal {
		return
	}

	l.logger.Fatal(args...)
}

func (l *structuredLogger) Fatalf(s string, args ...interface{}) {
	if l.level < Fatal {
		return
	}

	l.logger.Fatalf(s, args...)
}

func (l *structuredLogger) Error(args ...interface{}) {
	if l.level < Error {
		return
	}

	l.logger.Error(args...)
}

func (l *structuredLogger) Errorf(s string, args ...interface{}) {
	if l.level < Error {
		return
	}

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
