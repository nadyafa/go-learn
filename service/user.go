package service

import (
	"context"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
	"gorm.io/gorm"
)

type UserService interface {
	UserSignup(ctx context.Context, userSignup model.UserSignup) (res model.UserResponse, err error)
}

type UserServiceImpl struct {
	userRepo repository.UserRepo
	db       *gorm.DB
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) UserSignup(ctx context.Context, userSignup model.UserSignup) (res model.UserResponse, err error) {
	if err := middleware.ValidateUserSignup(s.db, userSignup); err != nil {
		return res, err
	}

	user := entity.User{
		Username:  userSignup.Username,
		Email:     userSignup.Email,
		Password:  userSignup.Password,
		FirstName: userSignup.FirstName,
		LastName:  userSignup.LastName,
		Role:      entity.Role(userSignup.Role),
	}

	_, userErr := s.userRepo.UserSignup(ctx, user)

	if userErr != nil {
		return res, err
	}

	res = model.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return res, nil
}
