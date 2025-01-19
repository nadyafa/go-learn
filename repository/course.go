package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type CourseRepo interface {
	CreateCourse(course *entity.Course) error
	GetCourses() ([]entity.Course, error)
	GetCourseByID(courseID string) (*entity.Course, error)
	UpdateCourseByID(courseID string, course *entity.Course) error
	DeleteUserByID(courseID string) error
}

type CourseRepoImpl struct {
	db *gorm.DB
}

func NewCourseRepo(db *gorm.DB) CourseRepo {
	return &CourseRepoImpl{
		db: db,
	}
}

func (r *CourseRepoImpl) CreateCourse(course *entity.Course) error {
	if err := r.db.Create(&course).Error; err != nil {
		return err
	}

	return nil
}

func (r *CourseRepoImpl) GetCourses() ([]entity.Course, error) {
	var courses []entity.Course

	if err := r.db.Find(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
}

func (r *CourseRepoImpl) GetCourseByID(courseID string) (*entity.Course, error) {
	var course entity.Course

	if err := r.db.First(&course, courseID).Error; err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepoImpl) UpdateCourseByID(courseID string, course *entity.Course) error {
	if err := r.db.Where("course_id = ?", courseID).Updates(&course).Error; err != nil {
		return err
	}

	return nil
}

func (r *CourseRepoImpl) DeleteUserByID(courseID string) error {
	var course entity.Course

	if err := r.db.Where("course_id = ?", courseID).Delete(&course).Error; err != nil {
		return err
	}

	return nil
}
