package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		log.Printf("panic recovered: method=%s path=%s error=%v\n%s", c.Request.Method, c.Request.URL.Path, recovered, debug.Stack())
		response.Error(c, http.StatusInternalServerError, 5000, "internal server error")
		c.Abort()
	})
}
