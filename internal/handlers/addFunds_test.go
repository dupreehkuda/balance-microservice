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
)

func TestHandlers_AddFunds(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		inputBody          addFunds
		actionsReturn      error
		expectedStatusCode int
	}{
		{
			name: "Success",
			inputBody: addFunds{
				AccountID: "test_user",
				Amount:    decimal.NewFromFloat(100.10),
			},
			actionsReturn:      nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "AccountID empty",
			inputBody: addFunds{
				AccountID: "",
				Amount:    decimal.NewFromFloat(100.10),
			},
			actionsReturn:      nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Server error",
			inputBody: addFunds{
				AccountID: "test_user_2",
				Amount:    decimal.NewFromFloat(100.10),
			},
			actionsReturn:      bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, _ := json.Marshal(testCase.inputBody)

			actions := &MockActions{
				AddFundsFunc: func(accountID string, funds decimal.Decimal) error {
					return testCase.actionsReturn
				},
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodPost, "/api/balance/add", bytes.NewReader(data))
			res := httptest.NewRecorder()

			server.AddFunds(res, req)

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
		})
	}
}
