package upload

import (
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

var allowedImageExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
}

func IsAllowedImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(strings.TrimSpace(filename)))
	_, ok := allowedImageExtensions[ext]
	return ok
}

func IsAllowedImageContent(contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/png", "image/webp":
		return true
	default:
		return false
	}
}

func DetectContentType(header []byte) string {
	return http.DetectContentType(header)
}

func BuildURL(baseURL, relativePath string) string {
	normalizedRelativePath := strings.TrimLeft(filepath.ToSlash(strings.TrimSpace(relativePath)), "/")
	if normalizedRelativePath == "" {
		return strings.TrimSpace(baseURL)
	}

	normalizedBaseURL := strings.TrimSpace(baseURL)
	if normalizedBaseURL == "" {
		return "/" + normalizedRelativePath
	}

	if parsed, err := url.Parse(normalizedBaseURL); err == nil && parsed.Scheme != "" && parsed.Host != "" {
		ref := &url.URL{Path: path.Join(strings.TrimRight(parsed.Path, "/"), normalizedRelativePath)}
		return parsed.ResolveReference(ref).String()
	}

	return path.Join(strings.TrimRight(normalizedBaseURL, "/"), normalizedRelativePath)
}
