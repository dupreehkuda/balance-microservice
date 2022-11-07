package internal

import "errors"

var (
	ErrWrongCredentials = errors.New("invalid order number")
	ErrNotEnoughFunds   = errors.New("not enough points")
	ErrNoSuchUser       = errors.New("there is no such user")
	ErrNoSuchOrder      = errors.New("there is no such order")
)
