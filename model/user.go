package model

import (
	"time"
)

type UserSignup struct {
	Username string `json:"username" validate:"required,alphanum,min=6,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,alphanum"`
	Role     string `json:"role" validate:"required,oneof=student mentor"`
}

type UserSignin struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"omitempty,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
