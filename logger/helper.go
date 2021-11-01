package logger

func Debug(args ...interface{}) {
	instance.Debug(args...)
}

func Debugf(s string, args ...interface{}) {
	instance.Debugf(s, args...)
}

func Info(args ...interface{}) {
	instance.Info(args...)
}

func Infof(s string, args ...interface{}) {
	instance.Infof(s, args...)
}

func Warn(args ...interface{}) {
	instance.Warn(args...)
}

func Warnf(s string, args ...interface{}) {
	instance.Warnf(s, args...)
}

func Fatal(args ...interface{}) {
	instance.Fatal(args...)
}

func Fatalf(s string, args ...interface{}) {
	instance.Fatalf(s, args...)
}

func Error(args ...interface{}) {
	instance.Error(args...)
}

func Errorf(s string, args ...interface{}) {
	instance.Errorf(s, args...)
}

func Flush() {
	instance.Flush()
}
