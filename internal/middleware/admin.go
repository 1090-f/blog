package middleware

import (
	"net/http"

	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			response.Error(c, http.StatusForbidden, 4030, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
