package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

func TestHandlers_GetHistory(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		inputBody          history
		outputBody         []byte
		actionsReturnErr   error
		expectedStatusCode int
	}{
		{
			name: "No such user",
			inputBody: history{
				User:      "test_user_2",
				SortBy:    "date",
				SortOrder: "desc",
				Quantity:  1,
			},
			outputBody:         nil,
			actionsReturnErr:   i.ErrNoSuchUser,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Server error",
			inputBody: history{
				User:      "test_user_2",
				SortBy:    "date",
				SortOrder: "desc",
				Quantity:  1,
			},
			outputBody:         nil,
			actionsReturnErr:   bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "No data",
			inputBody: history{
				User:      "test_user_1",
				SortBy:    "date",
				SortOrder: "desc",
				Quantity:  1,
			},
			outputBody:         nil,
			actionsReturnErr:   i.ErrNoData,
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Quantity equals 0",
			inputBody: history{
				User:      "test_user_1",
				SortBy:    "date",
				SortOrder: "desc",
				Quantity:  0,
			},
			outputBody:         nil,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Bad input sortBy",
			inputBody: history{
				User:      "test_user_1",
				SortBy:    "",
				SortOrder: "desc",
				Quantity:  0,
			},
			outputBody:         nil,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Bad input sortOrder",
			inputBody: history{
				User:      "test_user_1",
				SortBy:    "date",
				SortOrder: "des",
				Quantity:  0,
			},
			outputBody:         nil,
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, _ := json.Marshal(testCase.inputBody)

			actions := &MockActions{
				GetHistoryFunc: func(targetId, sortOrder, sortBy string, quantity int) ([]byte, error) {
					return testCase.outputBody, testCase.actionsReturnErr
				},
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodGet, "/api/history/", bytes.NewReader(data))
			res := httptest.NewRecorder()

			server.GetHistory(res, req)

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
			a.Equal(testCase.outputBody, res.Body.Bytes(), "Wrong body")
		})
	}
}
