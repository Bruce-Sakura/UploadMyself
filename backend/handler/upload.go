package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Bruce-Sakura/UploadMyself/backend/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadDir is the directory for stored uploads (set from main).
var UploadDir = "./uploads"

// UploadFile handles multipart file upload, saves with UUID filename.
func (h *Handler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	// Generate UUID filename preserving extension
	ext := filepath.Ext(file.Filename)
	newName := uuid.New().String() + ext
	savePath := filepath.Join(UploadDir, newName)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Detect mime type
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
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
