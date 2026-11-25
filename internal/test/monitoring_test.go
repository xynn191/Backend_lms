package test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
)

func TestSubjectTeachersInMonitoring(t *testing.T) {
	db := setupTestDB(t)

	// Create user & teacher
	user := models.User{ID: "teacher-1", RoleID: 2, Email: "guru1@school.id", Password: "x"}
	db.Create(&user)
	teacher := models.Teacher{UserID: user.ID, Name: "Pak Budi", NIP: "19800101"}
	db.Create(&teacher)

	// Create major & class
	major := models.Major{Code: "RPL", Name: "Rekayasa Perangkat Lunak"}
	db.Create(&major)
	class := models.Class{Name: "X RPL 1", MajorID: major.ID}
	db.Create(&class)

	// Create subject
	subj := models.Subject{Code: "MTK", Name: "Matematika"}
	db.Create(&subj)

	// Assign teacher to class_subject
	cs := models.ClassSubject{ClassID: class.ID, SubjectID: subj.ID, TeacherID: teacher.ID}
	db.Create(&cs)

	repo := repository.NewMonitoringRepository(db)
	data, err := repo.GetCurriculumMonitoring()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	b, _ := json.MarshalIndent(data.SubjectStats, "", "  ")
	fmt.Printf("SubjectStats Result:\n%s\n", string(b))

	if len(data.SubjectStats) == 0 {
		t.Fatalf("expected subjects, got 0")
	}

	s0 := data.SubjectStats[0]
	if len(s0.Teachers) == 0 {
		t.Fatalf("expected at least 1 teacher, got 0")
	}

	if s0.Teachers[0].TeacherName != "Pak Budi" {
		t.Fatalf("expected 'Pak Budi', got '%s'", s0.Teachers[0].TeacherName)
	}
}
