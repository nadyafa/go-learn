package controller

import (
	"gorm.io/gorm"
)

type EnrollController interface {
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
// func (c *EnrollControllerImpl) StudentEnroll(ctx *gin.Context) {
// 	// check if the currentUser is admin
// 	claims, _ := ctx.Get("currentUser")
// 	userClaims, ok := claims.(*middleware.UserClaims)
// 	if !ok || userClaims.Role == entity.Student {
// 		ctx.JSON(http.StatusForbidden, gin.H{
// 			"error": "Access Restricted",
// 			"code":  http.StatusForbidden,
// 		})
// 		return
// 	}

// 	// make sure courseID exist
// 	courseID := ctx.Param("course_id")
// 	var course entity.Course
// 	if err := c.db.First(&course, courseID).Error; err != nil {
// 		ctx.JSON(http.StatusNotFound, gin.H{
// 			"message": fmt.Sprintf("CourseID %s not found", courseID),
// 			"code":    http.StatusNotFound,
// 		})
// 		return
// 	}

// 	// bind json body with model
// 	var classReq model.CreateClass

// 	if err := ctx.ShouldBindJSON(&classReq); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": err.Error(),
// 			"code":  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	// validate input className, startDate, endDate
// 	isValid, _ := middleware.ValidateCourseName(classReq.ClassName)
// 	if !isValid {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": "Class name cannot be empty",
// 			"code":  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	isValid, validationMsg := middleware.ValidateCourseDate(classReq.StartDate.Format("02-01-2006 15:04"), classReq.EndDate.Format("02-01-2006 15:04"))
// 	if !isValid {
// 		ctx.JSON(http.StatusBadRequest, gin.H{
// 			"error": validationMsg,
// 			"code":  http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	// a mentor only able to create class for themselves
// 	if userClaims.Role == entity.Mentor {
// 		classReq.MentorID = userClaims.UserID
// 	}

// 	// add new course to db
// 	class := entity.Class{
// 		ClassName:   classReq.ClassName,
// 		Description: classReq.Description,
// 		MentorID:    classReq.MentorID,
// 		StartDate:   classReq.StartDate.Time,
// 		EndDate:     classReq.EndDate.Time,
// 		CourseID:    course.CourseID,
// 	}

// 	if err := c.db.Create(&class).Error; err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err,
// 			"code":  http.StatusInternalServerError,
// 			"msg":   course.CourseID,
// 		})
// 		return
// 	}

// 	// success response
// 	classResp := model.ClassResp{
// 		ClassID:     class.ClassID,
// 		CourseID:    class.CourseID,
// 		ClassName:   class.ClassName,
// 		Description: class.Description,
// 		MentorID:    class.MentorID,
// 		StartDate:   class.StartDate,
// 		EndDate:     class.EndDate,
// 		CreatedAt:   class.CreatedAt,
// 		UpdatedAt:   class.UpdatedAt,
// 	}

// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"message": fmt.Sprintf("Class %s created successfully", class.ClassName),
// 		"data":    classResp,
// 	})
// }
