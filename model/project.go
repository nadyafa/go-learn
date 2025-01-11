package model

import (
	"time"

	"github.com/nadyafa/go-learn/middleware"
)

type CreateProject struct {
	ProjectName string                `json:"project_name" validate:"required"`
	Description string                `json:"description"`
	Deadline    middleware.CustomTime `json:"deadline" validate:"required"`
}

type UpdateProject struct {
	ProjectName string                `json:"project_name"`
	Description string                `json:"description"`
	Deadline    middleware.CustomTime `json:"deadline"`
}

type ProjectResp struct {
	ProjectID   uint      `json:"project_id"`
	CourseID    uint      `json:"course_id"`
	ProjectName string    `json:"project_name"`
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
