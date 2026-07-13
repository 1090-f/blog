package controller

import (
	"errors"
	"net/http"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadService *service.UploadService
}

func NewUploadController(uploadService *service.UploadService) *UploadController {
	return &UploadController{uploadService: uploadService}
}

func (ctl *UploadController) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	fileURL, err := ctl.uploadService.SaveArticleImage(fileHeader)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUploadTypeNotAllowed), errors.Is(err, service.ErrInvalidUpload), errors.Is(err, service.ErrUploadTooLarge):
			response.Error(c, http.StatusBadRequest, 4006, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5018, "failed to upload file")
		}
		return
	}

	response.Success(c, gin.H{
		"url": fileURL,
	})
}
