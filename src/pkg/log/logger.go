package logger

import "go.uber.org/zap"

func InitLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	return logger.Sugar()
}
