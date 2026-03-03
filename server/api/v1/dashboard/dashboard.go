package dashboard

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DashboardApi struct{}

var dashboardService = service.ServiceGroupApp.DashboardServiceGroup.DashboardService

// GetDashboardStats 获取仪表盘统计数据
// @Tags Dashboard
// @Summary 获取仪表盘统计数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=model.DashboardStats,msg=string} "获取成功"
// @Router /dashboard/stats [get]
func (d *DashboardApi) GetDashboardStats(c *gin.Context) {
	stats, err := dashboardService.GetDashboardStats()
	if err != nil {
		global.GVA_LOG.Error("获取仪表盘统计数据失败!", zap.Error(err))
		response.FailWithMessage("获取仪表盘统计数据失败", c)
		return
	}
	response.OkWithDetailed(stats, "获取成功", c)
}