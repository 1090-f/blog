package service_test

import (
	"bytes"
	"mime/multipart"
	"os"
	"path/filepath"
	"testing"

	"blog/internal/service"
)

func createUploadHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file failed: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write form file failed: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close writer failed: %v", err)
	}

	reader := multipart.NewReader(bytes.NewReader(body.Bytes()), writer.Boundary())
	form, err := reader.ReadForm(int64(len(body.Bytes()) + 1024))
	if err != nil {
		t.Fatalf("read form failed: %v", err)
	}
	t.Cleanup(func() { _ = form.RemoveAll() })

	return form.File["file"][0]
}

func TestUploadRejectsTextFile(t *testing.T) {
	tempDir := t.TempDir()
	svc := service.NewUploadService(tempDir, "/uploads", 1024)

	fileHeader := createUploadHeader(t, "bad.txt", []byte("plain text"))
	_, err := svc.SaveArticleImage(fileHeader)
	if err != service.ErrUploadTypeNotAllowed {
		t.Fatalf("expected ErrUploadTypeNotAllowed, got %v", err)
	}
}

func TestUploadRejectsFileOverLimit(t *testing.T) {
	tempDir := t.TempDir()
	svc := service.NewUploadService(tempDir, "/uploads", 16)

	content := append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, bytes.Repeat([]byte("a"), 64)...)
	fileHeader := createUploadHeader(t, "large.png", content)
	_, err := svc.SaveArticleImage(fileHeader)
	if err != service.ErrUploadTooLarge {
		t.Fatalf("expected ErrUploadTooLarge, got %v", err)
	}
}

func TestUploadStoresImageAndReturnsURL(t *testing.T) {
	tempDir := t.TempDir()
	svc := service.NewUploadService(tempDir, "/uploads", 1024*1024)

	content := append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, bytes.Repeat([]byte{0}, 64)...)
	fileHeader := createUploadHeader(t, "ok.png", content)
	url, err := svc.SaveArticleImage(fileHeader)
	if err != nil {
		t.Fatalf("expected upload success, got error: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty url")
	}
	matches, err := filepath.Glob(filepath.Join(tempDir, "articles", "*", "*", "*.png"))
	if err != nil {
		t.Fatalf("glob failed: %v", err)
	}
	if len(matches) != 1 {
		t.Fatalf("expected one stored file, got %d", len(matches))
	}
	if _, err := os.Stat(matches[0]); err != nil {
		t.Fatalf("expected stored file to exist: %v", err)
	}
}
