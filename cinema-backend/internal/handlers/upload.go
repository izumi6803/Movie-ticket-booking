package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cinema-backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	UploadDir   = "uploads"
	MaxFileSize = 10 * 1024 * 1024 // 10MB
	AllowedExts = ".jpg,.jpeg,.png,.gif,.webp"
)

type UploadHandler struct {
	baseURL           string
	cloudinaryService *services.CloudinaryService
}

func NewUploadHandler(baseURL string, cloudinaryService *services.CloudinaryService) *UploadHandler {
	if err := os.MkdirAll(UploadDir, 0755); err != nil {
		panic(fmt.Sprintf("failed to create upload directory: %v", err))
	}

	return &UploadHandler{
		baseURL:           baseURL,
		cloudinaryService: cloudinaryService,
	}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "no file provided"})
		return
	}
	defer file.Close()

	if header.Size > MaxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "file too large (max 10MB)"})
		return
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !strings.Contains(AllowedExts, ext) {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid file type. Allowed: " + AllowedExts})
		return
	}

	if h.cloudinaryService != nil && h.cloudinaryService.IsEnabled() {
		url, err := h.cloudinaryService.Upload(c.Request.Context(), file, header)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to upload to cloudinary: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"url":      url,
				"filename": header.Filename,
				"size":     header.Size,
			},
		})
		return
	}

	// Fallback to local filesystem (dev mode)

	filename := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102_150405"), uuid.New().String()[:8], ext)
	filepath := filepath.Join(UploadDir, filename)

	out, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to save file"})
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to save file"})
		return
	}

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

	if h.cloudinaryService != nil && h.cloudinaryService.IsEnabled() {
		if err := h.cloudinaryService.Delete(c.Request.Context(), filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to delete file from cloudinary"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "file deleted"})
		return
	}

	filepath := filepath.Join(UploadDir, filename)
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
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "filename required"})
		return
	}

	absUploadDir, err := filepath.Abs(UploadDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "server error"})
		return
	}

	filePath := filepath.Join(absUploadDir, filename)
	if !strings.HasPrefix(filePath, absUploadDir) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "invalid filename"})
		return
	}

	fileInfo, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "file not found"})
		return
	}

	ext := strings.ToLower(filepath.Ext(filename))
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(filePath)
}
