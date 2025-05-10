package util

import (
	"fmt"
	"strconv"
	"time"
)

// ConvertStringToInt32 converts a string to int32.
func ConvertStringToInt32(src string) (int32, error) {
	parsed, err := strconv.ParseInt(src, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(parsed), nil
}

// validateDate checks if the input string is a valid date in "YYYY-MM-DD" format.
func ValidateDate(dateStr string) bool {
	const layout = "2006-01-02" // Go's date layout pattern
	_, err := time.Parse(layout, dateStr)
	return err == nil
}

func GetYearFromQueryParam(year string) (int32, error) {
	if year == "" {
		return int32(time.Now().Year()), nil // this year
	}
	
	if len(year) != 4 {
		return 0, fmt.Errorf("invalid year length in query param")
	}
	y, err := ConvertStringToInt32(year)
	if err != nil {
		return 0, fmt.Errorf("invalid year in query param")
	}
	return y, nil
}

func GetMonthFromQueryParam(month string) (int32, error) {
	if month == "" {
		return int32(time.Now().Month()), nil // this month
	}

	m, err := ConvertStringToInt32(month)
	if err != nil {
		return 0, fmt.Errorf("invalid month in query param")
	}
	return m, nil
}
