package storage

import (
	"context"

	"go.uber.org/zap"
)

// WriteReport writes report to the database
func (s storage) WriteReport(repID, report string) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "insert into reports (report_id, report) values ($1, $2);", repID, report)
	if err != nil {
		s.logger.Error("Error while inserting new report", zap.Error(err))
		return err
	}

	return nil
}

// ReadReport reads and returns needed report
func (s storage) ReadReport(repID string) (string, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return "", err
	}
	defer conn.Release()

	var result string

	conn.QueryRow(context.Background(), "select report from reports where report_id = $1", repID).Scan(&result)

	return result, nil
}
