package controller

import (
	"errors"
	"net/http"

	"blog/internal/dto"
	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	user, err := ctl.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) || errors.Is(err, service.ErrInvalidRegister) {
			response.Error(c, http.StatusBadRequest, 4002, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, 5001, "failed to register user")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"role":     user.Role,
	})
}

func (ctl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid request params")
		return
	}

	loginResponse, err := ctl.authService.Login(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			response.Error(c, http.StatusUnauthorized, 4011, err.Error())
		case errors.Is(err, service.ErrUserDisabled):
			response.Error(c, http.StatusForbidden, 4031, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5002, "failed to login")
		}
		return
	}

	response.Success(c, loginResponse)
}

func (ctl *AuthController) Profile(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := ctl.authService.Profile(userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			response.Error(c, http.StatusNotFound, 4041, "user not found")
		case errors.Is(err, service.ErrUserDisabled):
			response.Error(c, http.StatusForbidden, 4031, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5003, "failed to fetch profile")
		}
		return
	}

	response.Success(c, user)
}
