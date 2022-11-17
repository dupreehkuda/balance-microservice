package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetHistory gets user`s operations history
func (h handlers) GetHistory(w http.ResponseWriter, r *http.Request) {
	var data history

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		h.logger.Error("Unable to decode JSON", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.User == "" || data.SortOrder == "" || data.SortBy == "" || data.Quantity == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if data.SortOrder != "desc" && data.SortOrder != "asc" || data.SortBy != "amount" && data.SortBy != "date" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := h.actions.GetHistory(data.User, data.SortOrder, data.SortBy, data.Quantity)

	switch err {
	case i.ErrNoSuchUser:
		w.WriteHeader(http.StatusBadRequest)
		return
	case i.ErrNoData:
		w.WriteHeader(http.StatusNoContent)
		return
	case nil:
		w.Write(resp)
		return
	default:
		h.logger.Error("Error call to actions for reserve", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
