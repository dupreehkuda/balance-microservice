package handlers

import (
	"time"

	"github.com/shopspring/decimal"

	intl "github.com/dupreehkuda/balance-microservice/internal"
)

type Users struct {
	UserID  string          `json:"user_id"`
	Current decimal.Decimal `json:"current"`
	OnHold  decimal.Decimal `json:"on_hold"`
}

type Orders struct {
	OrderID       int
	ServiceID     string
	AccountID     string
	Amount        decimal.Decimal
	Processed     bool
	CreationDate  time.Time
	ProcessedDate time.Time
	Canceled      bool
}

type MockActions struct {
	AddFundsFunc        func(accountID string, funds decimal.Decimal) error
	GetBalanceFunc      func(accountID string) ([]byte, error)
	TransferFundsFunc   func(senderID, receiverID string, funds decimal.Decimal) error
	WithdrawBalanceFunc func(orderID int) error
	ReserveFundsFunc    func(targetID, serviceID string, orderID int, funds decimal.Decimal) error
	CancelReserveFunc   func(orderID int) error
	GetReportLinkFunc   func(month, year string) (string, error)
	GetReportFunc       func(reportID string) (string, error)
	GetHistoryFunc      func(targetId, sortOrder, sortBy string, quantity int) ([]byte, error)
	Users               []Users
	Orders              []Orders
}

func (m *MockActions) AddFunds(accountID string, funds decimal.Decimal) error {
	for _, val := range m.Users {
		if val.UserID == accountID {
			val.Current = val.Current.Add(funds)
			return m.AddFundsFunc(accountID, funds)
		}
	}

	m.Users = append(m.Users, Users{
		UserID:  accountID,
		Current: funds,
		OnHold:  decimal.Zero,
	})

	return m.AddFundsFunc(accountID, funds)
}

func (m *MockActions) GetBalance(accountID string) ([]byte, error) {
	return m.GetBalanceFunc(accountID)
}

func (m *MockActions) TransferFunds(senderID, receiverID string, funds decimal.Decimal) error {
	for _, val := range m.Users {
		if val.UserID == senderID {
			if val.Current.LessThan(funds) {
				return m.TransferFundsFunc(senderID, receiverID, funds)
			}

			val.Current = val.Current.Sub(funds)
		}
	}

	for _, val := range m.Users {
		if val.UserID == receiverID {
			val.Current = val.Current.Add(funds)
			return m.TransferFundsFunc(senderID, receiverID, funds)
		}
	}

	m.Users = append(m.Users, Users{
		UserID:  receiverID,
		Current: funds,
		OnHold:  decimal.Zero,
	})

	return m.TransferFundsFunc(senderID, receiverID, funds)
}

func (m *MockActions) ReserveFunds(targetID, serviceID string, orderID int, funds decimal.Decimal) error {
	if valid := intl.LuhnValid(orderID); !valid {
		return m.ReserveFundsFunc(targetID, serviceID, orderID, funds)
	}

	for _, val := range m.Orders {
		if val.OrderID == orderID {
			return m.ReserveFundsFunc(targetID, serviceID, orderID, funds)
		}
	}

	for i, val := range m.Users {
		if val.UserID == targetID {
			if val.Current.LessThan(funds) {
				return m.ReserveFundsFunc(targetID, serviceID, orderID, funds)
			}

			m.Users[i].Current = val.Current.Sub(funds)
			m.Users[i].OnHold = val.OnHold.Add(funds)

			m.Orders = append(m.Orders, Orders{
				OrderID:      orderID,
				ServiceID:    serviceID,
				AccountID:    targetID,
				Amount:       funds,
				Processed:    false,
				CreationDate: time.Now(),
			})

			return m.ReserveFundsFunc(targetID, serviceID, orderID, funds)
		}
	}

	return m.ReserveFundsFunc(targetID, serviceID, orderID, funds)
}

func (m *MockActions) WithdrawBalance(orderID int) error {
	if valid := intl.LuhnValid(orderID); !valid {
		return m.WithdrawBalanceFunc(orderID)
	}

	for i, oval := range m.Orders {
		if oval.OrderID == orderID {
			if oval.Processed || oval.Canceled {
				return m.WithdrawBalanceFunc(orderID)
			}

			for j, uval := range m.Users {
				if uval.UserID == oval.AccountID {
					m.Users[j].OnHold.Sub(oval.Amount)
				}

				m.Orders[i].Processed = true
				m.Orders[i].ProcessedDate = time.Now()

				return m.WithdrawBalanceFunc(orderID)
			}
		}
	}

	return m.WithdrawBalanceFunc(orderID)
}

func (m *MockActions) CancelReserve(orderID int) error {
	if valid := intl.LuhnValid(orderID); !valid {
		return m.CancelReserveFunc(orderID)
	}

	for i, oval := range m.Orders {
		if oval.OrderID == orderID {
			if oval.Processed || oval.Canceled {
				return m.CancelReserveFunc(orderID)
			}

			for j, uval := range m.Users {
				if uval.UserID == oval.AccountID {
					m.Users[j].OnHold.Sub(oval.Amount)
					m.Users[j].Current.Add(oval.Amount)
				}
			}

			m.Orders[i].Canceled = true

			return m.CancelReserveFunc(orderID)
		}
	}

	return m.CancelReserveFunc(orderID)
}

func (m *MockActions) GetReportLink(month, year string) (string, error) {
	return m.GetReportLinkFunc(month, year)
}

func (m *MockActions) GetReport(reportID string) (string, error) {
	return m.GetReportFunc(reportID)
}

func (m *MockActions) GetHistory(targetId, sortOrder, sortBy string, quantity int) ([]byte, error) {
	return m.GetHistoryFunc(targetId, sortOrder, sortBy, quantity)
}

func (m *MockActions) RunDeletion() {
}
