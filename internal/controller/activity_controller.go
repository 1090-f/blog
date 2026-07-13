package controller

import (
	"net/http"
	"strconv"

	"blog/internal/service"
	"blog/pkg/response"
	"github.com/gin-gonic/gin"
)

type ActivityController struct {
	service *service.ActivityService
}

func NewActivityController(activityService *service.ActivityService) *ActivityController {
	return &ActivityController{service: activityService}
}

func (ctl *ActivityController) Get(c *gin.Context) {
	year, err := parseOptionalInt(c.Query("year"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid year")
		return
	}
	month, err := parseOptionalInt(c.Query("month"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, "invalid month")
		return
	}

	activity, err := ctl.service.Get(year, month)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 4001, err.Error())
		return
	}
	response.Success(c, activity)
}

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}
