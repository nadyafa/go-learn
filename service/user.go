package service

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/repository"
)

type UserService interface {
	GenerateAdmin() error
	GetUsers(userClaims *middleware.UserClaims) ([]entity.User, error)
	GetUserByID(userID string, userClaims *middleware.UserClaims) (*entity.User, error)
	UpdateUserRoleByID(userClaims *middleware.UserClaims, userID string, role string) (*entity.User, error)
	DeleteUserByID(userClaims *middleware.UserClaims, userID string) error
}

type UserServiceImpl struct {
	userRepo repository.UserRepo
}

func NewUserService(userRepo repository.UserRepo) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

// generate super admin
func (s *UserServiceImpl) GenerateAdmin() error {
	// plan password
	password := os.Getenv("ADMIN_PASSWORD")

	// hashing password before storing
	hashedPassword, err := middleware.HashPassword(password)
	if err != nil {
		log.Fatal("Error hashing password:", err)
		return err
	}

	// create super admin
	admin := entity.User{
		UserID:   0,
		Username: os.Getenv("ADMIN_USERNAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: hashedPassword,
		Role:     entity.Admin,
	}

	// call repo layer to save super admin
	if err := s.userRepo.GenerateAdmin(&admin); err != nil {
		log.Println("Error creating super admin:", err)
		return err
	}

	log.Println("Super admin created successfully")
	return nil
}

func (s *UserServiceImpl) GetUsers(userClaims *middleware.UserClaims) ([]entity.User, error) {
	// check if userRole is admin
	if userClaims.Role != entity.Admin {
		return nil, fmt.Errorf("only admin can get users data")
	}

	// fetch users from repo
	users, err := s.userRepo.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("unable to get list of users")
	}

	return users, nil
}

func (s *UserServiceImpl) GetUserByID(userID string, userClaims *middleware.UserClaims) (*entity.User, error) {
	// validate user role
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("student have no access")
	}

	// fetch user from repo
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserServiceImpl) UpdateUserRoleByID(userClaims *middleware.UserClaims, userID, role string) (*entity.User, error) {
	if userClaims.Role != entity.Admin {
		return nil, fmt.Errorf("only admin can update user data")
	}

	// validate role input
	if strings.ToLower(role) != "student" && strings.ToLower(role) != "mentor" {
		return &entity.User{}, fmt.Errorf("role must be student or mentor")
	}

	// check if userID exist
	_, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// update user role
	_, err = s.userRepo.UpdateUserRoleByID(userID, role)
	if err != nil {
		return nil, fmt.Errorf("unable to update user role")
	}

	// return user value
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *UserServiceImpl) DeleteUserByID(userClaims *middleware.UserClaims, userID string) error {
	if userClaims.Role != entity.Admin {
		return fmt.Errorf("only admin can delete user")
	}

	// check if userID exist
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return fmt.Errorf("user not found")
	}

	// update user role
	if err := s.userRepo.DeleteUserByID(userID); err != nil {
		return fmt.Errorf("unable to update user role")
	}

	return nil
}
