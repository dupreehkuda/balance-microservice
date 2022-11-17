package actions

import (
	"github.com/shopspring/decimal"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// ReserveFunds checks for luhn and enough funds, then creates order
func (a actions) ReserveFunds(targetID, serviceID string, orderID int, funds decimal.Decimal) error {
	valid := i.LuhnValid(orderID)
	if !valid {
		return i.ErrWrongCredentials
	}

	exists, enoughFunds := a.storage.CheckAccountBalance(targetID, funds)

	if !exists {
		return i.ErrNoSuchOrder
	}

	if !enoughFunds {
		return i.ErrNotEnoughFunds
	}

	exists = a.storage.CheckOrderExistence(orderID)
	if exists {
		return i.ErrOrderAlreadyExists
	}

	err := a.storage.ReserveFunds(targetID, serviceID, orderID, funds)
	if err != nil {
		return err
	}

	return nil
}
