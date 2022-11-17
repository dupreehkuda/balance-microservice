package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

func TestHandlers_GetBalance(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		outputBody         Users
		akey               i.CtxKey
		accountID          string
		actionsReturnErr   error
		expectedStatusCode int
	}{
		{
			name:               "Get nonexistent user",
			outputBody:         Users{},
			akey:               "account",
			accountID:          "test_user",
			actionsReturnErr:   i.ErrNoSuchUser,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Get existing user",
			outputBody: Users{
				UserID:  "test_user",
				Current: decimal.Decimal{},
				OnHold:  decimal.Decimal{},
			},
			akey:               "account",
			accountID:          "test_user",
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "AccountID is empty",
			outputBody:         Users{},
			akey:               "account",
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Server error",
			outputBody:         Users{},
			akey:               "account",
			accountID:          "test_user",
			actionsReturnErr:   bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			body, _ := json.Marshal(testCase.outputBody)

			if testCase.outputBody.UserID == "" {
				body = nil
			}

			actions := &MockActions{
				GetBalanceFunc: func(accountID string) ([]byte, error) {
					return body, testCase.actionsReturnErr
				},
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodGet, "/api/balance/get", nil)
			res := httptest.NewRecorder()

			ctx := context.WithValue(req.Context(), testCase.akey, testCase.accountID)

			server.GetBalance(res, req.WithContext(ctx))

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
			a.Equal(body, res.Body.Bytes(), "Wrong body")
		})
	}
}
