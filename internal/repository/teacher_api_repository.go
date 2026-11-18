package repository

import (
	"errors"
	"time"

	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type TeacherDashboardStats struct {
	TotalClasses       int64 `json:"total_classes"`
	TotalSubjects      int64 `json:"total_subjects"`
	TotalMaterials     int64 `json:"total_materials"`
	TotalTasks         int64 `json:"total_tasks"`
	TotalAssessments   int64 `json:"total_assessments"`
	PendingSubmissions int64 `json:"pending_submissions"`
}

type TeacherGradeItem struct {
	StudentID        uint     `json:"student_id"`
	StudentName      string   `json:"student_name"`
	StudentNISN      string   `json:"student_nisn"`
	ClassID          uint     `json:"class_id"`
	ClassName        string   `json:"class_name"`
	SubjectName      string   `json:"subject_name"`
	TaskAverage      float64  `json:"task_average"`
	AssessmentAverage float64 `json:"assessment_average"`
	FinalScore       float64  `json:"final_score"`
}

type TeacherApiRepository interface {
	GetTeacherByUserID(userID string) (*models.Teacher, error)
	GetDashboardStats(teacherID uint) (*TeacherDashboardStats, error)
	GetClassSubjects(teacherID uint) ([]models.ClassSubject, error)
	GetClassSubjectByID(teacherID uint, id uint) (*models.ClassSubject, error)
	
	// Materials
	GetMaterials(teacherID uint, classSubjectID uint) ([]models.Material, error)
	GetMaterialByID(teacherID uint, id uint) (*models.Material, error)
	CreateMaterial(m *models.Material) error
	UpdateMaterial(m *models.Material) error
	DeleteMaterial(teacherID uint, id uint) error
	PublishMaterial(teacherID uint, id uint, publish bool) error

	// Tasks
	GetTasks(teacherID uint, classSubjectID uint) ([]models.Task, error)
	GetTaskByID(teacherID uint, id uint) (*models.Task, error)
	CreateTask(t *models.Task) error
	UpdateTask(t *models.Task) error
	DeleteTask(teacherID uint, id uint) error
	PublishTask(teacherID uint, id uint, publish bool) error
	GetTaskSubmissions(teacherID uint, taskID uint) ([]models.TaskSubmission, error)
	GradeSubmission(teacherID uint, submissionID uint, score float64, feedback string) (*models.TaskSubmission, error)

	// Assessments
	GetAssessments(teacherID uint, classSubjectID uint) ([]models.Assessment, error)
	GetAssessmentByID(teacherID uint, id uint) (*models.Assessment, error)
	CreateAssessment(a *models.Assessment) error
	UpdateAssessment(a *models.Assessment) error
	DeleteAssessment(teacherID uint, id uint) error
	PublishAssessment(teacherID uint, id uint, publish bool) error
	GetAssessmentResults(teacherID uint, assessmentID uint) ([]models.AssessmentSubmission, error)

	// Grades
	GetGradesOverview(teacherID uint, classSubjectID uint) ([]TeacherGradeItem, error)
	GetGradesDetail(teacherID uint, classSubjectID uint) (*TeacherGradesDetail, error)
}

type teacherApiRepository struct {
	db *gorm.DB
}

func NewTeacherApiRepository(db *gorm.DB) TeacherApiRepository {
	return &teacherApiRepository{db}
}

func (r *teacherApiRepository) GetTeacherByUserID(userID string) (*models.Teacher, error) {
	var teacher models.Teacher
	if err := r.db.Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return nil, err
	}
	return &teacher, nil
}

func (r *teacherApiRepository) GetDashboardStats(teacherID uint) (*TeacherDashboardStats, error) {
	var stats TeacherDashboardStats

	// Count subjects assigned
	r.db.Model(&models.ClassSubject{}).Where("teacher_id = ?", teacherID).Count(&stats.TotalSubjects)

	// Count unique classes taught
	r.db.Model(&models.ClassSubject{}).Where("teacher_id = ?", teacherID).Select("count(distinct class_id)").Scan(&stats.TotalClasses)

	// Count materials created across assigned class-subjects
	r.db.Model(&models.Material{}).
		Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Where("class_subjects.teacher_id = ?", teacherID).
		Count(&stats.TotalMaterials)

	// Count tasks created
	r.db.Model(&models.Task{}).
		Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Where("class_subjects.teacher_id = ?", teacherID).
		Count(&stats.TotalTasks)

	// Count assessments created
	r.db.Model(&models.Assessment{}).
		Joins("JOIN class_subjects ON class_subjects.id = assessments.class_subject_id").
		Where("class_subjects.teacher_id = ?", teacherID).
		Count(&stats.TotalAssessments)

	// Count pending task submissions (score is null)
	r.db.Model(&models.TaskSubmission{}).
		Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
		Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Where("class_subjects.teacher_id = ? AND task_submissions.score IS NULL", teacherID).
		Count(&stats.PendingSubmissions)

	return &stats, nil
}

func (r *teacherApiRepository) GetClassSubjects(teacherID uint) ([]models.ClassSubject, error) {
	var list []models.ClassSubject
	err := r.db.Preload("Class.Major").Preload("Subject").Where("teacher_id = ?", teacherID).Find(&list).Error
	return list, err
}

func (r *teacherApiRepository) GetClassSubjectByID(teacherID uint, id uint) (*models.ClassSubject, error) {
	var cs models.ClassSubject
	err := r.db.Preload("Class.Major").Preload("Subject").Where("id = ? AND teacher_id = ?", id, teacherID).First(&cs).Error
	if err != nil {
		return nil, errors.New("penugasan kelas tidak ditemukan atau Anda tidak memiliki akses")
	}
	return &cs, nil
}

// Materials
func (r *teacherApiRepository) GetMaterials(teacherID uint, classSubjectID uint) ([]models.Material, error) {
	list := []models.Material{}
	query := r.db.Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Preload("ClassSubject.Class").Preload("ClassSubject.Subject").
		Where("class_subjects.teacher_id = ?", teacherID)

	if classSubjectID > 0 {
		query = query.Where("materials.class_subject_id = ?", classSubjectID)
	}

	err := query.Order("materials.created_at DESC").Find(&list).Error
	return list, err
}

func (r *teacherApiRepository) GetMaterialByID(teacherID uint, id uint) (*models.Material, error) {
	var m models.Material
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
		Preload("ClassSubject.Class").Preload("ClassSubject.Subject").
		Where("materials.id = ? AND class_subjects.teacher_id = ?", id, teacherID).
		First(&m).Error
	if err != nil {
		return nil, errors.New("materi tidak ditemukan atau Anda tidak memiliki akses")
	}
	return &m, nil
}

func (r *teacherApiRepository) CreateMaterial(m *models.Material) error {
	// Verify class subject belongs to teacher
	var cs models.ClassSubject
	if err := r.db.First(&cs, m.ClassSubjectID).Error; err != nil {
		return errors.New("kelas-mata pelajaran tidak valid")
	}
	return r.db.Create(m).Error
}

func (r *teacherApiRepository) UpdateMaterial(m *models.Material) error {
	return r.db.Save(m).Error
}

func (r *teacherApiRepository) DeleteMaterial(teacherID uint, id uint) error {
	m, err := r.GetMaterialByID(teacherID, id)
	if err != nil {
		return err
	}
	return r.db.Delete(m).Error
}

func (r *teacherApiRepository) PublishMaterial(teacherID uint, id uint, publish bool) error {
	m, err := r.GetMaterialByID(teacherID, id)
	if err != nil {
		return err
	}
	return r.db.Model(m).Update("is_published", publish).Error
}

// Tasks
func (r *teacherApiRepository) GetTasks(teacherID uint, classSubjectID uint) ([]models.Task, error) {
	list := []models.Task{}
	query := r.db.Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Preload("ClassSubject.Class").Preload("ClassSubject.Subject").
		Where("class_subjects.teacher_id = ?", teacherID)

	if classSubjectID > 0 {
		query = query.Where("tasks.class_subject_id = ?", classSubjectID)
	}

	err := query.Order("tasks.created_at DESC").Find(&list).Error
	return list, err
}

func (r *teacherApiRepository) GetTaskByID(teacherID uint, id uint) (*models.Task, error) {
	var t models.Task
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Preload("ClassSubject.Class").Preload("ClassSubject.Subject").
		Where("tasks.id = ? AND class_subjects.teacher_id = ?", id, teacherID).
		First(&t).Error
	if err != nil {
		return nil, errors.New("tugas tidak ditemukan atau Anda tidak memiliki akses")
	}
	return &t, nil
}

func (r *teacherApiRepository) CreateTask(t *models.Task) error {
	return r.db.Create(t).Error
}

func (r *teacherApiRepository) UpdateTask(t *models.Task) error {
	return r.db.Save(t).Error
}

func (r *teacherApiRepository) DeleteTask(teacherID uint, id uint) error {
	t, err := r.GetTaskByID(teacherID, id)
	if err != nil {
		return err
	}
	return r.db.Delete(t).Error
}

func (r *teacherApiRepository) PublishTask(teacherID uint, id uint, publish bool) error {
	t, err := r.GetTaskByID(teacherID, id)
	if err != nil {
		return err
	}
	return r.db.Model(t).Update("is_published", publish).Error
}

func (r *teacherApiRepository) GetTaskSubmissions(teacherID uint, taskID uint) ([]models.TaskSubmission, error) {
	// Verify task ownership
	_, err := r.GetTaskByID(teacherID, taskID)
	if err != nil {
		return nil, err
	}

	var list []models.TaskSubmission
	err = r.db.Preload("Student.Class").Where("task_id = ?", taskID).Order("submitted_at DESC").Find(&list).Error
	return list, err
}

func (r *teacherApiRepository) GradeSubmission(teacherID uint, submissionID uint, score float64, feedback string) (*models.TaskSubmission, error) {
	var sub models.TaskSubmission
	err := r.db.Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
		Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
		Preload("Student").Preload("Task").
		Where("task_submissions.id = ? AND class_subjects.teacher_id = ?", submissionID, teacherID).
		First(&sub).Error
	if err != nil {
		return nil, errors.New("pengumpulan tugas tidak ditemukan atau Anda tidak memiliki akses")
	}

	now := time.Now()
	sub.Score = &score
	sub.Feedback = feedback
	sub.GradedAt = &now

	if err := r.db.Save(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// Assessments
func (r *teacherApiRepository) GetAssessments(teacherID uint, classSubjectID uint) ([]models.Assessment, error) {
	list := []models.Assessment{}
	query := r.db.Joins("JOIN class_subjects ON class_subjects.id = assessments.class_subject_id").
		Preload("ClassSubject.Class").Preload("ClassSubject.Subject").
		Preload("Questions.Options").
		Where("class_subjects.teacher_id = ?", teacherID)

	if classSubjectID > 0 {
		query = query.Where("assessments.class_subject_id = ?", classSubjectID)
	}

	err := query.Order("assessments.created_at DESC").Find(&list).Error
	return list, err
}

func (r *teacherApiRepository) GetAssessmentByID(teacherID uint, id uint) (*models.Assessment, error) {
	var a models.Assessment
	err := r.db.Joins("JOIN class_subjects ON class_subjects.id = assessments.class_subject_id").
		Preload("ClassSubject.Class").Preload("ClassSubject.Subject").
		Preload("Questions.Options").
		Where("assessments.id = ? AND class_subjects.teacher_id = ?", id, teacherID).
		First(&a).Error
	if err != nil {
		return nil, errors.New("assessment tidak ditemukan atau Anda tidak memiliki akses")
	}
	return &a, nil
}

func (r *teacherApiRepository) CreateAssessment(a *models.Assessment) error {
	return r.db.Create(a).Error
}

func (r *teacherApiRepository) UpdateAssessment(a *models.Assessment) error {
	// Delete old questions & options then recreate for clean edit
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingQuestions []models.AssessmentQuestion
		tx.Where("assessment_id = ?", a.ID).Find(&existingQuestions)
		for _, q := range existingQuestions {
			tx.Where("question_id = ?", q.ID).Delete(&models.AssessmentOption{})
		}
		tx.Where("assessment_id = ?", a.ID).Delete(&models.AssessmentQuestion{})

		return tx.Save(a).Error
	})
}

func (r *teacherApiRepository) DeleteAssessment(teacherID uint, id uint) error {
	a, err := r.GetAssessmentByID(teacherID, id)
	if err != nil {
		return err
	}
	return r.db.Delete(a).Error
}

func (r *teacherApiRepository) PublishAssessment(teacherID uint, id uint, publish bool) error {
	a, err := r.GetAssessmentByID(teacherID, id)
	if err != nil {
		return err
	}
	return r.db.Model(a).Update("is_published", publish).Error
}

func (r *teacherApiRepository) GetAssessmentResults(teacherID uint, assessmentID uint) ([]models.AssessmentSubmission, error) {
	_, err := r.GetAssessmentByID(teacherID, assessmentID)
	if err != nil {
		return nil, err
	}

	var list []models.AssessmentSubmission
	err = r.db.Preload("Student.Class").Preload("Answers").Where("assessment_id = ?", assessmentID).Order("submitted_at DESC").Find(&list).Error
	return list, err
}

func (r *teacherApiRepository) GetGradesOverview(teacherID uint, classSubjectID uint) ([]TeacherGradeItem, error) {
	// Find class subjects taught by teacher
	var csList []models.ClassSubject
	query := r.db.Preload("Class").Preload("Subject").Where("teacher_id = ?", teacherID)
	if classSubjectID > 0 {
		query = query.Where("id = ?", classSubjectID)
	}
	if err := query.Find(&csList).Error; err != nil {
		return nil, err
	}

	results := []TeacherGradeItem{}

	for _, cs := range csList {
		var students []models.Student
		r.db.Where("class_id = ?", cs.ClassID).Find(&students)

		for _, st := range students {
			// Calculate task average
			var taskAvg float64
			r.db.Table("task_submissions").
				Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
				Where("tasks.class_subject_id = ? AND task_submissions.student_id = ? AND task_submissions.score IS NOT NULL", cs.ID, st.ID).
				Select("COALESCE(AVG(task_submissions.score), 0)").
				Scan(&taskAvg)

			// Calculate assessment average
			var asmAvg float64
			r.db.Table("assessment_submissions").
				Joins("JOIN assessments ON assessments.id = assessment_submissions.assessment_id").
				Where("assessments.class_subject_id = ? AND assessment_submissions.student_id = ? AND assessment_submissions.score IS NOT NULL", cs.ID, st.ID).
				Select("COALESCE(AVG(assessment_submissions.score), 0)").
				Scan(&asmAvg)

			finalScore := (taskAvg + asmAvg) / 2.0
			if taskAvg == 0 && asmAvg != 0 {
				finalScore = asmAvg
			} else if asmAvg == 0 && taskAvg != 0 {
				finalScore = taskAvg
			}

			results = append(results, TeacherGradeItem{
				StudentID:         st.ID,
				StudentName:       st.Name,
				StudentNISN:       st.NISN,
				ClassID:           cs.ClassID,
				ClassName:         cs.Class.Name,
				SubjectName:       cs.Subject.Name,
				TaskAverage:       taskAvg,
				AssessmentAverage: asmAvg,
				FinalScore:        finalScore,
			})
		}
	}

	return results, nil
}

// ── Detail Grades (per assessment & per task) ──────────────────────────────────

type StudentScoreRow struct {
	StudentID   uint    `json:"student_id"`
	StudentName string  `json:"student_name"`
	StudentNISN string  `json:"student_nisn"`
	Score       *float64 `json:"score"`
	SubmittedAt *string `json:"submitted_at,omitempty"`
	SubmissionID *uint  `json:"submission_id,omitempty"`
	Feedback    string  `json:"feedback,omitempty"`
}

type AssessmentGrade struct {
	AssessmentID    uint              `json:"assessment_id"`
	AssessmentTitle string            `json:"assessment_title"`
	Students        []StudentScoreRow `json:"students"`
}

type TaskGrade struct {
	TaskID       uint              `json:"task_id"`
	TaskTitle    string            `json:"task_title"`
	DueDate      string            `json:"due_date"`
	Students     []StudentScoreRow `json:"students"`
}

type ClassSubjectGradesDetail struct {
	ClassSubjectID uint              `json:"class_subject_id"`
	ClassName      string            `json:"class_name"`
	SubjectName    string            `json:"subject_name"`
	Assessments    []AssessmentGrade `json:"assessments"`
	Tasks          []TaskGrade       `json:"tasks"`
}

type TeacherGradesDetail struct {
	ClassSubjects []ClassSubjectGradesDetail `json:"class_subjects"`
}

func (r *teacherApiRepository) GetGradesDetail(teacherID uint, classSubjectID uint) (*TeacherGradesDetail, error) {
	var csList []models.ClassSubject
	query := r.db.Preload("Class").Preload("Subject").Where("teacher_id = ?", teacherID)
	if classSubjectID > 0 {
		query = query.Where("id = ?", classSubjectID)
	}
	if err := query.Find(&csList).Error; err != nil {
		return nil, err
	}

	detail := &TeacherGradesDetail{}

	for _, cs := range csList {
		var students []models.Student
		r.db.Where("class_id = ?", cs.ClassID).Find(&students)


		csDetail := ClassSubjectGradesDetail{
			ClassSubjectID: cs.ID,
			ClassName:      cs.Class.Name,
			SubjectName:    cs.Subject.Name,
		}

		// --- Assessments ---
		var assessments []models.Assessment
		r.db.Where("class_subject_id = ?", cs.ID).Find(&assessments)

		for _, asm := range assessments {
			asmGrade := AssessmentGrade{
				AssessmentID:    asm.ID,
				AssessmentTitle: asm.Title,
			}
			for _, st := range students {
				row := StudentScoreRow{
					StudentID:   st.ID,
					StudentName: st.Name,
					StudentNISN: st.NISN,
				}
				var sub models.AssessmentSubmission
				err := r.db.Where("assessment_id = ? AND student_id = ?", asm.ID, st.ID).First(&sub).Error
				if err == nil {
					row.Score = sub.Score
					ts := sub.SubmittedAt.Format(time.RFC3339)
					row.SubmittedAt = &ts
					row.SubmissionID = &sub.ID
				}
				asmGrade.Students = append(asmGrade.Students, row)
			}
			csDetail.Assessments = append(csDetail.Assessments, asmGrade)
		}

		// --- Tasks ---
		var tasks []models.Task
		r.db.Where("class_subject_id = ?", cs.ID).Find(&tasks)

		for _, task := range tasks {
			taskGrade := TaskGrade{
				TaskID:    task.ID,
				TaskTitle: task.Title,
				DueDate:   task.DueDate.Format(time.RFC3339),
			}
			for _, st := range students {
				row := StudentScoreRow{
					StudentID:   st.ID,
					StudentName: st.Name,
					StudentNISN: st.NISN,
				}
				var sub models.TaskSubmission
				err := r.db.Where("task_id = ? AND student_id = ?", task.ID, st.ID).First(&sub).Error
				if err == nil {
					row.Score = sub.Score
					row.Feedback = sub.Feedback
					row.SubmissionID = &sub.ID
					ts := sub.SubmittedAt.Format(time.RFC3339)
					row.SubmittedAt = &ts
				}
				taskGrade.Students = append(taskGrade.Students, row)
			}
			csDetail.Tasks = append(csDetail.Tasks, taskGrade)
		}

		detail.ClassSubjects = append(detail.ClassSubjects, csDetail)
	}

	return detail, nil
}

