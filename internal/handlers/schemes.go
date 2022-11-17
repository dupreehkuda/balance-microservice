package handlers

import "github.com/shopspring/decimal"

type reserve struct {
	TargetID  string          `json:"target_id"`
	ServiceID string          `json:"service_id"`
	OrderID   int             `json:"order_id"`
	Amount    decimal.Decimal `json:"amount"`
}

type transfer struct {
	SenderID   string          `json:"sender_id"`
	ReceiverID string          `json:"receiver_id"`
	Amount     decimal.Decimal `json:"amount"`
}

type addFunds struct {
	AccountID string          `json:"account_id"`
	Amount    decimal.Decimal `json:"amount"`
}

type withdrawal struct {
	OrderID int `json:"order_id"`
}

type accounting struct {
	Month string `json:"month"`
	Year  string `json:"year"`
}

type linkResponce struct {
	Date string `json:"date"`
	Link string `json:"link"`
}

type history struct {
	User      string `json:"userID"`
	SortBy    string `json:"sortBy"`
	SortOrder string `json:"sortOrder"`
	Quantity  int    `json:"quantity"`
}
