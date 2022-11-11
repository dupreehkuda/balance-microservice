package handlers

import (
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetBalance gets current accounts balance
func (h handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	var key i.CtxKey = "account"
	accountID := r.Context().Value(key).(string)

	if accountID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := h.actions.GetBalance(accountID)

	switch err {
	case i.ErrNoSuchUser:
		w.WriteHeader(http.StatusBadRequest)
		return
	case nil:
		w.Write(res)
		return
	default:
		h.logger.Error("Error call to actions for getting", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
