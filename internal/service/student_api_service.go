package service

import (
	"errors"

	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type StudentApiService interface {
	GetDashboardStats(userID string) (*repository.StudentDashboardStats, error)
	GetMaterials(userID string) ([]models.Material, error)
	GetTasks(userID string) ([]models.Task, error)
}

type studentApiService struct {
	repo repository.StudentApiRepository
}

func NewStudentApiService(repo repository.StudentApiRepository) StudentApiService {
	return &studentApiService{repo}
}

func (s *studentApiService) getStudentClassID(userID string) (uint, error) {
	student, err := s.repo.GetStudentProfile(userID)
	if err != nil {
		return 0, errors.New("profil siswa tidak ditemukan")
	}
	return student.ClassID, nil
}

func (s *studentApiService) GetDashboardStats(userID string) (*repository.StudentDashboardStats, error) {
	classID, err := s.getStudentClassID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetDashboardStats(classID)
}

func (s *studentApiService) GetMaterials(userID string) ([]models.Material, error) {
	classID, err := s.getStudentClassID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMaterials(classID)
}

func (s *studentApiService) GetTasks(userID string) ([]models.Task, error) {
	classID, err := s.getStudentClassID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTasks(classID)
}
