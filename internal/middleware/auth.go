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

type UserReader interface {
	FindByID(id uint) (*model.User, error)
}

func Auth(secret string, userDAO UserReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 4010, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, http.StatusUnauthorized, 4010, "invalid authorization header")
			c.Abort()
			return
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := jwtpkg.ParseToken(tokenString, secret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 4010, "invalid or expired token")
			c.Abort()
			return
		}

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

// OptionalAuth identifies a logged-in visitor when a bearer token is present,
// while allowing unauthenticated requests to continue.
func OptionalAuth(secret string, userDAO UserReader) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") == "" {
			c.Next()
			return
		}
		Auth(secret, userDAO)(c)
	}
}
