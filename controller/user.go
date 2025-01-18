package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type UserController interface {
	GenerateAdmin()
	GetUsers(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	UpdateUserRoleByID(ctx *gin.Context)
	DeleteUserByID(ctx *gin.Context)
}

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}

// generate super admin
func (c *UserControllerImpl) GenerateAdmin() {
	if err := c.userService.GenerateAdmin(); err != nil {
		log.Println("Error generating admin:", err)
		return
	}

	fmt.Println("Super admin created successfully")
}

func (c *UserControllerImpl) GetUsers(ctx *gin.Context) {
	// check if the user is signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "User must sign in to get users list",
			"code":    http.StatusForbidden,
		})
		return
	}

	// validate userRole
	users, err := c.userService.GetUsers(userClaims)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
			"code":    http.StatusForbidden,
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
func (c *UserControllerImpl) GetUserByID(ctx *gin.Context) {
	// check if the user is signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "User must sign in to get user details",
			"code":    http.StatusForbidden,
		})
		return
	}

	// get userId
	userID := ctx.Param("user_id")

	// validate userRole
	user, err := c.userService.GetUserByID(userID, userClaims)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
			"code":    http.StatusForbidden,
		})
		return
	}

	// succeed response
	userRes := model.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User fetch successfully",
		"data":    userRes,
	})
}

// update user role by their ID
func (c *UserControllerImpl) UpdateUserRoleByID(ctx *gin.Context) {
	// check if the user is signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "User must sign in to get user details",
			"code":    http.StatusForbidden,
		})
		return
	}

	// get userId
	userID := ctx.Param("user_id")

	// validate role input
	var role struct {
		Role string `json:"role" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&role); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	userUpdate, err := c.userService.UpdateUserRoleByID(userClaims, userID, role.Role)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
			"code":    http.StatusForbidden,
		})
		return
	}

	// succeed response
	user := model.UserResponse{
		UserID:    userUpdate.UserID,
		Username:  userUpdate.Username,
		Email:     userUpdate.Email,
		Role:      string(userUpdate.Role),
		CreatedAt: userUpdate.CreatedAt,
		UpdatedAt: userUpdate.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("UserID %s role updated successfully", userID),
		"data":    user,
	})
}

func (c *UserControllerImpl) DeleteUserByID(ctx *gin.Context) {
	// check if the user is signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": "User must sign in to get user details",
			"code":    http.StatusForbidden,
		})
		return
	}

	// get userId
	userID := ctx.Param("user_id")

	// call service layer to delete user
	if err := c.userService.DeleteUserByID(userClaims, userID); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"message": err.Error(),
			"code":    http.StatusForbidden,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("UserID %s has been deleted", userID),
	})
}
