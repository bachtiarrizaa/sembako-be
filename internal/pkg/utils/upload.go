package utils

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type UploadConfig struct {
	FieldName   string
	UploadDir   string
	AllowedExts []string
	MaxSize     int64
}

type UploadResult struct {
	FileName string
	FilePath string
}

func DefaultImageConfig(uploadDir string) UploadConfig {
	return UploadConfig{
		FieldName:   "image",
		UploadDir:   uploadDir,
		AllowedExts: []string{".jpg", ".jpeg", ".png", ".webp"},
		MaxSize:     2 * 1024 * 1024,
	}
}

func HandleFileUpload(c *gin.Context, cfg UploadConfig) (*UploadResult, error) {
	file, err := c.FormFile(cfg.FieldName)
	if err != nil {
		return nil, nil
	}

	if err := validateFileExtension(file, cfg.AllowedExts); err != nil {
		return nil, err
	}

	if file.Size > cfg.MaxSize {
		maxMB := cfg.MaxSize / (1024 * 1024)
		return nil, NewUploadError(fmt.Sprintf("file size exceeds maximum limit of %dMB", maxMB))
	}

	if err := os.MkdirAll(cfg.UploadDir, os.ModePerm); err != nil {
		return nil, NewUploadError("failed to create upload directory")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	baseName := strings.TrimSuffix(file.Filename, filepath.Ext(file.Filename))
	sanitized := sanitizeFileName(baseName)
	fileName := fmt.Sprintf("%d_%s%s", time.Now().UnixMilli(), sanitized, ext)
	targetPath := filepath.Join(cfg.UploadDir, fileName)

	if err := c.SaveUploadedFile(file, targetPath); err != nil {
		return nil, NewUploadError("failed to save uploaded file")
	}

	cleanPath := filepath.ToSlash(filepath.Clean(filepath.Join(cfg.UploadDir, fileName)))
	cleanPath = strings.TrimPrefix(cleanPath, "./")
	cleanPath = strings.TrimPrefix(cleanPath, "/")
	relPath := "/" + cleanPath

	return &UploadResult{
		FileName: fileName,
		FilePath: relPath,
	}, nil
}

func validateFileExtension(file *multipart.FileHeader, allowedExts []string) error {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	for _, allowed := range allowedExts {
		if ext == allowed {
			return nil
		}
	}

	extsStr := strings.Join(allowedExts, ", ")
	return NewUploadError(fmt.Sprintf("invalid file format. Allowed formats: %s", extsStr))
}

type UploadError struct {
	Message string
}

func (e *UploadError) Error() string {
	return e.Message
}

func NewUploadError(message string) *UploadError {
	return &UploadError{Message: message}
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	re := regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	name = re.ReplaceAllString(name, "")
	re = regexp.MustCompile(`-{2,}`)
	name = re.ReplaceAllString(name, "-")
	name = strings.Trim(name, "-")
	if name == "" {
		name = "file"
	}
	return strings.ToLower(name)
}

func DeleteFile(relativePath string) {
	if relativePath == "" {
		return
	}
	localPath := "." + relativePath
	_ = os.Remove(localPath)
}
