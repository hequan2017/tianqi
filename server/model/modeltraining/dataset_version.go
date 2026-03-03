// 数据集版本管理
package modeltraining

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// DatasetVersion 数据集版本
type DatasetVersion struct {
	global.GVA_MODEL
	DatasetId   uint   `json:"datasetId" form:"datasetId" gorm:"comment:数据集ID;not null;index"`
	Version     string `json:"version" form:"version" gorm:"comment:版本号;size:20;not null"` // V1, V2, V3...
	DataCount   int64  `json:"dataCount" form:"dataCount" gorm:"comment:数据量;default:0"`
	StoragePath string `json:"storagePath" form:"storagePath" gorm:"comment:存储路径;size:255"`
	Description string `json:"description" form:"description" gorm:"comment:版本说明;size:500"`
	FileSize    int64  `json:"fileSize" form:"fileSize" gorm:"comment:文件大小(字节)"`
	Status      string `json:"status" form:"status" gorm:"comment:状态;size:20;default:'pending'"` // uploading/success/failed
	FileName    string `json:"fileName" form:"fileName" gorm:"comment:文件名称;size:255"`           // 上传的文件名
	FilePath    string `json:"filePath" form:"filePath" gorm:"comment:文件路径;size:500"`             // 实际文件存储路径
}

// TableName 数据集版本表名
func (DatasetVersion) TableName() string {
	return "dataset_version"
}