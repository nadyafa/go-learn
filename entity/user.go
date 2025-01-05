package entity

import (
	"time"
)

type Role string

const (
	Student Role = "Student"
	Admin   Role = "Admin"
	Mentor  Role = "Mentor"
)

type User struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"size:100;unique;notNull"`
	Email    string `json:"email" gorm:"notNull;unique"`
	// FirstName string    `json:"first_name" gorm:"size:50;notNull"`
	// LastName  string    `json:"last_name" gorm:"size:50;notNull"`
	Password         string `json:"password" gorm:"notNull"`
	EmailVerified    bool   `json:"email_verified"`
	VerificationCode string `json:"verification_code"`
	// Role      Role      `json:"role" gorm:"type:varchar(20);notNull"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Enrollments []Enrollment `gorm:"foreignKey:StudentID"`
	// Classes     []Class      `gorm:"foreignKey:MentorID"`
	// Courses     []Course     `gorm:"many2many:course_enrollments"`
	// Projects    []Project    `gorm:"foreignKey:CourseID"`
	// Attendances []Attendance `gorm:"foreignKey:StudentID"`
}
