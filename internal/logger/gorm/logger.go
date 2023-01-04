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

package gormlog

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

//nolint:nolintlint, gochecknoglobals
var instance ILogger

type ILogger interface {
	LogMode(gormlogger.LogLevel) gormlogger.Interface
	Info(ctx context.Context, s string, args ...interface{})
	Warn(ctx context.Context, s string, args ...interface{})
	Error(ctx context.Context, s string, args ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
}

type gormLogger struct {
	logger                    *logrus.Logger
	SlowThreshold             time.Duration
	SourceField               string
	IgnoreRecordNotFoundError bool
	Colorful                  bool
	LogLevel                  gormlogger.LogLevel
}

func initialize(l ILogger) {
	instance = l
}

func Logger() ILogger {
	if instance == nil {
		resp := &gormLogger{
			logger:                    logrus.New(),
			SlowThreshold:             time.Second,
			SourceField:               "",
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
			LogLevel:                  gormlogger.Info,
		}
		initialize(resp)
	}

	return instance
}

func (l *gormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logger.WithContext(ctx).Infof(s, args...)
}

func (l *gormLogger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logger.WithContext(ctx).Warnf(s, args...)
}

func (l *gormLogger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logger.WithContext(ctx).Errorf(s, args...)
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}

	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}

	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) {
		fields[logrus.ErrorKey] = err
		l.logger.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)

		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logger.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)

		return
	}

	l.logger.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}
