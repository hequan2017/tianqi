package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

// DatasetSearch 数据集搜索请求
type DatasetSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Name           *string     `json:"name" form:"name"`
	Type           *string     `json:"type" form:"type"`
	TrainMethod    *string     `json:"trainMethod" form:"trainMethod"`
	ImportStatus   *string     `json:"importStatus" form:"importStatus"`
	PublishStatus  *bool       `json:"publishStatus" form:"publishStatus"`
	UserId         *uint       `json:"userId" form:"userId"`
	request.PageInfo
}

// DatasetVersionSearch 数据集版本搜索请求
type DatasetVersionSearch struct {
	DatasetId uint `json:"datasetId" form:"datasetId" binding:"required"`
	request.PageInfo
}

// CreateDatasetReq 创建数据集请求
type CreateDatasetReq struct {
	Name        string `json:"name" form:"name" binding:"required"`
	Type        string `json:"type" form:"type" binding:"required"` // training/evaluation
	Format      string `json:"format" form:"format"`
	TrainMethod string `json:"trainMethod" form:"trainMethod"`
	Description string `json:"description" form:"description"`
}

// CreateVersionReq 创建版本请求
type CreateVersionReq struct {
	DatasetId   uint   `json:"datasetId" form:"datasetId" binding:"required"`
	Version     string `json:"version" form:"version" binding:"required"`
	DataCount   int64  `json:"dataCount" form:"dataCount"`
	StoragePath string `json:"storagePath" form:"storagePath"`
	Description string `json:"description" form:"description"`
	FileSize    int64  `json:"fileSize" form:"fileSize"`
}

// PublishDatasetReq 发布数据集请求
type PublishDatasetReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}