package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

// InitLogger initializes the logger.
// - `isProduction`: Controls whether to use production or development settings.
// - `logPath`: If non-empty, logs are written to the specified file in addition to stdout.
func InitLogger(isProduction bool, logPath string) {
	var core zapcore.Core

	// Encoder settings based on the environment
	var encoder zapcore.Encoder
	if isProduction {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// Default core writes to stdout
	consoleSyncer := zapcore.AddSync(os.Stdout)
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, consoleSyncer, zapcore.InfoLevel),
	}

	// If a logPath is provided, add file logging
	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic("failed to open log file: " + err.Error())
		}
		fileSyncer := zapcore.AddSync(file)
		cores = append(cores, zapcore.NewCore(encoder, fileSyncer, zapcore.InfoLevel))
	}

	// Combine cores
	core = zapcore.NewTee(cores...)

	// Build the logger
	zapLogger := zap.New(core)
	logger = zapLogger.Sugar()
}

// Sync flushes any buffered log entries.
// Should be called before the application exits.
func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

// Info logs an informational message.
func Info(format string, args ...interface{}) {
	// Use Sprintf to format the message if arguments are provided
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}
	logger.Info(message)
}

// Warn logs a warning message.
func Warn(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

// Error logs an error message.
func Error(format string, args ...interface{}) {
	// Use Sprintf to format the message if arguments are provided
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}
	logger.Errorw(message)
}

// Debug logs a debug message.
func Debug(format string, args ...interface{}) {
	// Use Sprintf to format the message if arguments are provided
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}
	logger.Debugw(message)
}

// Fatal logs a fatal message and exits the application.
func Fatal(format string, args ...interface{}) {
	// Use Sprintf to format the message if arguments are provided
	message := format
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	}
	logger.Fatalw(message)
}
