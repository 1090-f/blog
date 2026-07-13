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

type CommentController struct {
	commentService *service.CommentService
}

func NewCommentController(commentService *service.CommentService) *CommentController {
	return &CommentController{commentService: commentService}
}

func (ctl *CommentController) ListByArticle(c *gin.Context) {
	articleID, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	comments, err := ctl.commentService.ListPublicByArticleID(articleID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrArticleNotFound):
			response.Error(c, http.StatusNotFound, 4043, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5015, "failed to fetch comments")
		}
		return
	}

	response.Success(c, toCommentResponses(comments))
}

func (ctl *CommentController) Create(c *gin.Context) {
	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
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
			response.Error(c, http.StatusInternalServerError, 5016, "failed to create comment")
		}
		return
	}

	response.Success(c, toCommentResponse(*comment))
}

func (ctl *CommentController) DeleteMine(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	if err := ctl.commentService.DeleteByUser(id, c.GetUint("userID")); err != nil {
		switch {
		case errors.Is(err, service.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, 4044, err.Error())
		case errors.Is(err, service.ErrCommentForbidden):
			response.Error(c, http.StatusForbidden, 4030, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5017, "failed to delete comment")
		}
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func toCommentResponses(comments []model.Comment) []dto.CommentResponse {
	resp := make([]dto.CommentResponse, 0, len(comments))
	for _, comment := range comments {
		resp = append(resp, toCommentResponse(comment))
	}
	return resp
}

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
