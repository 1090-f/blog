package controller

import (
	"net/http"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type SiteStatsController struct {
	service *service.SiteStatsService
}

func NewSiteStatsController(service *service.SiteStatsService) *SiteStatsController {
	return &SiteStatsController{service: service}
}

func (ctl *SiteStatsController) Get(c *gin.Context) {
	stats, err := ctl.service.Get()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 5022, "failed to fetch site stats")
		return
	}
	response.Success(c, stats)
}
