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

// CategoryController 提供文章分类的 CRUD 操作。
type CategoryController struct {
	categoryService *service.CategoryService
}

// NewCategoryController 创建并初始化分类接口实例。
func NewCategoryController(categoryService *service.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

// List 查询分类接口列表。
func (ctl *CategoryController) List(c *gin.Context) {
	categories, err := ctl.categoryService.List()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5005, "获取分类列表失败")
		return
	}

	response.Success(c, toCategoryResponses(categories))
}

// Create 创建分类接口记录。
func (ctl *CategoryController) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	category, err := ctl.categoryService.Create(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCategory), errors.Is(err, service.ErrCategoryExists):
			response.Error(c, http.StatusBadRequest, 4003, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5006, "创建分类失败")
		}
		return
	}

	response.Success(c, toCategoryResponse(*category))
}

// Update 更新分类接口记录。
func (ctl *CategoryController) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	category, err := ctl.categoryService.Update(id, req)
	if err != nil {
		switch {
		// 分类不存在
		case errors.Is(err, service.ErrCategoryNotFound):
			response.Error(c, http.StatusNotFound, 4042, err.Error())
			// 参数无效或名称重复
		case errors.Is(err, service.ErrInvalidCategory), errors.Is(err, service.ErrCategoryExists):
			response.Error(c, http.StatusBadRequest, 4003, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5007, "更新分类失败")
		}
		return
	}

	response.Success(c, toCategoryResponse(*category))
}

// Delete 删除分类接口记录。
func (ctl *CategoryController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	if err := ctl.categoryService.Delete(id); err != nil {
		switch {
		// 分类不存在
		case errors.Is(err, service.ErrCategoryNotFound):
			response.Error(c, http.StatusNotFound, 4042, err.Error())
			// 分类下还有文章，不能删除
		case errors.Is(err, service.ErrCategoryInUse):
			response.Error(c, http.StatusBadRequest, 4003, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5008, "删除分类失败")
		}
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

// 将内部数据转换为接口响应结构。
func toCategoryResponses(categories []model.Category) []dto.CategoryResponse {
	resp := make([]dto.CategoryResponse, 0, len(categories))
	for _, category := range categories {
		resp = append(resp, toCategoryResponse(category))
	}
	return resp
}

// 将内部数据转换为接口响应结构。
func toCategoryResponse(category model.Category) dto.CategoryResponse {
	return dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}
}
