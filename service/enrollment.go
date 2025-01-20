package service

import (
	"fmt"
	"time"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/repository"
)

type EnrollService interface {
	StudentEnroll(userClaims *middleware.UserClaims, courseID, studentID string) (*entity.Enrollment, error)
	UpdateStudentEnroll(userClaims *middleware.UserClaims, courseID, studentID string, enrollStatus entity.Status) (*entity.Enrollment, error)
}

type EnrollServiceImpl struct {
	courseRepo repository.CourseRepo
	enrollRepo repository.EnrollRepo
	userRepo   repository.UserRepo
}

func NewEnrollService(courseRepo repository.CourseRepo, enrollRepo repository.EnrollRepo, userRepo repository.UserRepo) EnrollService {
	return &EnrollServiceImpl{
		courseRepo: courseRepo,
		enrollRepo: enrollRepo,
		userRepo:   userRepo,
	}
}

func (s *EnrollServiceImpl) StudentEnroll(userClaims *middleware.UserClaims, courseID, studentID string) (*entity.Enrollment, error) {
	// user & admin only
	if userClaims.Role == entity.Mentor {
		return nil, fmt.Errorf("you can't perform this action")
	}

	// check if courseID exist
	course, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found")
	}

	if userClaims.Role == entity.Student {
		studentID = fmt.Sprint(userClaims.UserID)
	} else {
		if studentID == "" {
			return nil, fmt.Errorf("user_id is required")
		}
	}

	// check if user exist
	userExist, err := s.userRepo.GetUserByID(studentID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// check if student already enroll to a course
	existingEnroll, err := s.enrollRepo.GetStudentCourseEnroll(studentID, courseID)
	if err == nil {
		return nil, fmt.Errorf("student has enroll with enrollment_id %d", existingEnroll.EnrollmentID)
	}

	// create student enrollment entity
	enroll := entity.Enrollment{
		StudentID:    userExist.UserID,
		CourseID:     course.CourseID,
		EnrollStatus: entity.Pending,
	}

	// enroll student
	newEnroll, err := s.enrollRepo.StudentEnroll(enroll)
	if err != nil {
		return nil, fmt.Errorf("student unable to enroll")
	}

	// // notify student
	// if userClaims.Role == entity.Student {
	// 	if err := middleware.SendMail(
	// 		userExist.Email,
	// 		"Go-Learn: Course Enrollment",
	// 		fmt.Sprintf("You have successfully signed to a courseID %s. You enrollment status currently on pending. We will soon notified you once it verified. Good luck!", courseID),
	// 	); err != nil {
	// 		return nil, fmt.Errorf("failed to send notification to student: %v", err)
	// 	}
	// }

	// // notify admin
	// if userClaims.Role == entity.Admin {
	// 	if err := middleware.SendMail(
	// 		os.Getenv("ADMIN_EMAIL"),
	// 		"A New User Course Enrollment",
	// 		fmt.Sprintf("UserID %d has signed to a courseID %s. Please validate their enrollment status.", userExist.UserID, courseID),
	// 	); err != nil {
	// 		return nil, fmt.Errorf("failed to send notification to admin: %v", err)
	// 	}
	// }

	return newEnroll, nil
}

func (s *EnrollServiceImpl) UpdateStudentEnroll(userClaims *middleware.UserClaims, courseID, studentID string, enrollStatus entity.Status) (*entity.Enrollment, error) {
	// admin only
	if userClaims.Role != entity.Admin {
		return nil, fmt.Errorf("only admin can update enrollment data")
	}

	// check if courseID exist
	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return nil, fmt.Errorf("course not found")
	}

	// check if user exist
	_, err := s.userRepo.GetUserByID(studentID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// check if student already enroll to a course
	existingEnroll, err := s.enrollRepo.GetStudentCourseEnroll(studentID, courseID)
	if err != nil {
		return nil, fmt.Errorf("student hasn't enroll to a course. courseid: %s, studentid: %s. err: %v", courseID, studentID, err.Error())
	}

	// create student enrollment entity
	updateEnroll := entity.Enrollment{
		EnrollmentID:   existingEnroll.EnrollmentID,
		StudentID:      existingEnroll.StudentID,
		CourseID:       existingEnroll.CourseID,
		EnrollmentDate: time.Now(),
		EnrollStatus:   enrollStatus,
		CreatedAt:      existingEnroll.CreatedAt,
		UpdatedAt:      time.Now(),
	}

	// update enrollmentStatus
	enroll, err := s.enrollRepo.UpdateStudentEnroll(courseID, studentID, updateEnroll)
	if err != nil {
		return nil, fmt.Errorf("unable to update student status enrollment")
	}

	// // notify student
	// if userClaims.Role == entity.Student {
	// 	if err := middleware.SendMail(
	// 		userExist.Email,
	// 		"Go-Learn: Course Enrollment",
	// 		fmt.Sprintf("You have successfully signed to a courseID %s. You enrollment status currently on pending. We will soon notified you once it verified. Good luck!", courseID),
	// 	); err != nil {
	// 		return nil, fmt.Errorf("failed to send notification to student: %v", err)
	// 	}
	// }

	return enroll, nil
}
