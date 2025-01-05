package repository

// import (
// 	"context"
// 	"net/http"

// 	"github.com/nadyafa/go-learn/config/helper"
// 	"github.com/nadyafa/go-learn/entity"
// 	"gorm.io/gorm"
// )

// type UserRepo interface {
// 	UserSignup(ctx context.Context, user entity.User) (entity.User, *helper.ErrorStruct)
// }

// type UserRepoImpl struct {
// 	db *gorm.DB
// }

// func NewUserRepo(db *gorm.DB) UserRepo {
// 	return &UserRepoImpl{
// 		db: db,
// 	}
// }

// func (r *UserRepoImpl) UserSignup(ctx context.Context, user entity.User) (entity.User, *helper.ErrorStruct) {
// 	// checking if username already exist
// 	var exitingUser entity.User
// 	if err := r.db.Where("username = ?", user.Username).First(&exitingUser).Error; err == nil {
// 		return error
// 	}

// 	if err := r.db.Debug().Create(&user).WithContext(ctx).Error; err != nil {
// 		helper.Logger(helper.LoggerLevelError, "Error during user signup", err)

// 		return entity.User{}, &helper.ErrorStruct{
// 			Err:  err,
// 			Code: http.StatusInternalServerError,
// 		}
// 	}

// 	return user, nil
// }
