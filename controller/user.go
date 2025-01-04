package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type UserController interface {
	UserSignup(ctx *gin.Context)
}

type UserControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		userService: userService,
	}
}

func (c *UserControllerImpl) UserSignup(ctx *gin.Context) {
	var userSignup model.UserSignup

	if err := ctx.ShouldBindJSON(&userSignup); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}

	res, err := c.userService.UserSignup(ctx, userSignup)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, res)
}
