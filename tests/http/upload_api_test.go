package http_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"blog/internal/controller"
	"blog/internal/service"
	"github.com/gin-gonic/gin"
)

func uploadRequest(t *testing.T, filename string, content []byte) (*http.Request, string) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file failed: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("write form file failed: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("close multipart writer failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, writer.FormDataContentType()
}

func TestUploadRejectsTextFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctl := controller.NewUploadController(service.NewUploadService(t.TempDir(), "/uploads", 1024))
	r := gin.New()
	r.POST("/api/upload", ctl.Upload)

	req, _ := uploadRequest(t, "bad.txt", []byte("plain text"))
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUploadRejectsLargeFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctl := controller.NewUploadController(service.NewUploadService(t.TempDir(), "/uploads", 16))
	r := gin.New()
	r.POST("/api/upload", ctl.Upload)

	content := append([]byte{0x89, 'P', 'N', 'G', '\r', '\n', 0x1a, '\n'}, bytes.Repeat([]byte("a"), 64)...)
	req, _ := uploadRequest(t, "large.png", content)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
