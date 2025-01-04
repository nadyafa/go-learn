package entity

import "time"

type Attendance struct {
	AttendID uint `json:"attend_id" gorm:"primaryKey;autoIncrement"`

	StudentID uint `json:"student_id" gorm:"notNull"`
	Student   User `gorm:"foreignKey:StudentID"`

	ClassID      uint      `json:"class_id" gorm:"index;notNull"`
	CourseID     uint      `json:"course_id" gorm:"index;notNull"`
	JoinDate     time.Time `json:"join_date"`
	CheckoutDate time.Time `json:"checkout_date"`
	Attended     bool      `json:"attended" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
