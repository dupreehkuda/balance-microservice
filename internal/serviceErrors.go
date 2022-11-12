package internal

import "errors"

var (
	ErrWrongCredentials   = errors.New("invalid order number")
	ErrNotEnoughFunds     = errors.New("not enough points")
	ErrNoSuchUser         = errors.New("there is no such user")
	ErrNoSuchOrder        = errors.New("there is no such order")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrOrderProcessed     = errors.New("order already processed")
	ErrNoData             = errors.New("no data for selection")
	ErrNoSuchReport       = errors.New("not existing report")
)
