package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/dto"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/entity"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/mapper"
	"github.com/Bruce-Sakura/UploadMyself/backend/pkg/file_uploads/service"
	"github.com/google/uuid"
)

// MaxUploadSize is the max file size in bytes (100MB).
const MaxUploadSize = 100 << 20

var allowedMIME = map[string]bool{
	"audio/wav": true, "audio/x-wav": true, "audio/mpeg": true, "audio/mp3": true,
	"audio/flac": true, "audio/ogg": true, "audio/webm": true,
	"image/png": true, "image/jpeg": true, "image/webp": true, "image/gif": true,
	"application/pdf": true, "text/plain": true, "text/markdown": true, "application/json": true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

var allowedExt = map[string]bool{
	".wav": true, ".mp3": true, ".flac": true, ".ogg": true, ".webm": true,
	".png": true, ".jpg": true, ".jpeg": true, ".webp": true, ".gif": true,
	".pdf": true, ".txt": true, ".json": true, ".md": true, ".doc": true, ".docx": true,
}

// Config holds runtime settings for storage and ML scripts.
type Config struct {
	UploadDir    string
	MLScriptsDir string
	PythonBin    string
}

type FileUploadServiceImpl struct {
	mapper *mapper.FileUploadMapper
	cfg    Config
}

func NewFileUploadService(m *mapper.FileUploadMapper, cfg Config) service.FileUploadService {
	return &FileUploadServiceImpl{mapper: m, cfg: cfg}
}

func (s *FileUploadServiceImpl) SaveUpload(ctx context.Context, fh *multipart.FileHeader) (*dto.UploadResultVO, error) {
	if fh.Size > MaxUploadSize {
		return nil, fmt.Errorf("file too large (max %d bytes)", int64(MaxUploadSize))
	}

	mime := fh.Header.Get("Content-Type")
	if mime == "" {
		mime = "application/octet-stream"
	}
	if !allowedMIME[mime] {
		return nil, fmt.Errorf("file type not allowed: %s", mime)
	}

	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if ext != "" && !allowedExt[ext] {
		return nil, fmt.Errorf("file extension not allowed: %s", ext)
	}

	newName := uuid.NewString() + ext
	savePath := filepath.Join(s.cfg.UploadDir, newName)
	if err := saveMultipart(fh, savePath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	f := &entity.FileUpload{
		ID:           uuid.NewString(),
		OriginalName: fh.Filename,
		StoredPath:   savePath,
		Size:         fh.Size,
		MimeType:     mime,
	}
	if err := s.mapper.Insert(ctx, f); err != nil {
		return nil, err
	}

	return &dto.UploadResultVO{ID: f.ID, Filename: newName, Path: savePath, Size: fh.Size}, nil
}

func (s *FileUploadServiceImpl) Get(ctx context.Context, id string) (*dto.FileUploadVO, error) {
	f, err := s.mapper.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.FileUploadVO{
		ID:           f.ID,
		OriginalName: f.OriginalName,
		StoredPath:   f.StoredPath,
		Size:         f.Size,
		MimeType:     f.MimeType,
		CreatedAt:    f.CreatedAt,
	}, nil
}

func (s *FileUploadServiceImpl) ExtractCorpus(ctx context.Context, fh *multipart.FileHeader) (*dto.CorpusResultVO, error) {
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	tmpPath := filepath.Join(os.TempDir(), "uploadmyself_corpus_"+uuid.NewString()+ext)
	if err := saveMultipart(fh, tmpPath); err != nil {
		return nil, fmt.Errorf("failed to save temp file: %w", err)
	}
	defer os.Remove(tmpPath)

	script := fmt.Sprintf("%s/extract_text.py", s.cfg.MLScriptsDir)
	cmd := exec.CommandContext(ctx, s.cfg.PythonBin, script, "--input", tmpPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("extraction failed: %s", string(out))
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, fmt.Errorf("invalid extraction output")
	}
	if e, ok := result["error"]; ok {
		return nil, fmt.Errorf("%v", e)
	}

	text, _ := result["text"].(string)
	method, _ := result["method"].(string)
	return &dto.CorpusResultVO{Text: text, Method: method, Name: fh.Filename}, nil
}

// saveMultipart copies an uploaded file header to dest on disk.
func saveMultipart(fh *multipart.FileHeader, dest string) error {
	src, err := fh.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}
