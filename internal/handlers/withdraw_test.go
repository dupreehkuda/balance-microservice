package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

func TestHandlers_WithdrawBalance(t *testing.T) {
	a := assert.New(t)

	users := []Users{{
		UserID:  "test_user_1",
		Current: decimal.NewFromFloat(301),
		OnHold:  decimal.Zero,
	}}
	orders := []Orders{
		{
			OrderID:      12345678903,
			ServiceID:    "test_purpose",
			AccountID:    "test_user_1",
			Amount:       decimal.NewFromFloat(100),
			Processed:    false,
			Canceled:     true,
			CreationDate: time.Now(),
		},
		{
			OrderID:      2657960791,
			ServiceID:    "test_purpose",
			AccountID:    "test_user_1",
			Amount:       decimal.NewFromFloat(100),
			Processed:    false,
			Canceled:     false,
			CreationDate: time.Now(),
		},
	}
	expected := users

	testCases := []struct {
		name               string
		inputBody          withdrawal
		expectedUsers      []Users
		actionsReturnErr   error
		expectedStatusCode int
	}{
		{
			name:               "Canceled order",
			inputBody:          withdrawal{OrderID: 12345678903},
			expectedUsers:      users,
			actionsReturnErr:   i.ErrOrderProcessed,
			expectedStatusCode: http.StatusNotAcceptable,
		},
		{
			name:               "Order does not exist",
			inputBody:          withdrawal{OrderID: 8636204631},
			expectedUsers:      users,
			actionsReturnErr:   i.ErrNoSuchOrder,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Successful withdraw",
			inputBody:          withdrawal{OrderID: 2657960791},
			expectedUsers:      expected,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "OrderID empty",
			inputBody:          withdrawal{OrderID: 0},
			expectedUsers:      users,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Server error",
			inputBody:          withdrawal{OrderID: 12345678903},
			expectedUsers:      users,
			actionsReturnErr:   bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, _ := json.Marshal(testCase.inputBody)

			actions := &MockActions{
				WithdrawBalanceFunc: func(orderID int) error {
					return testCase.actionsReturnErr
				},
				Users:  users,
				Orders: orders,
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodPost, "/api/order/withdraw", bytes.NewReader(data))
			res := httptest.NewRecorder()

			server.WithdrawBalance(res, req)

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
			a.Equal(testCase.expectedUsers, actions.Users, "Data doesn't match")
		})
	}
}
