package middleware

import (
	"fmt"
	"strings"
	"time"
)

func ValidateCourseName(courseName string) (bool, string) {
	// if course name empty
	if strings.TrimSpace(courseName) == "" {
		return false, "Course name cannot be empty"
	}

	// if all validation passed
	return true, ""
}

func ValidateCourseDate(startDateStr, endDateStr string) (bool, string) {
	// validate startDate & endDate format
	layout := "02-01-2006 15:04"
	startDate, err := time.Parse(layout, startDateStr)
	if err != nil {
		return false, fmt.Sprintf("Invalid start date format. Expected format: %s", layout)
	}

	endDate, err := time.Parse(layout, endDateStr)
	if err != nil {
		return false, fmt.Sprintf("Invalid end date format. Expected format: %s", layout)
	}

	// validate startDate has to be earlier than endDate
	if !startDate.Before(endDate) {
		return false, "Start date must be earlier than end date"
	}

	// if all validation passed
	return true, ""
}

// func ValidateBoolean(input string) (bool, error) {
// 	input = strings.ToLower(input)

// 	if input == "true" || input == "false" {
// 		return true, nil
// 	}

// 	return false, errors.New("input must be 'true' or 'false'")
// }
