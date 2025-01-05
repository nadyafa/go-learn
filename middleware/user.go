package middleware

import (
	"errors"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func ValidateUsername(username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(username) {
		return errors.New("username must be alphanumeric")
	}

	return nil
}

func ValidatePassword(password string) error {
	// Check for at least one lowercase letter
	lowercaseRegex := regexp.MustCompile(`[a-z]`)
	if !lowercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	// Check for at least one uppercase letter
	uppercaseRegex := regexp.MustCompile(`[A-Z]`)
	if !uppercaseRegex.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	// Check for at least one digit
	numberRegex := regexp.MustCompile(`\d`)
	if !numberRegex.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	// Check for at least one special character
	specialCharRegex := regexp.MustCompile(`[\W_]+`)
	if !specialCharRegex.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	// Check for minimum length of 8 characters
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

// func ValidateRole(role string) error {
// 	switch entity.Role(role) {
// 	case entity.Student, entity.Mentor:
// 		return nil
// 	default:
// 		return errors.New("invalid role: must be 'student' or 'mentor'")
// 	}
// }

// func ValidateUserSignup(db *gorm.DB, userSignup model.UserSignup) error {
// 	validate := validator.New()

// 	if err := validate.Struct(userSignup); err != nil {
// 		return err
// 	}

// 	if err := ValidateUsername(db, userSignup.Username); err != nil {
// 		return err
// 	}

// 	if err := ValidateRole(userSignup.Role); err != nil {
// 		return err
// 	}

// 	return nil
// }

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
