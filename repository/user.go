package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type UserRepo interface {
	GenerateAdmin(admin *entity.User) error
	GetUsers() ([]entity.User, error)
	GetUserByID(userID string) (*entity.User, error)
	UpdateUserRoleByID(userID string, role string) (*entity.User, error)
	DeleteUserByID(userID string) error
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

// generate super admin
func (r *UserRepoImpl) GenerateAdmin(admin *entity.User) error {
	if err := r.db.Create(admin).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepoImpl) GetUsers() ([]entity.User, error) {
	var users []entity.User

	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepoImpl) GetUserByID(userID string) (*entity.User, error) {
	var user entity.User

	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepoImpl) UpdateUserRoleByID(userID string, role string) (*entity.User, error) {
	var userUpdate entity.User

	if err := r.db.Model(userUpdate).Where("user_id = ?", userID).Update("role", role).Error; err != nil {
		return nil, err
	}

	return &userUpdate, nil
}

func (r *UserRepoImpl) DeleteUserByID(userID string) error {
	var user entity.User

	if err := r.db.Where("user_id = ?", userID).Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
