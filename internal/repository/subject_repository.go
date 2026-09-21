package repository

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type SubjectRepository interface {
	GetAll() ([]models.Subject, error)
	FindByID(id uint) (*models.Subject, error)
	Create(subject *models.Subject) error
	Update(subject *models.Subject) error
	Delete(id uint) error
}

type subjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) SubjectRepository {
	return &subjectRepository{db}
}

func (r *subjectRepository) GetAll() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Find(&subjects).Error
	return subjects, err
}

func (r *subjectRepository) FindByID(id uint) (*models.Subject, error) {
	var subject models.Subject
	err := r.db.First(&subject, id).Error
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *subjectRepository) Create(subject *models.Subject) error {
	return r.db.Create(subject).Error
}

func (r *subjectRepository) Update(subject *models.Subject) error {
	return r.db.Save(subject).Error
}

func (r *subjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.Subject{}, id).Error
}
