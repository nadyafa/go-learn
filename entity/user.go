package entity

import (
	"time"
)

type Role string

const (
	Student Role = "student"
	Admin   Role = "admin"
	Mentor  Role = "mentor"
)

type User struct {
	UserID   uint   `json:"user_id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"size:100;unique;notNull"`
	Email    string `json:"email" gorm:"notNull;unique"`
	// FirstName string    `json:"first_name" gorm:"size:50;notNull"`
	// LastName  string    `json:"last_name" gorm:"size:50;notNull"`
	Password  string    `json:"password" gorm:"notNull"`
	Role      Role      `json:"role" gorm:"default:student"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Enrollments []Enrollment `gorm:"foreignKey:StudentID"`
	Classes []Class `gorm:"foreignKey:MentorID;constrain:OnUpdate:CASCADE"`
	// Courses     []Course     `gorm:"many2many:course_enrollments"`
	Projects    []Project    `gorm:"foreignKey:CourseID;constrain:OnUpdate:CASCADE"`
	Attendances []Attendance `gorm:"foreignKey:StudentID;constrain:OnUpdate:CASCADE"`
}
