package controller

import (
	"net/http"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// SiteStatsController 提供站点公开统计数据（文章数、字数等）。
type SiteStatsController struct {
	service *service.SiteStatsService
}

// NewSiteStatsController 创建并初始化站点统计接口实例。
func NewSiteStatsController(service *service.SiteStatsService) *SiteStatsController {
	return &SiteStatsController{service: service}
}

// Get 获取站点统计接口数据。
func (ctl *SiteStatsController) Get(c *gin.Context) {
	stats, err := ctl.service.Get()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5022, "获取站点统计失败")
		return
	}
	response.Success(c, stats)
}
