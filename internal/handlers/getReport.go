package handlers

import (
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetReport gets reportID from context and passes data to get report
func (h handlers) GetReport(w http.ResponseWriter, r *http.Request) {
	var key i.CtxKey = "report"
	reportID := r.Context().Value(key).(string)

	if reportID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.actions.GetReport(reportID)

	switch err {
	case i.ErrNoSuchReport:
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		w.Write([]byte(res))
		return
	default:
		h.logger.Error("Error call to actions for getting", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
