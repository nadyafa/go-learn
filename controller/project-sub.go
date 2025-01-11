package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"gorm.io/gorm"
)

type ProjectSubController interface {
	StudentSubmitProject(ctx *gin.Context)
}

type ProjectSubControllerImpl struct {
	db *gorm.DB
}

func NewProjectSubController(db *gorm.DB) ProjectSubController {
	return &ProjectSubControllerImpl{
		db: db,
	}
}

// submit project (student only)
func (c *ProjectSubControllerImpl) StudentSubmitProject(ctx *gin.Context) {
	// check if the currentUser is admin
	claims, _ := ctx.Get("currentUser")
	userClaims, ok := claims.(*middleware.UserClaims)
	if !ok || userClaims.Role != entity.Student {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Access Restricted",
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

	// validate courseID
	projectID := ctx.Param("project_id")
	var project entity.Project
	if err := c.db.First(&project, projectID).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("ProjectID %s not found", projectID),
			"code":    http.StatusNotFound,
		})
		return
	}

	// validate current date with course startDate & endDate
	currentTime := time.Now()
	if currentTime.Before(course.StartDate) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Submissions are not allowed before the course starts",
		})
		return
	}

	if currentTime.After(course.EndDate) {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "Project submission period has ended",
		})
		return
	}

	// bind json body with model
	var projectSubReq model.ProjectSubStudent

	// if projectSubReq.Score == 0 {
	// 	projectSubReq.Score = 0
	// }

	if err := ctx.ShouldBind(&projectSubReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"code":  http.StatusBadRequest,
		})
		return
	}

	// retrieve file from req input
	file, err := ctx.FormFile("project_path")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to retrieve file",
		})
		return
	}

	// open file to check file size and validate MIME type
	openFile, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}
	defer openFile.Close()

	// check file size
	if err := middleware.CheckFileSize(openFile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	// reset file pointer for MIME type detection
	openFile.Seek(0, io.SeekStart)

	// validate MIME type
	if err := middleware.CheckMimeType(ctx, openFile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		ctx.Abort()
		return
	}

	// generate file name based on username & time upload
	newFileName := middleware.GenerateFileName(project.ProjectName, file.Filename)
	ctx.Set("newFileName", newFileName)

	// create directory to save file if it not exist
	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create directory",
		})
		return
	}

	// move file to directory uploads
	filePath := fmt.Sprintf("uploads/%s", file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to save file",
		})
		return
	}

	// submit projectSub to db
	projectSub := entity.ProjectSub{
		ProjectID:      project.ProjectID,
		StudentID:      userClaims.UserID,
		SubmissionDate: time.Now(),
		ProjectPath:    filePath,
		// Score:          projectSubReq.Score,
	}

	if err := c.db.Create(&projectSub).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create course",
			"code":  http.StatusInternalServerError,
		})
		return
	}

	// success response
	projectSubResp := model.StudentSubmitResp{
		ProjectSubID:   projectSub.ProjectSubID,
		ProjectID:      projectSub.ProjectID,
		StudentID:      projectSub.StudentID,
		SubmissionDate: projectSub.SubmissionDate,
		ProjectPath:    projectSub.ProjectPath,
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Course %s created successfully", course.CourseName),
		"user":    projectSubResp,
	})
}
