package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type MonitoringHandler struct {
	svc service.MonitoringService
}

func NewMonitoringHandler(svc service.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{svc: svc}
}

// GET /api/curriculum/monitoring
func (h *MonitoringHandler) CurriculumMonitoring(c *gin.Context) {
	data, err := h.svc.GetCurriculumMonitoring()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data monitoring kurikulum"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

// GET /api/principal/monitoring
func (h *MonitoringHandler) PrincipalMonitoring(c *gin.Context) {
	data, err := h.svc.GetPrincipalMonitoring()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data monitoring kepala sekolah"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}
