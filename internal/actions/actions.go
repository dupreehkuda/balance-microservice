package actions

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal/interfaces"
)

type actions struct {
	storage      i.Stored
	logger       *zap.Logger
	StopDeletion chan struct{}
}

// New creates new instance of actions
func New(storage i.Stored, logger *zap.Logger) *actions {
	return &actions{
		storage:      storage,
		logger:       logger,
		StopDeletion: make(chan struct{}),
	}
}
