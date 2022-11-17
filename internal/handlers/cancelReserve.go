package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// CancelReserve cancels early created reserve request
func (h handlers) CancelReserve(w http.ResponseWriter, r *http.Request) {
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

	err = h.actions.CancelReserve(data.OrderID)

	switch err {
	case i.ErrNoSuchOrder:
		w.WriteHeader(http.StatusBadRequest)
	case i.ErrOrderProcessed:
		w.WriteHeader(http.StatusNotAcceptable)
	case nil:
		return
	default:
		h.logger.Error("Error call to actions for cancel", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
