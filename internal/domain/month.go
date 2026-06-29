package domain

import (
	"fmt"
	"time"
)

const MonthLayout = "01-2006"

func ParseMonth(value string) (time.Time, error) {
	parsed, err := time.Parse(MonthLayout, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse month %q: %w", value, err)
	}

	return parsed, nil
}

func FormatMonth(value time.Time) string {
	return value.Format(MonthLayout)
}
