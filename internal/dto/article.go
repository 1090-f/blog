package dto

import "time"

type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required,max=150"`
	Summary    string `json:"summary" binding:"max=255"`
	Content    string `json:"content" binding:"required"`
	CoverImage string `json:"coverImage" binding:"max=255"`
	Status     string `json:"status" binding:"required"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	TagIDs     []uint `json:"tagIds"`
}

type UpdateArticleRequest struct {
	Title      string `json:"title" binding:"required,max=150"`
	Summary    string `json:"summary" binding:"max=255"`
	Content    string `json:"content" binding:"required"`
	CoverImage string `json:"coverImage" binding:"max=255"`
	Status     string `json:"status" binding:"required"`
	CategoryID uint   `json:"categoryId" binding:"required"`
	TagIDs     []uint `json:"tagIds"`
}

type ArticleListQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	CategoryID uint   `form:"categoryId"`
	TagID      uint   `form:"tagId"`
	Keyword    string `form:"keyword"`
}

type AdminArticleListQuery struct {
	Page       int    `form:"page"`
	PageSize   int    `form:"pageSize"`
	CategoryID uint   `form:"categoryId"`
	Keyword    string `form:"keyword"`
	Status     string `form:"status"`
	TagID      uint   `form:"tagId"`
}

type ArticleAuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type ArticleCategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ArticleSummaryResponse struct {
	ID         uint                    `json:"id"`
	Title      string                  `json:"title"`
	Summary    string                  `json:"summary"`
	CoverImage string                  `json:"coverImage"`
	Status     string                  `json:"status"`
	ViewCount  int                     `json:"viewCount"`
	CategoryID uint                    `json:"categoryId"`
	Category   ArticleCategoryResponse `json:"category"`
	Tags       []TagResponse           `json:"tags"`
	Author     ArticleAuthorResponse   `json:"author"`
	CreatedAt  time.Time               `json:"createdAt"`
	UpdatedAt  time.Time               `json:"updatedAt"`
}

type ArticleDetailResponse struct {
	ID         uint                    `json:"id"`
	Title      string                  `json:"title"`
	Summary    string                  `json:"summary"`
	Content    string                  `json:"content"`
	CoverImage string                  `json:"coverImage"`
	Status     string                  `json:"status"`
	ViewCount  int                     `json:"viewCount"`
	CategoryID uint                    `json:"categoryId"`
	Category   ArticleCategoryResponse `json:"category"`
	Tags       []TagResponse           `json:"tags"`
	Author     ArticleAuthorResponse   `json:"author"`
	CreatedAt  time.Time               `json:"createdAt"`
	UpdatedAt  time.Time               `json:"updatedAt"`
}

type ArticleFullDetailResponse struct {
	ID           uint                    `json:"id"`
	Title        string                  `json:"title"`
	Summary      string                  `json:"summary"`
	Content      string                  `json:"content"`
	CoverImage   string                  `json:"coverImage"`
	Status       string                  `json:"status"`
	ViewCount    int                     `json:"viewCount"`
	CategoryID   uint                    `json:"categoryId"`
	Category     ArticleCategoryResponse `json:"category"`
	Tags         []TagResponse           `json:"tags"`
	Author       ArticleAuthorResponse   `json:"author"`
	CommentCount int                     `json:"commentCount"`
	Comments     []CommentResponse       `json:"comments"`
	CreatedAt    time.Time               `json:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt"`
}

type ArticleListResponse struct {
	List     []ArticleSummaryResponse `json:"list"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	PageSize int                      `json:"pageSize"`
}
