package service

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type MajorDTO struct {
	Code string `json:"code" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type MajorService interface {
	GetAll() ([]models.Major, error)
	FindByID(id uint) (*models.Major, error)
	Create(dto *MajorDTO) (*models.Major, error)
	Update(id uint, dto *MajorDTO) (*models.Major, error)
	Delete(id uint) error
}

type majorService struct {
	repo repository.MajorRepository
}

func NewMajorService(repo repository.MajorRepository) MajorService {
	return &majorService{repo}
}

func (s *majorService) GetAll() ([]models.Major, error) {
	return s.repo.GetAll()
}

func (s *majorService) FindByID(id uint) (*models.Major, error) {
	return s.repo.FindByID(id)
}

func (s *majorService) Create(dto *MajorDTO) (*models.Major, error) {
	major := &models.Major{
		Code: dto.Code,
		Name: dto.Name,
	}
	if err := s.repo.Create(major); err != nil {
		return nil, err
	}
	return major, nil
}

func (s *majorService) Update(id uint, dto *MajorDTO) (*models.Major, error) {
	major, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	major.Code = dto.Code
	major.Name = dto.Name
	if err := s.repo.Update(major); err != nil {
		return nil, err
	}
	return major, nil
}

func (s *majorService) Delete(id uint) error {
	return s.repo.Delete(id)
}
