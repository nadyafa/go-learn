package service

import (
	"fmt"
	"time"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
)

type CourseService interface {
	CreateCourse(userClaims *middleware.UserClaims, courseReq model.CourseReq) (*entity.Course, error)
	GetCourses() ([]entity.Course, error)
	GetCourseByID(courseID string) (*entity.Course, error)
	UpdateCourseByID(userClaims *middleware.UserClaims, courseReq model.CourseReq, courseID string) (*entity.Course, error)
	DeleteCourseByID(userClaims *middleware.UserClaims, courseID string) error
}

type CourseServiceImpl struct {
	courseRepo repository.CourseRepo
}

func NewCourseService(courseRepo repository.CourseRepo) CourseService {
	return &CourseServiceImpl{
		courseRepo: courseRepo,
	}
}

func (s *CourseServiceImpl) CreateCourse(userClaims *middleware.UserClaims, courseReq model.CourseReq) (*entity.Course, error) {
	// userClaim role must be admin & mentor
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can create a new course")
	}

	// validate course input
	isValid, errMsg := middleware.ValidateCourseName(courseReq.CourseName)
	if !isValid {
		return nil, errMsg
	}

	isValid, errMsg = middleware.ValidateCourseDate(courseReq.StartDate.Format("02-01-2006 15:04"), courseReq.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		return nil, errMsg
	}

	if courseReq.MentorID == 0 {
		if userClaims.Role == entity.Mentor {
			courseReq.MentorID = userClaims.UserID
		} else {
			return nil, fmt.Errorf("mentor_id is required")
		}
	}

	// create course entity to repo layer
	course := entity.Course{
		CourseName:  courseReq.CourseName,
		Description: courseReq.Description,
		MentorID:    courseReq.MentorID,
		StartDate:   courseReq.StartDate.Time,
		EndDate:     courseReq.EndDate.Time,
	}

	if err := s.courseRepo.CreateCourse(&course); err != nil {
		return nil, fmt.Errorf("unable to create a new course")
	}

	return &course, nil
}

func (s *CourseServiceImpl) GetCourses() ([]entity.Course, error) {
	courses, err := s.courseRepo.GetCourses()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch list of courses")
	}

	return courses, nil
}

func (s *CourseServiceImpl) GetCourseByID(courseID string) (*entity.Course, error) {
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch a course")
	}

	return course, nil
}

func (s *CourseServiceImpl) UpdateCourseByID(userClaims *middleware.UserClaims, courseReq model.CourseReq, courseID string) (*entity.Course, error) {
	// only admin and mentor can update course
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can update a course")
	}

	// check if courseID exist
	existingCourse, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found")
	}

	// update course with value
	isValid, errMsg := middleware.ValidateCourseDate(existingCourse.StartDate.Format("02-01-2006 15:04"), existingCourse.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		return nil, errMsg
	}

	// update fields if not empty
	if courseReq.CourseName != "" {
		existingCourse.CourseName = courseReq.CourseName
	}

	if courseReq.Description != "" {
		existingCourse.Description = courseReq.Description
	}

	if courseReq.MentorID != 0 {
		existingCourse.MentorID = courseReq.MentorID
	}

	if !courseReq.StartDate.IsZero() {
		existingCourse.StartDate = courseReq.StartDate.Time
	}

	if !courseReq.EndDate.IsZero() {
		existingCourse.EndDate = courseReq.EndDate.Time
	}

	existingCourse.UpdatedAt = time.Now()

	// update course to db
	if err := s.courseRepo.UpdateCourseByID(courseID, existingCourse); err != nil {
		return nil, err
	}

	return existingCourse, nil
}

func (s *CourseServiceImpl) DeleteCourseByID(userClaims *middleware.UserClaims, courseID string) error {
	// only admin can have access
	if userClaims.Role == entity.Student {
		return fmt.Errorf("only admin can delete a course")
	}

	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return fmt.Errorf("course not found")
	}

	if err := s.courseRepo.DeleteUserByID(courseID); err != nil {
		return fmt.Errorf("unable to delete course")
	}

	return nil
}
