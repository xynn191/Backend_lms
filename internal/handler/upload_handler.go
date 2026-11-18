package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

const maxUploadSize = 25 * 1024 * 1024 // 25 MB

var allowedExtensions = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".xls":  true,
	".xlsx": true,
	".ppt":  true,
	".pptx": true,
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
	".gif":  true,
	".mp3":  true,
	".wav":  true,
	".ogg":  true,
	".mp4":  true,
	".mkv":  true,
	".webm": true,
}

func (h *UploadHandler) UploadFile(c *gin.Context) {
	// Parse multipart form with 25MB max memory
	if c.Request.ContentLength > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file melebihi batas maksimum 25 MB"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File tidak ditemukan dalam request"})
		return
	}

	if fileHeader.Size > maxUploadSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file melebihi batas maksimum 25 MB"})
		return
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !allowedExtensions[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Tipe file %s tidak diizinkan", ext)})
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat direktori upload"})
		return
	}

	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dst := filepath.Join(uploadDir, newFilename)

	if err := c.SaveUploadedFile(fileHeader, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
		return
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFilename)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload berhasil",
		"url":      fileURL,
		"filename": fileHeader.Filename,
		"size":     fileHeader.Size,
	})
}
