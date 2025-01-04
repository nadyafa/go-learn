package model

import "time"

type UserSignup struct {
	Username  string `json:"username" validate:"required,alphanum,min=6,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=3,max=50"`
	LastName  string `json:"last_name" validate:"required,min=3,max=50"`
	Role      string `json:"role" validate:"required,oneof=student admin mentor"`
}

type UserSignin struct {
	Username string `json:"username" validate:"required,alphanum,min=6,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserReqUpdate struct {
	FirstName string `json:"first_name" gorm:"size:50;notNull"`
	LastName  string `json:"last_name" gorm:"size:50;notNull"`
	Password  string `json:"password" gorm:"notNull"`
	// Description string `json:"description" gorm:"omitempty"`
	// Image_path  string `json:"image_path" gorm:"omitempty"`
}

type UserResponse struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	// Description string    `json:"description"`
	// Image_path  string    `json:"image_path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
