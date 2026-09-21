package models

import (
	"gorm.io/gorm"
)

// AutoMigrate all models
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Role{},
		&User{},
		&Admin{},
		&Teacher{},
		&Major{},
		&Class{}, // Dependent on Major and Teacher
		&Student{}, // Dependent on Class
		&Curriculum{},
		&Principal{},
		&Subject{},
		&ClassSubject{},
		&Material{},
		&Task{},
		&TaskSubmission{},
		&Assessment{},
		&AssessmentSubmission{},
	)
}
