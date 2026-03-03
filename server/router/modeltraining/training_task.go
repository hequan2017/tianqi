package modeltraining

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TrainingTaskRouter struct{}

// InitTrainingTaskRouter 初始化训练任务路由
func (r *TrainingTaskRouter) InitTrainingTaskRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	taskRouter := Router.Group("modeltraining/trainingTask").Use(middleware.OperationRecord())
	taskRouterWithoutRecord := Router.Group("modeltraining/trainingTask")
	taskRouterWithoutAuth := PublicRouter.Group("modeltraining/trainingTask")
	{
		taskRouter.POST("createTask", trainingTaskApi.CreateTask)             // 创建训练任务
		taskRouter.DELETE("deleteTask", trainingTaskApi.DeleteTask)           // 删除训练任务
		taskRouter.DELETE("deleteTaskByIds", trainingTaskApi.DeleteTaskByIds) // 批量删除训练任务
		taskRouter.PUT("updateTask", trainingTaskApi.UpdateTask)              // 更新训练任务
		taskRouter.POST("startTask", trainingTaskApi.StartTask)               // 启动训练
		taskRouter.POST("stopTask", trainingTaskApi.StopTask)                 // 停止训练
		taskRouter.POST("markCompleted", trainingTaskApi.MarkCompleted)       // 手动标记完成
		taskRouter.POST("startService", trainingTaskApi.StartService)         // 启动推理服务
		taskRouter.POST("stopService", trainingTaskApi.StopService)           // 停止推理服务
		taskRouter.POST("chatCompletion", trainingTaskApi.ChatCompletion)     // 模型对话测试
	}
	{
		taskRouterWithoutRecord.GET("findTask", trainingTaskApi.FindTask)                  // 查询任务详情
		taskRouterWithoutRecord.GET("getTaskList", trainingTaskApi.GetTaskList)            // 获取任务列表
		taskRouterWithoutRecord.GET("getTaskLogs", trainingTaskApi.GetTaskLogs)            // 获取训练日志
		taskRouterWithoutRecord.GET("getDefaultParams", trainingTaskApi.GetDefaultParams)  // 获取默认参数
	}
	{
		taskRouterWithoutAuth.GET("getTaskDataSource", trainingTaskApi.GetTaskDataSource) // 获取数据源
	}
}