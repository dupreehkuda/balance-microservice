package handlers

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal/interfaces"
)

type handlers struct {
	storage i.Stored
	actions i.Actions
	logger  *zap.Logger
}

// New creates new instance of handlers
func New(storage i.Stored, processor i.Actions, logger *zap.Logger) *handlers {
	return &handlers{
		storage: storage,
		actions: processor,
		logger:  logger,
	}
}

// todo: write unit-tests
