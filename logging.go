package disgord

import (
	"github.com/ImVexed/disgord/internal/logger"
	"go.uber.org/zap"
)

// Logger super basic logging interface
type Logger = logger.Logger

func DefaultLogger(debug bool) *logger.LoggerZap {
	return logger.DefaultLogger(debug)
}

func DefaultLoggerWithInstance(log *zap.Logger) *logger.LoggerZap {
	return logger.DefaultLoggerWithInstance(log)
}
