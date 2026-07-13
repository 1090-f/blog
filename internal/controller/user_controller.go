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

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (ctl *UserController) List(c *gin.Context) {
	var query dto.AdminUserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	users, total, page, pageSize, err := ctl.userService.List(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5022, "failed to fetch users")
		return
	}

	response.Success(c, dto.AdminUserListResponse{
		List:     toAdminUserResponses(users),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func (ctl *UserController) UpdateStatus(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	user, err := ctl.userService.UpdateStatus(id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, 4041, err.Error())
		case errors.Is(err, service.ErrInvalidUserStatus):
			response.Error(c, http.StatusBadRequest, 4007, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5023, "failed to update user status")
		}
		return
	}

	response.Success(c, toAdminUserResponse(*user))
}

func (ctl *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("userID")
	if userID == 0 {
		response.Error(c, http.StatusUnauthorized, 4010, "unauthorized")
		return
	}

	var req dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	user, err := ctl.userService.UpdateProfile(userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, 4041, err.Error())
		case errors.Is(err, service.ErrInvalidProfile):
			response.Error(c, http.StatusBadRequest, 4007, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5024, "failed to update profile")
		}
		return
	}

	response.Success(c, user)
}

func toAdminUserResponses(users []model.User) []dto.AdminUserResponse {
	resp := make([]dto.AdminUserResponse, 0, len(users))
	for _, user := range users {
		resp = append(resp, toAdminUserResponse(user))
	}
	return resp
}

func toAdminUserResponse(user model.User) dto.AdminUserResponse {
	return dto.AdminUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Role:      user.Role,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
