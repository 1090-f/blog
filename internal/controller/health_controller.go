package controller

import (
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// HealthController 提供服务健康检查端点。
type HealthController struct{}

// NewHealthController 创建并初始化健康检查接口实例。
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Check 返回服务健康状态。
func (ctl *HealthController) Check(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}
