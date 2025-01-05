package db

import (
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	db.AutoMigrate(
		&entity.User{},
		// &entity.Course{},
		// &entity.Class{},
		// &entity.Project{},
		// &entity.Test{},
		// &entity.Enrollment{},
		// &entity.ProjectSub{},
		// &entity.TestSub{},
		// &entity.Attendance{},
	)
}
