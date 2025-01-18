package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type EnrollRepo interface {
	StudentEnroll(enroll entity.Enrollment) (*entity.Enrollment, error)
	StudentCourseEnroll(courseID, userID string) (*entity.Enrollment, error)
}

type EnrollRepoImpl struct {
	db *gorm.DB
}

func NewEnrollRepo(db *gorm.DB) EnrollRepo {
	return &EnrollRepoImpl{
		db: db,
	}
}

func (r *EnrollRepoImpl) StudentEnroll(enroll entity.Enrollment) (*entity.Enrollment, error) {
	if err := r.db.Create(&enroll).Error; err != nil {
		return nil, err
	}

	return &enroll, nil
}

func (r *EnrollRepoImpl) StudentCourseEnroll(courseID, userID string) (*entity.Enrollment, error) {
	var studentEnroll entity.Enrollment

	if err := r.db.Where("user_id = ? AND course_id = ?", userID, courseID).First(&studentEnroll).Error; err != nil {
		return nil, err
	}

	return &studentEnroll, nil
}
