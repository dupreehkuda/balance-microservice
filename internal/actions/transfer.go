package actions

import (
	i "github.com/dupreehkuda/balance-microservice/internal"
	"github.com/shopspring/decimal"
)

func (a actions) TransferFunds(senderID, receiverID string, funds decimal.Decimal) error {
	exists := a.storage.CheckAccountExistence(senderID)
	if !exists {
		return i.ErrNoSuchUser
	}

	err := a.storage.TransferFunds(senderID, receiverID, funds)
	if err != nil {
		return err
	}

	return nil
}
