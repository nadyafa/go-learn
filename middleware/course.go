package middleware

import (
	"fmt"
	"strings"
	"time"
)

func ValidateCourseName(courseName string) (bool, error) {
	// if course name empty
	if strings.TrimSpace(courseName) == "" {
		return false, fmt.Errorf("course name cannot be empty")
	}

	// if all validation passed
	return true, nil
}

func ValidateCourseDate(startDateStr, endDateStr string) (bool, error) {
	// validate startDate & endDate format
	layout := "02-01-2006 15:04"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return false, fmt.Errorf("invalid start date format. Expected format: %s", layout)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return false, fmt.Errorf("invalid end date format. Expected format: %s", layout)
	}

	// validate startDate has to be earlier than endDate
	if !startDate.Before(endDate) {
		return false, fmt.Errorf("start date must be earlier than end date")
	}

	// if all validation passed
	return true, nil
}
