package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type AttendController interface {
	StudentAttendClass(ctx *gin.Context)
	GetClassAttendances(ctx *gin.Context)
	DeleteAttendanceByID(ctx *gin.Context)
}

type AttendControllerImpl struct {
	attendService service.AttendService
}

func NewAttendController(attendService service.AttendService) AttendController {
	return &AttendControllerImpl{
		attendService: attendService,
	}
}

// create course (admin & student)
func (c *AttendControllerImpl) StudentAttendClass(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to attend class",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get classID param
	classID := ctx.Param("class_id")

	// bind json body with model
	var attendReq model.AttendReq

	if err := ctx.ShouldBindJSON(&attendReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// create attendance
	attend, err := c.attendService.StudentAttendClass(userClaims, courseID, classID, attendReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// success response
	attendResp := model.AttendResp{
		AttendID:  attend.AttendID,
		StudentID: attend.StudentID,
		ClassID:   attend.ClassID,
		CourseID:  attend.CourseID,
		Attended:  attend.Attended,
		AttendAt:  attend.AttendAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Student successfully attend class",
		"user":    attendResp,
	})
}

// get list of students attend class (mentor & admin)
func (c *AttendControllerImpl) GetClassAttendances(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to get attendance lists",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// make sure courseID exist
	courseID := ctx.Param("course_id")

	// check if the class exist
	classID := ctx.Param("class_id")

	attendances, err := c.attendService.GetClassAttendances(userClaims, courseID, classID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// create attendance list response
	var attendResponses []model.AttendResp

	for _, attend := range attendances {
		attendResp := model.AttendResp{
			AttendID:  attend.AttendID,
			StudentID: attend.StudentID,
			ClassID:   attend.ClassID,
			CourseID:  attend.CourseID,
			Attended:  attend.Attended,
			AttendAt:  attend.AttendAt,
		}

		attendResponses = append(attendResponses, attendResp)
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "attendance lists fetch successfully",
		"code":    http.StatusOK,
		"data":    attendResponses,
	})
}

// delete student attendance (admin only)
func (c *AttendControllerImpl) DeleteAttendanceByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to delete attendance",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get classID param
	classID := ctx.Param("class_id")

	// get attendID param
	attendID := ctx.Param("attend_id")

	if err := c.attendService.DeleteAttendanceByID(userClaims, courseID, classID, attendID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Student attendance successfully deleted",
		"code":    http.StatusOK,
	})
}
