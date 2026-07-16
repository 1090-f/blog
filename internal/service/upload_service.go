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

// 上传服务相关的错误值。
var (
	ErrUploadTypeNotAllowed = errors.New("upload file type not allowed")     // 文件类型不允许
	ErrInvalidUpload        = errors.New("invalid upload file")              // 上传文件无效
	ErrUploadNotConfigured  = errors.New("upload service is not configured") // 上传服务未配置
	ErrUploadTooLarge       = errors.New("upload file is too large")         // 文件超过大小限制
)

// UploadService 图片上传业务逻辑层，处理文件校验和存储。
type UploadService struct {
	baseDir      string
	baseURL      string
	maxSizeBytes int64
}

// NewUploadService 创建并初始化图片上传实例，baseDir 为本地存储目录，baseURL 为对外访问地址前缀。
func NewUploadService(baseDir, baseURL string, maxSizeBytes int64) *UploadService {
	return &UploadService{
		baseDir:      strings.TrimSpace(baseDir),
		baseURL:      strings.TrimSpace(baseURL),
		maxSizeBytes: maxSizeBytes,
	}
}

// SaveArticleImage 校验上传文件的类型、大小和内容魔数后，按年月目录存储并返回可访问的 URL。
func (s *UploadService) SaveArticleImage(fileHeader *multipart.FileHeader) (string, error) {
	// 基础校验：文件不能为空
	if fileHeader == nil || strings.TrimSpace(fileHeader.Filename) == "" {
		return "", ErrInvalidUpload
	}
	// 上传服务必须配置存储目录和访问地址
	if strings.TrimSpace(s.baseDir) == "" || strings.TrimSpace(s.baseURL) == "" {
		return "", ErrUploadNotConfigured
	}
	// 文件大小校验（multipart 头部声明的大小）
	if s.maxSizeBytes > 0 && fileHeader.Size > s.maxSizeBytes {
		return "", ErrUploadTooLarge
	}
	// 文件扩展名校验（只允许常见图片格式）
	if !uploadpkg.IsAllowedImageExtension(fileHeader.Filename) {
		return "", ErrUploadTypeNotAllowed
	}

	source, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer source.Close()

	// 读取文件头 512 字节，用于检测真实 MIME 类型（防止伪造扩展名）
	header := make([]byte, 512)
	n, err := source.Read(header)
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	contentType := uploadpkg.DetectContentType(header[:n])
	// 文件内容必须是允许的图片类型
	if !uploadpkg.IsAllowedImageContent(contentType) {
		return "", ErrUploadTypeNotAllowed
	}

	// 将已读取的头部和剩余内容拼接成完整 reader
	reader := io.MultiReader(bytes.NewReader(header[:n]), source)
	now := time.Now()
	// 按 articles/年/月 目录结构存储
	relativeDir := filepath.Join("articles", now.Format("2006"), now.Format("01"))
	targetDir := filepath.Join(s.baseDir, relativeDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	filename, err := buildUniqueFileName(fileHeader.Filename, now)
	if err != nil {
		return "", err
	}

	// 拼接完整路径并创建目标文件
	relativePath := filepath.Join(relativeDir, filename)
	targetPath := filepath.Join(s.baseDir, relativePath)
	targetFile, err := os.Create(targetPath)
	if err != nil {
		return "", err
	}
	defer targetFile.Close()

	// 写入时二次限制文件大小，防止 multipart 头部声明与实际内容不一致
	copyReader := io.Reader(reader)
	if s.maxSizeBytes > 0 {
		copyReader = io.LimitReader(reader, s.maxSizeBytes+1)
	}
	written, err := io.Copy(targetFile, copyReader)
	if err != nil {
		// 写入失败时清理已创建的文件
		_ = os.Remove(targetPath)
		return "", err
	}
	// 多 +1 的 LimitReader 用于检测是否超限：实际写入超过限制则拒绝
	if s.maxSizeBytes > 0 && written > s.maxSizeBytes {
		_ = os.Remove(targetPath)
		return "", ErrUploadTooLarge
	}

	return uploadpkg.BuildURL(s.baseURL, relativePath), nil
}

// buildUniqueFileName 生成唯一文件名：时间戳 + 12 位随机十六进制串 + 原始扩展名，避免文件名冲突。
func buildUniqueFileName(originalName string, now time.Time) (string, error) {
	randomBytes := make([]byte, 6)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(originalName)))
	return now.Format("20060102150405") + "-" + hex.EncodeToString(randomBytes) + ext, nil
}
