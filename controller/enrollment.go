package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type EnrollController interface {
	StudentEnroll(ctx *gin.Context)
}

type EnrollControllerImpl struct {
	enrollService service.EnrollService
}

func NewEnrollController(enrollService service.EnrollService) EnrollController {
	return &EnrollControllerImpl{
		enrollService: enrollService,
	}
}

// create class (admin & mentor)
func (c *EnrollControllerImpl) StudentEnroll(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to enroll to a course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// bind json body with model
	var studentID struct {
		StudentID uint `json:"student_id"`
	}

	// validate with model req
	if err := ctx.ShouldBindJSON(&studentID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// enroll to a course
	enroll, err := c.enrollService.StudentEnroll(userClaims, courseID, fmt.Sprint(studentID.StudentID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// success response
	enrollResp := model.EnrollResp{
		EnrollmentID:   enroll.EnrollmentID,
		StudentID:      enroll.StudentID,
		CourseID:       enroll.CourseID,
		EnrollmentDate: nil,
		EnrollStatus:   entity.Pending,
		CreatedAt:      enroll.CreatedAt,
		UpdatedAt:      enroll.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("UserID %d enrollment request has been sent", enroll.StudentID),
		"data":    enrollResp,
	})
}

// admin only
func (c *EnrollControllerImpl) UpdateStudentEnroll(ctx *gin.Context) {

}

// mentor & admin || get a student by studentID & couseID
func (c *EnrollControllerImpl) StudentCourseEnroll(ctx *gin.Context) {

}

// mentor & admin || get students enroll to a course by courseID
