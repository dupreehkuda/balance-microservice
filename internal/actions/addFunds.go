package actions

import (
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

func (a actions) AddFunds(accountID string, funds decimal.Decimal) error {
	err := a.storage.AddFunds(accountID, funds)
	if err != nil {
		a.logger.Error("Error occurred in AddFunds call", zap.Error(err))
		return err
	}

	return nil
}
