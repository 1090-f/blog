package controller

import (
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (ctl *HealthController) Check(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}
