package interfaces

import (
	"net/http"

	"github.com/shopspring/decimal"
)

type Handlers interface {
	AddFunds(w http.ResponseWriter, r *http.Request)
	GetBalance(w http.ResponseWriter, r *http.Request)
	CancelReserve(w http.ResponseWriter, r *http.Request)
	ReserveFunds(w http.ResponseWriter, r *http.Request)
	TransferFunds(w http.ResponseWriter, r *http.Request)
	WithdrawBalance(w http.ResponseWriter, r *http.Request)
}

type Stored interface {
	AddFunds(accountID string, funds decimal.Decimal) error
	GetBalance(accountID string) ([]byte, error)
	CheckAccountExistence(accountID string) bool
	CheckAccountBalance(accountID string, needed decimal.Decimal) (bool, bool)
	ReserveFunds(targetID, serviceID string, orderID int, funds decimal.Decimal) error
	TransferFunds(senderID, receiverID string, funds decimal.Decimal) error
	WithdrawBalance(orderID int) error
	CancelReserve(orderID int) error
}

type Actions interface {
	AddFunds(accountID string, funds decimal.Decimal) error
	GetBalance(accountID string) ([]byte, error)
	TransferFunds(senderID, receiverID string, funds decimal.Decimal) error
	WithdrawBalance(orderID int) error
	ReserveFunds(targetID, serviceID string, orderID int, funds decimal.Decimal) error
	CancelReserve(orderID int) error
}

type Middleware interface {
	CheckCompression(next http.Handler) http.Handler
	WriteCompressed(next http.Handler) http.Handler
	AccountCtx(next http.Handler) http.Handler
}
