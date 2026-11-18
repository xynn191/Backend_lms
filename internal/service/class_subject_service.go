package service

import (
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type CreateClassSubjectRequest struct {
	ClassID   uint `json:"class_id" binding:"required"`
	SubjectID uint `json:"subject_id" binding:"required"`
	TeacherID uint `json:"teacher_id" binding:"required"`
}

type UpdateClassSubjectRequest struct {
	ClassID   uint `json:"class_id" binding:"required"`
	SubjectID uint `json:"subject_id" binding:"required"`
	TeacherID uint `json:"teacher_id" binding:"required"`
}

type ClassSubjectService interface {
	GetAll() ([]models.ClassSubject, error)
	GetByID(id uint) (*models.ClassSubject, error)
	Create(req *CreateClassSubjectRequest) (*models.ClassSubject, error)
	Update(id uint, req *UpdateClassSubjectRequest) (*models.ClassSubject, error)
	Delete(id uint) error
}

type classSubjectService struct {
	repo repository.ClassSubjectRepository
}

func NewClassSubjectService(repo repository.ClassSubjectRepository) ClassSubjectService {
	return &classSubjectService{repo}
}

func (s *classSubjectService) GetAll() ([]models.ClassSubject, error) {
	return s.repo.GetAll()
}

func (s *classSubjectService) GetByID(id uint) (*models.ClassSubject, error) {
	return s.repo.GetByID(id)
}

func (s *classSubjectService) Create(req *CreateClassSubjectRequest) (*models.ClassSubject, error) {
	cs := &models.ClassSubject{
		ClassID:   req.ClassID,
		SubjectID: req.SubjectID,
		TeacherID: req.TeacherID,
	}
	if err := s.repo.Create(cs); err != nil {
		return nil, err
	}
	return s.repo.GetByID(cs.ID)
}

func (s *classSubjectService) Update(id uint, req *UpdateClassSubjectRequest) (*models.ClassSubject, error) {
	cs, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	cs.ClassID = req.ClassID
	cs.SubjectID = req.SubjectID
	cs.TeacherID = req.TeacherID

	if err := s.repo.Update(cs); err != nil {
		return nil, err
	}
	return s.repo.GetByID(cs.ID)
}

func (s *classSubjectService) Delete(id uint) error {
	return s.repo.Delete(id)
}
