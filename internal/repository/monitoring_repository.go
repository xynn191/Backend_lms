package repository

import (
	"math"

	"github.com/xynn191/Backend_lms/internal/models"
	"gorm.io/gorm"
)

type ClassProgressItem struct {
	ClassID         uint    `json:"class_id"`
	ClassName       string  `json:"class_name"`
	MajorName       string  `json:"major_name"`
	StudentCount    int64   `json:"student_count"`
	MaterialCount   int64   `json:"material_count"`
	TaskCount       int64   `json:"task_count"`
	SubmissionCount int64   `json:"submission_count"`
	AverageScore    float64 `json:"average_score"`
}

type SubjectMonitoringItem struct {
	SubjectID     uint    `json:"subject_id"`
	SubjectCode   string  `json:"subject_code"`
	SubjectName   string  `json:"subject_name"`
	TeacherCount  int64   `json:"teacher_count"`
	MaterialCount int64   `json:"material_count"`
	TaskCount     int64   `json:"task_count"`
	AverageScore  float64 `json:"average_score"`
}

type CurriculumMonitoringData struct {
	TotalClasses           int64                   `json:"total_classes"`
	TotalSubjects          int64                   `json:"total_subjects"`
	TotalTeachers          int64                   `json:"total_teachers"`
	TotalMaterials         int64                   `json:"total_materials"`
	TotalTasks             int64                   `json:"total_tasks"`
	TotalAssessments       int64                   `json:"total_assessments"`
	TotalSubmissions       int64                   `json:"total_submissions"`
	SyllabusCompletionRate float64                 `json:"syllabus_completion_rate"`
	TaskSubmissionRate     float64                 `json:"task_submission_rate"`
	AverageSchoolScore     float64                 `json:"average_school_score"`
	ClassProgress          []ClassProgressItem     `json:"class_progress"`
	SubjectStats           []SubjectMonitoringItem `json:"subject_stats"`
}

type MajorStatItem struct {
	MajorID      uint    `json:"major_id"`
	MajorName    string  `json:"major_name"`
	ClassCount   int64   `json:"class_count"`
	StudentCount int64   `json:"student_count"`
	AverageScore float64 `json:"average_score"`
}

type TeacherActivityItem struct {
	TeacherID     uint   `json:"teacher_id"`
	TeacherName   string `json:"teacher_name"`
	NIP           string `json:"nip"`
	SubjectCount  int64  `json:"subject_count"`
	MaterialCount int64  `json:"material_count"`
	TaskCount     int64  `json:"task_count"`
}

type PrincipalMonitoringData struct {
	TotalStudents         int64                 `json:"total_students"`
	TotalTeachers         int64                 `json:"total_teachers"`
	TotalClasses          int64                 `json:"total_classes"`
	TotalMajors           int64                 `json:"total_majors"`
	TotalMaterials        int64                 `json:"total_materials"`
	TotalTasks            int64                 `json:"total_tasks"`
	TotalAssessments      int64                 `json:"total_assessments"`
	AcademicParticipation float64               `json:"academic_participation"`
	TeacherEffectiveness  float64               `json:"teacher_effectiveness"`
	AverageSchoolScore    float64               `json:"average_school_score"`
	MajorDistribution     []MajorStatItem       `json:"major_distribution"`
	TopTeacherActivity    []TeacherActivityItem `json:"top_teacher_activity"`
}

type MonitoringRepository interface {
	GetCurriculumMonitoring() (*CurriculumMonitoringData, error)
	GetPrincipalMonitoring() (*PrincipalMonitoringData, error)
}

type monitoringRepository struct {
	db *gorm.DB
}

func NewMonitoringRepository(db *gorm.DB) MonitoringRepository {
	return &monitoringRepository{db: db}
}

func round(val float64, precision int) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (r *monitoringRepository) GetCurriculumMonitoring() (*CurriculumMonitoringData, error) {
	data := &CurriculumMonitoringData{}

	r.db.Model(&models.Class{}).Count(&data.TotalClasses)
	r.db.Model(&models.Subject{}).Count(&data.TotalSubjects)
	r.db.Model(&models.Teacher{}).Count(&data.TotalTeachers)
	r.db.Model(&models.Material{}).Where("is_published = ?", true).Count(&data.TotalMaterials)
	r.db.Model(&models.Task{}).Where("is_published = ?", true).Count(&data.TotalTasks)
	r.db.Model(&models.Assessment{}).Where("is_published = ?", true).Count(&data.TotalAssessments)
	r.db.Model(&models.TaskSubmission{}).Count(&data.TotalSubmissions)

	// Hitung keterisian silabus: class_subject yang sudah memiliki materi terbit vs total class_subject
	var totalClassSubjects int64
	r.db.Model(&models.ClassSubject{}).Count(&totalClassSubjects)
	if totalClassSubjects > 0 {
		var activeClassSubjects int64
		r.db.Model(&models.Material{}).
			Where("is_published = ?", true).
			Distinct("class_subject_id").
			Count(&activeClassSubjects)
		data.SyllabusCompletionRate = round((float64(activeClassSubjects)/float64(totalClassSubjects))*100, 1)
	}

	// Hitung task submission rate
	var totalStudents int64
	r.db.Model(&models.Student{}).Count(&totalStudents)
	if data.TotalTasks > 0 && totalStudents > 0 {
		expectedSubmissions := float64(data.TotalTasks * totalStudents)
		submissionRate := (float64(data.TotalSubmissions) / expectedSubmissions) * 100
		if submissionRate > 100 {
			submissionRate = 100
		}
		data.TaskSubmissionRate = round(submissionRate, 1)
	} else if data.TotalSubmissions > 0 {
		data.TaskSubmissionRate = 100
	}

	// Hitung rerata nilai sekolah (gabungan task submission & assessment submission)
	var avgTaskScore *float64
	r.db.Model(&models.TaskSubmission{}).Where("score IS NOT NULL").Select("AVG(score)").Scan(&avgTaskScore)
	var avgAssessScore *float64
	r.db.Model(&models.AssessmentSubmission{}).Where("score IS NOT NULL").Select("AVG(score)").Scan(&avgAssessScore)

	var sumScore float64
	var scoreCount float64
	if avgTaskScore != nil {
		sumScore += *avgTaskScore
		scoreCount++
	}
	if avgAssessScore != nil {
		sumScore += *avgAssessScore
		scoreCount++
	}
	if scoreCount > 0 {
		data.AverageSchoolScore = round(sumScore/scoreCount, 1)
	}

	// Progress per Kelas
	var classes []models.Class
	r.db.Preload("Major").Order("name ASC").Find(&classes)
	for _, c := range classes {
		var stCount int64
		r.db.Model(&models.Student{}).Where("class_id = ?", c.ID).Count(&stCount)

		var matCount int64
		r.db.Model(&models.Material{}).
			Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
			Where("class_subjects.class_id = ? AND materials.is_published = ?", c.ID, true).
			Count(&matCount)

		var tskCount int64
		r.db.Model(&models.Task{}).
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Where("class_subjects.class_id = ? AND tasks.is_published = ?", c.ID, true).
			Count(&tskCount)

		var subCount int64
		r.db.Model(&models.TaskSubmission{}).
			Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Where("class_subjects.class_id = ?", c.ID).
			Count(&subCount)

		var classAvg *float64
		r.db.Model(&models.TaskSubmission{}).
			Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Where("class_subjects.class_id = ? AND task_submissions.score IS NOT NULL", c.ID).
			Select("AVG(task_submissions.score)").
			Scan(&classAvg)

		avg := 0.0
		if classAvg != nil {
			avg = round(*classAvg, 1)
		}

		majorName := "-"
		if c.Major.Name != "" {
			majorName = c.Major.Name
		}

		data.ClassProgress = append(data.ClassProgress, ClassProgressItem{
			ClassID:         c.ID,
			ClassName:       c.Name,
			MajorName:       majorName,
			StudentCount:    stCount,
			MaterialCount:   matCount,
			TaskCount:       tskCount,
			SubmissionCount: subCount,
			AverageScore:    avg,
		})
	}

	// Statistik per Mata Pelajaran
	var subjects []models.Subject
	r.db.Order("name ASC").Find(&subjects)
	for _, s := range subjects {
		var teacherCount int64
		r.db.Model(&models.ClassSubject{}).Where("subject_id = ?", s.ID).Distinct("teacher_id").Count(&teacherCount)

		var matCount int64
		r.db.Model(&models.Material{}).
			Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
			Where("class_subjects.subject_id = ? AND materials.is_published = ?", s.ID, true).
			Count(&matCount)

		var tskCount int64
		r.db.Model(&models.Task{}).
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Where("class_subjects.subject_id = ? AND tasks.is_published = ?", s.ID, true).
			Count(&tskCount)

		var subjAvg *float64
		r.db.Model(&models.TaskSubmission{}).
			Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Where("class_subjects.subject_id = ? AND task_submissions.score IS NOT NULL", s.ID).
			Select("AVG(task_submissions.score)").
			Scan(&subjAvg)

		avg := 0.0
		if subjAvg != nil {
			avg = round(*subjAvg, 1)
		}

		data.SubjectStats = append(data.SubjectStats, SubjectMonitoringItem{
			SubjectID:     s.ID,
			SubjectCode:   s.Code,
			SubjectName:   s.Name,
			TeacherCount:  teacherCount,
			MaterialCount: matCount,
			TaskCount:     tskCount,
			AverageScore:  avg,
		})
	}

	return data, nil
}

func (r *monitoringRepository) GetPrincipalMonitoring() (*PrincipalMonitoringData, error) {
	data := &PrincipalMonitoringData{}

	r.db.Model(&models.Student{}).Count(&data.TotalStudents)
	r.db.Model(&models.Teacher{}).Count(&data.TotalTeachers)
	r.db.Model(&models.Class{}).Count(&data.TotalClasses)
	r.db.Model(&models.Major{}).Count(&data.TotalMajors)
	r.db.Model(&models.Material{}).Where("is_published = ?", true).Count(&data.TotalMaterials)
	r.db.Model(&models.Task{}).Where("is_published = ?", true).Count(&data.TotalTasks)
	r.db.Model(&models.Assessment{}).Where("is_published = ?", true).Count(&data.TotalAssessments)

	// Hitung rata-rata skor sekolah
	var avgTaskScore *float64
	r.db.Model(&models.TaskSubmission{}).Where("score IS NOT NULL").Select("AVG(score)").Scan(&avgTaskScore)
	var avgAssessScore *float64
	r.db.Model(&models.AssessmentSubmission{}).Where("score IS NOT NULL").Select("AVG(score)").Scan(&avgAssessScore)

	var sumScore float64
	var scoreCount float64
	if avgTaskScore != nil {
		sumScore += *avgTaskScore
		scoreCount++
	}
	if avgAssessScore != nil {
		sumScore += *avgAssessScore
		scoreCount++
	}
	if scoreCount > 0 {
		data.AverageSchoolScore = round(sumScore/scoreCount, 1)
	}

	// Partisipasi Akademik
	var totalSubmissions int64
	r.db.Model(&models.TaskSubmission{}).Count(&totalSubmissions)
	if data.TotalTasks > 0 && data.TotalStudents > 0 {
		expected := float64(data.TotalTasks * data.TotalStudents)
		part := (float64(totalSubmissions) / expected) * 100
		if part > 100 {
			part = 100
		}
		data.AcademicParticipation = round(part, 1)
	} else if totalSubmissions > 0 {
		data.AcademicParticipation = 95.0
	} else {
		data.AcademicParticipation = 0.0
	}

	// Efektivitas Pendidik (% guru yang sudah aktif membuat materi/tugas)
	if data.TotalTeachers > 0 {
		var activeTeachersCount int64
		r.db.Model(&models.ClassSubject{}).
			Joins("JOIN materials ON materials.class_subject_id = class_subjects.id").
			Distinct("class_subjects.teacher_id").
			Count(&activeTeachersCount)
		data.TeacherEffectiveness = round((float64(activeTeachersCount)/float64(data.TotalTeachers))*100, 1)
	}

	// Distribusi Jurusan
	var majors []models.Major
	r.db.Order("name ASC").Find(&majors)
	for _, m := range majors {
		var classCount int64
		r.db.Model(&models.Class{}).Where("major_id = ?", m.ID).Count(&classCount)

		var stCount int64
		r.db.Model(&models.Student{}).
			Joins("JOIN classes ON classes.id = students.class_id").
			Where("classes.major_id = ?", m.ID).
			Count(&stCount)

		var majorAvg *float64
		r.db.Model(&models.TaskSubmission{}).
			Joins("JOIN tasks ON tasks.id = task_submissions.task_id").
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Joins("JOIN classes ON classes.id = class_subjects.class_id").
			Where("classes.major_id = ? AND task_submissions.score IS NOT NULL", m.ID).
			Select("AVG(task_submissions.score)").
			Scan(&majorAvg)

		avg := 0.0
		if majorAvg != nil {
			avg = round(*majorAvg, 1)
		}

		data.MajorDistribution = append(data.MajorDistribution, MajorStatItem{
			MajorID:      m.ID,
			MajorName:    m.Name,
			ClassCount:   classCount,
			StudentCount: stCount,
			AverageScore: avg,
		})
	}

	// Aktivitas Guru
	var teachers []models.Teacher
	r.db.Order("name ASC").Limit(10).Find(&teachers)
	for _, t := range teachers {
		var subjCount int64
		r.db.Model(&models.ClassSubject{}).Where("teacher_id = ?", t.ID).Distinct("subject_id").Count(&subjCount)

		var matCount int64
		r.db.Model(&models.Material{}).
			Joins("JOIN class_subjects ON class_subjects.id = materials.class_subject_id").
			Where("class_subjects.teacher_id = ?", t.ID).
			Count(&matCount)

		var tskCount int64
		r.db.Model(&models.Task{}).
			Joins("JOIN class_subjects ON class_subjects.id = tasks.class_subject_id").
			Where("class_subjects.teacher_id = ?", t.ID).
			Count(&tskCount)

		data.TopTeacherActivity = append(data.TopTeacherActivity, TeacherActivityItem{
			TeacherID:     t.ID,
			TeacherName:   t.Name,
			NIP:           t.NIP,
			SubjectCount:  subjCount,
			MaterialCount: matCount,
			TaskCount:     tskCount,
		})
	}

	return data, nil
}
