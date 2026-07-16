package middleware

import (
	"net/http"

	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// Admin 返回管理端访问控制中间件，要求当前用户角色为 admin，否则返回 403。
// 必须在 Auth 中间件之后使用，因为需要从 gin.Context 中读取 Auth 中间件写入的 role。
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		// 角色不是 admin 时拒绝访问
		if role != "admin" {
			response.Error(c, http.StatusForbidden, 4030, "admin access required")
			// 中止后续中间件和handler的执行
			c.Abort()
			return
		}
		c.Next()
	}
}
