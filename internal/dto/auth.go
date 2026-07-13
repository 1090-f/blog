package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Nickname string `json:"nickname" binding:"required,min=2,max=50"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string            `json:"token"`
	User  LoginUserResponse `json:"user"`
}

type LoginUserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
	Status   int8   `json:"status"`
}
