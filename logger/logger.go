package logger

import (
	structuredlog "samwang0723/jarvis/logger/structured"
)

var (
	instance structuredlog.ILogger
)

func Initialize(l structuredlog.ILogger) {
	if instance == nil {
		instance = l
	}
}
