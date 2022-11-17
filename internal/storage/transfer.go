package storage

import (
	"context"
	"fmt"
	"time"

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

	forSender := fmt.Sprintf("Sent %s to %s", funds.String(), receiverID)
	forReceiver := fmt.Sprintf("Recieved %s from %s", funds.String(), senderID)

	ctx := context.Background()
	tx, err := conn.Begin(ctx)
	if err != nil {
		s.logger.Error("Error occurred creating tx", zap.Error(err))
		return err
	}

	tx.Exec(ctx, "update accounts set funds = funds - $2 where account_id = $1", senderID, funds)
	tx.Exec(ctx, `insert into accounts (account_id, funds) values ($1, $2) on conflict (account_id)
					do update set funds = accounts.funds + $2 where accounts.account_id = $1;`, receiverID, funds)
	tx.Exec(ctx, "insert into history (operation, account_id, correspondent, funds, comment, processed_at) values ('transfer', $1, $2, $3, $4, $5) on conflict do nothing;", receiverID, senderID, funds, forReceiver, time.Now())
	tx.Exec(ctx, "insert into history (operation, account_id, correspondent, funds, comment, processed_at) values ('transfer', $1, $2, $3, $4, $5) on conflict do nothing;", senderID, receiverID, funds, forSender, time.Now())

	err = tx.Commit(ctx)
	if err != nil {
		s.logger.Error("Error occurred making reservation in db", zap.Error(err))
		return err
	}

	return nil
}
