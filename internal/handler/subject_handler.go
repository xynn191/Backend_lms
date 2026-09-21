package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type SubjectHandler struct {
	svc service.SubjectService
}

func NewSubjectHandler(svc service.SubjectService) *SubjectHandler {
	return &SubjectHandler{svc}
}

func (h *SubjectHandler) GetAll(c *gin.Context) {
	data, err := h.svc.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SubjectHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, err := h.svc.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mata pelajaran tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *SubjectHandler) Create(c *gin.Context) {
	var dto service.SubjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.Create(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Mata pelajaran berhasil dibuat", "data": data})
}

func (h *SubjectHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.SubjectDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.Update(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil diupdate", "data": data})
}

func (h *SubjectHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mata pelajaran berhasil dihapus"})
}
