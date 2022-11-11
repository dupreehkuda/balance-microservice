package storage

import (
	"context"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// TransferFunds transfers funds from sender to receiver
func (s storage) TransferFunds(senderID, receiverID string, funds decimal.Decimal) error {
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

	tx.Exec(ctx, "update accounts set funds = funds - $2 where account_id = $1", senderID, funds)
	tx.Exec(ctx, `insert into accounts (account_id, funds) values ($1, $2) on conflict (account_id)
					do update set funds = accounts.funds + $2 where accounts.account_id = $1;`, receiverID, funds)

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Error occurred making reservation in db", zap.Error(err))
		return err
	}

	return nil
}
