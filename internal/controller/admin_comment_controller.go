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

// AdminCommentController 提供管理端评论管理：列表、审核状态更新、删除。
type AdminCommentController struct {
	commentService *service.AdminCommentService
}

// NewAdminCommentController 创建并初始化管理端评论实例。
func NewAdminCommentController(commentService *service.AdminCommentService) *AdminCommentController {
	return &AdminCommentController{commentService: commentService}
}

// List 查询管理端评论列表。
func (ctl *AdminCommentController) List(c *gin.Context) {
	var query dto.AdminCommentListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	comments, total, page, pageSize, err := ctl.commentService.List(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5025, "获取评论列表失败")
		return
	}

	response.Success(c, dto.AdminCommentListResponse{
		List:     toAdminCommentResponses(comments),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// UpdateStatus 更新管理端评论记录。
func (ctl *AdminCommentController) UpdateStatus(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	var req dto.UpdateCommentStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
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
			response.Error(c, http.StatusInternalServerError, 5026, "更新评论状态失败")
		}
		return
	}

	response.Success(c, toAdminCommentResponse(*comment))
}

// Delete 删除管理端评论记录。
func (ctl *AdminCommentController) Delete(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	if err := ctl.commentService.Delete(id); err != nil {
		if errors.Is(err, service.ErrCommentNotFound) {
			response.Error(c, http.StatusNotFound, 4044, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, 5027, "删除评论失败")
		return
	}

	response.Success(c, gin.H{"deleted": true})
}

// 将内部数据转换为接口响应结构。
func toAdminCommentResponses(comments []model.Comment) []dto.AdminCommentResponse {
	resp := make([]dto.AdminCommentResponse, 0, len(comments))
	for _, comment := range comments {
		resp = append(resp, toAdminCommentResponse(comment))
	}
	return resp
}

// 将内部数据转换为接口响应结构。
func toAdminCommentResponse(comment model.Comment) dto.AdminCommentResponse {
	username := ""
	nickname := comment.GuestName
	if comment.User != nil {
		username = comment.User.Username
		nickname = comment.User.Nickname
	}

	return dto.AdminCommentResponse{
		ID:           comment.ID,
		ArticleID:    comment.ArticleID,
		ArticleTitle: comment.Article.Title,
		UserID:       comment.UserID,
		Username:     username,
		Nickname:     nickname,
		ParentID:     comment.ParentID,
		ReplyToID:    comment.ReplyToID,
		Content:      comment.Content,
		Status:       comment.Status,
		CreatedAt:    comment.CreatedAt,
		UpdatedAt:    comment.UpdatedAt,
	}
}
