package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/service"
)

type ProjectController interface {
	CreateProject(ctx *gin.Context)
	GetProjects(ctx *gin.Context)
	GetProjectByID(ctx *gin.Context)
	UpdateProjectByID(ctx *gin.Context)
	DeleteProjectByID(ctx *gin.Context)
}

type ProjectControllerImpl struct {
	projectService service.ProjectService
}

func NewProjectController(projectService service.ProjectService) ProjectController {
	return &ProjectControllerImpl{
		projectService: projectService,
	}
}

// create project (admin & mentor)
func (c *ProjectControllerImpl) CreateProject(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to create a new project",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// validate courseID
	courseID := ctx.Param("course_id")

	// bind json body with model
	var projectReq model.CreateProject

	if err := ctx.ShouldBindJSON(&projectReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// call service layer createProject
	project, err := c.projectService.CreateProject(userClaims, courseID, projectReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
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
		"message": fmt.Sprintf("ProjectID %s created successfully", fmt.Sprint(project.ProjectID)),
		"user":    projectResp,
	})
}

// get all projects (for all)
func (c *ProjectControllerImpl) GetProjects(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "User must sign in to get project lists",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// make sure courseID exist
	courseID := ctx.Param("course_id")

	// get projects
	projects, err := c.projectService.GetProjects(courseID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User must sign in to get project lists",
			"code":  http.StatusBadRequest,
		})
		return
	}

	// create response
	var projectResponses []model.ProjectResp

	for _, project := range projects {
		projectResp := model.ProjectResp{
			ProjectID:   project.ProjectID,
			CourseID:    project.CourseID,
			ProjectName: project.ProjectName,
			Description: project.Description,
			Deadline:    project.Deadline,
			CreatedAt:   project.CreatedAt,
			UpdatedAt:   project.UpdatedAt,
		}

		projectResponses = append(projectResponses, projectResp)
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Projects fetch successfully",
		"data":    projectResponses,
	})
}

// get project by id (for all)
func (c *ProjectControllerImpl) GetProjectByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	_, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You have to sign in to look at detail project",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get projectID param
	projectID := ctx.Param("project_id")

	// get projectID
	project, err := c.projectService.GetProjectByID(courseID, projectID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
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
		"message": fmt.Sprintf("ProjectID %s fetch successfully", projectID),
		"data":    projectResp,
	})
}

// update project (admin & mentor)
func (c *ProjectControllerImpl) UpdateProjectByID(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You have to sign in to update a project",
			"code":  http.StatusForbidden,
		})
		return
	}

	// get courseID param
	courseID := ctx.Param("course_id")

	// get projectID param
	projectID := ctx.Param("project_id")

	// bind json body with model
	var projectReq model.UpdateProject

	if err := ctx.ShouldBindJSON(&projectReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// update project
	project, err := c.projectService.UpdatedProjectByID(userClaims, courseID, projectID, projectReq)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
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

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("ProjectID %s updated successfully", projectID),
		"user":    projectResp,
	})
}

// delete project (admin only)
func (c *ProjectControllerImpl) DeleteProjectByID(ctx *gin.Context) {
	// make sure user has signed in
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You have to sign in to delete a project",
			"code":  http.StatusUnauthorized,
		})
		return
	}

	courseID := ctx.Param("course_id")

	projectID := ctx.Param("project_id")

	if err := c.projectService.DeleteProjectByID(userClaims, courseID, projectID); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// succeed response
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Project has been deleted",
		"code":    http.StatusOK,
	})
}
