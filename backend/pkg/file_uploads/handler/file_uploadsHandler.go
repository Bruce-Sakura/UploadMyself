package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/service"
	"github.com/gin-gonic/gin"
)

type FileUploadHandler struct {
	svc service.FileUploadService
}

func NewFileUploadHandler(svc service.FileUploadService) *FileUploadHandler {
	return &FileUploadHandler{svc: svc}
}

// Register mounts upload/serve routes under the given group (e.g. /api/v1).
func (h *FileUploadHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/upload", h.Upload)
	rg.GET("/files/:id", h.Serve)
	rg.POST("/upload-corpus", h.UploadCorpus)
}

func (h *FileUploadHandler) Upload(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	vo, err := h.svc.SaveUpload(c.Request.Context(), fh)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, vo)
}

func (h *FileUploadHandler) Serve(c *gin.Context) {
	vo, err := h.svc.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	if _, err := os.Stat(vo.StoredPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file missing from disk"})
		return
	}
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%q", vo.OriginalName))
	c.File(vo.StoredPath)
}

func (h *FileUploadHandler) UploadCorpus(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	vo, err := h.svc.ExtractCorpus(c.Request.Context(), fh)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, vo)
}
