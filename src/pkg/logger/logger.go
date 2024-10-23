package logger

import "go.uber.org/zap"

func InitLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer func() {
		if err := logger.Sync(); err != nil {
			panic(err)
		}
	}()
	return logger.Sugar()
}
