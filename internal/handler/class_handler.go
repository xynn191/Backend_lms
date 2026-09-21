package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type ClassHandler struct {
	svc service.ClassService
}

func NewClassHandler(svc service.ClassService) *ClassHandler {
	return &ClassHandler{svc}
}

func (h *ClassHandler) GetAll(c *gin.Context) {
	data, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *ClassHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kelas tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *ClassHandler) Create(c *gin.Context) {
	var dto service.ClassDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.Create(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Kelas berhasil dibuat", "data": data})
}

func (h *ClassHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.ClassDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.Update(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kelas berhasil diupdate", "data": data})
}

func (h *ClassHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kelas berhasil dihapus"})
}
