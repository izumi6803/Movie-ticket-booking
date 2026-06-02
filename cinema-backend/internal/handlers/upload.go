package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UploadDir   = "uploads"
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	AllowedExts = ".jpg,.jpeg,.png,.gif,.webp"
)

type UploadHandler struct {
	baseURL string
}

func NewUploadHandler(baseURL string) *UploadHandler {
	// Create uploads directory if not exists
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}

	return &UploadHandler{baseURL: baseURL}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	// Get file from request
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "no file provided"})
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "file too large (max 10MB)"})
		return
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !strings.Contains(AllowedExts, ext) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid file type. Allowed: " + AllowedExts})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102_150405"), uuid.New().String()[:8], ext)
	filepath := filepath.Join(UploadDir, filename)

	// Create file
	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to save file"})
		return
	}
	defer out.Close()

	// Copy file content
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to save file"})
		return
	}

	// Return file URL
	fileURL := fmt.Sprintf("%s/uploads/%s", h.baseURL, filename)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url":      fileURL,
			"filename": filename,
			"size":     header.Size,
		},
	})
}

func (h *UploadHandler) DeleteImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "filename required"})
		return
	}

	filepath := filepath.Join(UploadDir, filename)

	// Security check: prevent directory traversal
	if !strings.HasPrefix(filepath, UploadDir) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid filename"})
		return
	}

	if err := os.Remove(filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "file deleted"})
}

func (h *UploadHandler) ServeImage(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.Status(http.StatusNotFound)
		return
	}

	filepath := filepath.Join(UploadDir, filename)

	// Security check
	if !strings.HasPrefix(filepath, UploadDir) {
		c.Status(http.StatusNotFound)
		return
	}

	c.File(filepath)
}
