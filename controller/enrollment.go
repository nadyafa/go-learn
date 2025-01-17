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

type EnrollController interface {
	StudentEnroll(ctx *gin.Context)
}

type EnrollControllerImpl struct {
	db *gorm.DB
}

func NewEnrollController(db *gorm.DB) EnrollController {
	return &EnrollControllerImpl{
		db: db,
	}
}

// create class (admin & mentor)
func (c *EnrollControllerImpl) StudentEnroll(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role == entity.Mentor {
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

	// bind json body with model
	var enrollReq model.CreateEnrollment

	// validate UserID
	if enrollReq.UserID == 0 {
		if userClaims.Role != entity.Admin {
			enrollReq.UserID = userClaims.UserID
			enrollReq.UserRole = userClaims.Role
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "UserID and UserRole cannot be empty",
				"code":  http.StatusBadRequest,
			})
			return
		}
	}

	// validate with model req
	if err := ctx.ShouldBindJSON(&enrollReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
			"msg":   "bad request",
		})
		return
	}

	// add new user enroll to db
	enroll := entity.Enrollment{
		UserID:   enrollReq.UserID,
		UserRole: enrollReq.UserRole,
		CourseID: course.CourseID,
	}

	if err := c.db.Create(&enroll).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"code":  http.StatusInternalServerError,
			"msg":   course.CourseID,
		})
		return
	}

	// success response
	enrollResp := model.EnrollResp{
		EnrollmentID:   enroll.EnrollmentID,
		UserID:         enroll.UserID,
		UserRole:       enroll.UserRole,
		CourseID:       enroll.CourseID,
		EnrollmentDate: nil,
		EnrollStatus:   entity.Pending,
		CreatedAt:      enroll.CreatedAt,
		UpdatedAt:      enroll.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("UserID %d enrollment request has been sent", enroll.UserID),
		"data":    enrollResp,
	})
}
