package storage

import (
	"context"

	"go.uber.org/zap"
)

// DeleteUnprocessed deletes unprocessed orders reserved more than half an hour ago
func (s storage) DeleteUnprocessed() {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "DELETE FROM orders WHERE processed = false AND creation_date + INTERVAL '30 minutes' > CURRENT_TIMESTAMP;")
	if err != nil {
		s.logger.Error("Error while deleting unprocessed", zap.Error(err))
		return
	}
}
