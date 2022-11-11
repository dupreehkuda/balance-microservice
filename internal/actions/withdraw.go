package actions

import i "github.com/dupreehkuda/balance-microservice/internal"

// WithdrawBalance checks if order is processed and withdraws order
func (a actions) WithdrawBalance(orderID int) error {
	exists, processed := a.storage.CheckOrderProcessed(orderID)
	if !exists {
		return i.ErrNoSuchOrder
	}

	if processed {
		return i.ErrOrderProcessed
	}

	err := a.storage.WithdrawBalance(orderID)
	if err != nil {
		return err
	}

	return nil
}
