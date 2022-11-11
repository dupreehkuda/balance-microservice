package storage

import (
	"context"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// AddFunds adds funds to account
func (s storage) AddFunds(accountID string, funds decimal.Decimal) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()
	pgxdecimal.Register(conn.Conn().TypeMap())

	_, err = conn.Exec(context.Background(), `insert into accounts (account_id, funds) values ($1, $2) on conflict (account_id)
	do update set funds = accounts.funds + $2 where accounts.account_id = $1;`, accountID, funds)
	if err != nil {
		s.logger.Error("Error occurred adding funds to account", zap.Error(err))
		return err
	}

	return nil
}
