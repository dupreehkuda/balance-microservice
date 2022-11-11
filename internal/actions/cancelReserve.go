package actions

import i "github.com/dupreehkuda/balance-microservice/internal"

// CancelReserve checks if order is processed and cancels order
func (a actions) CancelReserve(orderID int) error {
	exists, processed := a.storage.CheckOrderProcessed(orderID)
	if !exists {
		return i.ErrNoSuchOrder
	}

	if processed {
		return i.ErrOrderProcessed
	}

	err := a.storage.CancelReserve(orderID)
	if err != nil {
		return err
	}

	return nil
}
