package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/config"
	"github.com/xynn191/Backend_lms/internal/handler"
	"github.com/xynn191/Backend_lms/internal/middleware"
	"github.com/xynn191/Backend_lms/internal/models"
	"github.com/xynn191/Backend_lms/internal/repository"
	"github.com/xynn191/Backend_lms/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()
	_ = r.SetTrustedProxies(nil)
	r.Use(cors())

	// Static route for file uploads
	r.Static("/uploads", "./uploads")

	// Dependency injection
	authRepo := repository.NewAuthRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	majorRepo := repository.NewMajorRepository(db)
	classRepo := repository.NewClassRepository(db)
	subjectRepo := repository.NewSubjectRepository(db)
	userRepo := repository.NewUserRepository(db)
	classSubjectRepo := repository.NewClassSubjectRepository(db)
	teacherApiRepo := repository.NewTeacherApiRepository(db)
	studentApiRepo := repository.NewStudentApiRepository(db)
	monitoringRepo := repository.NewMonitoringRepository(db)

	authSvc := service.NewAuthService(authRepo, cfg)
	adminSvc := service.NewAdminService(adminRepo)
	majorSvc := service.NewMajorService(majorRepo)
	classSvc := service.NewClassService(classRepo)
	subjectSvc := service.NewSubjectService(subjectRepo)
	userSvc := service.NewUserService(userRepo)
	classSubjectSvc := service.NewClassSubjectService(classSubjectRepo)
	teacherApiSvc := service.NewTeacherApiService(teacherApiRepo)
	studentApiSvc := service.NewStudentApiService(studentApiRepo)
	monitoringSvc := service.NewMonitoringService(monitoringRepo)

	authH := handler.NewAuthHandler(authSvc)
	adminH := handler.NewAdminHandler(adminSvc)
	majorH := handler.NewMajorHandler(majorSvc)
	classH := handler.NewClassHandler(classSvc)
	subjectH := handler.NewSubjectHandler(subjectSvc)
	userH := handler.NewUserHandler(userSvc)
	classSubjectH := handler.NewClassSubjectHandler(classSubjectSvc)
	teacherApiH := handler.NewTeacherApiHandler(teacherApiSvc)
	studentApiH := handler.NewStudentApiHandler(studentApiSvc)
	monitoringH := handler.NewMonitoringHandler(monitoringSvc)
	uploadH := handler.NewUploadHandler()

	api := r.Group("/api")

	// ── Public Routes ──────────────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/login", authH.Login)
	}

	// ── Protected Routes (All authenticated users) ──────────────────────────────
	protected := api.Group("/")
	protected.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		protected.GET("/auth/me", authH.Me)
		protected.POST("/upload", uploadH.UploadFile)
	}

	// ── Admin Routes ───────────────────────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.AuthRequired(cfg.JWTSecret))
	{
		// Strict Admin Only (Dashboard & Master Data)
		adminOnly := admin.Group("/")
		adminOnly.Use(middleware.RoleRequired(models.RoleAdmin))
		{
			// Dashboard
			adminOnly.GET("/dashboard", adminH.Dashboard)

			// Master Data: Jurusan
			adminOnly.GET("/majors", majorH.GetAll)
			adminOnly.POST("/majors", majorH.Create)
			adminOnly.GET("/majors/:id", majorH.GetByID)
			adminOnly.PUT("/majors/:id", majorH.Update)
			adminOnly.DELETE("/majors/:id", majorH.Delete)

			// Master Data: Kelas (Write operations)
			adminOnly.POST("/classes", classH.Create)
			adminOnly.PUT("/classes/:id", classH.Update)
			adminOnly.DELETE("/classes/:id", classH.Delete)

			// Master Data: Mata Pelajaran
			adminOnly.GET("/subjects", subjectH.GetAll)
			adminOnly.POST("/subjects", subjectH.Create)
			adminOnly.GET("/subjects/:id", subjectH.GetByID)
			adminOnly.PUT("/subjects/:id", subjectH.Update)
			adminOnly.DELETE("/subjects/:id", subjectH.Delete)

			// Penugasan Guru (ClassSubject)
			adminOnly.GET("/class-subjects", classSubjectH.GetAll)
			adminOnly.POST("/class-subjects", classSubjectH.Create)
			adminOnly.GET("/class-subjects/:id", classSubjectH.GetByID)
			adminOnly.PUT("/class-subjects/:id", classSubjectH.Update)
			adminOnly.DELETE("/class-subjects/:id", classSubjectH.Delete)
		}

		// User Management & Shared Lookups (Admin & Superadmin)
		userMgmt := admin.Group("/")
		userMgmt.Use(middleware.RoleRequired(models.RoleAdmin, models.RoleSuperadmin))
		{
			// Read classes for student form dropdowns
			userMgmt.GET("/classes", classH.GetAll)
			userMgmt.GET("/classes/:id", classH.GetByID)

			// Manajemen User: Siswa
			userMgmt.GET("/students", userH.GetAllStudents)
			userMgmt.POST("/students", userH.CreateStudent)
			userMgmt.PUT("/students/:id", userH.UpdateStudent)
			userMgmt.DELETE("/students/:id", userH.DeleteStudent)

			// Manajemen User: Guru
			userMgmt.GET("/teachers", userH.GetAllTeachers)
			userMgmt.POST("/teachers", userH.CreateTeacher)
			userMgmt.PUT("/teachers/:id", userH.UpdateTeacher)
			userMgmt.DELETE("/teachers/:id", userH.DeleteTeacher)

			// Manajemen User: Kurikulum
			userMgmt.GET("/curriculums", userH.GetAllCurriculums)
			userMgmt.POST("/curriculums", userH.CreateCurriculum)
			userMgmt.PUT("/curriculums/:id", userH.UpdateCurriculum)
			userMgmt.DELETE("/curriculums/:id", userH.DeleteCurriculum)

			// Manajemen User: Kepala Sekolah
			userMgmt.GET("/principals", userH.GetAllPrincipals)
			userMgmt.POST("/principals", userH.CreatePrincipal)
			userMgmt.PUT("/principals/:id", userH.UpdatePrincipal)
			userMgmt.DELETE("/principals/:id", userH.DeletePrincipal)
		}
	}

	// ── Superadmin Routes (Strictly User Governance & Admin CRUD) ──────────────
	superadmin := api.Group("/superadmin")
	superadmin.Use(middleware.AuthRequired(cfg.JWTSecret))
	superadmin.Use(middleware.RoleRequired(models.RoleSuperadmin))
	{
		superadmin.GET("/dashboard", userH.GetSuperadminStats)

		// Manajemen User: Administrator Sekolah
		superadmin.GET("/admins", userH.GetAllAdmins)
		superadmin.POST("/admins", userH.CreateAdmin)
		superadmin.PUT("/admins/:id", userH.UpdateAdmin)
		superadmin.DELETE("/admins/:id", userH.DeleteAdmin)
	}

	// ── Teacher Routes ─────────────────────────────────────────────────────────
	teacher := api.Group("/teacher")
	teacher.Use(middleware.AuthRequired(cfg.JWTSecret))
	teacher.Use(middleware.RoleRequired(models.RoleTeacher))
	{
		teacher.GET("/dashboard", teacherApiH.Dashboard)
		teacher.GET("/class-subjects", teacherApiH.ClassSubjects)
		teacher.GET("/class-subjects/:id", teacherApiH.ClassSubjectDetail)

		// Materi
		teacher.GET("/materials", teacherApiH.GetMaterials)
		teacher.POST("/materials", teacherApiH.CreateMaterial)
		teacher.GET("/materials/:id", teacherApiH.GetMaterialByID)
		teacher.PUT("/materials/:id", teacherApiH.UpdateMaterial)
		teacher.DELETE("/materials/:id", teacherApiH.DeleteMaterial)
		teacher.PATCH("/materials/:id/publish", teacherApiH.PublishMaterial)

		// Tugas
		teacher.GET("/tasks", teacherApiH.GetTasks)
		teacher.POST("/tasks", teacherApiH.CreateTask)
		teacher.GET("/tasks/:id", teacherApiH.GetTaskByID)
		teacher.PUT("/tasks/:id", teacherApiH.UpdateTask)
		teacher.DELETE("/tasks/:id", teacherApiH.DeleteTask)
		teacher.PATCH("/tasks/:id/publish", teacherApiH.PublishTask)
		teacher.GET("/tasks/:id/submissions", teacherApiH.TaskSubmissions)
		teacher.PUT("/submissions/:submission_id/grade", teacherApiH.GradeSubmission)

		// Assessment
		teacher.GET("/assessments", teacherApiH.GetAssessments)
		teacher.POST("/assessments", teacherApiH.CreateAssessment)
		teacher.GET("/assessments/:id", teacherApiH.GetAssessmentByID)
		teacher.PUT("/assessments/:id", teacherApiH.UpdateAssessment)
		teacher.PATCH("/assessments/:id/schedule", teacherApiH.UpdateAssessmentSchedule)
		teacher.DELETE("/assessments/:id", teacherApiH.DeleteAssessment)
		teacher.PATCH("/assessments/:id/publish", teacherApiH.PublishAssessment)
		teacher.GET("/assessments/:id/results", teacherApiH.AssessmentResults)

		// Nilai
		teacher.GET("/grades", teacherApiH.GradesOverview)
		teacher.GET("/grades/detail", teacherApiH.GradesDetail)
	}

	// ── Student Routes ─────────────────────────────────────────────────────────
	student := api.Group("/student")
	student.Use(middleware.AuthRequired(cfg.JWTSecret))
	student.Use(middleware.RoleRequired(models.RoleStudent))
	{
		student.GET("/dashboard", studentApiH.Dashboard)
		student.GET("/materials", studentApiH.Materials)
		student.GET("/materials/:id", studentApiH.MaterialDetail)
		student.GET("/tasks", studentApiH.Tasks)
		student.GET("/tasks/:id", studentApiH.TaskDetail)
		student.POST("/tasks/:id/submit", studentApiH.SubmitTask)
		student.GET("/assessments", studentApiH.Assessments)
		student.GET("/assessments/:id", studentApiH.AssessmentDetail)
		student.POST("/assessments/:id/submit", studentApiH.SubmitAssessment)
	}

	// ── Curriculum Routes ──────────────────────────────────────────────────────
	curriculum := api.Group("/curriculum")
	curriculum.Use(middleware.AuthRequired(cfg.JWTSecret))
	curriculum.Use(middleware.RoleRequired(models.RoleCurriculum, models.RoleAdmin))
	{
		curriculum.GET("/monitoring", monitoringH.CurriculumMonitoring)
		curriculum.GET("/grades", monitoringH.SchoolGrades)
		curriculum.GET("/classes", classH.GetAll)
		curriculum.GET("/subjects", subjectH.GetAll)
	}

	// ── Principal Routes ───────────────────────────────────────────────────────
	principal := api.Group("/principal")
	principal.Use(middleware.AuthRequired(cfg.JWTSecret))
	principal.Use(middleware.RoleRequired(models.RolePrincipal, models.RoleAdmin))
	{
		principal.GET("/monitoring", monitoringH.PrincipalMonitoring)
		principal.GET("/grades", monitoringH.SchoolGrades)
		principal.GET("/classes", classH.GetAll)
		principal.GET("/subjects", subjectH.GetAll)
	}

	return r
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "http://localhost:3000" || origin == "http://127.0.0.1:3000" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
