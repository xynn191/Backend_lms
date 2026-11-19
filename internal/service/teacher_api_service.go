package service

import (
	"errors"
	"time"

	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

type CreateMaterialRequest struct {
	ClassSubjectID uint   `json:"class_subject_id" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Content        string `json:"content"`
	AttachmentURL  string `json:"attachment_url"`
	IsPublished    bool   `json:"is_published"`
}

type UpdateMaterialRequest struct {
	Title         string `json:"title" binding:"required"`
	Content       string `json:"content"`
	AttachmentURL string `json:"attachment_url"`
	IsPublished   bool   `json:"is_published"`
}

type CreateTaskRequest struct {
	ClassSubjectID uint      `json:"class_subject_id" binding:"required"`
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description"`
	AttachmentURL  string    `json:"attachment_url"`
	AttachmentLink string    `json:"attachment_link"`
	DueDate        time.Time `json:"due_date" binding:"required"`
	IsPublished    bool      `json:"is_published"`
}

type UpdateTaskRequest struct {
	Title          string    `json:"title" binding:"required"`
	Description    string    `json:"description"`
	AttachmentURL  string    `json:"attachment_url"`
	AttachmentLink string    `json:"attachment_link"`
	DueDate        time.Time `json:"due_date" binding:"required"`
	IsPublished    bool      `json:"is_published"`
}

type GradeSubmissionRequest struct {
	Score    float64 `json:"score" binding:"required"`
	Feedback string  `json:"feedback"`
}

type AssessmentOptionRequest struct {
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

type AssessmentQuestionRequest struct {
	QuestionText string                    `json:"question_text" binding:"required"`
	Options      []AssessmentOptionRequest `json:"options" binding:"required,min=2"`
}

type CreateAssessmentRequest struct {
	ClassSubjectID  uint                         `json:"class_subject_id" binding:"required"`
	Title           string                       `json:"title" binding:"required"`
	Description     string                       `json:"description"`
	StartTime       time.Time                    `json:"start_time" binding:"required"`
	EndTime         time.Time                    `json:"end_time" binding:"required"`
	DurationMinutes int                          `json:"duration_minutes" binding:"required"`
	IsPublished     bool                         `json:"is_published"`
	Questions       []AssessmentQuestionRequest `json:"questions" binding:"required,min=1"`
}

type UpdateAssessmentRequest struct {
	Title           string                       `json:"title" binding:"required"`
	Description     string                       `json:"description"`
	StartTime       time.Time                    `json:"start_time" binding:"required"`
	EndTime         time.Time                    `json:"end_time" binding:"required"`
	DurationMinutes int                          `json:"duration_minutes" binding:"required"`
	IsPublished     bool                         `json:"is_published"`
	Questions       []AssessmentQuestionRequest `json:"questions" binding:"required,min=1"`
}

type UpdateAssessmentScheduleRequest struct {
	StartTime       time.Time `json:"start_time" binding:"required"`
	EndTime         time.Time `json:"end_time" binding:"required"`
	DurationMinutes *int      `json:"duration_minutes"`
}

type PublishRequest struct {
	IsPublished bool `json:"is_published"`
}

type TeacherApiService interface {
	GetDashboardStats(userID string) (*repository.TeacherDashboardStats, error)
	GetClassSubjects(userID string) ([]models.ClassSubject, error)
	GetClassSubjectByID(userID string, id uint) (*models.ClassSubject, error)

	// Material
	GetMaterials(userID string, classSubjectID uint) ([]models.Material, error)
	GetMaterialByID(userID string, id uint) (*models.Material, error)
	CreateMaterial(userID string, req *CreateMaterialRequest) (*models.Material, error)
	UpdateMaterial(userID string, id uint, req *UpdateMaterialRequest) (*models.Material, error)
	DeleteMaterial(userID string, id uint) error
	PublishMaterial(userID string, id uint, publish bool) error

	// Task
	GetTasks(userID string, classSubjectID uint) ([]models.Task, error)
	GetTaskByID(userID string, id uint) (*models.Task, error)
	CreateTask(userID string, req *CreateTaskRequest) (*models.Task, error)
	UpdateTask(userID string, id uint, req *UpdateTaskRequest) (*models.Task, error)
	DeleteTask(userID string, id uint) error
	PublishTask(userID string, id uint, publish bool) error
	GetTaskSubmissions(userID string, taskID uint) ([]models.TaskSubmission, error)
	GradeSubmission(userID string, submissionID uint, req *GradeSubmissionRequest) (*models.TaskSubmission, error)

	// Assessment
	GetAssessments(userID string, classSubjectID uint) ([]models.Assessment, error)
	GetAssessmentByID(userID string, id uint) (*models.Assessment, error)
	CreateAssessment(userID string, req *CreateAssessmentRequest) (*models.Assessment, error)
	UpdateAssessment(userID string, id uint, req *UpdateAssessmentRequest) (*models.Assessment, error)
	UpdateAssessmentSchedule(userID string, id uint, req *UpdateAssessmentScheduleRequest) (*models.Assessment, error)
	DeleteAssessment(userID string, id uint) error
	PublishAssessment(userID string, id uint, publish bool) error
	GetAssessmentResults(userID string, assessmentID uint) ([]models.AssessmentSubmission, error)

	// Grades
	GetGradesOverview(userID string, classSubjectID uint) ([]repository.TeacherGradeItem, error)
	GetGradesDetail(userID string, classSubjectID uint) (*repository.TeacherGradesDetail, error)
}

type teacherApiService struct {
	repo repository.TeacherApiRepository
}

func NewTeacherApiService(repo repository.TeacherApiRepository) TeacherApiService {
	return &teacherApiService{repo}
}

func (s *teacherApiService) getTeacherID(userID string) (uint, error) {
	teacher, err := s.repo.GetTeacherByUserID(userID)
	if err != nil {
		return 0, errors.New("profil guru tidak ditemukan")
	}
	return teacher.ID, nil
}

func (s *teacherApiService) GetDashboardStats(userID string) (*repository.TeacherDashboardStats, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetDashboardStats(teacherID)
}

func (s *teacherApiService) GetClassSubjects(userID string) ([]models.ClassSubject, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetClassSubjects(teacherID)
}

func (s *teacherApiService) GetClassSubjectByID(userID string, id uint) (*models.ClassSubject, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetClassSubjectByID(teacherID, id)
}

func (s *teacherApiService) GetMaterials(userID string, classSubjectID uint) ([]models.Material, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMaterials(teacherID, classSubjectID)
}

func (s *teacherApiService) GetMaterialByID(userID string, id uint) (*models.Material, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetMaterialByID(teacherID, id)
}

func (s *teacherApiService) CreateMaterial(userID string, req *CreateMaterialRequest) (*models.Material, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	// Verify class subject belongs to teacher
	_, err = s.repo.GetClassSubjectByID(teacherID, req.ClassSubjectID)
	if err != nil {
		return nil, errors.New("kelas-mata pelajaran ini tidak ditugaskan kepada Anda")
	}

	m := &models.Material{
		ClassSubjectID: req.ClassSubjectID,
		Title:          req.Title,
		Content:        req.Content,
		AttachmentURL:  req.AttachmentURL,
		IsPublished:    req.IsPublished,
	}
	if err := s.repo.CreateMaterial(m); err != nil {
		return nil, err
	}
	return s.repo.GetMaterialByID(teacherID, m.ID)
}

func (s *teacherApiService) UpdateMaterial(userID string, id uint, req *UpdateMaterialRequest) (*models.Material, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	m, err := s.repo.GetMaterialByID(teacherID, id)
	if err != nil {
		return nil, err
	}

	m.Title = req.Title
	m.Content = req.Content
	m.AttachmentURL = req.AttachmentURL
	m.IsPublished = req.IsPublished

	if err := s.repo.UpdateMaterial(m); err != nil {
		return nil, err
	}
	return s.repo.GetMaterialByID(teacherID, id)
}

func (s *teacherApiService) DeleteMaterial(userID string, id uint) error {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteMaterial(teacherID, id)
}

func (s *teacherApiService) PublishMaterial(userID string, id uint, publish bool) error {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return err
	}
	return s.repo.PublishMaterial(teacherID, id, publish)
}

func (s *teacherApiService) GetTasks(userID string, classSubjectID uint) ([]models.Task, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTasks(teacherID, classSubjectID)
}

func (s *teacherApiService) GetTaskByID(userID string, id uint) (*models.Task, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTaskByID(teacherID, id)
}

func (s *teacherApiService) CreateTask(userID string, req *CreateTaskRequest) (*models.Task, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.GetClassSubjectByID(teacherID, req.ClassSubjectID)
	if err != nil {
		return nil, errors.New("kelas-mata pelajaran ini tidak ditugaskan kepada Anda")
	}

	t := &models.Task{
		ClassSubjectID: req.ClassSubjectID,
		Title:          req.Title,
		Description:    req.Description,
		AttachmentURL:  req.AttachmentURL,
		AttachmentLink: req.AttachmentLink,
		DueDate:        req.DueDate,
		IsPublished:    req.IsPublished,
	}
	if err := s.repo.CreateTask(t); err != nil {
		return nil, err
	}
	return s.repo.GetTaskByID(teacherID, t.ID)
}

func (s *teacherApiService) UpdateTask(userID string, id uint, req *UpdateTaskRequest) (*models.Task, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	t, err := s.repo.GetTaskByID(teacherID, id)
	if err != nil {
		return nil, err
	}

	t.Title = req.Title
	t.Description = req.Description
	t.AttachmentURL = req.AttachmentURL
	t.AttachmentLink = req.AttachmentLink
	t.DueDate = req.DueDate
	t.IsPublished = req.IsPublished

	if err := s.repo.UpdateTask(t); err != nil {
		return nil, err
	}
	return s.repo.GetTaskByID(teacherID, id)
}

func (s *teacherApiService) DeleteTask(userID string, id uint) error {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteTask(teacherID, id)
}

func (s *teacherApiService) PublishTask(userID string, id uint, publish bool) error {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return err
	}
	return s.repo.PublishTask(teacherID, id, publish)
}

func (s *teacherApiService) GetTaskSubmissions(userID string, taskID uint) ([]models.TaskSubmission, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTaskSubmissions(teacherID, taskID)
}

func (s *teacherApiService) GradeSubmission(userID string, submissionID uint, req *GradeSubmissionRequest) (*models.TaskSubmission, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	if req.Score < 0 || req.Score > 100 {
		return nil, errors.New("nilai harus berada pada rentang 0 hingga 100")
	}
	return s.repo.GradeSubmission(teacherID, submissionID, req.Score, req.Feedback)
}

func (s *teacherApiService) GetAssessments(userID string, classSubjectID uint) ([]models.Assessment, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAssessments(teacherID, classSubjectID)
}

func (s *teacherApiService) GetAssessmentByID(userID string, id uint) (*models.Assessment, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAssessmentByID(teacherID, id)
}

func (s *teacherApiService) CreateAssessment(userID string, req *CreateAssessmentRequest) (*models.Assessment, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.GetClassSubjectByID(teacherID, req.ClassSubjectID)
	if err != nil {
		return nil, errors.New("kelas-mata pelajaran ini tidak ditugaskan kepada Anda")
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("waktu selesai harus setelah waktu mulai")
	}

	questions := make([]models.AssessmentQuestion, len(req.Questions))
	for i, q := range req.Questions {
		opts := make([]models.AssessmentOption, len(q.Options))
		hasCorrect := false
		for j, o := range q.Options {
			if o.IsCorrect {
				hasCorrect = true
			}
			opts[j] = models.AssessmentOption{
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
			}
		}
		if !hasCorrect {
			return nil, errors.New("setiap soal harus memiliki setidaknya satu opsi jawaban yang benar")
		}
		questions[i] = models.AssessmentQuestion{
			QuestionText: q.QuestionText,
			Options:      opts,
		}
	}

	a := &models.Assessment{
		ClassSubjectID:  req.ClassSubjectID,
		Title:           req.Title,
		Description:     req.Description,
		StartTime:       req.StartTime,
		EndTime:         req.EndTime,
		DurationMinutes: req.DurationMinutes,
		IsPublished:     req.IsPublished,
		Questions:       questions,
	}

	if err := s.repo.CreateAssessment(a); err != nil {
		return nil, err
	}
	return s.repo.GetAssessmentByID(teacherID, a.ID)
}

func (s *teacherApiService) UpdateAssessment(userID string, id uint, req *UpdateAssessmentRequest) (*models.Assessment, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	a, err := s.repo.GetAssessmentByID(teacherID, id)
	if err != nil {
		return nil, err
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("waktu selesai harus setelah waktu mulai")
	}

	questions := make([]models.AssessmentQuestion, len(req.Questions))
	for i, q := range req.Questions {
		opts := make([]models.AssessmentOption, len(q.Options))
		hasCorrect := false
		for j, o := range q.Options {
			if o.IsCorrect {
				hasCorrect = true
			}
			opts[j] = models.AssessmentOption{
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
			}
		}
		if !hasCorrect {
			return nil, errors.New("setiap soal harus memiliki setidaknya satu opsi jawaban yang benar")
		}
		questions[i] = models.AssessmentQuestion{
			AssessmentID: a.ID,
			QuestionText: q.QuestionText,
			Options:      opts,
		}
	}

	a.Title = req.Title
	a.Description = req.Description
	a.StartTime = req.StartTime
	a.EndTime = req.EndTime
	a.DurationMinutes = req.DurationMinutes
	a.IsPublished = req.IsPublished
	a.Questions = questions

	if err := s.repo.UpdateAssessment(a); err != nil {
		return nil, err
	}
	return s.repo.GetAssessmentByID(teacherID, id)
}

func (s *teacherApiService) UpdateAssessmentSchedule(userID string, id uint, req *UpdateAssessmentScheduleRequest) (*models.Assessment, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.GetAssessmentByID(teacherID, id)
	if err != nil {
		return nil, err
	}

	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("waktu selesai harus setelah waktu mulai")
	}

	if err := s.repo.UpdateAssessmentSchedule(id, req.StartTime, req.EndTime, req.DurationMinutes); err != nil {
		return nil, err
	}
	return s.repo.GetAssessmentByID(teacherID, id)
}

func (s *teacherApiService) DeleteAssessment(userID string, id uint) error {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return err
	}
	return s.repo.DeleteAssessment(teacherID, id)
}

func (s *teacherApiService) PublishAssessment(userID string, id uint, publish bool) error {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return err
	}
	return s.repo.PublishAssessment(teacherID, id, publish)
}

func (s *teacherApiService) GetAssessmentResults(userID string, assessmentID uint) ([]models.AssessmentSubmission, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAssessmentResults(teacherID, assessmentID)
}

func (s *teacherApiService) GetGradesOverview(userID string, classSubjectID uint) ([]repository.TeacherGradeItem, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetGradesOverview(teacherID, classSubjectID)
}

func (s *teacherApiService) GetGradesDetail(userID string, classSubjectID uint) (*repository.TeacherGradesDetail, error) {
	teacherID, err := s.getTeacherID(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetGradesDetail(teacherID, classSubjectID)
}
