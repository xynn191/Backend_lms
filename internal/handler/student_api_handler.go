package handler

import (
	"net/http"

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

func (h *StudentApiHandler) Tasks(c *gin.Context) {
	userID := c.GetString("user_id")
	data, err := h.svc.GetTasks(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
