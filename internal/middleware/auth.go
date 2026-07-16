package middleware

import (
	"errors"
	"net/http"
	"strings"

	"blog/internal/model"
	jwtpkg "blog/pkg/jwt"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserReader 中间件所需的用户查询抽象，用于 JWT 鉴权时根据 ID 查用户。
type UserReader interface {
	FindByID(id uint) (*model.User, error)
}

// Auth 返回校验登录令牌的认证中间件。
func Auth(secret string, userDAO UserReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 4010, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		// 必须由空格分成两部分；第一部分必须是 Bearer(不区分大小写); 第二部分（令牌）不能为空
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, http.StatusUnauthorized, 4010, "invalid authorization header")
			c.Abort()
			return
		}

		// 验证签名是否有效或过期
		tokenString := strings.TrimSpace(parts[1])
		claims, err := jwtpkg.ParseToken(tokenString, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 4010, "invalid or expired token")
			c.Abort()
			return
		}

		// 查询数据库确认用户是否存在；
		user, err := userDAO.FindByID(claims.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				response.Error(c, http.StatusUnauthorized, 4010, "invalid or expired token")
			} else {
				response.Error(c, http.StatusInternalServerError, 5004, "failed to verify current user")
			}
			c.Abort()
			return
		}
		if user.Status != 1 {
			response.Error(c, http.StatusForbidden, 4031, "user is disabled")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", user.Role)
		c.Next()
	}
}

// OptionalAuth 在请求携带 Bearer 令牌时识别已登录用户，主要用于评论
// 未登录请求仍可继续访问。
func OptionalAuth(secret string, userDAO UserReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.Next()
			return
		}
		Auth(secret, userDAO)(c)
	}
}
