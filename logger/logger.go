package logger

import (
	"io"
	"time"

	"github.com/sirupsen/logrus"
	gormlogger "gorm.io/gorm/logger"
)

var (
	logImpl *logger
)

type logger struct {
	logrus                    *logrus.Logger
	SlowThreshold             time.Duration
	SourceField               string
	IgnoreRecordNotFoundError bool
	Colorful                  bool
	LogLevel                  gormlogger.LogLevel
}

func initialize(l *logger) {
	logImpl = l
	Info("initialized logger")
}

func Logger() *logger {
	if logImpl == nil {
		resp := &logger{
			SlowThreshold:             time.Second,      // Slow SQL threshold
			LogLevel:                  gormlogger.Error, // Log level
			IgnoreRecordNotFoundError: true,             // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Disable color
			logrus:                    logrus.New(),
		}
		initialize(resp)
	}
	return logImpl
}

func UpdateConfig(output io.Writer, level logrus.Level, caller bool) {
	l := Logger().logrus
	l.SetOutput(output)
	l.SetLevel(level)
	l.SetReportCaller(caller)
}
