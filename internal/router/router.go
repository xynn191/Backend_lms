package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/config"
	"github.com/xynn191/Backend_lms/internal/handler"
	"github.com/xynn191/Backend_lms/internal/middleware"
	"github.com/xynn191/Backend_lms/internal/repository"
	"github.com/xynn191/Backend_lms/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	r.Use(cors())

	// Dependency injection
	authRepo := repository.NewAuthRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	majorRepo := repository.NewMajorRepository(db)
	classRepo := repository.NewClassRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	userRepo := repository.NewUserRepository(db)

	authSvc := service.NewAuthService(authRepo, cfg)
	adminSvc := service.NewAdminService(adminRepo)
	majorSvc := service.NewMajorService(majorRepo)
	classSvc := service.NewClassService(classRepo)
	subjectSvc := service.NewSubjectService(subjectRepo)
	userSvc := service.NewUserService(userRepo)

	authH := handler.NewAuthHandler(authSvc)
	adminH := handler.NewAdminHandler(adminSvc)
	majorH := handler.NewMajorHandler(majorSvc)
	classH := handler.NewClassHandler(classSvc)
	subjectH := handler.NewSubjectHandler(subjectSvc)
	userH := handler.NewUserHandler(userSvc)

	api := r.Group("/api")

	// ── Public Routes ──────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/login", authH.Login)
	}

	// ── Protected Routes ───────────────────────────────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		protected.GET("/auth/me", authH.Me)
	}

	// ── Admin Routes ───────────────────────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecret))
	admin.Use(middleware.RoleRequired("admin"))
	{
		// Dashboard
		admin.GET("/dashboard", adminH.Dashboard)

		// Master Data: Jurusan
		admin.GET("/majors", majorH.GetAll)
		admin.POST("/majors", majorH.Create)
		admin.GET("/majors/:id", majorH.GetByID)
		admin.PUT("/majors/:id", majorH.Update)
		admin.DELETE("/majors/:id", majorH.Delete)

		// Master Data: Kelas
		admin.GET("/classes", classH.GetAll)
		admin.POST("/classes", classH.Create)
		admin.GET("/classes/:id", classH.GetByID)
		admin.PUT("/classes/:id", classH.Update)
		admin.DELETE("/classes/:id", classH.Delete)

		// Master Data: Mata Pelajaran
		admin.GET("/subjects", subjectH.GetAll)
		admin.POST("/subjects", subjectH.Create)
		admin.GET("/subjects/:id", subjectH.GetByID)
		admin.PUT("/subjects/:id", subjectH.Update)
		admin.DELETE("/subjects/:id", subjectH.Delete)

		// Manajemen User: Siswa
		admin.GET("/students", userH.GetAllStudents)
		admin.POST("/students", userH.CreateStudent)
		admin.PUT("/students/:id", userH.UpdateStudent)
		admin.DELETE("/students/:id", userH.DeleteStudent)

		// Manajemen User: Guru
		admin.GET("/teachers", userH.GetAllTeachers)
		admin.POST("/teachers", userH.CreateTeacher)
		admin.PUT("/teachers/:id", userH.UpdateTeacher)
		admin.DELETE("/teachers/:id", userH.DeleteTeacher)

		// Manajemen User: Kurikulum
		admin.GET("/curriculums", userH.GetAllCurriculums)
		admin.POST("/curriculums", userH.CreateCurriculum)
		admin.PUT("/curriculums/:id", userH.UpdateCurriculum)
		admin.DELETE("/curriculums/:id", userH.DeleteCurriculum)

		// Manajemen User: Kepala Sekolah
		admin.GET("/principals", userH.GetAllPrincipals)
		admin.POST("/principals", userH.CreatePrincipal)
		admin.PUT("/principals/:id", userH.UpdatePrincipal)
		admin.DELETE("/principals/:id", userH.DeletePrincipal)
	}

	return r
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
