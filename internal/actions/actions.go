package actions

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal/interfaces"
)

type actions struct {
	storage i.Stored
	logger  *zap.Logger
}

// New creates new instance of actions
func New(storage i.Stored, logger *zap.Logger) *actions {
	return &actions{
		storage: storage,
		logger:  logger,
	}
}
