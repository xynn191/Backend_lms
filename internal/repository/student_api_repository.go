package repository

import (
	"errors"
	"time"

	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type StudentDashboardStats struct {
	StudentName      string `json:"student_name"`
	StudentNISN      string `json:"student_nisn"`
	ClassName        string `json:"class_name"`
	TotalSubjects    int64  `json:"total_subjects"`
	TotalMaterials   int64  `json:"total_materials"`
	TotalTasks       int64  `json:"total_tasks"`
	PendingTasks     int64  `json:"pending_tasks"`
	TotalAssessments int64  `json:"total_assessments"`
}

type StudentTaskItem struct {
	models.Task
	Submission *models.TaskSubmission `json:"submission"`
}

type StudentAssessmentItem struct {
	models.Assessment
	Submission *models.AssessmentSubmission `json:"submission"`
}

type StudentApiRepository interface {
	GetStudentProfile(userID string) (*models.Student, error)
	GetMaterials(classID uint) ([]models.Material, error)
	GetMaterialByID(classID uint, id uint) (*models.Material, error)
	GetTasks(classID uint, studentID uint) ([]StudentTaskItem, error)
	GetTaskByID(classID uint, studentID uint, id uint) (*StudentTaskItem, error)
	SubmitTask(studentID uint, taskID uint, fileURL string, linkURL string) (*models.TaskSubmission, error)
	GetAssessments(classID uint, studentID uint) ([]StudentAssessmentItem, error)
	GetAssessmentByID(classID uint, studentID uint, id uint) (*StudentAssessmentItem, error)
	SubmitAssessment(studentID uint, assessmentID uint, answers map[uint]uint) (*models.AssessmentSubmission, error)
	GetDashboardStats(classID uint, studentID uint, student *models.Student) (*StudentDashboardStats, error)
}

type studentApiRepository struct {
	db *gorm.DB
}

func NewStudentApiRepository(db *gorm.DB) StudentApiRepository {
	return &studentApiRepository{db}
}

func (r *studentApiRepository) GetStudentProfile(userID string) (*models.Student, error) {
	var student models.Student
	if err := r.db.Preload("Class.Major").Where("user_id = ?", userID).First(&student).Error; err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentApiRepository) GetMaterials(classID uint) ([]models.Material, error) {
	var materials []models.Material
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Preload("ClassSubject.Subject").Preload("ClassSubject.Teacher").
		Where("class_subjects.class_id = ? AND materials.is_published = ?", classID, true).
		Order("materials.created_at DESC").
		Find(&materials).Error
	return materials, err
}

func (r *studentApiRepository) GetMaterialByID(classID uint, id uint) (*models.Material, error) {
	var m models.Material
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Preload("ClassSubject.Subject").Preload("ClassSubject.Teacher").
		Where("materials.id = ? AND class_subjects.class_id = ? AND materials.is_published = ?", id, classID, true).
		First(&m).Error
	if err != nil {
		return nil, errors.New("materi tidak ditemukan atau belum diterbitkan")
	}
	return &m, nil
}

func (r *studentApiRepository) GetTasks(classID uint, studentID uint) ([]StudentTaskItem, error) {
	var tasks []models.Task
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Preload("ClassSubject.Subject").Preload("ClassSubject.Teacher").
		Where("class_subjects.class_id = ? AND tasks.is_published = ?", classID, true).
		Order("tasks.due_date ASC").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	result := make([]StudentTaskItem, len(tasks))
	for i, t := range tasks {
		var sub models.TaskSubmission
		err := r.db.Where("task_id = ? AND student_id = ?", t.ID, studentID).First(&sub).Error
		item := StudentTaskItem{Task: t}
		if err == nil {
			item.Submission = &sub
		}
		result[i] = item
	}
	return result, nil
}

func (r *studentApiRepository) GetTaskByID(classID uint, studentID uint, id uint) (*StudentTaskItem, error) {
	var t models.Task
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Preload("ClassSubject.Subject").Preload("ClassSubject.Teacher").
		Where("tasks.id = ? AND class_subjects.class_id = ? AND tasks.is_published = ?", id, classID, true).
		First(&t).Error
	if err != nil {
		return nil, errors.New("tugas tidak ditemukan atau belum diterbitkan")
	}

	var sub models.TaskSubmission
	err = r.db.Where("task_id = ? AND student_id = ?", t.ID, studentID).First(&sub).Error
	item := StudentTaskItem{Task: t}
	if err == nil {
		item.Submission = &sub
	}
	return &item, nil
}

func (r *studentApiRepository) SubmitTask(studentID uint, taskID uint, fileURL string, linkURL string) (*models.TaskSubmission, error) {
	var task models.Task
	if err := r.db.Where("id = ? AND is_published = ?", taskID, true).First(&task).Error; err != nil {
		return nil, errors.New("tugas tidak ditemukan atau belum diterbitkan")
	}

	now := time.Now()
	var sub models.TaskSubmission
	err := r.db.Where("task_id = ? AND student_id = ?", taskID, studentID).First(&sub).Error
	if err == nil {
		// Update existing submission
		sub.FileURL = fileURL
		sub.LinkURL = linkURL
		sub.SubmittedAt = now
		if err := r.db.Save(&sub).Error; err != nil {
			return nil, err
		}
		return &sub, nil
	}

	sub = models.TaskSubmission{
		TaskID:      taskID,
		StudentID:   studentID,
		FileURL:     fileURL,
		LinkURL:     linkURL,
		SubmittedAt: now,
	}
	if err := r.db.Create(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *studentApiRepository) GetAssessments(classID uint, studentID uint) ([]StudentAssessmentItem, error) {
	var list []models.Assessment
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = assessments.class_subject_id").
		Preload("ClassSubject.Subject").Preload("ClassSubject.Teacher").
		Where("class_subjects.class_id = ? AND assessments.is_published = ?", classID, true).
		Order("assessments.start_time ASC").
		Find(&list).Error
	if err != nil {
		return nil, err
	}

	result := make([]StudentAssessmentItem, len(list))
	for i, a := range list {
		var sub models.AssessmentSubmission
		err := r.db.Where("assessment_id = ? AND student_id = ?", a.ID, studentID).First(&sub).Error
		item := StudentAssessmentItem{Assessment: a}
		if err == nil {
			item.Submission = &sub
		}
		result[i] = item
	}
	return result, nil
}

func (r *studentApiRepository) GetAssessmentByID(classID uint, studentID uint, id uint) (*StudentAssessmentItem, error) {
	var a models.Assessment
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = assessments.class_subject_id").
		Preload("ClassSubject.Subject").Preload("ClassSubject.Teacher").
		Preload("Questions.Options", func(db *gorm.DB) *gorm.DB {
			// Don't leak IsCorrect flag to students in option output query
			return db.Select("id, question_id, option_text, created_at, updated_at")
		}).
		Where("assessments.id = ? AND class_subjects.class_id = ? AND assessments.is_published = ?", id, classID, true).
		First(&a).Error
	if err != nil {
		return nil, errors.New("assessment tidak ditemukan atau belum diterbitkan")
	}

	var sub models.AssessmentSubmission
	err = r.db.Where("assessment_id = ? AND student_id = ?", a.ID, studentID).First(&sub).Error
	item := StudentAssessmentItem{Assessment: a}
	if err == nil {
		item.Submission = &sub
	}
	return &item, nil
}

func (r *studentApiRepository) SubmitAssessment(studentID uint, assessmentID uint, answers map[uint]uint) (*models.AssessmentSubmission, error) {
	var a models.Assessment
	err := r.db.Preload("Questions.Options").
		Where("id = ? AND is_published = ?", assessmentID, true).
		First(&a).Error
	if err != nil {
		return nil, errors.New("assessment tidak ditemukan atau belum diterbitkan")
	}

	now := time.Now()
	if now.Before(a.StartTime) || now.After(a.EndTime) {
		return nil, errors.New("assessment tidak dalam jadwal aktif")
	}

	// Check existing submission
	var existing models.AssessmentSubmission
	if err := r.db.Where("assessment_id = ? AND student_id = ?", assessmentID, studentID).First(&existing).Error; err == nil {
		return nil, errors.New("Anda sudah mengerjakan assessment ini")
	}

	// Calculate score automatically
	var correctCount int
	totalQuestions := len(a.Questions)
	var answerRecords []models.AssessmentAnswer

	for _, q := range a.Questions {
		selectedOptID, answered := answers[q.ID]
		isCorrect := false
		if answered {
			for _, opt := range q.Options {
				if opt.ID == selectedOptID && opt.IsCorrect {
					isCorrect = true
					correctCount++
					break
				}
			}
		}
		answerRecords = append(answerRecords, models.AssessmentAnswer{
			QuestionID:       q.ID,
			SelectedOptionID: selectedOptID,
			IsCorrect:        isCorrect,
		})
	}

	var score float64
	if totalQuestions > 0 {
		score = (float64(correctCount) / float64(totalQuestions)) * 100.0
	}

	sub := models.AssessmentSubmission{
		AssessmentID: assessmentID,
		StudentID:    studentID,
		Score:        &score,
		SubmittedAt:  now,
		Answers:      answerRecords,
	}

	if err := r.db.Create(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *studentApiRepository) GetDashboardStats(classID uint, studentID uint, student *models.Student) (*StudentDashboardStats, error) {
	var stats StudentDashboardStats

	if student != nil {
		stats.StudentName = student.Name
		stats.StudentNISN = student.NISN
		if student.Class != nil {
			stats.ClassName = student.Class.Name
		}
	}

	r.db.Model(&models.ClassSubject{}).Where("class_id = ?", classID).Count(&stats.TotalSubjects)

	r.db.Model(&models.Material{}).Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Where("class_subjects.class_id = ? AND materials.is_published = ?", classID, true).
		Count(&stats.TotalMaterials)

	r.db.Model(&models.Task{}).Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Where("class_subjects.class_id = ? AND tasks.is_published = ?", classID, true).
		Count(&stats.TotalTasks)

	var submittedTaskCount int64
	r.db.Model(&models.TaskSubmission{}).
		Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
		Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Where("class_subjects.class_id = ? AND tasks.is_published = ? AND task_submissions.student_id = ?", classID, true, studentID).
		Count(&submittedTaskCount)

	stats.PendingTasks = stats.TotalTasks - submittedTaskCount
	if stats.PendingTasks < 0 {
		stats.PendingTasks = 0
	}

	r.db.Model(&models.Assessment{}).Joins("JOIN class_subjects ON class_subjects.id = assessments.class_subject_id").
		Where("class_subjects.class_id = ? AND assessments.is_published = ?", classID, true).
		Count(&stats.TotalAssessments)

	return &stats, nil
}
