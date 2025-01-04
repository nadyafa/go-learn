package entity

import "time"

type Project struct {
	ProjectID   uint      `json:"project_id" gorm:"primaryKey;autoIncrement"`
	CourseID    uint      `json:"course_id" gorm:"indexl;notNull"`
	ProjectName string    `json:"project_name" gorm:"notNull"`
	Description string    `json:"description" gorm:"omitempty"`
	Deadline    time.Time `json:"deadline" gorm:"notNull"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// one-to-one
	ProjectSub ProjectSub `gorm:"constraint:OnDelete:CASCADE;"`
}

type ProjectSub struct {
	ProjectSubID   uint      `json:"project_sub_id" gorm:"primaryKey;autoIncrement"`
	ProjectID      uint      `json:"project_id" gorm:"index;notNull"`
	StudentID      uint      `json:"student_id" gorm:"index;notNull"`
	SubmissionDate time.Time `json:"submission_date" gorm:"notNull"`
	Score          int       `json:"score" gorm:"default:0"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
