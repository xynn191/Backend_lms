package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type StudentDashboardStats struct {
	TotalSubjects   int64 `json:"total_subjects"`
	TotalMaterials  int64 `json:"total_materials"`
	TotalTasks      int64 `json:"total_tasks"`
	PendingTasks    int64 `json:"pending_tasks"`
}

type StudentApiRepository interface {
	GetStudentProfile(userID string) (*models.Student, error)
	GetMaterials(classID uint) ([]models.Material, error)
	GetTasks(classID uint) ([]models.Task, error)
	GetDashboardStats(classID uint) (*StudentDashboardStats, error)
}

type studentApiRepository struct {
	db *gorm.DB
}

func NewStudentApiRepository(db *gorm.DB) StudentApiRepository {
	return &studentApiRepository{db}
}

func (r *studentApiRepository) GetStudentProfile(userID string) (*models.Student, error) {
	var student models.Student
	if err := r.db.Preload("Class").Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentApiRepository) GetMaterials(classID uint) ([]models.Material, error) {
	var materials []models.Material
	// Join with class_subjects to filter by classID and preload Subject for badge
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Preload("ClassSubject.Subject").
		Where("class_subjects.class_id = ?", classID).
		Order("materials.created_at DESC").
		Find(&materials).Error
	return materials, err
}

func (r *studentApiRepository) GetTasks(classID uint) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Preload("ClassSubject.Subject").
		Where("class_subjects.class_id = ?", classID).
		Order("tasks.due_date ASC").
		Find(&tasks).Error
	return tasks, err
}

func (r *studentApiRepository) GetDashboardStats(classID uint) (*StudentDashboardStats, error) {
	var stats StudentDashboardStats
	
	r.db.Model(&models.ClassSubject{}).Where("class_id = ?", classID).Count(&stats.TotalSubjects)
	
	r.db.Model(&models.Material{}).Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Where("class_subjects.class_id = ?", classID).Count(&stats.TotalMaterials)
		
	r.db.Model(&models.Task{}).Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Where("class_subjects.class_id = ?", classID).Count(&stats.TotalTasks)
		
	// Assuming pending is same as total for simplicity right now
	stats.PendingTasks = stats.TotalTasks 

	return &stats, nil
}
