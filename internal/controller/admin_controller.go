package controller

import (
	"net/http"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminController 提供管理端仪表盘聚合统计数据。
type AdminController struct {
	adminService *service.AdminService
}

// NewAdminController 创建并初始化管理端仪表盘实例。
func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

// Dashboard 获取管理端仪表盘统计数据。
func (ctl *AdminController) Dashboard(c *gin.Context) {
	stats, err := ctl.adminService.Dashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5019, "获取仪表盘统计失败")
		return
	}

	response.Success(c, stats)
}
