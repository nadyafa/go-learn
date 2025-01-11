package entity

import "time"

type Class struct {
	ClassID     uint   `json:"class_id" gorm:"primaryKey;autoIncrement"`
	CourseID    uint   `json:"course_id" gorm:"index;notNull"`
	ClassName   string `json:"class_name" gorm:"notNull"`
	Description string `json:"description" gorm:"omitempty"`

	MentorID uint `json:"mentor_id" gorm:"notNull"`
	Mentor   User `gorm:"foreignKey:MentorID"`

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Attendances []Attendance `gorm:"foreignKey:ClassID;constrain:OnUpdate:CASCADE"` //ori one2many
	// Attendance Attendance `gorm:"foreignKey:ClassID"`
}
