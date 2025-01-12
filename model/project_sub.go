package model

import "time"

type ProjectSubMentor struct {
	Description string `json:"description"`
	Score       int    `json:"score" validate:"required"`
}

type ProjectSubStudent struct {
	ProjectPath string `json:"project_path" validate:"required"`
	// Score       int    `json:"score"`
}

type StudentSubmitResp struct {
	ProjectSubID   uint      `json:"project_sub_id"`
	ProjectID      uint      `json:"project_id"`
	StudentID      uint      `json:"student_id"`
	SubmissionDate time.Time `json:"submission_date"`
	ProjectPath    string    `json:"project_path"`
}

type MentorSubmitResp struct {
	ProjectSubID   uint      `json:"project_sub_id"`
	ProjectID      uint      `json:"project_id"`
	StudentID      uint      `json:"student_id"`
	SubmissionDate time.Time `json:"submission_date"`
	ProjectPath    string    `json:"project_path"`
	Score          int       `json:"score"`
	Description    string    `json:"description"`
}
