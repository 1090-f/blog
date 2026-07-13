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

type ArticleController struct {
	articleService *service.ArticleService
}

func NewArticleController(articleService *service.ArticleService) *ArticleController {
	return &ArticleController{articleService: articleService}
}

func (ctl *ArticleController) List(c *gin.Context) {
	var query dto.ArticleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	articles, total, page, pageSize, err := ctl.articleService.ListPublished(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5009, "failed to fetch articles")
		return
	}

	response.Success(c, dto.ArticleListResponse{
		List:     toArticleSummaryResponses(articles),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (ctl *ArticleController) Detail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	article, err := ctl.articleService.GetPublishedDetail(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5010, "failed to fetch article")
		}
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

func (ctl *ArticleController) FullDetail(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	detail, err := ctl.articleService.GetPublishedFullDetail(id)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5019, "failed to fetch article detail")
		}
		return
	}

	response.Success(c, toArticleFullDetailResponse(detail))
}

func (ctl *ArticleController) Latest(c *gin.Context) {
	articles, err := ctl.articleService.ListLatest(parseLimitQuery(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5020, "failed to fetch latest articles")
		return
	}

	response.Success(c, toArticleSummaryResponses(articles))
}

func (ctl *ArticleController) Popular(c *gin.Context) {
	articles, err := ctl.articleService.ListPopular(parseLimitQuery(c))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5021, "failed to fetch popular articles")
		return
	}

	response.Success(c, toArticleSummaryResponses(articles))
}

func (ctl *ArticleController) AdminList(c *gin.Context) {
	var query dto.AdminArticleListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	articles, total, page, pageSize, err := ctl.articleService.ListAdmin(query)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidArticleStatus):
			response.Error(c, http.StatusBadRequest, 4004, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5011, "failed to fetch admin articles")
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

func (ctl *ArticleController) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	article, err := ctl.articleService.Create(req, c.GetUint("userID"))
	if err != nil {
		respondArticleWriteError(c, err, "failed to create article", 5012)
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

func (ctl *ArticleController) Update(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	article, err := ctl.articleService.Update(id, req, c.GetUint("userID"))
	if err != nil {
		respondArticleWriteError(c, err, "failed to update article", 5013)
		return
	}

	response.Success(c, toArticleDetailResponse(*article))
}

func (ctl *ArticleController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	if err := ctl.articleService.Delete(id); err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		case errors.Is(err, service.ErrArticleInUse):
			response.Error(c, http.StatusBadRequest, 4004, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5014, "failed to delete article")
		}
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

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

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	value, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || value == 0 {
		return 0, false
	}
	return uint(value), true
}

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
