package dto

import "time"

// AdminCommentListQuery 管理端评论分页查询参数。
type AdminCommentListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"pageSize"`
	Keyword   string `form:"keyword"`
	ArticleID uint   `form:"articleId"`
	Status    *int8  `form:"status"`
}

// UpdateCommentStatusRequest 更新评论审核状态请求体。
type UpdateCommentStatusRequest struct {
	Status int8 `json:"status" binding:"oneof=0 1"`
}

// CreateCommentRequest 创建评论请求体，支持注册用户和游客评论。
type CreateCommentRequest struct {
	ArticleID    uint   `json:"articleId" binding:"required"`
	ReplyToID    *uint  `json:"replyToId"`
	Content      string `json:"content" binding:"required"`
	GuestName    string `json:"guestName"`
	GuestEmail   string `json:"guestEmail"`
	GuestWebsite string `json:"guestWebsite"`
}

// CommentAuthorResponse 评论作者信息（注册用户或游客）。
type CommentAuthorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Website  string `json:"website"`
}

// CommentResponse 公开评论数据响应。
type CommentResponse struct {
	ID            uint                   `json:"id"`
	ArticleID     uint                   `json:"articleId"`
	UserID        *uint                  `json:"userId"`
	ParentID      *uint                  `json:"parentId"`
	ReplyToID     *uint                  `json:"replyToId"`
	Content       string                 `json:"content"`
	CreatedAt     time.Time              `json:"createdAt"`
	Author        CommentAuthorResponse  `json:"author"`
	ReplyToAuthor *CommentAuthorResponse `json:"replyToAuthor,omitempty"`
}

// AdminCommentResponse 管理端评论数据，包含审核状态和文章标题。
type AdminCommentResponse struct {
	ID           uint      `json:"id"`
	ArticleID    uint      `json:"articleId"`
	ArticleTitle string    `json:"articleTitle"`
	UserID       *uint     `json:"userId"`
	Username     string    `json:"username"`
	Nickname     string    `json:"nickname"`
	ParentID     *uint     `json:"parentId"`
	ReplyToID    *uint     `json:"replyToId"`
	Content      string    `json:"content"`
	Status       int8      `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// AdminCommentListResponse 管理端分页评论列表响应。
type AdminCommentListResponse struct {
	List     []AdminCommentResponse `json:"list"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"pageSize"`
}
