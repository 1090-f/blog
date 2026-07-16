package upload

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

// 允许上传的图片文件扩展名集合。
var allowedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
}

// IsAllowedImageExtension 判断文件扩展名是否属于允许的图片类型（不区分大小写）。
func IsAllowedImageExtension(filename string) bool {
	// Ext返回文件的扩展名
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(filename)))
	_, ok := allowedImageExtensions[ext]
	return ok
}

// IsAllowedImageContent 判断 MIME 内容类型是否属于允许的图片类型（不区分大小写）。
func IsAllowedImageContent(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

// DetectContentType 根据文件头魔数识别实际内容类型，不依赖文件扩展名。
func DetectContentType(header []byte) string {
	return http.DetectContentType(header)
}

// BuildURL 根据基础 URL 和相对路径 拼接完整的文件访问地址。
// 支持绝对 URL（如 https://cdn.example.com）和相对路径（如 /uploads）两种基础地址格式。
func BuildURL(baseURL, relativePath string) string {
	// 统一为正斜杠并去除前导斜杠
	normalizedRelativePath := strings.TrimLeft(filepath.ToSlash(strings.TrimSpace(relativePath)), "/")
	if normalizedRelativePath == "" {
		return strings.TrimSpace(baseURL)
	}

	normalizedBaseURL := strings.TrimSpace(baseURL)
	if normalizedBaseURL == "" {
		return "/" + normalizedRelativePath
	}

	// 如果 baseURL 是完整的绝对 URL（含 scheme 和 host），使用 URL 解析拼接
	if parsed, err := url.Parse(normalizedBaseURL); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		ref := &url.URL{Path: path.Join(strings.TrimRight(parsed.Path, "/"), normalizedRelativePath)}
		return parsed.ResolveReference(ref).String()
	}

	// 否则当作相对路径直接拼接
	return path.Join(strings.TrimRight(normalizedBaseURL, "/"), normalizedRelativePath)
}
