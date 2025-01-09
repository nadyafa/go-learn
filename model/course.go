package model

import (
	"fmt"
	"time"
)

type CourseReq struct {
	CourseName  string     `json:"course_name" validate:"required"`
	Description string     `json:"description"`
	StartDate   customTime `json:"start_date" validate:"required"`
	EndDate     customTime `json:"end_date" validate:"required"`
}

type CourseResp struct {
	CourseID    uint      `json:"course_id"`
	CourseName  string    `json:"course_name"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type customTime struct {
	time.Time
}

func (c *customTime) UnmarshalJSON(b []byte) error {
	layout := "02-01-2006 15:04" //dd-mm-yyyy hour:minute
	parseTime, err := time.Parse(fmt.Sprintf("\"%s\"", layout), string(b))
	if err != nil {
		return fmt.Errorf("date time format input invalid. expected format: %s", layout)
	}

	c.Time = parseTime
	return nil
}
