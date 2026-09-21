package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type MajorHandler struct {
	svc service.MajorService
}

func NewMajorHandler(svc service.MajorService) *MajorHandler {
	return &MajorHandler{svc}
}

func (h *MajorHandler) GetAll(c *gin.Context) {
	data, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *MajorHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jurusan tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *MajorHandler) Create(c *gin.Context) {
	var dto service.MajorDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.Create(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Jurusan berhasil dibuat", "data": data})
}

func (h *MajorHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.MajorDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.Update(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jurusan berhasil diupdate", "data": data})
}

func (h *MajorHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jurusan berhasil dihapus"})
}
