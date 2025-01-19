package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
	"gorm.io/gorm"
)

type AttendController interface {
	StudentAttendClass(ctx *gin.Context)
	GetClassAttendances(ctx *gin.Context)
	DeleteAttendanceByID(ctx *gin.Context)
}

type AttendControllerImpl struct {
	db         *gorm.DB
	courseRepo repository.CourseRepo
}

func NewAttendController(db *gorm.DB, courseRepo repository.CourseRepo) AttendController {
	return &AttendControllerImpl{
		db:         db,
		courseRepo: courseRepo,
	}
}

// create course (admin & student)
func (c *AttendControllerImpl) StudentAttendClass(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role == entity.Mentor {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin & students can add attendance",
			"code":  http.StatusForbidden,
		})
		return
	}

	// make sure courseID exist
	courseIDStr := ctx.Param("course_id")
	courseID, err := strconv.ParseUint(courseIDStr, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "CourseID must be a positive number",
		})
	}

	var course entity.Course
	if err := c.db.First(&course, courseIDStr).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("CourseID %s not found", courseIDStr),
			"code":    http.StatusNotFound,
		})
		return
	}

	// validate classID
	classIDStr := ctx.Param("class_id")
	classID, err := strconv.ParseUint(classIDStr, 10, 0)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ClassID must be a positive number",
		})
	}

	var class entity.Class
	if err := c.db.First(&class, classIDStr).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("ClassID %s not found", classIDStr),
			"code":    http.StatusNotFound,
		})
		return
	}

	// bind json body with model
	var attendReq model.AttendReq

	if userClaims.Role == entity.Student {
		attendReq.StudentID = userClaims.UserID
	}

	if err := ctx.ShouldBindJSON(&attendReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// validate input
	// attend, err := middleware.ValidateBoolean(attendReq.Attended)
	// if err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err,
	// 	})
	// }

	// add new course to db
	attend := entity.Attendance{
		ClassID:   uint(classID),
		StudentID: userClaims.UserID,
		CourseID:  uint(courseID),
		Attended:  attendReq.Attended,
		AttendAt:  time.Now(),
	}

	if err := c.db.Create(&attend).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process attend class",
			"code":  http.StatusInternalServerError,
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
		"message": "Class attendance successfully recorded",
		"user":    attendResp,
	})
}

// get list of students attend class (mentor & admin)
func (c *AttendControllerImpl) GetClassAttendances(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role == entity.Student {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Access Restricted",
			"code":  http.StatusForbidden,
		})
		return
	}

	// make sure courseID exist
	courseID := ctx.Param("course_id")
	var course entity.Course
	if err := c.db.First(&course, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("CourseID %s not found", courseID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// check if the class exist
	classID := ctx.Param("class_id")
	var class entity.Class

	if err := c.db.Find(&class, classID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("ClassID %s not found", classID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// verify if the mentor own the course
	existingCourse, err := c.courseRepo.GetCourseByID(courseID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": fmt.Sprintf("CourseID %s not found", courseID),
			"code":  http.StatusNotFound,
		})
		return
	}

	if userClaims.Role == entity.Mentor {
		if existingCourse.MentorID != userClaims.UserID {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "you have no access to create new class to this course",
				"code":  http.StatusBadRequest,
			})
			return
		}
	}

	// get attendances model
	var attendances []entity.Attendance

	if err := c.db.Where("class_id = ?", classID).Find(&attendances).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve list class attendances",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Courses fetch successfully",
		"code":    http.StatusOK,
		"data":    attendances,
	})
}

// delete student attendance (admin only)
func (c *AttendControllerImpl) DeleteAttendanceByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Admin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Access Restricted",
			"code":  http.StatusForbidden,
		})
		return
	}

	// make sure courseID exist
	courseID := ctx.Param("course_id")
	var course entity.Course
	if err := c.db.First(&course, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("CourseID %s not found", courseID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// check if the class exist
	classID := ctx.Param("class_id")
	var class entity.Class

	if err := c.db.Find(&class, classID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("ClassID %s not found", classID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// get courseID
	attendID := ctx.Param("attend_id")

	// find courseID
	var attendance entity.Attendance
	if err := c.db.First(&attendance, attendID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// update course to db
	if err := c.db.Delete(&attendID).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete student with attendID %s", attendID),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Student with attendID %s has been deleted", attendID),
		"code":    http.StatusOK,
	})
}
