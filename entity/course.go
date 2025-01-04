package entity

import "time"

type Course struct {
	CourseID    uint      `json:"course_id" gorm:"primaryKey;autoIncrement"`
	CourseName  string    `json:"course_name" gorm:"notNull"`
	Description string    `json:"description" gorm:"omitempty"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Class      Class        `gorm:"foreignKey:CourseID"`
	Enrollment []Enrollment `gorm:"foreignKey:CourseID"`
	Project    []Project    `gorm:"foreignKey:CourseID"`
}
