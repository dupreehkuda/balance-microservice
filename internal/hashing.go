package internal

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// ReportHash returns short hash string
func ReportHash(month, year string) string {
	coupled := fmt.Sprintf("%s%s", month, year)

	hsha := sha1.Sum([]byte(coupled))

	return hex.EncodeToString(hsha[:5])
}
