package storage

import (
	"context"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CheckAccountExistence checks if account is present
func (s storage) CheckAccountExistence(accountID string) bool {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false
	}
	defer conn.Release()

	var aID string

	conn.QueryRow(context.Background(), "select account_id from accounts where account_id = $1", accountID).Scan(&aID)

	if aID == "" {
		return false
	}

	return true
}

// CheckAccountBalance checks if account is present and have enough funds
func (s storage) CheckAccountBalance(accountID string, needed decimal.Decimal) (bool, bool) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false, false
	}
	defer conn.Release()
	pgxdecimal.Register(conn.Conn().TypeMap())

	var aID string
	var current decimal.Decimal

	conn.QueryRow(context.Background(), "select account_id, funds from accounts where account_id = $1", accountID).Scan(&aID, &current)
	s.logger.Debug("in checkers", zap.String("account", aID), zap.String("balance", current.String()))

	if aID == "" {
		return false, false
	}

	if current.LessThan(needed) {
		return true, false
	}

	return true, true
}

// CheckOrderExistence checks if order is present
func (s storage) CheckOrderExistence(orderID int) bool {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false
	}
	defer conn.Release()

	var aID string

	conn.QueryRow(context.Background(), "select account_id from orders where order_id = $1", orderID).Scan(&aID)

	if aID == "" {
		return false
	}

	return true
}

// CheckOrderProcessed checks if order is present and processed
func (s storage) CheckOrderProcessed(orderID int) (bool, bool) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return false, false
	}
	defer conn.Release()

	var aID string
	var processed bool

	conn.QueryRow(context.Background(), "select account_id, processed from orders where order_id = $1", orderID).Scan(&aID, &processed)

	if aID == "" {
		return false, false
	}

	if !processed {
		return true, false
	}

	return true, true
}
