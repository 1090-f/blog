package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回记录请求日志的中间件，在请求完成后输出方法、路径、状态码和耗时。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		// 执行后续中间件和 handler
		c.Next()

		// 获取路由注册时的路径模板（如 /api/articles/:id），404 时为空则降级为实际请求路径
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 输出日志：方法、路径、HTTP 状态码、处理耗时
		log.Printf("method=%s path=%s status=%d latency=%s", c.Request.Method, path, c.Writer.Status(), time.Since(start))
	}
}
