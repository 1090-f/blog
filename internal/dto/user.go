package dto

import "time"

type AdminUserListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"pageSize"`
	Keyword  string `form:"keyword"`
	Role     string `form:"role"`
	Status   *int8  `form:"status"`
}

type UpdateUserStatusRequest struct {
	Status int8 `json:"status" binding:"oneof=0 1"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
	Avatar   string `json:"avatar" binding:"max=255"`
}

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

type AdminUserListResponse struct {
	List     []AdminUserResponse `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"pageSize"`
}
