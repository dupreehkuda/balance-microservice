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

func TestHandlers_GetReportLink(t *testing.T) {
	a := assert.New(t)

	testCases := []struct {
		name               string
		inputBody          accounting
		outputBody         linkResponce
		actionsReturnErr   error
		actionsReturnHash  string
		expectedStatusCode int
	}{
		{
			name: "Get empty report",
			inputBody: accounting{
				Month: "01",
				Year:  "2023",
			},
			outputBody:         linkResponce{},
			actionsReturnErr:   i.ErrNoData,
			actionsReturnHash:  "",
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "Month is empty",
			inputBody: accounting{
				Month: "",
				Year:  "2023",
			},
			outputBody:         linkResponce{},
			actionsReturnErr:   nil,
			actionsReturnHash:  "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Server error",
			inputBody: accounting{
				Month: "11",
				Year:  "2022",
			},
			outputBody:         linkResponce{},
			actionsReturnErr:   bytes.ErrTooLarge,
			actionsReturnHash:  "",
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name: "Successful make report",
			inputBody: accounting{
				Month: "11",
				Year:  "2022",
			},
			outputBody: linkResponce{
				Date: "11-2022",
				Link: "http://example.com/accounting/get/163bd0a74e",
			},
			actionsReturnErr:   nil,
			actionsReturnHash:  "163bd0a74e",
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			data, _ := json.Marshal(testCase.inputBody)
			resp, _ := json.Marshal(testCase.outputBody)

			if testCase.outputBody.Link == "" {
				resp = nil
			}

			actions := &MockActions{
				GetReportLinkFunc: func(month, year string) (string, error) {
					return testCase.actionsReturnHash, testCase.actionsReturnErr
				},
			}

			zp, _ := zap.NewDevelopment()
			server := New(actions, zp)

			req := httptest.NewRequest(http.MethodPost, "/api/accounting/add", bytes.NewReader(data))
			res := httptest.NewRecorder()

			server.GetReportLink(res, req)

			a.Equal(testCase.expectedStatusCode, res.Result().StatusCode, "Wrong status code")
			a.Equal(bytes.NewBuffer(resp), res.Body, "Wrong body")
		})
	}
}
