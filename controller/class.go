package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
	"gorm.io/gorm"
)

type ClassController interface {
	CreateClass(ctx *gin.Context)
	GetClasses(ctx *gin.Context)
	GetClassByID(ctx *gin.Context)
	UpdateClassByID(ctx *gin.Context)
	DeleteClassByID(ctx *gin.Context)
}

type ClassControllerImpl struct {
	db         *gorm.DB
	courseRepo repository.CourseRepo
}

func NewClassController(db *gorm.DB, courseRepo repository.CourseRepo) ClassController {
	return &ClassControllerImpl{
		db:         db,
		courseRepo: courseRepo,
	}
}

// create class (admin & mentor)
func (c *ClassControllerImpl) CreateClass(ctx *gin.Context) {
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

	// bind json body with model
	var classReq model.CreateClass

	if err := ctx.ShouldBindJSON(&classReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// validate input className, startDate, endDate
	isValid, _ := middleware.ValidateCourseName(classReq.ClassName)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Class name cannot be empty",
			"code":  http.StatusBadRequest,
		})
		return
	}

	isValid, validationMsg := middleware.ValidateCourseDate(classReq.StartDate.Format("02-01-2006 15:04"), classReq.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
		return
	}

	// // a mentor only able to create class for themselves
	// if userClaims.Role == entity.Mentor {
	// 	classReq.MentorID = userClaims.UserID
	// }

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

	// add new course to db
	class := entity.Class{
		ClassName:   classReq.ClassName,
		Description: classReq.Description,
		StartDate:   classReq.StartDate.Time,
		EndDate:     classReq.EndDate.Time,
		CourseID:    course.CourseID,
	}

	if err := c.db.Create(&class).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
			"code":  http.StatusInternalServerError,
			"msg":   course.CourseID,
		})
		return
	}

	// success response
	classResp := model.ClassResp{
		ClassID:     class.ClassID,
		CourseID:    class.CourseID,
		ClassName:   class.ClassName,
		Description: class.Description,
		StartDate:   class.StartDate,
		EndDate:     class.EndDate,
		CreatedAt:   class.CreatedAt,
		UpdatedAt:   class.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Class %s created successfully", class.ClassName),
		"data":    classResp,
	})
}

// get all classes (for all)
func (c *ClassControllerImpl) GetClasses(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You have to signin to open classes",
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

	// get classes
	var classes []entity.Class

	if err := c.db.Where("course_id = ?", courseID).Find(&classes).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve classes",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// create response
	var classResponses []model.ClassResp

	for _, class := range classes {
		classResp := model.ClassResp{
			ClassID:     class.ClassID,
			CourseID:    class.CourseID,
			ClassName:   class.ClassName,
			Description: class.Description,
			StartDate:   class.StartDate,
			EndDate:     class.EndDate,
			CreatedAt:   class.CreatedAt,
			UpdatedAt:   class.UpdatedAt,
		}

		classResponses = append(classResponses, classResp)
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Classes fetch successfully",
		"code":    http.StatusOK,
		"data":    classResponses,
	})
}

// get class by id (for all)
func (c *ClassControllerImpl) GetClassByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You have to sign in to look at classes",
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

	// get class
	classID := ctx.Param("class_id")
	var class entity.Class

	if err := c.db.Where("course_id = ?", courseID).First(&class, classID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("ClassID %s not found", classID),
			"code":    http.StatusNotFound,
		})
		return
	}

	classResp := model.ClassResp{
		ClassID:     class.ClassID,
		CourseID:    class.CourseID,
		ClassName:   class.ClassName,
		Description: class.Description,
		StartDate:   class.StartDate,
		EndDate:     class.EndDate,
		CreatedAt:   class.CreatedAt,
		UpdatedAt:   class.UpdatedAt,
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Class fetch successfully",
		"code":    http.StatusOK,
		"data":    classResp,
	})
}

// update class (admin & mentor)
func (c *ClassControllerImpl) UpdateClassByID(ctx *gin.Context) {
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

	// get clasID
	classID := ctx.Param("class_id")

	// find classID
	var existingClass entity.Class
	if err := c.db.First(&existingClass, classID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Class not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// get class body input
	var classReq model.UpdateClass
	if err := ctx.ShouldBindJSON(&classReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// validate date input
	isValid, validationMsg := middleware.ValidateCourseDate(classReq.StartDate.Format("02-01-2006 15:04"), classReq.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
		return
	}

	// update fields if not empty
	if classReq.ClassName != "" {
		existingClass.ClassName = classReq.ClassName
	}

	if classReq.Description != "" {
		existingClass.Description = classReq.Description
	}

	if classReq.StartDate.IsZero() {
		existingClass.StartDate = classReq.StartDate.Time
	}

	if classReq.EndDate.IsZero() {
		existingClass.EndDate = classReq.EndDate.Time
	}

	existingClass.UpdatedAt = time.Now()

	// a mentor only able to create class for themselves
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

	// update class to db
	if err := c.db.Save(&existingClass).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to update class with ID %s", classID),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	class := model.ClassResp{
		ClassID:     existingClass.ClassID,
		CourseID:    existingClass.CourseID,
		ClassName:   existingClass.ClassName,
		Description: existingClass.Description,
		StartDate:   existingClass.StartDate,
		EndDate:     existingClass.EndDate,
		CreatedAt:   existingClass.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ClassID %s updated successfully", classID),
		"code":    http.StatusOK,
		"data":    class,
	})
}

// delete class (admin only)
func (c *ClassControllerImpl) DeleteClassByID(ctx *gin.Context) {
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
	if err := c.db.Find(&course, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("CourseID %s not found", courseID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// get classID
	classID := ctx.Param("class_id")

	// find classID
	var class entity.Class
	if err := c.db.First(&class, classID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Class not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// update class to db
	if err := c.db.Delete(&class).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete class with ID %s", classID),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ClassID %s has been deleted", classID),
		"code":    http.StatusOK,
	})
}
