package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type ClassRepo interface {
	CreateClass(class *entity.Class) error
	GetClasses(courseID string) ([]entity.Class, error)
	GetClassByID(courseID, classID string) (*entity.Class, error)
	UpdateClassByID(courseID, classID string, class entity.Class) (*entity.Class, error)
	DeleteClassByID(courseID, classID string) error
}

type ClassRepoImpl struct {
	db *gorm.DB
}

func NewClassRepo(db *gorm.DB) ClassRepo {
	return &ClassRepoImpl{
		db: db,
	}
}

func (r *ClassRepoImpl) CreateClass(class *entity.Class) error {
	if err := r.db.Create(&class).Error; err != nil {
		return err
	}

	return nil
}

func (r *ClassRepoImpl) GetClasses(courseID string) ([]entity.Class, error) {
	var classes []entity.Class

	if err := r.db.Where("course_id = ?", courseID).Find(&classes).Error; err != nil {
		return nil, err
	}

	return classes, nil
}

func (r *ClassRepoImpl) GetClassByID(courseID, classID string) (*entity.Class, error) {
	var class entity.Class

	if err := r.db.Where("course_id = ? AND class_id = ?", courseID, classID).First(&class).Error; err != nil {
		return nil, err
	}

	return &class, nil
}

func (r *ClassRepoImpl) UpdateClassByID(courseID, classID string, class entity.Class) (*entity.Class, error) {
	if err := r.db.Where("course_id = ? AND class_id = ?", courseID, classID).Updates(class).Error; err != nil {
		return nil, err
	}

	return &class, nil
}

func (r *ClassRepoImpl) DeleteClassByID(courseID, classID string) error {
	var class entity.Class

	if err := r.db.Where("course_id = ? AND class_id = ?", courseID, classID).Delete(class).Error; err != nil {
		return err
	}

	return nil
}
