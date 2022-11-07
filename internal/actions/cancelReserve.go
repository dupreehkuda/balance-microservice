package actions

func (a actions) CancelReserve(orderID int) error {
	err := a.storage.CancelReserve(orderID)
	if err != nil {
		return err
	}

	return nil
}
