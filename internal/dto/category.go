package dto

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name" binding:"required,max=50"`
	Description string `json:"description" binding:"max=255"`
}

type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
