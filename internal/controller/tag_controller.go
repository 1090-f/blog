package controller

import (
	"errors"
	"net/http"

	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// TagController 提供文章标签的 CRUD 操作。
type TagController struct {
	tagService *service.TagService
}

// NewTagController 创建并初始化标签接口实例。
func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

// List 查询标签接口列表。
func (ctl *TagController) List(c *gin.Context) {
	tags, err := ctl.tagService.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5028, "获取标签列表失败")
		return
	}
	response.Success(c, toTagResponses(tags))
}

// Create 创建标签接口记录。
func (ctl *TagController) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	tag, err := ctl.tagService.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTag), errors.Is(err, service.ErrTagExists):
			response.Error(c, http.StatusBadRequest, 4008, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5029, "创建标签失败")
		}
		return
	}
	response.Success(c, toTagResponse(*tag))
}

// Update 更新标签接口记录。
func (ctl *TagController) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	tag, err := ctl.tagService.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTagNotFound):
			response.Error(c, http.StatusNotFound, 4045, err.Error())
		case errors.Is(err, service.ErrInvalidTag), errors.Is(err, service.ErrTagExists):
			response.Error(c, http.StatusBadRequest, 4008, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5030, "更新标签失败")
		}
		return
	}
	response.Success(c, toTagResponse(*tag))
}

// Delete 删除标签接口记录。
func (ctl *TagController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	if err := ctl.tagService.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrTagNotFound):
			response.Error(c, http.StatusNotFound, 4045, err.Error())
		case errors.Is(err, service.ErrTagInUse):
			response.Error(c, http.StatusBadRequest, 4008, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5031, "删除标签失败")
		}
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

// 将内部数据转换为接口响应结构。
func toTagResponses(tags []model.Tag) []dto.TagResponse {
	resp := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		resp = append(resp, toTagResponse(tag))
	}
	return resp
}

// 将内部数据转换为接口响应结构。
func toTagResponse(tag model.Tag) dto.TagResponse {
	return dto.TagResponse{ID: tag.ID, Name: tag.Name}
}
