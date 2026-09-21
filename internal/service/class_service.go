package service

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type ClassDTO struct {
	MajorID           uint   `json:"major_id" binding:"required"`
	Name              string `json:"name" binding:"required"`
	GradeLevel        string `json:"grade_level" binding:"required,oneof=X XI XII"`
	HomeroomTeacherID uint   `json:"homeroom_teacher_id"`
}

type ClassService interface {
	GetAll() ([]models.Class, error)
	FindByID(id uint) (*models.Class, error)
	Create(dto *ClassDTO) (*models.Class, error)
	Update(id uint, dto *ClassDTO) (*models.Class, error)
	Delete(id uint) error
}

type classService struct {
	repo repository.ClassRepository
}

func NewClassService(repo repository.ClassRepository) ClassService {
	return &classService{repo}
}

func (s *classService) GetAll() ([]models.Class, error) {
	return s.repo.GetAll()
}

func (s *classService) FindByID(id uint) (*models.Class, error) {
	return s.repo.FindByID(id)
}

func (s *classService) Create(dto *ClassDTO) (*models.Class, error) {
	class := &models.Class{
		MajorID:           dto.MajorID,
		Name:              dto.Name,
		GradeLevel:        dto.GradeLevel,
		HomeroomTeacherID: dto.HomeroomTeacherID,
	}
	if err := s.repo.Create(class); err != nil {
		return nil, err
	}
	return class, nil
}

func (s *classService) Update(id uint, dto *ClassDTO) (*models.Class, error) {
	class, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	class.MajorID = dto.MajorID
	class.Name = dto.Name
	class.GradeLevel = dto.GradeLevel
	class.HomeroomTeacherID = dto.HomeroomTeacherID
	if err := s.repo.Update(class); err != nil {
		return nil, err
	}
	return class, nil
}

func (s *classService) Delete(id uint) error {
	return s.repo.Delete(id)
}
