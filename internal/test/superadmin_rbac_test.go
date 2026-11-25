package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/config"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/router"
)

func TestSuperadminRBACAndUserManagement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB(t)

	cfg := &config.Config{
		JWTSecret: "superadmin_test_secret_key",
	}

	r := router.SetupRouter(db, cfg)

	// Seed roles & superadmin
	models.Seed(db)

	superadminToken := generateTestToken("00000000-0000-0000-0000-000000000099", models.RoleSuperadmin, cfg.JWTSecret)
	adminToken := generateTestToken("00000000-0000-0000-0000-000000000001", models.RoleAdmin, cfg.JWTSecret)

	// 1. Superadmin accesses /api/superadmin/dashboard -> 200 OK
	t.Run("Superadmin accesses dashboard", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/superadmin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+superadminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}
	})

	// 2. Regular Admin is FORBIDDEN from /api/superadmin/dashboard -> 403 Forbidden
	t.Run("Admin forbidden from superadmin dashboard", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/superadmin/dashboard", nil)
		req.Header.Set("Authorization", "Bearer "+adminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
	})

	// 3. Superadmin creates a new Admin via /api/superadmin/admins -> 201 Created
	t.Run("Superadmin creates an admin", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":        "newadmin@sekolah.id",
			"password":     "password123",
			"name":         "Admin Budi",
			"phone_number": "081234567890",
		}
		jsonBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/superadmin/admins", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Authorization", "Bearer "+superadminToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("expected status 201, got %d: %s", w.Code, w.Body.String())
		}
	})

	// 4. Regular Admin is FORBIDDEN from creating Admin -> 403 Forbidden
	t.Run("Regular Admin forbidden from creating Admin", func(t *testing.T) {
		payload := map[string]interface{}{
			"email":        "unauthorized@sekolah.id",
			"password":     "password123",
			"name":         "Hacker Admin",
			"phone_number": "081234567890",
		}
		jsonBytes, _ := json.Marshal(payload)
		req, _ := http.NewRequest(http.MethodPost, "/api/superadmin/admins", bytes.NewBuffer(jsonBytes))
		req.Header.Set("Authorization", "Bearer "+adminToken)
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
	})

	// 5. Superadmin is FORBIDDEN from accessing Master Data (e.g. /api/admin/majors) -> 403 Forbidden
	t.Run("Superadmin forbidden from master data majors", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/admin/majors", nil)
		req.Header.Set("Authorization", "Bearer "+superadminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusForbidden {
			t.Fatalf("expected status 403, got %d", w.Code)
		}
	})

	// 6. Superadmin can access student user management (/api/admin/students) -> 200 OK
	t.Run("Superadmin can access students list", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/api/admin/students", nil)
		req.Header.Set("Authorization", "Bearer "+superadminToken)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}
	})
}
