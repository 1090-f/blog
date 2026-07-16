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

// AuthController 处理用户注册、登录、管理端登录及会话查询。
type AuthController struct {
	authService *service.AuthService
}

// NewAuthController 创建并初始化认证接口实例。
func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Register 注册用户。
func (ctl *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
		return
	}

	user, err := ctl.authService.Register(req)
	if err != nil {
		if errors.Is(err, service.ErrUserExists) || errors.Is(err, service.ErrInvalidRegister) {
			response.Error(c, http.StatusBadRequest, 4002, err.Error())
			return
		}
		response.Error(c, http.StatusInternalServerError, 5001, "注册用户失败")
		return
	}

	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"role":     user.Role,
	})
}

// Login 验证用户凭据并签发登录令牌。
func (ctl *AuthController) Login(c *gin.Context) {
	ctl.login(c, false)
}

// AdminLogin 验证管理员凭据并签发令牌。
func (ctl *AuthController) AdminLogin(c *gin.Context) {
	ctl.login(c, true)
}

func (ctl *AuthController) login(c *gin.Context, adminOnly bool) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "请求参数无效")
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
			response.Error(c, http.StatusInternalServerError, 5002, "登录失败")
		}
		return
	}
	if adminOnly && loginResponse.User.Role != "admin" {
		response.Error(c, http.StatusForbidden, 4030, "需要管理员权限")
		return
	}

	response.Success(c, loginResponse)
}

// Session 获取当前登录会话的用户信息。
func (ctl *AuthController) Session(c *gin.Context) {
	userID := c.GetUint("userID")

	user, err := ctl.authService.CurrentUser(userID)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			response.Error(c, http.StatusNotFound, 4041, "用户不存在")
		case errors.Is(err, service.ErrUserDisabled):
			response.Error(c, http.StatusForbidden, 4031, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, 5003, "获取当前用户信息失败")
		}
		return
	}

	response.Success(c, user)
}
