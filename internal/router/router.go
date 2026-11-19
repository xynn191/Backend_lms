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
	admin.Use(middleware.RoleRequired(models.RoleAdmin))
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

		// Penugasan Guru (ClassSubject)
		admin.GET("/class-subjects", classSubjectH.GetAll)
		admin.POST("/class-subjects", classSubjectH.Create)
		admin.GET("/class-subjects/:id", classSubjectH.GetByID)
		admin.PUT("/class-subjects/:id", classSubjectH.Update)
		admin.DELETE("/class-subjects/:id", classSubjectH.Delete)

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
	}

	// ── Principal Routes ───────────────────────────────────────────────────────
	principal := api.Group("/principal")
	principal.Use(middleware.AuthRequired(cfg.JWTSecret))
	principal.Use(middleware.RoleRequired(models.RolePrincipal, models.RoleAdmin))
	{
		principal.GET("/monitoring", monitoringH.PrincipalMonitoring)
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
