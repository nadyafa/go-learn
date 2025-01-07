package controller

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/middleware"
)

func AuthMiddleware(ctx *gin.Context) {
	// get auth token jwt from header
	authHeader := ctx.GetHeader("Authorization")
	var tokenStr string

	if authHeader == "" {
		// if auth token exist in header
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	} else {
		// if no auth token in header, get from cookies
		cookie, err := ctx.Cookie("auth_token")
		if err == nil {
			tokenStr = cookie
		} else {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization token missing",
				"code":  http.StatusUnauthorized,
			})

			ctx.Abort()
			return
		}
	}

	// parse jwt user info
	user, err := middleware.ParseJWT(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization token",
			"code":  http.StatusUnauthorized,
		})

		ctx.Abort()
		return
	}

	// set user info in context
	ctx.Set("currentUser", user)
	ctx.Next()
}
