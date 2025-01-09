package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"gorm.io/gorm"
)

type CourseController interface {
	CreateCourse(ctx *gin.Context)
	GetCourses(ctx *gin.Context)
	GetCourseByID(ctx *gin.Context)
	UpdateCourseByID(ctx *gin.Context)
	DeleteCourseByID(ctx *gin.Context)
}

type CourseControllerImpl struct {
	db *gorm.DB
}

func NewCourseController(db *gorm.DB) CourseController {
	return &CourseControllerImpl{
		db: db,
	}
}

// create course (admin only)
func (c *CourseControllerImpl) CreateCourse(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Admin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin can create a new course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// bind json body with model
	var courseReq model.CourseReq

	if err := ctx.ShouldBindJSON(&courseReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// validate input courseName, startDate, endDate
	isValid, validationMsg := middleware.ValidateCourseName(courseReq.CourseName)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
	}

	isValid, validationMsg = middleware.ValidateCourseDate(courseReq.StartDate.Format("02-01-2006 15:04"), courseReq.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
	}

	// add new course to db
	course := entity.Course{
		CourseName:  courseReq.CourseName,
		Description: courseReq.Description,
		StartDate:   courseReq.StartDate.Time,
		EndDate:     courseReq.EndDate.Time,
	}

	if err := c.db.Create(&course).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create course",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// success response
	courseResp := model.CourseResp{
		CourseID:    course.CourseID,
		CourseName:  course.CourseName,
		Description: course.Description,
		StartDate:   course.StartDate,
		EndDate:     course.EndDate,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Course %s created successfully", course.CourseName),
		"user":    courseResp,
	})
}

// get all course (admin only)
func (c *CourseControllerImpl) GetCourses(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Admin {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin can create a new course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courses
	var courses []entity.Course

	if err := c.db.Find(&courses).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve courses",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Courses fetch successfully",
		"code":    http.StatusOK,
		"data":    courses,
	})
}

// get course by id (for all)
func (c *CourseControllerImpl) GetCourseByID(ctx *gin.Context) {
	// get courseID
	courseID := ctx.Param("course_id")

	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin can create a new course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courses
	var course []entity.Course

	if err := c.db.Find(&course, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("CourseID %s not found", courseID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course fetch successfully",
		"code":    http.StatusOK,
		"data":    course,
	})
}

// update course (admin & mentor)
func (c *CourseControllerImpl) UpdateCourseByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role == entity.Student {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Access Restricted",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID
	courseID := ctx.Param("course_id")

	// find courseID
	var existingCourse entity.Course
	if err := c.db.First(&existingCourse, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// get course body input
	var courseReq model.CourseReq
	if err := ctx.ShouldBindJSON(&courseReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// update course with value
	isValid, validationMsg := middleware.ValidateCourseDate(courseReq.StartDate.Format("02-01-2006 15:04"), courseReq.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
	}

	// update fields if not empty
	if courseReq.CourseName != "" {
		existingCourse.CourseName = courseReq.CourseName
	}

	if courseReq.Description != "" {
		existingCourse.Description = courseReq.Description
	}

	if courseReq.StartDate.IsZero() {
		existingCourse.StartDate = courseReq.StartDate.Time
	}

	if courseReq.EndDate.IsZero() {
		existingCourse.EndDate = courseReq.EndDate.Time
	}

	existingCourse.UpdatedAt = time.Now()

	// update course to db
	if err := c.db.Save(&existingCourse).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to update course with ID %s", courseID),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	course := model.CourseResp{
		CourseID:    existingCourse.CourseID,
		CourseName:  existingCourse.CourseName,
		Description: existingCourse.Description,
		StartDate:   existingCourse.StartDate,
		EndDate:     existingCourse.EndDate,
		CreatedAt:   existingCourse.CreatedAt,
		UpdatedAt:   existingCourse.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course updated successfully",
		"code":    http.StatusOK,
		"data":    course,
	})
}

// delete course (admin only)
func (c *CourseControllerImpl) DeleteCourseByID(ctx *gin.Context) {
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

	// get courseID
	courseID := ctx.Param("course_id")

	// find courseID
	var course entity.Course
	if err := c.db.First(&course, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// update course to db
	if err := c.db.Delete(&course).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete course with ID %s", courseID),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("CourseID %s has been deleted", courseID),
		"code":    http.StatusOK,
	})
}
