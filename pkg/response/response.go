package response

import "github.com/gin-gonic/gin"

func Success(c *gin.Context, data any) {
	c.JSON(200, gin.H{
		"code":    0,
		"message": "success",
		"data":    data,
	})
}

func Error(c *gin.Context, status int, code int, message string) {
	c.JSON(status, gin.H{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}
