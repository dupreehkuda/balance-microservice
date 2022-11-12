package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

// InitializeLogger initializes new logger instance
func InitializeLogger() *zap.Logger {
	Logger, _ = zap.NewProduction()
	return Logger
}
