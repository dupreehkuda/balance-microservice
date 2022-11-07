package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// WithdrawBalance commits early created reserve request
func (h handlers) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	var data withdrawal

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.OrderID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actions.WithdrawBalance(data.OrderID)
	if err != nil {
		h.logger.Error("Error call to actions for withdraw", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
