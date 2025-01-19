package entity

import "time"

type Course struct {
	CourseID    uint   `json:"course_id" gorm:"primaryKey;autoIncrement"`
	CourseName  string `json:"course_name" gorm:"notNull"`
	Description string `json:"description" gorm:"omitempty"`

	MentorID uint `json:"mentor_id" gorm:"notNull"`
	Mentor   User `gorm:"foreignKey:MentorID"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Classes []Class `gorm:"foreignKey:CourseID"` //ori one2many
	// Class Class `gorm:"foreignKey:CourseID"`
	Enrollments []Enrollment `gorm:"foreignKey:CourseID"`
	Projects    []Project    `gorm:"foreignKey:CourseID"`
	// Tests       []Test       `gorm:"foreignKey:CourseID"`
}
