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

type AdminCommentController struct {
	commentService *service.AdminCommentService
}

func NewAdminCommentController(commentService *service.AdminCommentService) *AdminCommentController {
	return &AdminCommentController{commentService: commentService}
}

func (ctl *AdminCommentController) List(c *gin.Context) {
	var query dto.AdminCommentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	comments, total, page, pageSize, err := ctl.commentService.List(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5025, "failed to fetch comments")
		return
	}

	response.Success(c, dto.AdminCommentListResponse{
		List:     toAdminCommentResponses(comments),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (ctl *AdminCommentController) UpdateStatus(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	comment, err := ctl.commentService.UpdateStatus(id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrCommentNotFound):
			response.Error(c, http.StatusNotFound, 4044, err.Error())
		case errors.Is(err, service.ErrInvalidCommentStatus):
			response.Error(c, http.StatusBadRequest, 4005, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5026, "failed to update comment status")
		}
		return
	}

	response.Success(c, toAdminCommentResponse(*comment))
}

func (ctl *AdminCommentController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	if err := ctl.commentService.Delete(id); err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			response.Error(c, http.StatusNotFound, 4044, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, 5027, "failed to delete comment")
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

func toAdminCommentResponses(comments []model.Comment) []dto.AdminCommentResponse {
	resp := make([]dto.AdminCommentResponse, 0, len(comments))
	for _, comment := range comments {
		resp = append(resp, toAdminCommentResponse(comment))
	}
	return resp
}

func toAdminCommentResponse(comment model.Comment) dto.AdminCommentResponse {
	return dto.AdminCommentResponse{
		ID:           comment.ID,
		ArticleID:    comment.ArticleID,
		ArticleTitle: comment.Article.Title,
		UserID:       comment.UserID,
		Username:     comment.User.Username,
		Nickname:     comment.User.Nickname,
		ParentID:     comment.ParentID,
		ReplyToID:    comment.ReplyToID,
		Content:      comment.Content,
		Status:       comment.Status,
		CreatedAt:    comment.CreatedAt,
		UpdatedAt:    comment.UpdatedAt,
	}
}
