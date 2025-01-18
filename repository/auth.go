package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type AuthRepo interface {
	UserSignup(user *entity.User) error
	FindByUsername(username string) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
}

type AuthRepoImpl struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &AuthRepoImpl{
		db: db,
	}
}

func (r *AuthRepoImpl) UserSignup(user *entity.User) error {
	if err := r.db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *AuthRepoImpl) FindByUsername(username string) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *AuthRepoImpl) FindByEmail(email string) (*entity.User, error) {
	var user entity.User

	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}
