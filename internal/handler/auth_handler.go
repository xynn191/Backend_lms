package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type AuthHandler struct {
	svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc}
}

// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.svc.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login berhasil",
		"data":    result,
	})
}

// GET /api/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.svc.GetProfile(userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role.Name,
		},
	})
}
