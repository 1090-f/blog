package dto

import "time"

// AdminUserListQuery 管理端用户分页查询参数。
type AdminUserListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Role     string `form:"role"`
	Status   *int8  `form:"status"`
}

// UpdateUserStatusRequest 更新用户状态（启用/禁用）请求体。
type UpdateUserStatusRequest struct {
	Status int8 `json:"status" binding:"oneof=0 1"`
}

// CreateAdminUserRequest 由管理员创建后台管理员账号的请求体。
type CreateAdminUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
}

// UpdateUserRoleRequest 更新用户角色请求体。
type UpdateUserRoleRequest struct {
	Role string `json:"role" binding:"required,oneof=admin user"`
}

// AdminUserResponse 管理端用户数据响应。
type AdminUserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Role      string    `json:"role"`
	Avatar    string    `json:"avatar"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// AdminUserListResponse 管理端分页用户列表响应。
type AdminUserListResponse struct {
	List     []AdminUserResponse `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}
