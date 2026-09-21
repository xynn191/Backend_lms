package service

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type SubjectDTO struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type SubjectService interface {
	GetAll() ([]models.Subject, error)
	FindByID(id uint) (*models.Subject, error)
	Create(dto *SubjectDTO) (*models.Subject, error)
	Update(id uint, dto *SubjectDTO) (*models.Subject, error)
	Delete(id uint) error
}

type subjectService struct {
	repo repository.SubjectRepository
}

func NewSubjectService(repo repository.SubjectRepository) SubjectService {
	return &subjectService{repo}
}

func (s *subjectService) GetAll() ([]models.Subject, error) {
	return s.repo.GetAll()
}

func (s *subjectService) FindByID(id uint) (*models.Subject, error) {
	return s.repo.FindByID(id)
}

func (s *subjectService) Create(dto *SubjectDTO) (*models.Subject, error) {
	subject := &models.Subject{
		Code: dto.Code,
		Name: dto.Name,
	}
	if err := s.repo.Create(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *subjectService) Update(id uint, dto *SubjectDTO) (*models.Subject, error) {
	subject, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	subject.Code = dto.Code
	subject.Name = dto.Name
	if err := s.repo.Update(subject); err != nil {
		return nil, err
	}
	return subject, nil
}

func (s *subjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}
