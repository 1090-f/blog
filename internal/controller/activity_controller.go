package controller

import (
	"net/http"
	"strconv"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

// ActivityController 提供访问活动日历数据（文章/评论发布统计）。
type ActivityController struct {
	service *service.ActivityService
}

// NewActivityController 创建并初始化访问活动数据实例。
func NewActivityController(activityService *service.ActivityService) *ActivityController {
	return &ActivityController{service: activityService}
}

// Get 获取访问活动数据数据。
func (ctl *ActivityController) Get(c *gin.Context) {
	year, err := parseOptionalInt(c.Query("year"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "年份参数无效")
		return
	}
	month, err := parseOptionalInt(c.Query("month"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "月份参数无效")
		return
	}

	activity, err := ctl.service.Get(year, month)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, err.Error())
		return
	}
	response.Success(c, activity)
}

// 解析并校验请求参数。
func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}
