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

type CategoryController struct {
	categoryService *service.CategoryService
}

func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (ctl *CategoryController) List(c *gin.Context) {
	categories, err := ctl.categoryService.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5005, "failed to fetch categories")
		return
	}

	response.Success(c, toCategoryResponses(categories))
}

func (ctl *CategoryController) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	category, err := ctl.categoryService.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCategory), errors.Is(err, service.ErrCategoryExists):
			response.Error(c, http.StatusBadRequest, 4003, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5006, "failed to create category")
		}
		return
	}

	response.Success(c, toCategoryResponse(*category))
}

func (ctl *CategoryController) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	category, err := ctl.categoryService.Update(id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			response.Error(c, http.StatusNotFound, 4042, err.Error())
		case errors.Is(err, service.ErrInvalidCategory), errors.Is(err, service.ErrCategoryExists):
			response.Error(c, http.StatusBadRequest, 4003, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5007, "failed to update category")
		}
		return
	}

	response.Success(c, toCategoryResponse(*category))
}

func (ctl *CategoryController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	if err := ctl.categoryService.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrCategoryNotFound):
			response.Error(c, http.StatusNotFound, 4042, err.Error())
		case errors.Is(err, service.ErrCategoryInUse):
			response.Error(c, http.StatusBadRequest, 4003, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5008, "failed to delete category")
		}
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func toCategoryResponses(categories []model.Category) []dto.CategoryResponse {
	resp := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		resp = append(resp, toCategoryResponse(category))
	}
	return resp
}

func toCategoryResponse(category model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}
