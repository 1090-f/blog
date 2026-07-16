package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recovery 返回一个用于捕获 panic 的 Gin 中间件。
// 当处理请求的过程中发生 panic 时，该中间件会：
//   - 记录 panic 的详细信息（包括请求方法、路径、错误内容和堆栈跟踪）
//   - 向客户端返回统一格式的 500 错误响应
//   - 阻止 panic 向上冒泡导致服务器进程崩溃
//
// 建议在所有 HTTP 服务入口处使用，防止未处理的异常导致服务不可用。
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Printf("panic recovered: method=%s path=%s error=%v\n%s", c.Request.Method, c.Request.URL.Path, recovered, debug.Stack())
		response.Error(c, http.StatusInternalServerError, 5000, "internal server error")
		c.Abort()
	})
}
