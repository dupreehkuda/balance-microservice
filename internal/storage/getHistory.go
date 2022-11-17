package storage

import (
	"context"
	"encoding/json"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"go.uber.org/zap"
)

// GetHistory gets user`s operations history
func (s storage) GetHistory(targetID, params string) ([]byte, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()
	pgxdecimal.Register(conn.Conn().TypeMap())

	resp := []historyData{}

	rows, err := conn.Query(context.Background(), "select operation, correspondent, funds, comment, processed_at from history where account_id = $1 "+params, targetID)
	if err != nil {
		s.logger.Error("Error while executing history query", zap.Error(err))
		return nil, err
	}

	for rows.Next() {
		r := historyData{}
		rows.Scan(&r.Operation, &r.Correspondent, &r.Funds, &r.Comment, &r.ProcessedAt)
		resp = append(resp, r)
	}

	body, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error occurred while marshalling body", zap.Error(err))
		return nil, err
	}

	return body, nil
}
