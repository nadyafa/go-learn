package model

import (
	"time"

	"github.com/nadyafa/go-learn/entity"
)

type CreateEnrollment struct {
	UserID   uint        `json:"user_id" validate:"required"`
	UserRole entity.Role `json:"user_role" validate:"required"`
}

type UpdateEnrollment struct {
	UserID uint `json:"user_id" validate:"required"`
	// UserRole entity.Role `json:"user_role" validate:"required"`
}

type EnrollResp struct {
	EnrollmentID uint `json:"enrollment_id" validate:"required"`
	StudentID    uint `json:"user_id" validate:"required"`
	// UserRole       entity.Role   `json:"user_role" validate:"required"`
	CourseID       uint          `json:"course_id" validate:"required"`
	EnrollmentDate *time.Time    `json:"enrollment_date"`
	EnrollStatus   entity.Status `json:"completion_status" gorm:"default:pending"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}
