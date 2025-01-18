package service

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
)

type AuthService interface {
	UserSignup(userSignup model.UserSignup) (*entity.User, error)
	UserSignin(user model.UserSignin) (*entity.User, string, error)
}

type AuthServiceImpl struct {
	authRepo  repository.AuthRepo
	validator *validator.Validate
}

func NewAuthService(authRepo repository.AuthRepo) AuthService {
	return &AuthServiceImpl{
		validator: validator.New(),
		authRepo:  authRepo,
	}
}

func (s *AuthServiceImpl) UserSignup(userSignup model.UserSignup) (*entity.User, error) {
	// validate user input
	errorMsg := middleware.ValidateUserInput(userSignup, true)
	if len(errorMsg) > 0 {
		return nil, fmt.Errorf("validation failed: %v", errorMsg)
	}

	// check if username is already exist
	existingUser, _ := s.authRepo.FindByUsername(userSignup.Username)
	if existingUser == nil {
		return nil, fmt.Errorf("username is taken")
	}

	// check if email is already exist
	existingUser, _ = s.authRepo.FindByEmail(userSignup.Email)
	if existingUser == nil {
		return nil, fmt.Errorf("email is taken")
	}

	// hash password
	hashedPassword, err := middleware.HashPassword(userSignup.Password)
	if err != nil {
		return nil, fmt.Errorf("unable to hash password")
	}

	// create user entity
	user := &entity.User{
		Username: userSignup.Username,
		Email:    userSignup.Email,
		Password: hashedPassword,
		Role:     entity.Student,
	}

	// save user to db
	if err := s.authRepo.UserSignup(user); err != nil {
		return nil, fmt.Errorf("user signup process is failed")
	}

	// notify user
	// if userSignup.Role == string(entity.Mentor) {
	// 	// notify admin
	// 	if err := middleware.SendMail(
	// 		os.Getenv("ADMIN_EMAIL"),
	// 		"New Mentor Sign Up Pending Validation",
	// 		fmt.Sprintf("A user has signed up with the role of mentor. Please validate user with UserID %s and Username %s", fmt.Sprint(user.UserID), user.Username),
	// 	); err != nil {
	// 		return nil, fmt.Errorf("failed to send notification to admin: %v", err)
	// 	}

	// 	// notify mentor
	// 	if err := middleware.SendMail(
	// 		userSignup.Email,
	// 		"Go-Learn Sign Up",
	// 		fmt.Sprintf("You have successfully sign up with UserID %s and Username %s. We will notify you as soon as your role as a mentor is being verified, but you still can sign in to account. Good luck!", fmt.Sprint(user.UserID), user.Username),
	// 	); err != nil {
	// 		return nil, fmt.Errorf("failed to send notification to mentor: %v", err)
	// 	}
	// }

	// if userSignup.Role == string(entity.Student) || userSignup.Role == "" {
	// 	// notify mentor
	// 	if err := middleware.SendMail(
	// 		userSignup.Email,
	// 		"Go-Learn Sign Up",
	// 		fmt.Sprintf("You have successfully sign up with UserID %s and Username %s. Please sign in to verify your login. Good luck!", fmt.Sprint(user.UserID), user.Username),
	// 	); err != nil {
	// 		return nil, fmt.Errorf("failed to send notification to student: %v", err)
	// 	}
	// }

	return user, nil
}

func (s *AuthServiceImpl) UserSignin(user model.UserSignin) (*entity.User, string, error) {
	// check if user or email exist
	var existingUser *entity.User
	var err error

	if user.Email != "" {
		existingUser, err = s.authRepo.FindByEmail(user.Email)

		if err != nil {
			return nil, "", fmt.Errorf("email not found")
		}
	} else if user.Username != "" {
		existingUser, err = s.authRepo.FindByUsername(user.Username)

		if err != nil {
			return nil, "", fmt.Errorf("username not found")
		}
	}

	// verify password
	if !middleware.CheckPasswordHash(user.Password, existingUser.Password) {
		return nil, "", fmt.Errorf("password incorrect")
	}

	// generate JWT token
	token, err := middleware.GenerateJWT(existingUser.Username, existingUser.Role, existingUser.UserID)
	if err != nil {
		return nil, "", fmt.Errorf("unable to generate token: %v", err)
	}

	return existingUser, token, nil
}
