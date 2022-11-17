package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

func TestHandlers_GetReport(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		outputBody         string
		reportID           string
		actionsReturnErr   error
		expectedStatusCode int
	}{
		{
			name:               "Get nonexistent report",
			outputBody:         "",
			reportID:           "qojw4b",
			actionsReturnErr:   i.ErrNoSuchReport,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "ReportID is empty",
			outputBody:         "",
			reportID:           "",
			actionsReturnErr:   nil,
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Server error",
			outputBody:         "",
			reportID:           "fghj",
			actionsReturnErr:   bytes.ErrTooLarge,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		t.Run("Get nonexistent report", func(t *testing.T) {
			actions := &MockActions{
				GetReportFunc: func(reportID string) (string, error) {
					return testCase.outputBody, testCase.actionsReturnErr
				},
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodGet, "/api/accounting/get/", nil)
			res := httptest.NewRecorder()

			var rkey i.CtxKey = "report"
			ctx := context.WithValue(req.Context(), rkey, testCase.reportID)

			server.GetReport(res, req.WithContext(ctx))

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
			a.Equal(testCase.outputBody, res.Body.String(), "Wrong body")
		})
	}
}
