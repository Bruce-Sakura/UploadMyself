package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadDir is the directory for stored uploads (set from main).
var UploadDir = "./uploads"

// MaxUploadSize is the max file size in bytes (100MB).
const MaxUploadSize = 100 << 20

// AllowedMIMETypes is the set of allowed MIME types for upload.
var AllowedMIMETypes = map[string]bool{
	"audio/wav":            true,
	"audio/x-wav":          true,
	"audio/mpeg":           true,
	"audio/mp3":            true,
	"audio/flac":           true,
	"audio/ogg":            true,
	"audio/webm":           true,
	"image/png":            true,
	"image/jpeg":           true,
	"image/webp":           true,
	"image/gif":            true,
	"application/pdf":      true,
	"text/plain":           true,
	"text/markdown":        true,
	"application/json":     true,
	"application/msword":   true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

// UploadFile handles multipart file upload, saves with UUID filename.
func (h *Handler) UploadFile(c *gin.Context) {
	// Limit request body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Validate MIME type
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if !AllowedMIMETypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file type not allowed: %s", mimeType)})
		return
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".wav": true, ".mp3": true, ".flac": true, ".ogg": true, ".webm": true,
		".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true,
		".pdf": true, ".txt": true, ".json": true, ".md": true,
		".doc": true, ".docx": true,
	}
	if ext != "" && !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file extension not allowed: %s", ext)})
		return
	}

	newName := uuid.New().String() + ext
	savePath := filepath.Join(UploadDir, newName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	fu := model.FileUpload{
		ID:           uuid.New().String(),
		OriginalName: file.Filename,
		StoredPath:   savePath,
		Size:         file.Size,
		MimeType:     mimeType,
	}
	h.db.Create(&fu)

	c.JSON(http.StatusCreated, gin.H{
		"id":       fu.ID,
		"filename": newName,
		"path":     savePath,
		"size":     file.Size,
	})
}

// ServeFile serves an uploaded file by its ID.
func (h *Handler) ServeFile(c *gin.Context) {
	id := c.Param("id")
	var fu model.FileUpload
	if err := h.db.First(&fu, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	if _, err := os.Stat(fu.StoredPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file missing from disk"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", fu.OriginalName))
	c.File(fu.StoredPath)
}
