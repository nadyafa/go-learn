package controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type ClassController interface {
	CreateClass(ctx *gin.Context)
	GetClasses(ctx *gin.Context)
	GetClassByID(ctx *gin.Context)
	UpdateClassByID(ctx *gin.Context)
	DeleteClassByID(ctx *gin.Context)
}

type ClassControllerImpl struct {
	classService service.ClassService
}

func NewClassController(classService service.ClassService) ClassController {
	return &ClassControllerImpl{
		classService: classService,
	}
}

// create class (admin & mentor)
func (c *ClassControllerImpl) CreateClass(ctx *gin.Context) {
	// check if the user is signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to create a new class",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// bind json body with model
	var classReq model.CreateClass

	if err := ctx.ShouldBindJSON(&classReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// call service layer to create class
	class, err := c.classService.CreateClass(userClaims, courseID, classReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to get classes",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	classes, err := c.classService.GetClasses(courseID)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": err.Error(),
			"code":  http.StatusForbidden,
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
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to get a class",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get classID param
	classID := ctx.Param("class_id")

	class, err := c.classService.GetClassByID(courseID, classID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
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
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "User must sign in to update a class",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get classID param
	classID := ctx.Param("class_id")

	// get class body input
	var classReq model.UpdateClass
	if err := ctx.ShouldBindJSON(&classReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// update class
	class, err := c.classService.UpdateClassByID(userClaims, courseID, classID, classReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// succeed response
	classResp := model.ClassResp{
		ClassID:     class.ClassID,
		CourseID:    class.CourseID,
		ClassName:   class.ClassName,
		Description: class.Description,
		StartDate:   class.StartDate,
		EndDate:     class.EndDate,
		CreatedAt:   class.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ClassID %s updated successfully", classID),
		"code":    http.StatusOK,
		"data":    classResp,
	})
}

// delete class (admin only)
func (c *ClassControllerImpl) DeleteClassByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to delete a class",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get classID param
	classID := ctx.Param("class_id")

	// delete class
	if err := c.classService.DeleteClassByID(userClaims, courseID, classID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ClassID %s has been deleted", classID),
		"code":    http.StatusOK,
	})
}
