package entity

import "time"

type Test struct {
	TestId      uint      `json:"test_id" gorm:"primaryKey;autoIncrement"`
	CourseID    uint      `json:"course_id" gorm:"indexl;notNull"`
	Title       string    `json:"title" gorm:"notNull"`
	Description string    `json:"description" gorm:"omitempty"`
	Deadline    time.Time `json:"deadline" gorm:"notNull"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// one-to-one
	TestSub TestSub `gorm:"constraint:OnDelete:CASCADE;"`
}

type TestSub struct {
	TestSubID      uint      `json:"test_sub_id" gorm:"primaryKey;autoIncrement"`
	TestID         uint      `json:"test_id" gorm:"index;autoIncrement"`
	StudentID      uint      `json:"student_id" gorm:"index;autoIncrement"`
	SubmissionDate time.Time `json:"submission_date" gorm:"notNull"`
	Score          int       `json:"score" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
