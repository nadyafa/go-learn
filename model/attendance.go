package model

import "time"

type AttendReq struct {
	StudentID uint `json:"student_id" validate:"required"`
	// ClassID   uint      `json:"class_id" validate:"required"`
	Attended bool `json:"attended" validate:"required"`
	// AttendAt  time.Time `json:"attend_at"`
}

type AttendResp struct {
	AttendID  uint      `json:"attend_id"`
	StudentID uint      `json:"student_id"`
	ClassID   uint      `json:"class_id"`
	CourseID  uint      `json:"course_id"`
	Attended  bool      `json:"attended"`
	AttendAt  time.Time `json:"attend_at"`
}
