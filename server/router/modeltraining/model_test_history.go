package modeltraining

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ModelTestHistoryRouter struct{}

// InitModelTestHistoryRouter 初始化模型测试历史路由
func (r *ModelTestHistoryRouter) InitModelTestHistoryRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	testRouter := Router.Group("modeltraining/modelTest").Use(middleware.OperationRecord())
	testRouterWithoutRecord := Router.Group("modeltraining/modelTest")
	_ = PublicRouter.Group("modeltraining/modelTest")
	{
		testRouter.POST("createTestHistory", modelTestHistoryApi.CreateTestHistory)       // 创建测试历史
		testRouter.DELETE("deleteTestHistory", modelTestHistoryApi.DeleteTestHistory)     // 删除测试历史
		testRouter.DELETE("clearTestHistory", modelTestHistoryApi.ClearTestHistory)       // 清空测试历史
	}
	{
		testRouterWithoutRecord.GET("getTestHistoryList", modelTestHistoryApi.GetTestHistoryList) // 获取测试历史列表
	}
}