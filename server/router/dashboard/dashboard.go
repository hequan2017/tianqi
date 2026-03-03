package dashboard

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DashboardRouter struct{}

// InitDashboardRouter 初始化 仪表盘 路由信息
func (s *DashboardRouter) InitDashboardRouter(Router *gin.RouterGroup) {
	dashboardPrivateRouter := Router.Group("dashboard")
	dashboardPrivateRouter.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		dashboardPrivateRouter.GET("stats", v1.ApiGroupApp.DashboardApiGroup.GetDashboardStats) // 获取仪表盘统计数据
	}
}