package dto

import "time"

type AdminCommentListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
	Keyword   string `form:"keyword"`
	ArticleID uint   `form:"articleId"`
	Status    *int8  `form:"status"`
}

type UpdateCommentStatusRequest struct {
	Status int8 `json:"status" binding:"oneof=0 1"`
}

type CreateCommentRequest struct {
	ArticleID uint   `json:"articleId" binding:"required"`
	ReplyToID *uint  `json:"replyToId"`
	Content   string `json:"content" binding:"required"`
}

type CommentAuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type CommentResponse struct {
	ID            uint                   `json:"id"`
	ArticleID     uint                   `json:"articleId"`
	UserID        uint                   `json:"userId"`
	ParentID      *uint                  `json:"parentId"`
	ReplyToID     *uint                  `json:"replyToId"`
	Content       string                 `json:"content"`
	CreatedAt     time.Time              `json:"createdAt"`
	Author        CommentAuthorResponse  `json:"author"`
	ReplyToAuthor *CommentAuthorResponse `json:"replyToAuthor,omitempty"`
}

type AdminCommentResponse struct {
	ID           uint      `json:"id"`
	ArticleID    uint      `json:"articleId"`
	ArticleTitle string    `json:"articleTitle"`
	UserID       uint      `json:"userId"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	ParentID     *uint     `json:"parentId"`
	ReplyToID    *uint     `json:"replyToId"`
	Content      string    `json:"content"`
	Status       int8      `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AdminCommentListResponse struct {
	List     []AdminCommentResponse `json:"list"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
}
