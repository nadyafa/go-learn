package model

import (
	"time"

	"github.com/nadyafa/go-learn/middleware"
)

type CourseReq struct {
	CourseName  string                `json:"course_name" validate:"required"`
	Description string                `json:"description"`
	StartDate   middleware.CustomTime `json:"start_date" validate:"required"`
	EndDate     middleware.CustomTime `json:"end_date" validate:"required"`
}
type UpdateCourse struct {
	CourseName  string                `json:"course_name" validate:"required"`
	Description string                `json:"description"`
	StartDate   middleware.CustomTime `json:"start_date" validate:"required"`
	EndDate     middleware.CustomTime `json:"end_date" validate:"required"`
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
