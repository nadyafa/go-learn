package repository

import (
	"context"

	"gorm.io/gorm"
)

type UsersRepo interface {
}

type UserRepoImpl struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UsersRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (r *UserRepoImpl) AddUser(ctx context.Context, userID uint)
