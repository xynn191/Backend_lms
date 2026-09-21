package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xynn191/Backend_lms/internal/service"
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc}
}

// ── Students ──────────────────────────────────────────────────────────────────

func (h *UserHandler) GetAllStudents(c *gin.Context) {
	data, err := h.svc.GetAllStudents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *UserHandler) CreateStudent(c *gin.Context) {
	var dto service.CreateStudentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.CreateStudent(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Siswa berhasil ditambahkan", "data": data})
}

func (h *UserHandler) UpdateStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.UpdateStudentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.UpdateStudent(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data siswa berhasil diupdate", "data": data})
}

func (h *UserHandler) DeleteStudent(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteStudent(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Siswa berhasil dihapus"})
}

// ── Teachers ──────────────────────────────────────────────────────────────────

func (h *UserHandler) GetAllTeachers(c *gin.Context) {
	data, err := h.svc.GetAllTeachers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *UserHandler) CreateTeacher(c *gin.Context) {
	var dto service.CreateTeacherDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.CreateTeacher(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Guru berhasil ditambahkan", "data": data})
}

func (h *UserHandler) UpdateTeacher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.UpdateTeacherDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.UpdateTeacher(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data guru berhasil diupdate", "data": data})
}

func (h *UserHandler) DeleteTeacher(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteTeacher(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Guru berhasil dihapus"})
}

// ── Curriculums ───────────────────────────────────────────────────────────────

func (h *UserHandler) GetAllCurriculums(c *gin.Context) {
	data, err := h.svc.GetAllCurriculums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *UserHandler) CreateCurriculum(c *gin.Context) {
	var dto service.CreateStaffDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.CreateCurriculum(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Tim kurikulum berhasil ditambahkan", "data": data})
}

func (h *UserHandler) UpdateCurriculum(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.UpdateStaffDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.UpdateCurriculum(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data kurikulum berhasil diupdate", "data": data})
}

func (h *UserHandler) DeleteCurriculum(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeleteCurriculum(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Tim kurikulum berhasil dihapus"})
}

// ── Principals ────────────────────────────────────────────────────────────────

func (h *UserHandler) GetAllPrincipals(c *gin.Context) {
	data, err := h.svc.GetAllPrincipals()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func (h *UserHandler) CreatePrincipal(c *gin.Context) {
	var dto service.CreateStaffDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.CreatePrincipal(&dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Kepala sekolah berhasil ditambahkan", "data": data})
}

func (h *UserHandler) UpdatePrincipal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var dto service.UpdateStaffDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	data, err := h.svc.UpdatePrincipal(uint(id), &dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Data kepala sekolah berhasil diupdate", "data": data})
}

func (h *UserHandler) DeletePrincipal(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.svc.DeletePrincipal(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Kepala sekolah berhasil dihapus"})
}
