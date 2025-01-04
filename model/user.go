package model

import "time"

type UserSignup struct {
	Username string `json:"username" gorm:"size:100;unique;notNull"`
	Email    string `json:"email" gorm:"notNull;unique"`
	Password string `json:"password" gorm:"notNull"`
	Role     string `json:"role" gorm:"notNull"`
}

type UserSignin struct {
	Username string `json:"username" gorm:"size:100;unique;notNull"`
	Email    string `json:"email" gorm:"notNull;unique"`
	Password string `json:"password" gorm:"notNull"`
	Role     string `json:"role" gorm:"notNull"`
}

type UserReqUpdate struct {
	FirstName string `json:"first_name" gorm:"size:50;notNull"`
	LastName  string `json:"last_name" gorm:"size:50;notNull"`
	Password  string `json:"password" gorm:"notNull"`
	// Description string `json:"description" gorm:"omitempty"`
	// Image_path  string `json:"image_path" gorm:"omitempty"`
}

type UserResponse struct {
	UserID      int       `json:"user_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Password    string    `json:"password"`
	Role        string    `json:"role"`
	Description string    `json:"description"`
	Image_path  string    `json:"image_path"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
