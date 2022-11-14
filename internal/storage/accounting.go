package storage

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"

	"go.uber.org/zap"

	i "github.com/dupreehkuda/balance-microservice/internal"
)

// GetReportLink gets data for the report and retutns a CSV
func (s storage) GetReportLink(month, year string) (string, error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		s.logger.Error("Error while acquiring connection", zap.Error(err))
		return "", err
	}
	defer conn.Release()

	var data = []accountingData{}

	period := fmt.Sprintf("%s-%s-%%", year, month)

	rows, err := conn.Query(context.Background(), `select service_id, sum(amount) from orders
		where cast(processed_date as varchar) like $1
		group by service_id;`, period)

	if err != nil {
		s.logger.Error("Error occurred getting accounting data", zap.Error(err))
		return "", err
	}

	for rows.Next() {
		var r accountingData
		err := rows.Scan(&r.ServiceID, &r.Sum)
		if err != nil {
			s.logger.Error("Error while scanning rows", zap.Error(err))
			return "", err
		}
		data = append(data, r)
	}

	if len(data) == 0 {
		return "", i.ErrNoData
	}

	return makeCSV(data), nil
}

// makeCSV converts a slice of accountingData and returns a CSV string
func makeCSV(data []accountingData) string {
	var dataToCSV = [][]string{{"Service", "Sum"}}

	for _, val := range data {
		dataToCSV = append(dataToCSV, []string{val.ServiceID, val.Sum.String()})
	}

	buf := new(bytes.Buffer)
	w := csv.NewWriter(buf)
	w.WriteAll(dataToCSV)
	csvString := buf.String()

	return csvString
}
