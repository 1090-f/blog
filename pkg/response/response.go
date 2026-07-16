package response

import "github.com/gin-gonic/gin"

// Success 返回统一的成功响应。
func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

// Error 返回统一的错误响应。
func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}
