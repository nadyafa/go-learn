package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nadyafa/go-learn/entity"
)

var sercretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserID,
		// "username": user.Username,
		// "email":    user.Email,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString(sercretKey)
	if err != nil {
		fmt.Println("Error generating JWT token:", err) //log error generate token
		return "", err
	}

	return tokenStr, nil
}

func ParseJWT(tokenStr string) (*entity.User, error) {
	// parse jwt token
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return sercretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// get claims and check if valid
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claimsMap := *claims

	// extract user info
	userIDFloat, ok := claimsMap["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("user_id claim is not a valid float64")
	}

	// extract current role and update role as the custom role provide
	roleStr, ok := claimsMap["role"].(string)
	if !ok {
		return nil, fmt.Errorf("role claim is not a valid string")
	}

	// convert role entity to str
	role := entity.Role(roleStr)

	user := &entity.User{
		UserID: uint(userIDFloat),
		Role:   role,
	}

	return user, nil
}
