package storage

import (
	"context"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// CancelReserve cancels order and unfreezes balance
func (s storage) CancelReserve(orderID int) error {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return err
	}
	defer conn.Release()
	pgxdecimal.Register(conn.Conn().TypeMap())

	var aID string
	var amount decimal.Decimal

	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		s.logger.Error("Error occurred creating tx", zap.Error(err))
		return err
	}

	tx.QueryRow(ctx, "select account_id, amount from orders where order_id = $1;", orderID).Scan(&aID, &amount)
	tx.Exec(ctx, "update accounts set on_hold = on_hold - $2, funds = funds + $2 where account_id = $1;", aID, amount)
	tx.Exec(ctx, "delete from orders where order_id = $1;", orderID)

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Error occurred committing tx", zap.Error(err))
		return err
	}

	return nil
}
