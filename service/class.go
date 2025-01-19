package service

import (
	"fmt"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
)

type ClassService interface {
	CreateClass(userClaims *middleware.UserClaims, courseID string, class model.CreateClass) (*entity.Class, error)
	GetClasses(courseID string) ([]entity.Class, error)
	GetClassByID(courseID, classID string) (*entity.Class, error)
	UpdateClassByID(userClaims *middleware.UserClaims, courseID, classID string, classReq model.UpdateClass) (*entity.Class, error)
	DeleteClassByID(userClaims *middleware.UserClaims, courseID, classID string) error
}

type ClassServiceImpl struct {
	classRepo  repository.ClassRepo
	courseRepo repository.CourseRepo
}

func NewClassService(classRepo repository.ClassRepo, courseRepo repository.CourseRepo) ClassService {
	return &ClassServiceImpl{
		classRepo:  classRepo,
		courseRepo: courseRepo,
	}
}

func (s *ClassServiceImpl) CreateClass(userClaims *middleware.UserClaims, courseID string, class model.CreateClass) (*entity.Class, error) {
	// only admin & mentor
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can create a new class")
	}

	// check if courseID exist
	existingCourse, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// check if user authorized
	if userClaims.Role == entity.Mentor {
		if existingCourse.MentorID != userClaims.UserID {
			return nil, fmt.Errorf("you have no access to create a new class in this course")
		}
	}

	// validate input className
	isValid, errMsg := middleware.ValidateCourseName(class.ClassName)
	if !isValid {
		return nil, errMsg
	}

	// validate input startDate, endDate
	isValid, errMsg = middleware.ValidateCourseDate(class.StartDate.Format("02-01-2006 15:04"), class.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		return nil, errMsg
	}

	// add new course to db
	newClass := entity.Class{
		ClassName:   class.ClassName,
		Description: class.Description,
		StartDate:   class.StartDate.Time,
		EndDate:     class.EndDate.Time,
		CourseID:    existingCourse.CourseID,
	}

	// create new class
	if err := s.classRepo.CreateClass(&newClass); err != nil {
		return nil, fmt.Errorf("unable to create a new class")
	}

	return &newClass, nil
}

func (s *ClassServiceImpl) GetClasses(courseID string) ([]entity.Class, error) {
	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	class, err := s.classRepo.GetClasses(courseID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch classes")
	}

	return class, nil
}

func (s *ClassServiceImpl) GetClassByID(courseID, classID string) (*entity.Class, error) {
	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	class, err := s.classRepo.GetClassByID(courseID, classID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch a class")
	}

	return class, nil
}

func (s *ClassServiceImpl) UpdateClassByID(userClaims *middleware.UserClaims, courseID, classID string, classReq model.UpdateClass) (*entity.Class, error) {
	// only admin & mentor can update course
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can create a new class")
	}

	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	existingClass, err := s.classRepo.GetClassByID(courseID, classID)
	if err != nil {
		return nil, fmt.Errorf("class_id %s not found", courseID)
	}

	if classReq.ClassName != "" {
		isValid, errMsg := middleware.ValidateCourseName(classReq.ClassName)
		if !isValid {
			return nil, errMsg
		}

		existingClass.ClassName = classReq.ClassName
	}

	isValid, errMsg := middleware.ValidateCourseDate(classReq.StartDate.Format("02-01-2006 15:04"), classReq.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		return nil, errMsg
	}

	if classReq.StartDate.IsZero() {
		existingClass.StartDate = classReq.StartDate.Time
	}

	if classReq.EndDate.IsZero() {
		existingClass.EndDate = classReq.EndDate.Time
	}

	if classReq.Description != "" {
		existingClass.Description = classReq.Description
	}

	// update class
	class, err := s.classRepo.UpdateClassByID(courseID, classID, *existingClass)
	if err != nil {
		return nil, fmt.Errorf("unable to update a class")
	}

	return class, nil
}

func (s *ClassServiceImpl) DeleteClassByID(userClaims *middleware.UserClaims, courseID, classID string) error {
	if userClaims.Role != entity.Admin {
		return fmt.Errorf("only admin can delete a class")
	}

	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return fmt.Errorf("course_id %s not found", courseID)
	}

	if _, err := s.classRepo.GetClassByID(courseID, classID); err != nil {
		return fmt.Errorf("class_id %s not found", classID)
	}

	if err := s.classRepo.DeleteClassByID(courseID, classID); err != nil {
		return fmt.Errorf("unable to delete a class")
	}

	return nil
}
