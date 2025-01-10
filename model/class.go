package model

import (
	"time"

	"github.com/nadyafa/go-learn/middleware"
)

type CreateClass struct {
	ClassName   string `json:"class_name" validate:"required"`
	Description string `json:"description"`
	// MentorID    uint                  `json:"mentor_id" validate:"required"`
	StartDate middleware.CustomTime `json:"start_date" validate:"required"`
	EndDate   middleware.CustomTime `json:"end_date" validate:"required"`
}

type UpdateClass struct {
	ClassName   string `json:"class_name"`
	Description string `json:"description"`
	// MentorID    uint                  `json:"mentor_id" validate:"required"`
	StartDate middleware.CustomTime `json:"start_date"`
	EndDate   middleware.CustomTime `json:"end_date"`
}

type ClassResp struct {
	ClassID     uint   `json:"class_id"`
	CourseID    uint   `json:"course_id"`
	ClassName   string `json:"class_name"`
	Description string `json:"description"`
	// MentorID    uint      `json:"mentor_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
