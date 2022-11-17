package actions

import (
	"fmt"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetHistory checks if user exists and then gets user`s operations history
func (a actions) GetHistory(targetId, sortOrder, sortBy string, quantity int) ([]byte, error) {
	exists := a.storage.CheckAccountExistence(targetId)
	if !exists {
		return nil, i.ErrNoSuchUser
	}

	if sortBy == "date" {
		sortBy = "processed_at"
	} else {
		sortBy = "funds"
	}

	param := fmt.Sprintf("order by %s %s limit %v;", sortBy, sortOrder, quantity)

	resp, err := a.storage.GetHistory(targetId, param)
	if err != nil {
		return nil, err
	}

	if len(resp) == 0 {
		return nil, i.ErrNoData
	}

	return resp, nil
}
