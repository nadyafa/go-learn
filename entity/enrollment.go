package entity

import (
	"time"
)

type Status string

const (
	Pending  Status = "pending"
	Enroll   Status = "enroll"
	Complete Status = "complete"
	Failed   Status = "failed"
	Cancel   Status = "cancel"
)

type Enrollment struct {
	EnrollmentID uint `json:"enrollment_id" gorm:"primaryKey;autoIncrement"`

	UserID uint `json:"user_id" gorm:"notNull"`
	// User     User `gorm:"foreignKey:UserID"`
	UserRole Role `json:"user_role"`

	CourseID       uint      `json:"course_id" gorm:"index;notNull"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	EnrollStatus   Status    `json:"completion_status" gorm:"default:pending"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// enrollment statusnya ada pending, enroll, passed, unfinished
