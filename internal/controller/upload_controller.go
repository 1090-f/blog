package controller

import (
	"errors"
	"net/http"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// UploadController 处理图片文件上传请求。
type UploadController struct {
	uploadService *service.UploadService
}

// NewUploadController 创建并初始化上传接口实例。
func NewUploadController(uploadService *service.UploadService) *UploadController {
	return &UploadController{uploadService: uploadService}
}

// Upload 接收并保存上传的图片。
func (ctl *UploadController) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	fileURL, err := ctl.uploadService.SaveArticleImage(fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUploadTypeNotAllowed), errors.Is(err, service.ErrInvalidUpload), errors.Is(err, service.ErrUploadTooLarge):
			response.Error(c, http.StatusBadRequest, 4006, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5018, "文件上传失败")
		}
		return
	}

	response.Success(c, gin.H{
		"url": fileURL,
	})
}
