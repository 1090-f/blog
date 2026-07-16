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

// UserController 提供管理端用户管理：列表查询、状态更新。
type UserController struct {
	userService *service.UserService
}

// NewUserController 创建并初始化用户管理接口实例。
func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

// List 查询用户管理接口列表。
func (ctl *UserController) List(c *gin.Context) {
	var query dto.AdminUserListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	users, total, page, pageSize, err := ctl.userService.List(query)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5022, "获取用户列表失败")
		return
	}

	response.Success(c, dto.AdminUserListResponse{
		List:     toAdminUserResponses(users),
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// CreateAdmin 由当前管理员创建新的后台管理员账号。
func (ctl *UserController) CreateAdmin(c *gin.Context) {
	var req dto.CreateAdminUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}
	user, err := ctl.userService.CreateAdmin(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserExists):
			response.Error(c, http.StatusBadRequest, 4002, "用户名已存在")
		case errors.Is(err, service.ErrInvalidRegister):
			response.Error(c, http.StatusBadRequest, 4001, "注册信息不符合要求")
		default:
			response.Error(c, http.StatusInternalServerError, 5024, "创建管理员失败")
		}
		return
	}
	response.Success(c, toAdminUserResponse(*user))
}

// UpdateStatus 更新用户管理接口记录。
func (ctl *UserController) UpdateStatus(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	user, err := ctl.userService.UpdateStatus(c.GetUint("userID"), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, 4041, err.Error())
		case errors.Is(err, service.ErrInvalidUserStatus):
			response.Error(c, http.StatusBadRequest, 4007, err.Error())
		case errors.Is(err, service.ErrCannotModifySelf):
			response.Error(c, http.StatusForbidden, 4032, "不能禁用当前登录的管理员账号")
		case errors.Is(err, service.ErrCannotModifyAdmin):
			response.Error(c, http.StatusForbidden, 4033, "管理员之间不能互相修改状态")
		default:
			response.Error(c, http.StatusInternalServerError, 5023, "更新用户状态失败")
		}
		return
	}

	response.Success(c, toAdminUserResponse(*user))
}

// UpdateRole 更新用户角色。
func (ctl *UserController) UpdateRole(c *gin.Context) {
	id, ok := parseUintParam(c, "id")
	if !ok {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}
	var req dto.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}
	user, err := ctl.userService.UpdateRole(c.GetUint("userID"), id, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			response.Error(c, http.StatusNotFound, 4041, err.Error())
		case errors.Is(err, service.ErrInvalidUserRole):
			response.Error(c, http.StatusBadRequest, 4008, err.Error())
		case errors.Is(err, service.ErrCannotModifySelf):
			response.Error(c, http.StatusForbidden, 4032, "不能降低当前登录管理员的权限")
		case errors.Is(err, service.ErrCannotModifyAdmin):
			response.Error(c, http.StatusForbidden, 4033, "管理员之间不能互相修改角色")
		default:
			response.Error(c, http.StatusInternalServerError, 5025, "更新用户角色失败")
		}
		return
	}
	response.Success(c, toAdminUserResponse(*user))
}

// 将内部数据转换为接口响应结构。
func toAdminUserResponses(users []model.User) []dto.AdminUserResponse {
	resp := make([]dto.AdminUserResponse, 0, len(users))
	for _, user := range users {
		resp = append(resp, toAdminUserResponse(user))
	}
	return resp
}

// 将内部数据转换为接口响应结构。
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
