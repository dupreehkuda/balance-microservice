package actions

import i "github.com/dupreehkuda/balance-microservice/internal"

// GetReport checks for report and gets it
func (a actions) GetReport(reportID string) (string, error) {
	exists := a.storage.CheckReportExistence(reportID)
	if !exists {
		return "", i.ErrNoSuchReport
	}

	resp, err := a.storage.ReadReport(reportID)
	if err != nil {
		return "", err
	}

	return resp, nil
}
