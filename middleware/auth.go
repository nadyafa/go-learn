package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nadyafa/go-learn/entity"
)

var sercretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UserClaims struct {
	UserID uint        `json:"user_id"`
	Role   entity.Role `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string, role entity.Role, userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	})

	// sign token and generate jwt string
	tokenStr, err := token.SignedString(sercretKey)
	if err != nil {
		fmt.Println("Error generating JWT token:", err) //log error generate token
		return "", err
	}

	return tokenStr, nil
}

func ParseJWT(tokenStr string) (*entity.User, error) {
	// parse jwt token with correct structure
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return sercretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	// get claims and check if token valid
	claims, ok := token.Claims.(*UserClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// check token expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// create user object based on claims
	user := &entity.User{
		UserID: claims.UserID,
		Role:   claims.Role,
	}
	// var user entity.User
	// if err := c.db.First(&user, claims.UserID).Error; err != nil {
	//     return nil, fmt.Errorf("user not found: %v", err)
	// }

	return user, nil
}

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
	user, err := ParseJWT(tokenStr)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authorization token",
			"code":  http.StatusUnauthorized,
		})

		ctx.Abort()
		return
	}

	// set user info in context
	ctx.Set("currentUser", &UserClaims{
		UserID: user.UserID,
		Role:   user.Role,
	})
	ctx.Next()
}
