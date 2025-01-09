package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"gorm.io/gorm"
)

type CourseController interface {
	CreateCourse(ctx *gin.Context)
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
	isValid, validationMsg := middleware.ValidateCourse(courseReq.CourseName, courseReq.StartDate.Format("02-01-2006 15:04"), courseReq.EndDate.Format("02-01-2006 15:04"))
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
// get course by id (admin & mentor)
// update course (admin & mentor)
// delete course (admin only)
