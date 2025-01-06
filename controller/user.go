package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"gorm.io/gorm"
)

type UserController interface {
	UserSignup(ctx *gin.Context)
	UserSignin(ctx *gin.Context)
	AdminLogin()
}

type UserControllerImpl struct {
	db        *gorm.DB
	validator *validator.Validate
}

func NewUserController(db *gorm.DB) UserController {
	return &UserControllerImpl{
		db:        db,
		validator: validator.New(),
	}
}

func (c *UserControllerImpl) UserSignup(ctx *gin.Context) {
	// binding input req
	var userSignup model.UserSignup

	if err := ctx.ShouldBindJSON(&userSignup); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  http.StatusBadRequest,
		})
	}

	// use validator to validate input with model struct
	errorMessage := middleware.ValidateUserInput(userSignup, true)
	if len(errorMessage) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errorMessage,
			"code":    http.StatusBadRequest,
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

	user := entity.User{
		Username: userSignup.Username,
		Email:    userSignup.Email,
		Password: userSignup.Password,
	}

	if err := c.db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to sign up",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// response succeed
	ctx.JSON(http.StatusCreated, gin.H{
		"user_id":  user.UserID,
		"username": userSignup.Username,
		"email":    userSignup.Email,
		"message":  "User successfully created",
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

	// validate user input
	errorMessage := middleware.ValidateUserInput(userSignin, false)
	if len(errorMessage) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errorMessage,
			"code":    http.StatusBadRequest,
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

// generate super admin
func (c *UserControllerImpl) AdminLogin() {
	var admin entity.User

	// check super admin existance in db
	var count int64
	c.db.Model(&admin).Where("role = ?", entity.Admin).Count(&count)
	if count > 0 {
		log.Println("Super admin already exist")
		return
	}

	// plan password
	password := os.Getenv("ADMIN_PASSWORD")

	// hashing password before storing
	hashedPassword, err := middleware.HashPassword(password)
	if err != nil {
		log.Fatal("Error hashing password:", err)
		return
	}

	// create super admin
	admin = entity.User{
		UserID:   0,
		Username: os.Getenv("ADMIN_USERNAME"),
		Email:    os.Getenv("ADMIN_EMAIL"),
		Password: hashedPassword,
		Role:     entity.Admin,
	}

	if err := c.db.Create(&admin).Error; err != nil {
		fmt.Println("Error creating super admin:", err)
	} else {
		fmt.Println("Super admin created successfully")
	}
}
