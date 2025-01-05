package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"gorm.io/gorm"
)

type UserController interface {
	UserSignup(ctx *gin.Context)
	UserSignin(ctx *gin.Context)
}

type UserControllerImpl struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return &UserControllerImpl{
		db: db,
	}
}

func (c *UserControllerImpl) UserSignup(ctx *gin.Context) {
	// binding input req
	var userSignup entity.User

	if err := ctx.ShouldBindJSON(&userSignup); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  http.StatusBadRequest,
		})
	}

	// username format validation
	if err := middleware.ValidateUsername(userSignup.Username); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Username must contain alphanumeric",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// password validation
	// password must contain lower & uppercase letter, number and special character and have at least 8 characters
	if err := middleware.ValidatePassword(userSignup.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be 8 character of combination aphanumeric and special character",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// checking username if exist in db
	var exitingUser entity.User

	if err := c.db.Where("username = ?", userSignup.Username).First(&exitingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Username is already exist",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// checking email
	if err := c.db.Where("email = ?", userSignup.Email).First(&exitingUser).Error; err == nil {
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "Email is already exist",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// hashing password
	hashedPassword, err := middleware.HashPassword(userSignup.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to hash password",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// create new user
	userSignup.Password = hashedPassword
	if err := c.db.Create(&userSignup).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to sign up",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// response succeed

	ctx.JSON(http.StatusCreated, gin.H{
		"user_id":  userSignup.UserID,
		"username": userSignup.Username,
		"email":    userSignup.Email,
	})
}

func (c *UserControllerImpl) UserSignin(ctx *gin.Context) {
	var userSignin entity.User

	// binding incoming req
	if err := ctx.ShouldBindJSON(&userSignin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// check if the user exist
	var exitingUser entity.User
	if err := c.db.Where("username = ? OR email = ?", userSignin.Username, userSignin.Email).First(&exitingUser).Error; err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid username or email",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	if !middleware.CheckPasswordHash(userSignin.Password, exitingUser.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid password",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user_id":  exitingUser.UserID,
		"username": exitingUser.Username,
		"email":    exitingUser.Email,
		"message":  "User signed in successfully",
	})
}
