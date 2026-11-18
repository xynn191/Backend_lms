package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type StudentApiHandler struct {
	svc service.StudentApiService
}

func NewStudentApiHandler(svc service.StudentApiService) *StudentApiHandler {
	return &StudentApiHandler{svc}
}

func (h *StudentApiHandler) Dashboard(c *gin.Context) {
	userID := c.GetString("user_id")
	stats, err := h.svc.GetDashboardStats(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}

func (h *StudentApiHandler) Materials(c *gin.Context) {
	userID := c.GetString("user_id")
	data, err := h.svc.GetMaterials(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StudentApiHandler) MaterialDetail(c *gin.Context) {
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

func (h *StudentApiHandler) Tasks(c *gin.Context) {
	userID := c.GetString("user_id")
	data, err := h.svc.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StudentApiHandler) TaskDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	item, err := h.svc.GetTaskByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *StudentApiHandler) SubmitTask(c *gin.Context) {
	userID := c.GetString("user_id")
	taskID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tugas tidak valid"})
		return
	}

	var req service.SubmitTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.SubmitTask(userID, uint(taskID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sub, "message": "Tugas berhasil dikumpulkan"})
}

func (h *StudentApiHandler) Assessments(c *gin.Context) {
	userID := c.GetString("user_id")
	data, err := h.svc.GetAssessments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *StudentApiHandler) AssessmentDetail(c *gin.Context) {
	userID := c.GetString("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}
	item, err := h.svc.GetAssessmentByID(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

func (h *StudentApiHandler) SubmitAssessment(c *gin.Context) {
	userID := c.GetString("user_id")
	assessmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID assessment tidak valid"})
		return
	}

	var req service.SubmitAssessmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.svc.SubmitAssessment(userID, uint(assessmentID), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sub, "message": "Assessment berhasil dikumpulkan"})
}
