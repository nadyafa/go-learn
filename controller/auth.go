package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type AuthController interface {
	UserSignup(ctx *gin.Context)
	UserSignin(ctx *gin.Context)
	UserSignout(ctx *gin.Context)
}

type AuthControllerImpl struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) AuthController {
	return &AuthControllerImpl{
		authService: authService,
	}
}

func (c *AuthControllerImpl) UserSignup(ctx *gin.Context) {
	// binding input req
	var userSignup model.UserSignup

	if err := ctx.ShouldBindJSON(&userSignup); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  http.StatusBadRequest,
		})
	}

	// call userSignup service
	user, err := c.authService.UserSignup(userSignup)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// response succeed
	userResp := model.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User successfully created",
		"user":    userResp,
	})
}

func (c *AuthControllerImpl) UserSignin(ctx *gin.Context) {
	var userSignin model.UserSignin

	// binding incoming req
	if err := ctx.ShouldBindJSON(&userSignin); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid input",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// call userSignin service
	user, token, err := c.authService.UserSignin(userSignin)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// set jwt token in cookie
	httpCookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   3600,
	}
	http.SetCookie(ctx.Writer, httpCookie)

	// send succeed response
	userResp := model.UserResponse{
		UserID:    user.UserID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User signed in successfully",
		"user":    userResp,
	})
}

func (c *AuthControllerImpl) UserSignout(ctx *gin.Context) {
	// clear jwt token
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0),
	})

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "User sign out successfully",
	})
}
