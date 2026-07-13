package service

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	uploadpkg "blog/pkg/upload"
)

var (
	ErrUploadTypeNotAllowed = errors.New("upload file type not allowed")
	ErrInvalidUpload        = errors.New("invalid upload file")
	ErrUploadNotConfigured  = errors.New("upload service is not configured")
	ErrUploadTooLarge       = errors.New("upload file is too large")
)

type UploadService struct {
	baseDir      string
	baseURL      string
	maxSizeBytes int64
}

func NewUploadService(baseDir, baseURL string, maxSizeBytes int64) *UploadService {
	return &UploadService{
		baseDir:      strings.TrimSpace(baseDir),
		baseURL:      strings.TrimSpace(baseURL),
		maxSizeBytes: maxSizeBytes,
	}
}

func (s *UploadService) SaveArticleImage(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader == nil || strings.TrimSpace(fileHeader.Filename) == "" {
		return "", ErrInvalidUpload
	}
	if strings.TrimSpace(s.baseDir) == "" || strings.TrimSpace(s.baseURL) == "" {
		return "", ErrUploadNotConfigured
	}
	if s.maxSizeBytes > 0 && fileHeader.Size > s.maxSizeBytes {
		return "", ErrUploadTooLarge
	}
	if !uploadpkg.IsAllowedImageExtension(fileHeader.Filename) {
		return "", ErrUploadTypeNotAllowed
	}

	source, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	header := make([]byte, 512)
	n, err := source.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	contentType := uploadpkg.DetectContentType(header[:n])
	if !uploadpkg.IsAllowedImageContent(contentType) {
		return "", ErrUploadTypeNotAllowed
	}

	reader := io.MultiReader(bytes.NewReader(header[:n]), source)
	now := time.Now()
	relativeDir := filepath.Join("articles", now.Format("2006"), now.Format("01"))
	targetDir := filepath.Join(s.baseDir, relativeDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	filename, err := buildUniqueFileName(fileHeader.Filename, now)
	if err != nil {
		return "", err
	}

	relativePath := filepath.Join(relativeDir, filename)
	targetPath := filepath.Join(s.baseDir, relativePath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	copyReader := io.Reader(reader)
	if s.maxSizeBytes > 0 {
		// Enforce the limit while reading as well as from the multipart metadata.
		copyReader = io.LimitReader(reader, s.maxSizeBytes+1)
	}
	written, err := io.Copy(targetFile, copyReader)
	if err != nil {
		_ = os.Remove(targetPath)
		return "", err
	}
	if s.maxSizeBytes > 0 && written > s.maxSizeBytes {
		_ = os.Remove(targetPath)
		return "", ErrUploadTooLarge
	}

	return uploadpkg.BuildURL(s.baseURL, relativePath), nil
}

func buildUniqueFileName(originalName string, now time.Time) (string, error) {
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(originalName)))
	return now.Format("20060102150405") + "-" + hex.EncodeToString(randomBytes) + ext, nil
}
