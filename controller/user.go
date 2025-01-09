package controller

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"gorm.io/gorm"
)

type AdminController interface {
	GenerateAdmin()
	GetUsers(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	UpdateUserRoleByID(ctx *gin.Context)
	DeleteUserByID(ctx *gin.Context)
}

type AdminControllerImpl struct {
	db *gorm.DB
}

func NewAdminController(db *gorm.DB) AdminController {
	return &AdminControllerImpl{
		db: db,
	}
}

// // create new user
// func (c *AdminControllerImpl) CreateNewUser(ctx *gin.Context) {
// 	// check if the currentUser is admin
// 	claims, _ := ctx.Get("currentUser")
// 	userClaims, ok := claims.(*middleware.UserClaims)
// 	if !ok || userClaims.Role != entity.Admin {
// 		ctx.JSON(http.StatusForbidden, gin.H{
// 			"error": "Only admin can create a new user",
// 			"code":  http.StatusForbidden,
// 		})
// 		return
// 	}

// 	// bind user input
// 	var user model.UserSignup
// 	if err := ctx.ShouldBindJSON(&user); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 			"code":  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	// hashing password input
// 	hashedPassword, err := middleware.HashPassword(user.Password)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Unable to hash password",
// 			"code":  http.StatusInternalServerError,
// 		})
// 		return
// 	}

// 	user.Password = hashedPassword

// 	// add new user to db
// 	if err := c.db.Create(&user).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": "Unable to create a new user",
// 			"code":  http.StatusInternalServerError,
// 		})
// 		return
// 	}

// 	// success response
// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"message": "User created successfully",
// 		"user":    user,
// 	})
// }

// generate super admin
func (c *AdminControllerImpl) GenerateAdmin() {
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

func (c *AdminControllerImpl) GetUsers(ctx *gin.Context) {
	// check if the current user is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Admin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Only admin can get users data",
			"code":    http.StatusForbidden,
		})
		return
	}

	// get users
	var users []entity.User

	if err := c.db.Find(&users).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve users",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Users fetch successfully",
		"code":    http.StatusOK,
		"data":    users,
	})
}

// only admin and mentor can see UserByID
func (c *AdminControllerImpl) GetUserByID(ctx *gin.Context) {
	// get userId
	userID := ctx.Param("user_id")

	// get currentUser info (userID and role)
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)

	// check if the current user is student
	if !ok || userClaims.Role == entity.Student {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Access Restricted",
			"code":    http.StatusForbidden,
		})
		return
	}

	// find user by input param userId
	var user entity.User
	if err := c.db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Users fetch successfully",
		"code":    http.StatusOK,
		"data":    user,
	})
}

// update user role by their ID
func (c *AdminControllerImpl) UpdateUserRoleByID(ctx *gin.Context) {
	// check if the current user is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Admin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Only admin can update user role",
			"code":    http.StatusForbidden,
		})
		return
	}

	// get userId
	userID := ctx.Param("user_id")

	var req struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// make sure new role input valid
	if strings.ToLower(req.Role) != "student" && strings.ToLower(req.Role) != "mentor" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid role",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// find user by input param userId
	var user entity.User
	if err := c.db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// update user role
	if err := c.db.Model(&user).Updates(entity.User{Role: entity.Role(req.Role)}).Where("user_id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update user role",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	userUpdate := model.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(req.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("UserID %s role updated successfully", userID),
		"code":    http.StatusOK,
		"data":    userUpdate,
	})
}

func (c *AdminControllerImpl) DeleteUserByID(ctx *gin.Context) {
	// check if the current user is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Admin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "Only admin can delete user",
			"code":    http.StatusForbidden,
		})
		return
	}

	// get userId
	userID := ctx.Param("user_id")

	// find user by input param userId
	var user entity.User
	if err := c.db.First(&user, userID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// delete user from db
	if err := c.db.Delete(&user).Where("user_id = ?", userID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Failed to delete user with ID %s", userID),
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("UserID %s has been deleted", userID),
	})
}
