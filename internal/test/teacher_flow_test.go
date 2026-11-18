package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/xynn191/Backend_lms/internal/config"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/router"
	"github.com/xynn191/Backend_lms/internal/service"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dbName := fmt.Sprintf("file:memdb_%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test db: %v", err)
	}

	err = models.AutoMigrate(db)
	if err != nil {
		t.Fatalf("Failed to auto migrate test db: %v", err)
	}

	return db
}

func generateTestToken(userID, role, secret string) string {
	claims := &service.JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString([]byte(secret))
	return tokenStr
}

func TestTeacherFlows(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	cfg := &config.Config{
		JWTSecret:    "test-secret",
		JWTExpiresIn: 1 * time.Hour,
	}

	// 1. Seed Roles
	roleAdmin := models.Role{Name: models.RoleAdmin}
	roleTeacher := models.Role{Name: models.RoleTeacher}
	roleStudent := models.Role{Name: models.RoleStudent}
	db.Create(&roleAdmin)
	db.Create(&roleTeacher)
	db.Create(&roleStudent)

	// 2. Seed Users
	userAdmin := models.User{ID: uuid.New().String(), Email: "admin@test.com", Password: "pass", RoleID: roleAdmin.ID}
	userTeacher1 := models.User{ID: uuid.New().String(), Email: "guru1@test.com", Password: "pass", RoleID: roleTeacher.ID}
	userTeacher2 := models.User{ID: uuid.New().String(), Email: "guru2@test.com", Password: "pass", RoleID: roleTeacher.ID}
	userStudent := models.User{ID: uuid.New().String(), Email: "siswa@test.com", Password: "pass", RoleID: roleStudent.ID}
	db.Create(&userAdmin)
	db.Create(&userTeacher1)
	db.Create(&userTeacher2)
	db.Create(&userStudent)

	// 3. Seed Profiles & Master Data
	teacher1 := models.Teacher{UserID: userTeacher1.ID, Name: "Guru Satu", NIP: "111"}
	teacher2 := models.Teacher{UserID: userTeacher2.ID, Name: "Guru Dua", NIP: "222"}
	db.Create(&teacher1)
	db.Create(&teacher2)

	major := models.Major{Code: "RPL", Name: "Rekayasa Perangkat Lunak"}
	db.Create(&major)

	class := models.Class{MajorID: major.ID, Name: "X RPL 1", GradeLevel: "X"}
	db.Create(&class)

	student := models.Student{UserID: userStudent.ID, NISN: "999", Name: "Siswa Satu", ClassID: class.ID}
	db.Create(&student)

	subj1 := models.Subject{Code: "WEB", Name: "Pemrograman Web"}
	subj2 := models.Subject{Code: "MTK", Name: "Matematika"}
	db.Create(&subj1)
	db.Create(&subj2)

	// Admin assigns ClassSubject
	cs1 := models.ClassSubject{ClassID: class.ID, SubjectID: subj1.ID, TeacherID: teacher1.ID}
	cs2 := models.ClassSubject{ClassID: class.ID, SubjectID: subj2.ID, TeacherID: teacher2.ID}
	db.Create(&cs1)
	db.Create(&cs2)

	r := router.SetupRouter(db, cfg)

	tokenAdmin := generateTestToken(userAdmin.ID, models.RoleAdmin, cfg.JWTSecret)
	tokenTeacher1 := generateTestToken(userTeacher1.ID, models.RoleTeacher, cfg.JWTSecret)
	tokenTeacher2 := generateTestToken(userTeacher2.ID, models.RoleTeacher, cfg.JWTSecret)
	tokenStudent := generateTestToken(userStudent.ID, models.RoleStudent, cfg.JWTSecret)

	// ── TEST 1: Admin prevent duplicate ClassSubject assignment ────────────────
	t.Run("Admin Prevent Duplicate ClassSubject", func(t *testing.T) {
		body, _ := json.Marshal(map[string]interface{}{
			"class_id":   class.ID,
			"subject_id": subj1.ID,
			"teacher_id": teacher2.ID,
		})
		req := httptest.NewRequest("POST", "/api/admin/class-subjects", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+tokenAdmin)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400 Bad Request for duplicate ClassSubject, got %d: %s", w.Code, w.Body.String())
		}
	})

	// ── TEST 2: Teacher Authorization Isolation ────────────────────────────────
	t.Run("Teacher Authorization Isolation", func(t *testing.T) {
		// Teacher 1 creates material for CS1
		body, _ := json.Marshal(map[string]interface{}{
			"class_subject_id": cs1.ID,
			"title":            "Materi Web 1",
			"content":          "Penjelasan HTML",
			"is_published":     false,
		})
		req := httptest.NewRequest("POST", "/api/teacher/materials", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("Expected 201 Created for Teacher 1 material, got %d: %s", w.Code, w.Body.String())
		}

		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		matID := uint(resp["data"]["ID"].(float64))

		// Teacher 2 attempts to view Teacher 1's material
		req2 := httptest.NewRequest("GET", fmt.Sprintf("/api/teacher/materials/%d", matID), nil)
		req2.Header.Set("Authorization", "Bearer "+tokenTeacher2)
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		if w2.Code != http.StatusNotFound {
			t.Fatalf("Expected 404/403 for Teacher 2 accessing Teacher 1 material, got %d", w2.Code)
		}
	})

	// ── TEST 3: Draft vs Published Visibility ──────────────────────────────────
	t.Run("Draft vs Published Visibility", func(t *testing.T) {
		// Teacher 1 creates draft material
		body, _ := json.Marshal(map[string]interface{}{
			"class_subject_id": cs1.ID,
			"title":            "Materi Draft",
			"content":          "Draft Content",
			"is_published":     false,
		})
		req := httptest.NewRequest("POST", "/api/teacher/materials", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		matID := uint(resp["data"]["ID"].(float64))

		// Student checks materials list (Draft should NOT appear)
		reqStu1 := httptest.NewRequest("GET", "/api/student/materials", nil)
		reqStu1.Header.Set("Authorization", "Bearer "+tokenStudent)
		wStu1 := httptest.NewRecorder()
		r.ServeHTTP(wStu1, reqStu1)

		var matList map[string][]map[string]interface{}
		json.Unmarshal(wStu1.Body.Bytes(), &matList)

		found := false
		for _, m := range matList["data"] {
			if uint(m["ID"].(float64)) == matID {
				found = true
			}
		}
		if found {
			t.Fatalf("Draft material should not be visible to student")
		}

		// Teacher publishes the material
		pubBody, _ := json.Marshal(map[string]interface{}{"is_published": true})
		reqPub := httptest.NewRequest("PATCH", fmt.Sprintf("/api/teacher/materials/%d/publish", matID), bytes.NewBuffer(pubBody))
		reqPub.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		reqPub.Header.Set("Content-Type", "application/json")
		wPub := httptest.NewRecorder()
		r.ServeHTTP(wPub, reqPub)

		if wPub.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for publish material, got %d", wPub.Code)
		}

		// Student checks again (Published SHOULD appear)
		reqStu2 := httptest.NewRequest("GET", "/api/student/materials", nil)
		reqStu2.Header.Set("Authorization", "Bearer "+tokenStudent)
		wStu2 := httptest.NewRecorder()
		r.ServeHTTP(wStu2, reqStu2)

		json.Unmarshal(wStu2.Body.Bytes(), &matList)
		found = false
		for _, m := range matList["data"] {
			if uint(m["ID"].(float64)) == matID {
				found = true
			}
		}
		if !found {
			t.Fatalf("Published material should be visible to student")
		}
	})

	// ── TEST 4: File Upload & Size Limits ──────────────────────────────────────
	t.Run("File Upload & Size Limits", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "test_document.pdf")
		part.Write([]byte("%PDF-1.4 test file content"))
		writer.Close()

		req := httptest.NewRequest("POST", "/api/upload", body)
		req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		req.Header.Set("Content-Type", writer.FormDataContentType())
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for valid upload, got %d: %s", w.Code, w.Body.String())
		}

		// Upload executable / invalid file
		body2 := &bytes.Buffer{}
		writer2 := multipart.NewWriter(body2)
		part2, _ := writer2.CreateFormFile("file", "script.exe")
		part2.Write([]byte("malicious content"))
		writer2.Close()

		req2 := httptest.NewRequest("POST", "/api/upload", body2)
		req2.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		req2.Header.Set("Content-Type", writer2.FormDataContentType())
		w2 := httptest.NewRecorder()
		r.ServeHTTP(w2, req2)

		if w2.Code != http.StatusBadRequest {
			t.Fatalf("Expected 400 Bad Request for invalid extension, got %d", w2.Code)
		}
	})

	// ── TEST 5: Task Submission & Teacher Grading ─────────────────────────────
	t.Run("Task Submission & Teacher Grading", func(t *testing.T) {
		// Teacher creates published task
		body, _ := json.Marshal(map[string]interface{}{
			"class_subject_id": cs1.ID,
			"title":            "Tugas HTML Form",
			"description":      "Buat form login",
			"due_date":         time.Now().Add(24 * time.Hour),
			"is_published":     true,
		})
		req := httptest.NewRequest("POST", "/api/teacher/tasks", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		var resp map[string]map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &resp)
		taskID := uint(resp["data"]["ID"].(float64))

		// Student submits task via link
		subBody, _ := json.Marshal(map[string]interface{}{
			"link_url": "https://github.com/student/html-form",
		})
		reqSub := httptest.NewRequest("POST", fmt.Sprintf("/api/student/tasks/%d/submit", taskID), bytes.NewBuffer(subBody))
		reqSub.Header.Set("Authorization", "Bearer "+tokenStudent)
		reqSub.Header.Set("Content-Type", "application/json")
		wSub := httptest.NewRecorder()
		r.ServeHTTP(wSub, reqSub)

		if wSub.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for student task submission, got %d: %s", wSub.Code, wSub.Body.String())
		}

		var subResp map[string]map[string]interface{}
		json.Unmarshal(wSub.Body.Bytes(), &subResp)
		submissionID := uint(subResp["data"]["ID"].(float64))

		// Teacher grades submission
		gradeBody, _ := json.Marshal(map[string]interface{}{
			"score":    95,
			"feedback": "Kerja bagus, tampilan rapi!",
		})
		reqGrade := httptest.NewRequest("PUT", fmt.Sprintf("/api/teacher/submissions/%d/grade", submissionID), bytes.NewBuffer(gradeBody))
		reqGrade.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		reqGrade.Header.Set("Content-Type", "application/json")
		wGrade := httptest.NewRecorder()
		r.ServeHTTP(wGrade, reqGrade)

		if wGrade.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for teacher grading, got %d: %s", wGrade.Code, wGrade.Body.String())
		}
	})

	// ── TEST 6: Multiple Choice Assessment Quiz & Auto-Grading ────────────────
	t.Run("Assessment Quiz & Auto Grading", func(t *testing.T) {
		// Teacher creates active assessment with 2 questions
		startTime := time.Now().Add(-10 * time.Minute)
		endTime := time.Now().Add(50 * time.Minute)

		asmBody, _ := json.Marshal(map[string]interface{}{
			"class_subject_id": cs1.ID,
			"title":            "Ulangan Harian 1",
			"description":      "Kuis Pilihan Ganda",
			"start_time":       startTime,
			"end_time":         endTime,
			"duration_minutes": 30,
			"is_published":     true,
			"questions": []map[string]interface{}{
				{
					"question_text": "Apa singkatan dari HTML?",
					"options": []map[string]interface{}{
						{"option_text": "HyperText Markup Language", "is_correct": true},
						{"option_text": "HighText Machine Language", "is_correct": false},
					},
				},
				{
					"question_text": "Protokol standar untuk web adalah?",
					"options": []map[string]interface{}{
						{"option_text": "HTTP", "is_correct": true},
						{"option_text": "FTP", "is_correct": false},
					},
				},
			},
		})

		reqAsm := httptest.NewRequest("POST", "/api/teacher/assessments", bytes.NewBuffer(asmBody))
		reqAsm.Header.Set("Authorization", "Bearer "+tokenTeacher1)
		reqAsm.Header.Set("Content-Type", "application/json")
		wAsm := httptest.NewRecorder()
		r.ServeHTTP(wAsm, reqAsm)

		if wAsm.Code != http.StatusCreated {
			t.Fatalf("Expected 201 Created for assessment, got %d: %s", wAsm.Code, wAsm.Body.String())
		}

		// Fetch assessment detail as student to get question & option IDs
		reqDetail := httptest.NewRequest("GET", fmt.Sprintf("/api/student/assessments/%s", "1"), nil)
		reqDetail.Header.Set("Authorization", "Bearer "+tokenStudent)
		wDetail := httptest.NewRecorder()
		r.ServeHTTP(wDetail, reqDetail)

		var asmData map[string]interface{}
		json.Unmarshal(wDetail.Body.Bytes(), &asmData)
		dataObj := asmData["data"].(map[string]interface{})
		var questions []interface{}
		if asmMap, ok := dataObj["Assessment"].(map[string]interface{}); ok {
			questions = asmMap["Questions"].([]interface{})
		} else {
			questions = dataObj["Questions"].([]interface{})
		}

		q1 := questions[0].(map[string]interface{})
		q1ID := uint(q1["ID"].(float64))
		q1Opts := q1["Options"].([]interface{})
		opt1ID := uint(q1Opts[0].(map[string]interface{})["ID"].(float64)) // Correct answer

		q2 := questions[1].(map[string]interface{})
		q2ID := uint(q2["ID"].(float64))
		q2Opts := q2["Options"].([]interface{})
		opt2ID := uint(q2Opts[0].(map[string]interface{})["ID"].(float64)) // Correct answer

		// Student submits assessment answers (2 out of 2 correct = score 100)
		ansBody, _ := json.Marshal(map[string]interface{}{
			"answers": []map[string]interface{}{
				{"question_id": q1ID, "selected_option_id": opt1ID},
				{"question_id": q2ID, "selected_option_id": opt2ID},
			},
		})

		reqSubAsm := httptest.NewRequest("POST", "/api/student/assessments/1/submit", bytes.NewBuffer(ansBody))
		reqSubAsm.Header.Set("Authorization", "Bearer "+tokenStudent)
		reqSubAsm.Header.Set("Content-Type", "application/json")
		wSubAsm := httptest.NewRecorder()
		r.ServeHTTP(wSubAsm, reqSubAsm)

		if wSubAsm.Code != http.StatusOK {
			t.Fatalf("Expected 200 OK for student assessment submission, got %d: %s", wSubAsm.Code, wSubAsm.Body.String())
		}

		var subRes map[string]map[string]interface{}
		json.Unmarshal(wSubAsm.Body.Bytes(), &subRes)
		score := subRes["data"]["Score"].(float64)
		if score != 100.0 {
			t.Fatalf("Expected score 100, got %v", score)
		}
	})
}
