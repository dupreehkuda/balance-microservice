package storage

import (
	"time"

	"github.com/shopspring/decimal"
)

type getResponse struct {
	UserId  string          `json:"user_id"`
	Current decimal.Decimal `json:"current"`
	OnHold  decimal.Decimal `json:"on_hold"`
}

type accountingData struct {
	ServiceID string
	Sum       decimal.Decimal
}

type historyData struct {
	Operation     string          `json:"operation"`
	Correspondent string          `json:"correspondent"`
	Funds         decimal.Decimal `json:"funds"`
	Comment       string          `json:"comment"`
	ProcessedAt   time.Time       `json:"processedAt"`
}
