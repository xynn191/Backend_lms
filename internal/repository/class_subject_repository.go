package repository

import (
	"errors"

	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type ClassSubjectRepository interface {
	GetAll() ([]models.ClassSubject, error)
	GetByID(id uint) (*models.ClassSubject, error)
	Create(cs *models.ClassSubject) error
	Update(cs *models.ClassSubject) error
	Delete(id uint) error
	FindByClassAndSubject(classID, subjectID uint) (*models.ClassSubject, error)
}

type classSubjectRepository struct {
	db *gorm.DB
}

func NewClassSubjectRepository(db *gorm.DB) ClassSubjectRepository {
	return &classSubjectRepository{db}
}

func (r *classSubjectRepository) GetAll() ([]models.ClassSubject, error) {
	var list []models.ClassSubject
	err := r.db.Preload("Class.Major").Preload("Subject").Preload("Teacher").Find(&list).Error
	return list, err
}

func (r *classSubjectRepository) GetByID(id uint) (*models.ClassSubject, error) {
	var cs models.ClassSubject
	if err := r.db.Preload("Class.Major").Preload("Subject").Preload("Teacher").First(&cs, id).Error; err != nil {
		return nil, err
	}
	return &cs, nil
}

func (r *classSubjectRepository) FindByClassAndSubject(classID, subjectID uint) (*models.ClassSubject, error) {
	var cs models.ClassSubject
	if err := r.db.Where("class_id = ? AND subject_id = ?", classID, subjectID).First(&cs).Error; err != nil {
		return nil, err
	}
	return &cs, nil
}

func (r *classSubjectRepository) Create(cs *models.ClassSubject) error {
	var existing models.ClassSubject
	if err := r.db.Where("class_id = ? AND subject_id = ?", cs.ClassID, cs.SubjectID).First(&existing).Error; err == nil {
		return errors.New("penugasan untuk kelas dan mata pelajaran ini sudah ada")
	}
	return r.db.Create(cs).Error
}

func (r *classSubjectRepository) Update(cs *models.ClassSubject) error {
	var existing models.ClassSubject
	if err := r.db.Where("class_id = ? AND subject_id = ? AND id != ?", cs.ClassID, cs.SubjectID, cs.ID).First(&existing).Error; err == nil {
		return errors.New("penugasan untuk kelas dan mata pelajaran ini sudah ada")
	}
	return r.db.Save(cs).Error
}

func (r *classSubjectRepository) Delete(id uint) error {
	return r.db.Delete(&models.ClassSubject{}, id).Error
}
