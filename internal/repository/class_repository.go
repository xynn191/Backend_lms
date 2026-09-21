package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type ClassRepository interface {
	GetAll() ([]models.Class, error)
	FindByID(id uint) (*models.Class, error)
	Create(class *models.Class) error
	Update(class *models.Class) error
	Delete(id uint) error
}

type classRepository struct {
	db *gorm.DB
}

func NewClassRepository(db *gorm.DB) ClassRepository {
	return &classRepository{db}
}

func (r *classRepository) GetAll() ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Preload("Major").Preload("HomeroomTeacher").Find(&classes).Error
	return classes, err
}

func (r *classRepository) FindByID(id uint) (*models.Class, error) {
	var class models.Class
	err := r.db.Preload("Major").Preload("HomeroomTeacher").First(&class, id).Error
	if err != nil {
		return nil, err
	}
	return &class, nil
}

func (r *classRepository) Create(class *models.Class) error {
	return r.db.Create(class).Error
}

func (r *classRepository) Update(class *models.Class) error {
	return r.db.Save(class).Error
}

func (r *classRepository) Delete(id uint) error {
	return r.db.Delete(&models.Class{}, id).Error
}
