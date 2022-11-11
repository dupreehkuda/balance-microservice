package storage

import (
	"context"
	"encoding/json"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// GetBalance return current info
func (s storage) GetBalance(accountID string) ([]byte, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return nil, err
	}
	defer conn.Release()
	pgxdecimal.Register(conn.Conn().TypeMap())

	var aID string
	var current decimal.Decimal
	var onHold decimal.Decimal

	conn.QueryRow(context.Background(), "select account_id, funds, on_hold from accounts where account_id = $1", accountID).Scan(&aID, &current, &onHold)

	resp := getResponse{
		UserId:  aID,
		Current: current,
		OnHold:  onHold,
	}

	body, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error("Error occurred while marshalling body", zap.Error(err))
		return nil, err
	}

	return body, nil
}
