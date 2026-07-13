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

type TagController struct {
	tagService *service.TagService
}

func NewTagController(tagService *service.TagService) *TagController {
	return &TagController{tagService: tagService}
}

func (ctl *TagController) List(c *gin.Context) {
	tags, err := ctl.tagService.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5028, "failed to fetch tags")
		return
	}
	response.Success(c, toTagResponses(tags))
}

func (ctl *TagController) Create(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	tag, err := ctl.tagService.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidTag), errors.Is(err, service.ErrTagExists):
			response.Error(c, http.StatusBadRequest, 4008, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5029, "failed to create tag")
		}
		return
	}
	response.Success(c, toTagResponse(*tag))
}

func (ctl *TagController) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
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
			response.Error(c, http.StatusInternalServerError, 5030, "failed to update tag")
		}
		return
	}
	response.Success(c, toTagResponse(*tag))
}

func (ctl *TagController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	if err := ctl.tagService.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrTagNotFound):
			response.Error(c, http.StatusNotFound, 4045, err.Error())
		case errors.Is(err, service.ErrTagInUse):
			response.Error(c, http.StatusBadRequest, 4008, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5031, "failed to delete tag")
		}
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func toTagResponses(tags []model.Tag) []dto.TagResponse {
	resp := make([]dto.TagResponse, 0, len(tags))
	for _, tag := range tags {
		resp = append(resp, toTagResponse(tag))
	}
	return resp
}

func toTagResponse(tag model.Tag) dto.TagResponse {
	return dto.TagResponse{ID: tag.ID, Name: tag.Name}
}
