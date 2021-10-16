package logger

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func (l *logger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *logger) Info(ctx context.Context, s string, args ...interface{}) {
	l.logrus.WithContext(ctx).Infof(s, args...)
}

func (l *logger) Warn(ctx context.Context, s string, args ...interface{}) {
	l.logrus.WithContext(ctx).Warnf(s, args...)
}

func (l *logger) Error(ctx context.Context, s string, args ...interface{}) {
	l.logrus.WithContext(ctx).Errorf(s, args...)
}

func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, _ := fc()
	fields := logrus.Fields{}
	if l.SourceField != "" {
		fields[l.SourceField] = utils.FileWithLineNum()
	}
	if err != nil && !(errors.Is(err, gorm.ErrRecordNotFound) && l.IgnoreRecordNotFoundError) {
		fields[logrus.ErrorKey] = err
		l.logrus.WithContext(ctx).WithFields(fields).Errorf("%s [%s]", sql, elapsed)
		return
	}

	if l.SlowThreshold != 0 && elapsed > l.SlowThreshold {
		l.logrus.WithContext(ctx).WithFields(fields).Warnf("%s [%s]", sql, elapsed)
		return
	}

	l.logrus.WithContext(ctx).WithFields(fields).Debugf("%s [%s]", sql, elapsed)
}

func Debug(args ...interface{}) {
	Logger().logrus.Debug(args...)
}

func Debugf(s string, args ...interface{}) {
	Logger().logrus.Debugf(s, args...)
}

func Info(args ...interface{}) {
	Logger().logrus.Info(args...)
}

func Infof(s string, args ...interface{}) {
	Logger().logrus.Infof(s, args...)
}

func Warn(args ...interface{}) {
	Logger().logrus.Warn(args...)
}

func Warnf(s string, args ...interface{}) {
	Logger().logrus.Warnf(s, args...)
}

func Fatal(args ...interface{}) {
	Logger().logrus.Fatal(args...)
}

func Fatalf(s string, args ...interface{}) {
	Logger().logrus.Fatalf(s, args...)
}

func Error(args ...interface{}) {
	Logger().logrus.Error(args...)
}

func Errorf(s string, args ...interface{}) {
	Logger().logrus.Errorf(s, args...)
}

func Panic(args ...interface{}) {
	Logger().logrus.Panic(args...)
}

func Panicf(s string, args ...interface{}) {
	Logger().logrus.Panicf(s, args...)
}
