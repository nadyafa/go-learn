package service

import (
	"fmt"

	"github.com/nadyafa/go-learn/entity"
	"github.com/nadyafa/go-learn/middleware"
	"github.com/nadyafa/go-learn/model"
	"github.com/nadyafa/go-learn/repository"
)

type ProjectService interface {
	CreateProject(userClaims *middleware.UserClaims, courseID string, projectReq model.CreateProject) (*entity.Project, error)
	GetProjects(courseID string) ([]entity.Project, error)
	GetProjectByID(courseID, projectID string) (*entity.Project, error)
	UpdatedProjectByID(userClaims *middleware.UserClaims, courseID, projectID string, projectReq model.UpdateProject) (*entity.Project, error)
	DeleteProjectByID(userClaims *middleware.UserClaims, courseID, projectID string) error
}

type ProjectServiceImpl struct {
	projectRepo repository.ProjectRepo
	courseRepo  repository.CourseRepo
}

func NewProjectService(projectRepo repository.ProjectRepo, courseRepo repository.CourseRepo) ProjectService {
	return &ProjectServiceImpl{
		projectRepo: projectRepo,
		courseRepo:  courseRepo,
	}
}

func (s *ProjectServiceImpl) CreateProject(userClaims *middleware.UserClaims, courseID string, projectReq model.CreateProject) (*entity.Project, error) {
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can create a new project")
	}

	// validate if the course exist
	courseExist, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// validate if mentor own the class
	if userClaims.Role == entity.Mentor {
		if userClaims.UserID != courseExist.MentorID {
			return nil, fmt.Errorf("user has no authority to this action")
		}
	}

	// validate course name input
	isValid, errMsg := middleware.ValidateCourseName(projectReq.ProjectName)
	if !isValid {
		return nil, errMsg
	}

	// validate deadline project >= course endDate
	isValid, errMsg = middleware.ValidateCourseDate(projectReq.Deadline.Format("02-01-2006 15:04"), courseExist.EndDate.Format("02-01-2006 15:04"))
	if !isValid {
		return nil, errMsg
	}

	newProject := entity.Project{
		CourseID:    courseExist.CourseID,
		ProjectName: projectReq.ProjectName,
		Description: projectReq.Description,
		Deadline:    projectReq.Deadline.Time,
	}

	// input project to db
	if err := s.projectRepo.CreateProject(&newProject); err != nil {
		return nil, fmt.Errorf("unable to create a new project")
	}

	return &newProject, nil
}

func (s *ProjectServiceImpl) GetProjects(courseID string) ([]entity.Project, error) {
	// check courseID exist
	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// get projects from repo/db
	projects, err := s.projectRepo.GetProjects(courseID)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch projects")
	}

	return projects, nil
}

func (s *ProjectServiceImpl) GetProjectByID(courseID, projectID string) (*entity.Project, error) {
	// check courseID exist
	if _, err := s.courseRepo.GetCourseByID(courseID); err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// get project from repo/db
	project, err := s.projectRepo.GetProjectByID(courseID, projectID)
	if err != nil {
		return nil, fmt.Errorf("project_id %s not found", projectID)
	}

	return project, nil
}

func (s *ProjectServiceImpl) UpdatedProjectByID(userClaims *middleware.UserClaims, courseID, projectID string, projectReq model.UpdateProject) (*entity.Project, error) {
	if userClaims.Role == entity.Student {
		return nil, fmt.Errorf("only admin & mentor can update a project")
	}

	// check courseID exist
	courseExist, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return nil, fmt.Errorf("course_id %s not found", courseID)
	}

	// get project from repo/db
	projectExist, err := s.projectRepo.GetProjectByID(courseID, projectID)
	if err != nil {
		return nil, fmt.Errorf("project_id %s not found", projectID)
	}

	// validate if mentor own the class
	if userClaims.Role == entity.Mentor {
		if userClaims.UserID != courseExist.MentorID {
			return nil, fmt.Errorf("user has no authority to this action")
		}
	}

	// validate course name input
	if projectReq.ProjectName != "" {
		isValid, errMsg := middleware.ValidateCourseName(projectReq.ProjectName)
		if !isValid {
			return nil, errMsg
		}

		projectExist.ProjectName = projectReq.ProjectName
	}

	// validate deadline project >= course endDate
	if !projectReq.Deadline.IsZero() {
		isValid, errMsg := middleware.ValidateCourseDate(projectReq.Deadline.Format("02-01-2006 15:04"), courseExist.EndDate.Format("02-01-2006 15:04"))
		if !isValid {
			return nil, errMsg
		}

		projectExist.Deadline = projectReq.Deadline.Time
	}

	// if desc not empty
	if projectReq.Description != "" {
		projectExist.Description = projectReq.Description
	}

	// update project
	project, err := s.projectRepo.UpdateProjectByID(courseID, projectID, *projectExist)
	if err != nil {
		return nil, fmt.Errorf("unable to update project_id %s", projectID)
	}

	return project, nil
}

func (s *ProjectServiceImpl) DeleteProjectByID(userClaims *middleware.UserClaims, courseID, projectID string) error {
	if userClaims.Role == entity.Student {
		return fmt.Errorf("only admin & mentor can delete a project")
	}

	// check courseID exist
	courseExist, err := s.courseRepo.GetCourseByID(courseID)
	if err != nil {
		return fmt.Errorf("course_id %s not found", courseID)
	}

	// get project from repo/db
	if _, err := s.projectRepo.GetProjectByID(courseID, projectID); err != nil {
		return fmt.Errorf("project_id %s not found", projectID)
	}

	// validate if mentor own the class
	if userClaims.Role == entity.Mentor {
		if userClaims.UserID != courseExist.MentorID {
			return fmt.Errorf("user has no authority to this action")
		}
	}

	// update project
	if err := s.projectRepo.DeleteProjectByID(courseID, projectID); err != nil {
		return fmt.Errorf("unable to delete project_id %s", projectID)
	}

	return nil
}
