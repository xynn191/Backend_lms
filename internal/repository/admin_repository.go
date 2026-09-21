package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type DashboardStats struct {
	TotalStudents  int64
	TotalTeachers  int64
	TotalClasses   int64
	TotalSubjects  int64
}

type AdminRepository interface {
	GetDashboardStats() (*DashboardStats, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db}
}

func (r *adminRepository) GetDashboardStats() (*DashboardStats, error) {
	var stats DashboardStats
	r.db.Model(&models.Student{}).Count(&stats.TotalStudents)
	r.db.Model(&models.Teacher{}).Count(&stats.TotalTeachers)
	r.db.Model(&models.Class{}).Count(&stats.TotalClasses)
	r.db.Model(&models.Subject{}).Count(&stats.TotalSubjects)
	return &stats, nil
}
