package dto

// CreateTagRequest 创建标签请求体。
type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// UpdateTagRequest 更新标签请求体。
type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

// TagResponse 标签数据响应。
type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
