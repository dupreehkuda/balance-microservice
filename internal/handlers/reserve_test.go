package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

func TestHandlers_ReserveFunds(t *testing.T) {
	a := assert.New(t)

	amount := decimal.NewFromFloat(50.50)
	users := []Users{{
		UserID:  "test_user_1",
		Current: decimal.NewFromFloat(100.50),
		OnHold:  decimal.Zero,
	}}
	expected := users
	expected[0] = Users{
		UserID:  users[0].UserID,
		Current: users[0].Current.Sub(amount),
		OnHold:  users[0].OnHold.Add(amount),
	}

	testCases := []struct {
		name               string
		inputBody          reserve
		expectedUsers      []Users
		actionsReturnErr   error
		expectedStatusCode int
	}{
		{
			name: "Wrong order num",
			inputBody: reserve{
				TargetID:  "test_user_1",
				ServiceID: "test_purpose",
				OrderID:   837429,
				Amount:    amount,
			},
			expectedUsers:      users,
			actionsReturnErr:   i.ErrWrongCredentials,
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name: "Successful reserve",
			inputBody: reserve{
				TargetID:  "test_user_1",
				ServiceID: "test_purpose",
				OrderID:   12345678903,
				Amount:    amount,
			},
			expectedUsers:      expected,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusCreated,
		},
		{
			name: "Duplicate order",
			inputBody: reserve{
				TargetID:  "test_user_1",
				ServiceID: "test_purpose",
				OrderID:   12345678903,
				Amount:    amount,
			},
			expectedUsers:      users,
			actionsReturnErr:   i.ErrOrderAlreadyExists,
			expectedStatusCode: http.StatusConflict,
		},
		{
			name: "Not enough funds",
			inputBody: reserve{
				TargetID:  "test_user_1",
				ServiceID: "test_purpose",
				OrderID:   9522120790,
				Amount:    decimal.NewFromFloat(200),
			},
			expectedUsers:      users,
			actionsReturnErr:   i.ErrNotEnoughFunds,
			expectedStatusCode: http.StatusPaymentRequired,
		},
		{
			name: "User not found",
			inputBody: reserve{
				TargetID:  "test_user_2",
				ServiceID: "test_purpose",
				OrderID:   3161292887,
				Amount:    decimal.NewFromFloat(200),
			},
			expectedUsers:      expected,
			actionsReturnErr:   i.ErrNoSuchUser,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "TargetID is empty",
			inputBody: reserve{
				TargetID:  "",
				ServiceID: "test_purpose",
				OrderID:   3161292887,
				Amount:    decimal.NewFromFloat(200),
			},
			expectedUsers:      expected,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Server error",
			inputBody: reserve{
				TargetID:  "test_user_1",
				ServiceID: "test_purpose",
				OrderID:   6805595219,
				Amount:    decimal.NewFromFloat(1),
			},
			expectedUsers:      expected,
			actionsReturnErr:   bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, _ := json.Marshal(testCase.inputBody)

			actions := &MockActions{
				ReserveFundsFunc: func(targetID, serviceID string, orderID int, funds decimal.Decimal) error {
					return testCase.actionsReturnErr
				},
				Users: users,
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodPost, "/api/order/reserve", bytes.NewReader(data))
			res := httptest.NewRecorder()

			server.ReserveFunds(res, req)

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code", testCase.name)
			a.Equal(testCase.expectedUsers, actions.Users, "Data doesn't match", testCase.name)
		})
	}
}
