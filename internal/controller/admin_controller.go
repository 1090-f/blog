package controller

import (
	"net/http"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	adminService *service.AdminService
}

func NewAdminController(adminService *service.AdminService) *AdminController {
	return &AdminController{adminService: adminService}
}

func (ctl *AdminController) Dashboard(c *gin.Context) {
	stats, err := ctl.adminService.Dashboard()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5019, "failed to fetch dashboard stats")
		return
	}

	response.Success(c, stats)
}
