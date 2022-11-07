package actions

import (
	i "github.com/dupreehkuda/balance-microservice/internal"
	"github.com/shopspring/decimal"
)

func (a actions) ReserveFunds(targetID, serviceID string, orderID int, funds decimal.Decimal) error {
	exists, enoughFunds := a.storage.CheckAccountBalance(targetID, funds)

	if !exists {
		return i.ErrNoSuchOrder
	}

	if !enoughFunds {
		return i.ErrNotEnoughFunds
	}

	err := a.storage.ReserveFunds(targetID, serviceID, orderID, funds)
	if err != nil {
		return err
	}

	return nil
}
