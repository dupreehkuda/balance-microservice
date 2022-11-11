package storage

import "github.com/shopspring/decimal"

type getResponse struct {
	UserId  string          `json:"user_id"`
	Current decimal.Decimal `json:"current"`
	OnHold  decimal.Decimal `json:"on_hold"`
}
