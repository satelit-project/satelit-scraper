package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A logger which writes given message to stderr.
//
// It's safe to call all methods on a nil receiver.
type Logger struct {
	inner *zap.SugaredLogger
}

// Creates and returns new logger instance
func NewLogger() (*Logger, error) {
	logger, err := makeLogger()
	if err != nil {
		return nil, err
	}

	return &Logger{logger.Sugar()}, nil
}

// Redirects standard logger's output to current logger instance.
func (l *Logger) CaptureSTDLog() error {
	if !l.canSafeExec() {
		return nil
	}

	l.Sync()
	logger, err := makeLogger()
	if err != nil {
		return err
	}

	minLevel := int8(zapcore.InfoLevel)
	maxLevel := int8(zapcore.FatalLevel)
	for i := minLevel; i <= maxLevel; i++ {
		if _, err := zap.RedirectStdLogAt(logger, zapcore.Level(i)); err != nil {
			return err
		}
	}

	l.inner = logger.Sugar()
	return nil
}

// Adds a variadic number of fields to the logging context. The first value
// will become a key and the second one will become a value.
func (l *Logger) With(args ...interface{}) *Logger {
	var pt *Logger
	l.safeExec(func() {
		inner := l.inner.With(args...)
		pt = &Logger{inner}

	})

	return pt
}

// Flushes all bufferred log entries.
func (l *Logger) Sync() error {
	if !l.canSafeExec() {
		return nil
	}

	return l.inner.Sync()
}

// Logs formatted debug message.
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Debugf(template, args...)
	})
}

// Logs formatted info message.
func (l *Logger) Infof(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Infof(template, args...)
	})
}

// Logs formatted warning message.
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Warnf(template, args...)
	})
}

// Logs formatted error message.
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Errorf(template, args...)
	})
}

// Logs formatted error message and kills current process.
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.safeExec(func() {
		l.inner.Fatalf(template, args...)
	})
}

func (l *Logger) canSafeExec() bool {
	return l != nil && l.inner != nil
}

// Executes given closure only if receiver and inner loggers are not nil.
func (l *Logger) safeExec(f func()) {
	if l.canSafeExec() {
		f()
	}
}

// Creates and returns new Uber's logger instance.
func makeLogger() (*zap.Logger, error) {
	return zap.NewDevelopment(zap.AddCallerSkip(3))
}
