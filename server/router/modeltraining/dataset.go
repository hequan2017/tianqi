package modeltraining

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DatasetRouter struct{}

// InitDatasetRouter 初始化数据集路由
func (r *DatasetRouter) InitDatasetRouter(Router *gin.RouterGroup, PublicRouter *gin.RouterGroup) {
	datasetRouter := Router.Group("modeltraining/dataset").Use(middleware.OperationRecord())
	datasetRouterWithoutRecord := Router.Group("modeltraining/dataset")
	datasetRouterWithoutAuth := PublicRouter.Group("modeltraining/dataset")
	{
		datasetRouter.POST("createDataset", datasetApi.CreateDataset)             // 创建数据集
		datasetRouter.DELETE("deleteDataset", datasetApi.DeleteDataset)           // 删除数据集
		datasetRouter.DELETE("deleteDatasetByIds", datasetApi.DeleteDatasetByIds) // 批量删除数据集
		datasetRouter.PUT("updateDataset", datasetApi.UpdateDataset)              // 更新数据集
		datasetRouter.POST("createVersion", datasetApi.CreateVersion)             // 创建版本
		datasetRouter.DELETE("deleteVersion", datasetApi.DeleteVersion)           // 删除版本
		datasetRouter.POST("publishDataset", datasetApi.PublishDataset)           // 发布数据集
		datasetRouter.POST("uploadFile", datasetApi.UploadDatasetFile)            // 上传数据集文件
		datasetRouter.POST("uploadVersionFile", datasetApi.UploadVersionFile)     // 上传版本文件
	}
	{
		datasetRouterWithoutRecord.GET("findDataset", datasetApi.FindDataset)           // 查询数据集详情
		datasetRouterWithoutRecord.GET("getDatasetList", datasetApi.GetDatasetList)     // 获取数据集列表
		datasetRouterWithoutRecord.GET("getVersionList", datasetApi.GetVersionList)     // 获取版本列表
	}
	{
		datasetRouterWithoutAuth.GET("getDatasetDataSource", datasetApi.GetDatasetDataSource) // 获取数据源
	}
}