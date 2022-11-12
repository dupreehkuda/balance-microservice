package actions

import (
	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetReportLink checks report by hash and creates new report
func (a actions) GetReportLink(month, year string) (string, error) {
	repID := i.ReportHash(month, year)

	exitsts := a.storage.CheckReportExistence(repID)
	if exitsts {
		return repID, nil
	}

	csv, err := a.storage.GetReportLink(month, year)
	if err != nil {
		a.logger.Error("Error in call to get report", zap.Error(err))
		return "", err
	}

	err = a.storage.WriteReport(repID, csv)
	if err != nil {
		a.logger.Error("Error in call to write report", zap.Error(err))
		return "", err
	}

	return repID, nil
}
