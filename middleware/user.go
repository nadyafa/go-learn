package middleware

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/model"
	"gorm.io/gorm"
)

func ValidateUsername(db *gorm.DB, username string) error {
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(username) {
		return errors.New("username must be alphanumeric")
	}

	username = strings.ToLower(username)

	var user entity.User
	if err := db.Where("LOWER(username) = ?", username).First(&user).Error; err == nil {
		return errors.New("username already taken")
	}

	return nil
}

func ValidateRole(role string) error {
	switch entity.Role(role) {
	case entity.Student, entity.Admin, entity.Mentor:
		return nil
	default:
		return errors.New("invalid role: must be 'student', 'admin', or 'mentor'")
	}
}

func ValidateUserSignup(db *gorm.DB, userSignup model.UserSignup) error {
	validate := validator.New()

	if err := validate.Struct(userSignup); err != nil {
		return err
	}

	if err := ValidateUsername(db, userSignup.Username); err != nil {
		return err
	}

	if err := ValidateRole(userSignup.Role); err != nil {
		return err
	}

	return nil
}
