package entity

import "time"

type Enrollment struct {
	EnrollmentID uint `json:"enrollment_id" gorm:"primaryKey;autoIncrement"`

	StudentID uint `json:"student_id" gorm:"notNull"`
	Student   User `gorm:"foreignKey:StudentID"`

	CourseID         uint      `json:"course_id" gorm:"index;notNull"`
	EnrollmentDate   time.Time `json:"enrollment_date"`
	CompletionStatus int       `json:"completion_status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
