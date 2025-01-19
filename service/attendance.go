package service

import (
	"fmt"
	"time"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
)

type AttendService interface {
	StudentAttendClass(userClaims *middleware.UserClaims, courseID, classID string, attendReq model.AttendReq) (*entity.Attendance, error)
	GetClassAttendances(userClaims *middleware.UserClaims, courseID, classID string) ([]entity.Attendance, error)
	DeleteAttendanceByID(userClaims *middleware.UserClaims, courseID, classID, attendID string) error
}

type AttendServiceImpl struct {
	attendRepo repository.AttendRepo
	courseRepo repository.CourseRepo
	classRepo  repository.ClassRepo
	enrollRepo repository.EnrollRepo
}

func NewAttendService(attendRepo repository.AttendRepo, courseRepo repository.CourseRepo, classRepo repository.ClassRepo, enrollRepo repository.EnrollRepo) AttendService {
	return &AttendServiceImpl{
		attendRepo: attendRepo,
		courseRepo: courseRepo,
		classRepo:  classRepo,
		enrollRepo: enrollRepo,
	}
}

func (s *AttendServiceImpl) StudentAttendClass(userClaims *middleware.UserClaims, courseID, classID string, attendReq model.AttendReq) (*entity.Attendance, error) {
	if userClaims.Role == entity.Mentor {
		return nil, fmt.Errorf("only admin & student can post attend class")
	}

	// check if course exist
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// check if class exist
	class, err := s.classRepo.GetClassByID(courseID, classID)
	if err != nil {
		return nil, fmt.Errorf("class_id %s not found", classID)
	}

	if attendReq.StudentID == 0 {
		if userClaims.Role == entity.Student {
			attendReq.StudentID = userClaims.UserID
		} else {
			return nil, fmt.Errorf("student_id is required")
		}
	}

	// make sure user is enrolled in course
	enroll, err := s.enrollRepo.GetStudentCourseEnroll(courseID, fmt.Sprint(attendReq.StudentID))
	if err != nil {
		return nil, fmt.Errorf("student not enroll to course")
	}

	if enroll.EnrollStatus != entity.Enroll {
		return nil, fmt.Errorf("student not enroll to course")
	}

	// create entity attendance
	newAttend := entity.Attendance{
		StudentID: attendReq.StudentID,
		ClassID:   class.ClassID,
		CourseID:  course.CourseID,
		Attended:  true,
		AttendAt:  time.Now(),
	}

	// create attendance
	attend, err := s.attendRepo.CreateAttendance(newAttend)
	if err != nil {
		return nil, fmt.Errorf("unable to create attendance")
	}

	return attend, nil
}

func (s *AttendServiceImpl) GetClassAttendances(userClaims *middleware.UserClaims, courseID, classID string) ([]entity.Attendance, error) {
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can see list of attendance")
	}

	// check if course exist
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// check if class exist
	if _, err := s.classRepo.GetClassByID(courseID, classID); err != nil {
		return nil, fmt.Errorf("class_id %s not found", classID)
	}

	// check if mentor own the course
	if userClaims.Role == entity.Mentor {
		if course.MentorID != userClaims.UserID {
			return nil, fmt.Errorf("you prohibit to do this action")
		}
	}

	// get list of attendances
	attendances, err := s.attendRepo.GetClassAttendances(courseID, classID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch attendance lists")
	}

	return attendances, nil
}

func (s *AttendServiceImpl) DeleteAttendanceByID(userClaims *middleware.UserClaims, courseID, classID, attendID string) error {
	if userClaims.Role != entity.Admin {
		return fmt.Errorf("only admin can delete attendance")
	}

	// check if course exist
	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return fmt.Errorf("course_id %s not found", courseID)
	}

	// check if class exist
	if _, err := s.classRepo.GetClassByID(courseID, classID); err != nil {
		return fmt.Errorf("class_id %s not found", classID)
	}

	// delete attendance
	if err := s.attendRepo.DeleteAttendanceByID(courseID, classID, attendID); err != nil {
		return fmt.Errorf("unable to delete attendance")
	}

	return nil
}
