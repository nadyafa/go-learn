package repository

import "gorm.io/gorm"

type UserRepo interface {
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (r *UserRepoImpl) GetUserByID() {

}
