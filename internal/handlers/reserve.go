package handlers

import (
	"encoding/json"
	"net/http"

	i "github.com/dupreehkuda/balance-microservice/internal"

	"go.uber.org/zap"
)

// ReserveFunds makes a request to reserve funds on account
func (h handlers) ReserveFunds(w http.ResponseWriter, r *http.Request) {
	var data reserve

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.TargetID != "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.actions.ReserveFunds(data.TargetID, data.ServiceID, data.OrderID, data.Amount)

	switch err {
	case i.ErrNoSuchUser:
		w.WriteHeader(http.StatusBadRequest)
		return
	case i.ErrWrongCredentials:
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	case i.ErrNotEnoughFunds:
		w.WriteHeader(http.StatusPaymentRequired)
		return
	case nil:
		return
	default:
		h.logger.Error("Error call to actions for reserve", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
