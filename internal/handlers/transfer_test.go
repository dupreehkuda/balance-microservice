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

func TestHandlers_TransferFunds(t *testing.T) {
	a := assert.New(t)

	users := []Users{{
		UserID:  "test_user_1",
		Current: decimal.NewFromFloat(100),
		OnHold:  decimal.Zero,
	}, {
		UserID:  "test_user_2",
		Current: decimal.Zero,
		OnHold:  decimal.Zero,
	}}

	want := users
	amount := decimal.NewFromFloat(100.10)
	want[0].Current = want[0].Current.Sub(amount)
	want[1].Current = want[1].Current.Add(amount)

	testCases := []struct {
		name               string
		inputBody          transfer
		actionsReturnErr   error
		expectedStatusCode int
		expectedUsers      []Users
	}{
		{
			name: "Transfer with lack of funds",
			inputBody: transfer{
				SenderID:   "test_user_1",
				ReceiverID: "test_user_2",
				Amount:     amount,
			},
			actionsReturnErr:   i.ErrNotEnoughFunds,
			expectedStatusCode: http.StatusPaymentRequired,
			expectedUsers:      users,
		},
		{
			name: "Successful transfer",
			inputBody: transfer{
				SenderID:   "test_user_1",
				ReceiverID: "test_user_2",
				Amount:     amount,
			},
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusOK,
			expectedUsers:      want,
		},
		{
			name: "Empty sender",
			inputBody: transfer{
				SenderID:   "",
				ReceiverID: "test_user_2",
				Amount:     amount,
			},
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedUsers:      users,
		},
		{
			name: "No user",
			inputBody: transfer{
				SenderID:   "test_user_3",
				ReceiverID: "test_user_2",
				Amount:     amount,
			},
			actionsReturnErr:   i.ErrNoSuchUser,
			expectedStatusCode: http.StatusBadRequest,
			expectedUsers:      users,
		},
		{
			name: "Server error",
			inputBody: transfer{
				SenderID:   "test_user_1",
				ReceiverID: "test_user_2",
				Amount:     amount,
			},
			actionsReturnErr:   bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
			expectedUsers:      users,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, _ := json.Marshal(testCase.inputBody)

			actions := &MockActions{
				TransferFundsFunc: func(senderID, receiverID string, funds decimal.Decimal) error {
					return testCase.actionsReturnErr
				},
				Users: users,
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodPost, "/api/balance/transfer", bytes.NewReader(data))
			res := httptest.NewRecorder()

			server.TransferFunds(res, req)

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
			a.Equal(testCase.expectedUsers, actions.Users, "Data doesn't match")
		})
	}
}
