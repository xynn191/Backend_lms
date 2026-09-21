package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type MajorRepository interface {
	GetAll() ([]models.Major, error)
	FindByID(id uint) (*models.Major, error)
	Create(major *models.Major) error
	Update(major *models.Major) error
	Delete(id uint) error
}

type majorRepository struct {
	db *gorm.DB
}

func NewMajorRepository(db *gorm.DB) MajorRepository {
	return &majorRepository{db}
}

func (r *majorRepository) GetAll() ([]models.Major, error) {
	var majors []models.Major
	err := r.db.Find(&majors).Error
	return majors, err
}

func (r *majorRepository) FindByID(id uint) (*models.Major, error) {
	var major models.Major
	err := r.db.First(&major, id).Error
	if err != nil {
		return nil, err
	}
	return &major, nil
}

func (r *majorRepository) Create(major *models.Major) error {
	return r.db.Create(major).Error
}

func (r *majorRepository) Update(major *models.Major) error {
	return r.db.Save(major).Error
}

func (r *majorRepository) Delete(id uint) error {
	return r.db.Delete(&models.Major{}, id).Error
}
