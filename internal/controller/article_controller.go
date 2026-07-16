package controller

import (
	"errors"
	"net/http"
	"strconv"

	"blog/internal/dto"
	"blog/internal/model"
	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// ArticleController 提供文章的完整 CRUD 及公开列表、详情、最新、热门查询。
type ArticleController struct {
	articleService *service.ArticleService
}

// NewArticleController 创建并初始化文章接口实例。
func NewArticleController(articleService *service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

// List 查询文章接口列表。
func (ctl *ArticleController) List(c *gin.Context) {
	var query dto.ArticleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	articles, total, page, pageSize, err := ctl.articleService.ListPublished(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5009, "获取文章列表失败")
		return
	}

	response.Success(c, dto.ArticleListResponse{
		List:     toArticleSummaryResponses(articles),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Detail 获取文章详情。
func (ctl *ArticleController) Detail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	article, err := ctl.articleService.GetPublishedDetail(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5010, "获取文章失败")
		}
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

// FullDetail 获取包含评论的完整文章详情。
func (ctl *ArticleController) FullDetail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	detail, err := ctl.articleService.GetPublishedFullDetail(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5019, "获取文章详情失败")
		}
		return
	}

	response.Success(c, toArticleFullDetailResponse(detail))
}

// Latest 获取最新发布的文章。
func (ctl *ArticleController) Latest(c *gin.Context) {
	articles, err := ctl.articleService.ListLatest(parseLimitQuery(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5020, "获取最新文章失败")
		return
	}

	response.Success(c, toArticleSummaryResponses(articles))
}

// Popular 获取热门文章。
func (ctl *ArticleController) Popular(c *gin.Context) {
	articles, err := ctl.articleService.ListPopular(parseLimitQuery(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5021, "获取热门文章失败")
		return
	}

	response.Success(c, toArticleSummaryResponses(articles))
}

// AdminList 获取管理端文章列表。
func (ctl *ArticleController) AdminList(c *gin.Context) {
	var query dto.AdminArticleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	articles, total, page, pageSize, err := ctl.articleService.ListAdmin(query)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidArticleStatus):
			response.Error(c, http.StatusBadRequest, 4004, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5011, "获取管理端文章列表失败")
		}
		return
	}

	response.Success(c, dto.ArticleListResponse{
		List:     toArticleSummaryResponses(articles),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// AdminDetail 获取管理端文章详情，不增加浏览量。
func (ctl *ArticleController) AdminDetail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	article, err := ctl.articleService.GetAdminDetail(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5010, "获取文章失败")
		}
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

// Create 创建文章接口记录。
func (ctl *ArticleController) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	article, err := ctl.articleService.Create(req, c.GetUint("userID"))
	if err != nil {
		respondArticleWriteError(c, err, "创建文章失败", 5012)
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

// Update 更新文章接口记录。
func (ctl *ArticleController) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	article, err := ctl.articleService.Update(id, req, c.GetUint("userID"))
	if err != nil {
		respondArticleWriteError(c, err, "更新文章失败", 5013)
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

// Delete 删除文章接口记录。
func (ctl *ArticleController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	if err := ctl.articleService.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		case errors.Is(err, service.ErrArticleInUse):
			response.Error(c, http.StatusBadRequest, 4004, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5014, "删除文章失败")
		}
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

// 将写入失败原因转换为 HTTP 错误响应。
func respondArticleWriteError(c *gin.Context, err error, message string, internalCode int) {
	switch {
	case errors.Is(err, service.ErrArticleNotFound):
		response.Error(c, http.StatusNotFound, 4043, err.Error())
	case errors.Is(err, service.ErrCategoryNotFound):
		response.Error(c, http.StatusBadRequest, 4003, err.Error())
	case errors.Is(err, service.ErrTagNotFound):
		response.Error(c, http.StatusBadRequest, 4008, err.Error())
	case errors.Is(err, service.ErrInvalidArticle), errors.Is(err, service.ErrInvalidArticleStatus):
		response.Error(c, http.StatusBadRequest, 4004, err.Error())
	default:
		response.Error(c, http.StatusInternalServerError, internalCode, message)
	}
}

// 解析并校验请求参数。
func parseUintParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		return 0, false
	}
	return uint(value), true
}

// 解析并校验请求参数。
func parseLimitQuery(c *gin.Context) int {
	limitText := c.Query("limit")
	if limitText == "" {
		return 0
	}

	limit, err := strconv.Atoi(limitText)
	if err != nil {
		return 0
	}

	return limit
}

// 将内部数据转换为接口响应结构。
func toArticleSummaryResponses(articles []model.Article) []dto.ArticleSummaryResponse {
	resp := make([]dto.ArticleSummaryResponse, 0, len(articles))
	for _, article := range articles {
		resp = append(resp, dto.ArticleSummaryResponse{
			ID:         article.ID,
			Title:      article.Title,
			Summary:    article.Summary,
			CoverImage: article.CoverImage,
			Status:     article.Status,
			ViewCount:  article.ViewCount,
			CategoryID: article.CategoryID,
			Category: dto.ArticleCategoryResponse{
				ID:   article.Category.ID,
				Name: article.Category.Name,
			},
			Tags: toTagResponses(article.Tags),
			Author: dto.ArticleAuthorResponse{
				ID:       article.User.ID,
				Username: article.User.Username,
				Nickname: article.User.Nickname,
			},
			CreatedAt: article.CreatedAt,
			UpdatedAt: article.UpdatedAt,
		})
	}
	return resp
}

// 将内部数据转换为接口响应结构。
func toArticleDetailResponse(article model.Article) dto.ArticleDetailResponse {
	return dto.ArticleDetailResponse{
		ID:         article.ID,
		Title:      article.Title,
		Summary:    article.Summary,
		Content:    article.Content,
		CoverImage: article.CoverImage,
		Status:     article.Status,
		ViewCount:  article.ViewCount,
		CategoryID: article.CategoryID,
		Category: dto.ArticleCategoryResponse{
			ID:   article.Category.ID,
			Name: article.Category.Name,
		},
		Tags: toTagResponses(article.Tags),
		Author: dto.ArticleAuthorResponse{
			ID:       article.User.ID,
			Username: article.User.Username,
			Nickname: article.User.Nickname,
		},
		CreatedAt: article.CreatedAt,
		UpdatedAt: article.UpdatedAt,
	}
}

// 将内部数据转换为接口响应结构。
func toArticleFullDetailResponse(detail *service.ArticleFullDetail) dto.ArticleFullDetailResponse {
	if detail == nil || detail.Article == nil {
		return dto.ArticleFullDetailResponse{}
	}

	article := detail.Article
	return dto.ArticleFullDetailResponse{
		ID:         article.ID,
		Title:      article.Title,
		Summary:    article.Summary,
		Content:    article.Content,
		CoverImage: article.CoverImage,
		Status:     article.Status,
		ViewCount:  article.ViewCount,
		CategoryID: article.CategoryID,
		Category: dto.ArticleCategoryResponse{
			ID:   article.Category.ID,
			Name: article.Category.Name,
		},
		Tags: toTagResponses(article.Tags),
		Author: dto.ArticleAuthorResponse{
			ID:       article.User.ID,
			Username: article.User.Username,
			Nickname: article.User.Nickname,
		},
		CommentCount: detail.CommentCount,
		Comments:     toCommentResponses(detail.Comments),
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
	}
}
