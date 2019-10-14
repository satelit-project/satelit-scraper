package logging

import "go.uber.org/zap"

var defaultLogger *zap.Logger

func init() {
	logger, err := newLogger()
	if err != nil {
		panic(err)
	}

	// TODO: for prod redirect to file
	defaultLogger = logger
}

func DefaultLogger() *zap.SugaredLogger {
	return defaultLogger.Sugar()
}

func newLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig() // TODO: get from env
	return cfg.Build()
}
