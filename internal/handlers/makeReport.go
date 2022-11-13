package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetReportLink passes time data to form and get accounting report
func (h handlers) GetReportLink(w http.ResponseWriter, r *http.Request) {
	var data accounting

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.Month == "" || data.Year == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result, err := h.actions.GetReportLink(data.Month, data.Year)

	switch err {
	case i.ErrNoData:
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		w.WriteHeader(http.StatusOK)

		body, err := json.Marshal(linkResponce{
			Date: fmt.Sprintf("%s-%s", data.Month, data.Year),
			Link: fmt.Sprintf("http://%s/accounting/get/%s", r.Host, result),
		})

		if err != nil {
			h.logger.Error("Unable to marshal JSON", zap.Error(err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write(body)
		return
	default:
		h.logger.Error("Unable to get report", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
