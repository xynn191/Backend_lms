package service

import (
	"errors"

	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type SubmitTaskRequest struct {
	FileURL string `json:"file_url"`
	LinkURL string `json:"link_url"`
}

type SubmitAssessmentAnswer struct {
	QuestionID       uint `json:"question_id" binding:"required"`
	SelectedOptionID uint `json:"selected_option_id" binding:"required"`
}

type SubmitAssessmentRequest struct {
	Answers []SubmitAssessmentAnswer `json:"answers" binding:"required"`
}

type StudentApiService interface {
	GetDashboardStats(userID string) (*repository.StudentDashboardStats, error)
	GetMaterials(userID string) ([]models.Material, error)
	GetMaterialByID(userID string, id uint) (*models.Material, error)
	GetTasks(userID string) ([]repository.StudentTaskItem, error)
	GetTaskByID(userID string, id uint) (*repository.StudentTaskItem, error)
	SubmitTask(userID string, taskID uint, req *SubmitTaskRequest) (*models.TaskSubmission, error)
	GetAssessments(userID string) ([]repository.StudentAssessmentItem, error)
	GetAssessmentByID(userID string, id uint) (*repository.StudentAssessmentItem, error)
	SubmitAssessment(userID string, assessmentID uint, req *SubmitAssessmentRequest) (*models.AssessmentSubmission, error)
}

type studentApiService struct {
	repo repository.StudentApiRepository
}

func NewStudentApiService(repo repository.StudentApiRepository) StudentApiService {
	return &studentApiService{repo}
}

func (s *studentApiService) getStudent(userID string) (*models.Student, error) {
	st, err := s.repo.GetStudentProfile(userID)
	if err != nil {
		return nil, errors.New("profil siswa tidak ditemukan")
	}
	return st, nil
}

func (s *studentApiService) GetDashboardStats(userID string) (*repository.StudentDashboardStats, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetDashboardStats(st.ClassID, st.ID, st)
}

func (s *studentApiService) GetMaterials(userID string) ([]models.Material, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMaterials(st.ClassID)
}

func (s *studentApiService) GetMaterialByID(userID string, id uint) (*models.Material, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMaterialByID(st.ClassID, id)
}

func (s *studentApiService) GetTasks(userID string) ([]repository.StudentTaskItem, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTasks(st.ClassID, st.ID)
}

func (s *studentApiService) GetTaskByID(userID string, id uint) (*repository.StudentTaskItem, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTaskByID(st.ClassID, st.ID, id)
}

func (s *studentApiService) SubmitTask(userID string, taskID uint, req *SubmitTaskRequest) (*models.TaskSubmission, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	if req.FileURL == "" && req.LinkURL == "" {
		return nil, errors.New("harap sertakan file atau tautan jawaban")
	}
	return s.repo.SubmitTask(st.ID, taskID, req.FileURL, req.LinkURL)
}

func (s *studentApiService) GetAssessments(userID string) ([]repository.StudentAssessmentItem, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAssessments(st.ClassID, st.ID)
}

func (s *studentApiService) GetAssessmentByID(userID string, id uint) (*repository.StudentAssessmentItem, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAssessmentByID(st.ClassID, st.ID, id)
}

func (s *studentApiService) SubmitAssessment(userID string, assessmentID uint, req *SubmitAssessmentRequest) (*models.AssessmentSubmission, error) {
	st, err := s.getStudent(userID)
	if err != nil {
		return nil, err
	}

	answerMap := make(map[uint]uint)
	for _, ans := range req.Answers {
		answerMap[ans.QuestionID] = ans.SelectedOptionID
	}

	return s.repo.SubmitAssessment(st.ID, assessmentID, answerMap)
}
