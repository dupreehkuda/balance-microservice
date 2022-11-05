package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// AddFunds passes account data to add funds
func (h handlers) AddFunds(w http.ResponseWriter, r *http.Request) {
	var data addFunds

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.AccountID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actions.AddFunds(data.AccountID, data.Amount)
	if err != nil {
		h.logger.Error("Error occurred in call to actions for adding", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
