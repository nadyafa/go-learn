package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type CourseController interface {
	CreateCourse(ctx *gin.Context)
	GetCourses(ctx *gin.Context)
	GetCourseByID(ctx *gin.Context)
	UpdateCourseByID(ctx *gin.Context)
	DeleteCourseByID(ctx *gin.Context)
}

type CourseControllerImpl struct {
	courseService service.CourseService
}

func NewCourseController(courseService service.CourseService) CourseController {
	return &CourseControllerImpl{
		courseService: courseService,
	}
}

// create course (admin only & mentor)
func (c *CourseControllerImpl) CreateCourse(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to create a new course",
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

	// call createCourse service layer
	course, err := c.courseService.CreateCourse(userClaims, courseReq)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
			"code":  http.StatusForbidden,
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

// get all course (for all)
func (c *CourseControllerImpl) GetCourses(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to get list of courses",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courses
	courses, err := c.courseService.GetCourses()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch courses",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Courses fetch successfully",
		"data":    courses,
	})
}

// get course by id (for all)
func (c *CourseControllerImpl) GetCourseByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to get a course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID
	courseID := ctx.Param("course_id")

	// get courses
	course, err := c.courseService.GetCourseByID(courseID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch courses",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course fetch successfully",
		"data":    course,
	})
}

// update course (admin & mentor)
func (c *CourseControllerImpl) UpdateCourseByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to update a course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID
	courseID := ctx.Param("course_id")

	// get course body input
	var courseReq model.CourseReq
	if err := ctx.ShouldBindJSON(&courseReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// update course to db
	course, err := c.courseService.UpdateCourseByID(userClaims, courseReq, courseID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// succeed response
	courseResp := model.CourseResp{
		CourseID:    course.CourseID,
		CourseName:  course.CourseName,
		Description: course.Description,
		StartDate:   course.StartDate,
		EndDate:     course.EndDate,
		CreatedAt:   course.CreatedAt,
		UpdatedAt:   course.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Course updated successfully",
		"data":    courseResp,
	})
}

// delete course (admin only)
func (c *CourseControllerImpl) DeleteCourseByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to create a new course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID
	courseID := ctx.Param("course_id")

	// delete course by courseID
	if err := c.courseService.DeleteCourseByID(userClaims, courseID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("CourseID %s has been deleted", courseID),
		"code":    http.StatusOK,
	})
}
