package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type AttendRepo interface {
	CreateAttendance(attendClass entity.Attendance) (*entity.Attendance, error)
	GetClassAttendances(courseID, classID string) ([]entity.Attendance, error)
	DeleteAttendanceByID(courseID, classID, attendID string) error
}

type AttendRepoImpl struct {
	db *gorm.DB
}

func NewAttendRepo(db *gorm.DB) AttendRepo {
	return &AttendRepoImpl{
		db: db,
	}
}

func (r *AttendRepoImpl) CreateAttendance(attendClass entity.Attendance) (*entity.Attendance, error) {
	if err := r.db.Create(&attendClass).Error; err != nil {
		return nil, err
	}

	return &attendClass, nil
}

func (r *AttendRepoImpl) GetClassAttendances(courseID, classID string) ([]entity.Attendance, error) {
	var attendances []entity.Attendance

	if err := r.db.Where("course_id = ? AND class_id = ?", courseID, classID).Find(&attendances).Error; err != nil {
		return nil, err
	}

	return attendances, nil
}

func (r *AttendRepoImpl) DeleteAttendanceByID(courseID, classID, attendID string) error {
	var attendance entity.Attendance

	if err := r.db.Where("course_id = ? AND class_id = ?", courseID, classID).Delete(attendance).Error; err != nil {
		return err
	}

	return nil
}
