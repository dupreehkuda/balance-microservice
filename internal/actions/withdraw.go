package actions

func (a actions) WithdrawBalance(orderID int) error {
	err := a.storage.WithdrawBalance(orderID)
	if err != nil {
		return err
	}

	return nil
}
