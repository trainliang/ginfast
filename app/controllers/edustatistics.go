package controllers

import (
	"gin-fast/app/models"
	"gin-fast/app/service"

	"github.com/gin-gonic/gin"
)

// EduStatisticsController 教培统计控制器
type EduStatisticsController struct {
	Common
	EduStatisticsService *service.EduStatisticsService
}

func NewEduStatisticsController() *EduStatisticsController {
	return &EduStatisticsController{Common: Common{}, EduStatisticsService: service.NewEduStatisticsService()}
}

func (ctl *EduStatisticsController) Schedule(c *gin.Context) {
	var req models.EduStatisticsRangeRequest
	if err := c.ShouldBind(&req); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	resp, err := ctl.EduStatisticsService.ScheduleStatistics(c, &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, resp)
}

func (ctl *EduStatisticsController) Benefits(c *gin.Context) {
	var req models.EduStatisticsRangeRequest
	if err := c.ShouldBind(&req); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	resp, err := ctl.EduStatisticsService.BenefitStatistics(c, &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, resp)
}

func (ctl *EduStatisticsController) ExternalSync(c *gin.Context) {
	var req models.EduStatisticsRangeRequest
	if err := c.ShouldBind(&req); err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	resp, err := ctl.EduStatisticsService.ExternalSyncStatistics(c, &req)
	if err != nil {
		ctl.FailAndAbort(c, err.Error(), err)
	}
	ctl.Success(c, resp)
}
