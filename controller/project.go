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

type ProjectController interface {
	CreateProject(ctx *gin.Context)
	GetProjects(ctx *gin.Context)
	GetProjectByID(ctx *gin.Context)
	UpdateProject(ctx *gin.Context)
	DeleteProjectByID(ctx *gin.Context)
}

type ProjectControllerImpl struct {
	db *gorm.DB
}

func NewProjectController(db *gorm.DB) ProjectController {
	return &ProjectControllerImpl{
		db: db,
	}
}

// create project (admin & mentor)
func (c *ProjectControllerImpl) CreateProject(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role == entity.Student {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin & mentor can create a new course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// validate courseID
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
	var projectReq model.CreateProject

	if err := ctx.ShouldBindJSON(&projectReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// validate input projectName
	isValid, validationMsg := middleware.ValidateCourseName(projectReq.ProjectName)
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
		return
	}

	// validate deadline project < endDate course
	isValid, validationMsg = middleware.ValidateCourseDate(projectReq.Deadline.Format("02-01-2006 15:04"), course.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": validationMsg,
			"code":  http.StatusBadRequest,
		})
		return
	}

	// add new course to db
	project := entity.Project{
		CourseID:    course.CourseID,
		ProjectName: projectReq.ProjectName,
		Description: projectReq.Description,
		Deadline:    projectReq.Deadline.Time,
	}

	if err := c.db.Create(&project).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create course",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// success response
	projectResp := model.ProjectResp{
		ProjectID:   project.ProjectID,
		CourseID:    project.CourseID,
		ProjectName: project.ProjectName,
		Description: project.Description,
		Deadline:    project.Deadline,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Course %s created successfully", course.CourseName),
		"user":    projectResp,
	})
}

// get all projects (for all)
func (c *ProjectControllerImpl) GetProjects(ctx *gin.Context) {
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
	var projects []entity.Project

	if err := c.db.Where("course_id = ?", courseID).Find(&projects).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to retrieve projects",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Classes fetch successfully",
		"code":    http.StatusOK,
		"data":    projects,
	})
}

// get project by id (for all)
func (c *ProjectControllerImpl) GetProjectByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You have to sign in to look at detail project",
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
	projectID := ctx.Param("project_id")
	var project entity.Project

	if err := c.db.Where("course_id = ?", courseID).First(&project, projectID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("ProjectID %s not found", projectID),
			"code":    http.StatusNotFound,
		})
		return
	}

	projectResp := model.ProjectResp{
		ProjectID:   project.ProjectID,
		CourseID:    project.CourseID,
		ProjectName: project.ProjectName,
		Description: project.Description,
		Deadline:    project.Deadline,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Class fetch successfully",
		"code":    http.StatusOK,
		"data":    projectResp,
	})
}

// update project (admin & mentor)
func (c *ProjectControllerImpl) UpdateProject(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role == entity.Student {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Only admin & mentor can create a new course",
			"code":  http.StatusForbidden,
		})
		return
	}

	// validate courseID
	courseID := ctx.Param("course_id")
	var course entity.Course
	if err := c.db.First(&course, courseID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("CourseID %s not found", courseID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// get projectID
	projectID := ctx.Param("project_id")
	var existingProject entity.Project

	if err := c.db.First(&existingProject, projectID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Project not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// bind json body with model
	var projectReq model.UpdateProject

	if err := ctx.ShouldBindJSON(&projectReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// update fields if not empty
	if projectReq.ProjectName != "" {
		existingProject.ProjectName = projectReq.ProjectName
	}

	if projectReq.Description != "" {
		existingProject.Description = projectReq.Description
	}

	if !projectReq.Deadline.IsZero() {
		// validate deadline project < endDate course
		isValid, validationMsg := middleware.ValidateCourseDate(projectReq.Deadline.Format("02-01-2006 15:04"), course.EndDate.Format("02-01-2006 15:04"))
		if !isValid {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": validationMsg,
				"code":  http.StatusBadRequest,
			})
			return
		}

		existingProject.Deadline = projectReq.Deadline.Time
	}

	// add new course to db
	project := entity.Project{
		CourseID:    course.CourseID,
		ProjectName: existingProject.ProjectName,
		Description: existingProject.Description,
		Deadline:    existingProject.Deadline,
	}

	if err := c.db.Create(&project).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create course",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// success response
	projectResp := model.ProjectResp{
		ProjectID:   project.ProjectID,
		CourseID:    project.CourseID,
		ProjectName: project.ProjectName,
		Description: project.Description,
		Deadline:    project.Deadline,
		CreatedAt:   project.CreatedAt,
		UpdatedAt:   project.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Course %s created successfully", course.CourseName),
		"user":    projectResp,
	})
}

// delete project (admin only)
func (c *ProjectControllerImpl) DeleteProjectByID(ctx *gin.Context) {
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

	// get classID
	projectID := ctx.Param("project_id")

	// find classID
	var project entity.Project
	if err := c.db.First(&project, projectID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "Project not found",
			"code":  http.StatusNotFound,
		})
		return
	}

	// update class to db
	if err := c.db.Delete(&project).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("Failed to delete project with ID %s", projectID),
			"code":    http.StatusInternalServerError,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ProjectID %s has been deleted", projectID),
		"code":    http.StatusOK,
	})
}
