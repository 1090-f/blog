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

// CommentController 提供公开评论操作：按文章查询、创建评论、删除自己的评论。
type CommentController struct {
	commentService *service.CommentService
}

// NewCommentController 创建并初始化评论接口实例。
func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

// ListByArticle 查询评论接口列表。
func (ctl *CommentController) ListByArticle(c *gin.Context) {
	articleID, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	comments, err := ctl.commentService.ListPublicByArticleID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5015, "获取评论列表失败")
		}
		return
	}

	response.Success(c, toCommentResponses(comments))
}

// Create 创建评论接口记录。
func (ctl *CommentController) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	comment, err := ctl.commentService.Create(req, c.GetUint("userID"))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		case errors.Is(err, service.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, 4044, err.Error())
		case errors.Is(err, service.ErrInvalidComment),
			errors.Is(err, service.ErrGuestNameRequired),
			errors.Is(err, service.ErrGuestEmailRequired),
			errors.Is(err, service.ErrInvalidGuestEmail),
			errors.Is(err, service.ErrInvalidGuestWebsite):
			response.Error(c, http.StatusBadRequest, 4005, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5016, "创建评论失败")
		}
		return
	}

	response.Success(c, toCommentResponse(*comment))
}

// DeleteMine 删除评论接口记录。
func (ctl *CommentController) DeleteMine(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	if err := ctl.commentService.DeleteByUser(id, c.GetUint("userID")); err != nil {
		switch {
		case errors.Is(err, service.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, 4044, err.Error())
		case errors.Is(err, service.ErrCommentForbidden):
			response.Error(c, http.StatusForbidden, 4030, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5017, "删除评论失败")
		}
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

// 将内部数据转换为接口响应结构。
func toCommentResponses(comments []model.Comment) []dto.CommentResponse {
	resp := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		resp = append(resp, toCommentResponse(comment))
	}
	return resp
}

// 将内部数据转换为接口响应结构。
func toCommentResponse(comment model.Comment) dto.CommentResponse {
	resp := dto.CommentResponse{
		ID:        comment.ID,
		ArticleID: comment.ArticleID,
		UserID:    comment.UserID,
		ParentID:  comment.ParentID,
		ReplyToID: comment.ReplyToID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt,
		Author:    commentAuthor(comment),
	}

	if comment.ReplyTo != nil {
		author := commentAuthor(*comment.ReplyTo)
		resp.ReplyToAuthor = &author
	}

	return resp
}

// commentAuthor 根据登录用户和游客返回对应的评论作者信息
func commentAuthor(comment model.Comment) dto.CommentAuthorResponse {
	if comment.User != nil {
		return dto.CommentAuthorResponse{
			ID:       comment.User.ID,
			Username: comment.User.Username,
			Nickname: comment.User.Nickname,
			Avatar:   comment.User.Avatar,
		}
	}

	return dto.CommentAuthorResponse{
		Nickname: comment.GuestName,
		Website:  comment.GuestWebsite,
	}
}
