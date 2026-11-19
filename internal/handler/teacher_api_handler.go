package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type TeacherApiHandler struct {
	svc service.TeacherApiService
}

func NewTeacherApiHandler(svc service.TeacherApiService) *TeacherApiHandler {
	return &TeacherApiHandler{svc}
}

func (h *TeacherApiHandler) Dashboard(c *gin.Context) {
	userID := c.GetString("user_id")
	stats, err := h.svc.GetDashboardStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *TeacherApiHandler) ClassSubjects(c *gin.Context) {
	userID := c.GetString("user_id")
	list, err := h.svc.GetClassSubjects(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TeacherApiHandler) ClassSubjectDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	cs, err := h.svc.GetClassSubjectByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": cs})
}

// Materials
func (h *TeacherApiHandler) GetMaterials(c *gin.Context) {
	userID := c.GetString("user_id")
	classSubjectID, _ := strconv.ParseUint(c.Query("class_subject_id"), 10, 32)
	list, err := h.svc.GetMaterials(userID, uint(classSubjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TeacherApiHandler) GetMaterialByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	m, err := h.svc.GetMaterialByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *TeacherApiHandler) CreateMaterial(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.CreateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.svc.CreateMaterial(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": m})
}

func (h *TeacherApiHandler) UpdateMaterial(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.UpdateMaterialRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.svc.UpdateMaterial(userID, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": m})
}

func (h *TeacherApiHandler) DeleteMaterial(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.svc.DeleteMaterial(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Materi berhasil dihapus"})
}

func (h *TeacherApiHandler) PublishMaterial(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.PublishMaterial(userID, uint(id), req.IsPublished); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status materi berhasil diperbarui"})
}

// Tasks
func (h *TeacherApiHandler) GetTasks(c *gin.Context) {
	userID := c.GetString("user_id")
	classSubjectID, _ := strconv.ParseUint(c.Query("class_subject_id"), 10, 32)
	list, err := h.svc.GetTasks(userID, uint(classSubjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TeacherApiHandler) GetTaskByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	t, err := h.svc.GetTaskByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

func (h *TeacherApiHandler) CreateTask(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.svc.CreateTask(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": t})
}

func (h *TeacherApiHandler) UpdateTask(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := h.svc.UpdateTask(userID, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": t})
}

func (h *TeacherApiHandler) DeleteTask(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.svc.DeleteTask(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tugas berhasil dihapus"})
}

func (h *TeacherApiHandler) PublishTask(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.PublishTask(userID, uint(id), req.IsPublished); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status tugas berhasil diperbarui"})
}

func (h *TeacherApiHandler) TaskSubmissions(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tugas tidak valid"})
		return
	}
	list, err := h.svc.GetTaskSubmissions(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TeacherApiHandler) GradeSubmission(c *gin.Context) {
	userID := c.GetString("user_id")
	subID, err := strconv.ParseUint(c.Param("submission_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pengumpulan tidak valid"})
		return
	}
	var req service.GradeSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sub, err := h.svc.GradeSubmission(userID, uint(subID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sub, "message": "Penilaian berhasil disimpan"})
}

// Assessments
func (h *TeacherApiHandler) GetAssessments(c *gin.Context) {
	userID := c.GetString("user_id")
	classSubjectID, _ := strconv.ParseUint(c.Query("class_subject_id"), 10, 32)
	list, err := h.svc.GetAssessments(userID, uint(classSubjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TeacherApiHandler) GetAssessmentByID(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	a, err := h.svc.GetAssessmentByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

func (h *TeacherApiHandler) CreateAssessment(c *gin.Context) {
	userID := c.GetString("user_id")
	var req service.CreateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.svc.CreateAssessment(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": a})
}

func (h *TeacherApiHandler) UpdateAssessment(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.UpdateAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.svc.UpdateAssessment(userID, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a})
}

func (h *TeacherApiHandler) UpdateAssessmentSchedule(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.UpdateAssessmentScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	a, err := h.svc.UpdateAssessmentSchedule(userID, uint(id), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": a, "message": "Jadwal assessment berhasil diperbarui"})
}

func (h *TeacherApiHandler) DeleteAssessment(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	if err := h.svc.DeleteAssessment(userID, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Assessment berhasil dihapus"})
}

func (h *TeacherApiHandler) PublishAssessment(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	var req service.PublishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.PublishAssessment(userID, uint(id), req.IsPublished); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Status assessment berhasil diperbarui"})
}

func (h *TeacherApiHandler) AssessmentResults(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID assessment tidak valid"})
		return
	}
	results, err := h.svc.GetAssessmentResults(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": results})
}

func (h *TeacherApiHandler) GradesOverview(c *gin.Context) {
	userID := c.GetString("user_id")
	classSubjectID, _ := strconv.ParseUint(c.Query("class_subject_id"), 10, 32)
	list, err := h.svc.GetGradesOverview(userID, uint(classSubjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *TeacherApiHandler) GradesDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	classSubjectID, _ := strconv.ParseUint(c.Query("class_subject_id"), 10, 32)
	detail, err := h.svc.GetGradesDetail(userID, uint(classSubjectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": detail})
}
