package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type AdminHandler struct {
	svc service.AdminService
}

func NewAdminHandler(svc service.AdminService) *AdminHandler {
	return &AdminHandler{svc}
}

// GET /api/admin/dashboard
func (h *AdminHandler) Dashboard(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil statistik"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": stats})
}
