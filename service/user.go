package service

// import (
// 	"context"
// 	"fmt"

// 	"github.com/nadyafa/go-learn/config/helper"
// 	"github.com/nadyafa/go-learn/entity"
// 	"github.com/nadyafa/go-learn/middleware"
// 	"github.com/nadyafa/go-learn/model"
// 	"github.com/nadyafa/go-learn/repository"
// 	"gorm.io/gorm"
// )

// type UserService interface {
// 	UserSignup(ctx context.Context, userSignup model.UserSignup) (model.UserResponse, error)
// }

// type UserServiceImpl struct {
// 	userRepo repository.UserRepo
// 	db       *gorm.DB
// }

// func NewUserService(userRepo repository.UserRepo) UserService {
// 	return &UserServiceImpl{
// 		userRepo: userRepo,
// 		db:       nil,
// 	}
// }

// func (s *UserServiceImpl) UserSignup(ctx context.Context, userSignup model.UserSignup) (model.UserResponse, error) {
// 	if err := middleware.ValidateUserSignup(s.db, userSignup); err != nil {
// 		helper.Logger(helper.LoggerLevelError, "Validation error during user signup", err)

// 		return model.UserResponse{}, fmt.Errorf("signup validation failed: %w", err)
// 	}

// 	user := entity.User{
// 		Username:  userSignup.Username,
// 		Email:     userSignup.Email,
// 		Password:  userSignup.Password,
// 		FirstName: userSignup.FirstName,
// 		LastName:  userSignup.LastName,
// 		Role:      entity.Role(userSignup.Role),
// 	}

// 	_, err := s.userRepo.UserSignup(ctx, user)

// 	if err != nil {
// 		helper.Logger(helper.LoggerLevelError, "Database error during user signup", err.Err)

// 		return model.UserResponse{}, fmt.Errorf("failed to create user: %v", err)
// 	}

// 	res := model.UserResponse{
// 		UserID:    user.UserID,
// 		Username:  user.Username,
// 		Email:     user.Email,
// 		FirstName: user.FirstName,
// 		LastName:  user.LastName,
// 		Password:  user.Password,
// 		Role:      string(user.Role),
// 		CreatedAt: user.CreatedAt,
// 		UpdatedAt: user.UpdatedAt,
// 	}

// 	return res, nil
// }
