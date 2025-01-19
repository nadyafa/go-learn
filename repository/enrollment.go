package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type EnrollRepo interface {
	StudentEnroll(enroll entity.Enrollment) (*entity.Enrollment, error)
	GetStudentCourseEnroll(courseID, userID string) (*entity.Enrollment, error)
	UpdateStudentEnroll(courseID, userID string, updateEnroll entity.Enrollment) (*entity.Enrollment, error)
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

func (r *EnrollRepoImpl) GetStudentCourseEnroll(courseID, studentID string) (*entity.Enrollment, error) {
	var studentEnroll entity.Enrollment

	if err := r.db.First(&studentEnroll, studentID, courseID).Error; err != nil {
		return nil, err
	}

	return &studentEnroll, nil
}

func (r *EnrollRepoImpl) UpdateStudentEnroll(courseID, studentID string, updateEnroll entity.Enrollment) (*entity.Enrollment, error) {
	if err := r.db.Where("student_id = ? AND course_id = ?", studentID, courseID).Updates(updateEnroll).Error; err != nil {
		return nil, err
	}

	return &updateEnroll, nil
}
