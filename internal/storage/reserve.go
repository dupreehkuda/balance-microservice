package storage

import (
	"context"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// ReserveFunds creates order and freezes balance
func (s storage) ReserveFunds(targetID, serviceID string, orderID int, funds decimal.Decimal) error {
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

	_, err = tx.Exec(ctx, "insert into orders (order_id, service_id, account_id, amount, creation_date) values ($1, $2, $3, $4, $5) on conflict do nothing;", orderID, serviceID, targetID, funds, time.Now())
	if err != nil {
		s.logger.Debug("first reserve exec", zap.Error(err))
	}

	_, err = tx.Exec(ctx, "update accounts set funds = funds - $2, on_hold = on_hold + $2 where account_id = $1;", targetID, funds)
	if err != nil {
		s.logger.Debug("second reserve exec", zap.Error(err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Error occurred making reservation in db", zap.Error(err))
		return err
	}

	return nil
}
