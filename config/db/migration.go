package db

import (
	"github.com/nadyafa/go-learn/config/helper"
	"github.com/nadyafa/go-learn/entity"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Course{},
		&entity.Class{},
		&entity.Project{},
		&entity.Test{},
		&entity.Enrollment{},
		&entity.ProjectSub{},
		&entity.TestSub{},
		&entity.Attendance{},
	)

	if err != nil {
		helper.Logger(helper.LoggerLevelError, "Failed to migrate tables", err)
	}

	helper.Logger(helper.LoggerLevelInfo, "Database migrated successfully", nil)
}
