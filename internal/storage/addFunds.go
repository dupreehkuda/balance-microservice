package storage

import (
	"context"
	"fmt"
	"time"

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

	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		s.logger.Error("Error occurred creating tx", zap.Error(err))
		return err
	}

	tx.Exec(context.Background(), `insert into accounts (account_id, funds) values ($1, $2) on conflict (account_id)
	do update set funds = accounts.funds + $2 where accounts.account_id = $1;`, accountID, funds)

	comment := fmt.Sprintf("Top-up %s", funds.String())
	tx.Exec(ctx, "insert into history (operation, account_id, correspondent, funds, comment, processed_at) values ('top-up', $1, $2, $3, $4, $5) on conflict do nothing;", accountID, "service", funds, comment, time.Now())

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Error occurred making reservation in db", zap.Error(err))
		return err
	}

	if err != nil {
		s.logger.Error("Error occurred adding funds to account", zap.Error(err))
		return err
	}

	return nil
}
