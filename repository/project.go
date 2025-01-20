package repository

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

type ProjectRepo interface {
	CreateProject(project *entity.Project) error
	GetProjects(courseID string) ([]entity.Project, error)
	GetProjectByID(courseID, projectID string) (*entity.Project, error)
	UpdateProjectByID(courseID, projectID string, project entity.Project) (*entity.Project, error)
	DeleteProjectByID(courseID, projectID string) error
}

type ProjectRepoImpl struct {
	db *gorm.DB
}

func NewProjectRepo(db *gorm.DB) ProjectRepo {
	return &ProjectRepoImpl{
		db: db,
	}
}

func (r *ProjectRepoImpl) CreateProject(project *entity.Project) error {
	if err := r.db.Create(&project).Error; err != nil {
		return err
	}

	return nil
}

func (r *ProjectRepoImpl) GetProjects(courseID string) ([]entity.Project, error) {
	var projects []entity.Project

	if err := r.db.Where("course_id = ?", courseID).Find(&projects).Error; err != nil {
		return nil, err
	}

	return projects, nil
}

func (r *ProjectRepoImpl) GetProjectByID(courseID, projectID string) (*entity.Project, error) {
	var project entity.Project

	if err := r.db.Where("course_id = ? AND project_id = ?", courseID, projectID).First(&project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepoImpl) UpdateProjectByID(courseID, projectID string, project entity.Project) (*entity.Project, error) {
	if err := r.db.Where("course_id = ? AND project_id = ?", courseID, projectID).Updates(project).Error; err != nil {
		return nil, err
	}

	return &project, nil
}

func (r *ProjectRepoImpl) DeleteProjectByID(courseID, projectID string) error {
	var project entity.Project

	if err := r.db.Where("course_id = ? AND project_id = ?", courseID, projectID).Delete(project).Error; err != nil {
		return err
	}

	return nil
}
