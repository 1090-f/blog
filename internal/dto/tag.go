package dto

type CreateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required,max=50"`
}

type TagResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
