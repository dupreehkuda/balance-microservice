package actions

import i "github.com/dupreehkuda/balance-microservice/internal"

// GetBalance checks if account exists and returns current info
func (a actions) GetBalance(accountID string) ([]byte, error) {
	exists := a.storage.CheckAccountExistence(accountID)
	if !exists {
		return nil, i.ErrNoSuchUser
	}

	ans, err := a.storage.GetBalance(accountID)
	if err != nil {
		return nil, err
	}

	return ans, nil
}
