package utils

import (
	"fmt"
	"time"
)

func ParseDateYYYYMMDD(dateStr string) (time.Time, error) {
	// format: 2006-01-02 (Go reference time)
	layout := "2006-01-02"
	parsed, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
	}
	return parsed, nil
}
