package middleware

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// hashing password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

// check if the hashed password match plain password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// validate input based on struct model
func ValidateUserInput(input interface{}, isSignup bool) map[string]string {
	errorMessage := make(map[string]string)

	// perform validation
	if err := validator.New().Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				if isSignup {
					errorMessage["username"] = "Username must be at least 6 characters alphanumerical"
				} else {
					errorMessage["username"] = "Username is required"
				}
			case "Email":
				if isSignup {
					errorMessage["email"] = "Invalid email input"
				} else {
					errorMessage["email"] = "Email is required"
				}
			case "Password":
				if isSignup {
					errorMessage["password"] = "Password must be at least 8 characters alphanumerical"
				} else {
					errorMessage["password"] = "Password is required"
				}
			}
		}
	}

	return errorMessage
}
