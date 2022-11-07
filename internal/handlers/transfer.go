package handlers

import (
	"encoding/json"
	"net/http"

	i "github.com/dupreehkuda/balance-microservice/internal"

	"go.uber.org/zap"
)

// TransferFunds transfer funds from one account to another
func (h handlers) TransferFunds(w http.ResponseWriter, r *http.Request) {
	var data transfer

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.SenderID != "" && data.ReceiverID != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actions.TransferFunds(data.SenderID, data.ReceiverID, data.Amount)

	switch err {
	case i.ErrNoSuchUser:
		w.WriteHeader(http.StatusBadRequest)
		return
	case i.ErrNotEnoughFunds:
		w.WriteHeader(http.StatusPaymentRequired)
		return
	case nil:
		return
	default:
		h.logger.Error("Error call to actions for transfer", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
