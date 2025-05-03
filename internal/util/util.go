package util

import (
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
