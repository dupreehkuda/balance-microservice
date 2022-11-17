package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// WithdrawBalance gets order done and gets frozen funds from account
func (s storage) WithdrawBalance(orderID int) error {
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
	tx.Exec(ctx, "update accounts set on_hold = on_hold - $2 where account_id = $1;", aID, amount)
	tx.Exec(ctx, "update orders set processed = true, processed_date = $2 where order_id = $1;", orderID, time.Now())

	comment := fmt.Sprintf("Payed %s for order №%v", amount.String(), orderID)
	tx.Exec(ctx, "insert into history (operation, account_id, correspondent, funds, comment, processed_at) values ('payment', $1, $2, $3, $4, $5) on conflict do nothing;", aID, strconv.Itoa(orderID), amount, comment, time.Now())

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Error occurred committing tx", zap.Error(err))
		return err
	}

	return nil
}
