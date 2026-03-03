// 数据集管理
package modeltraining

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Dataset 数据集模型
type Dataset struct {
	global.GVA_MODEL
	Name          string  `json:"name" form:"name" gorm:"comment:数据集名称;size:100;not null"`
	Type          string  `json:"type" form:"type" gorm:"comment:数据集类型;size:20;not null"` // training/evaluation
	Format        string  `json:"format" form:"format" gorm:"comment:数据格式;size:50"`          // 文本生成、图片理解等
	TrainMethod   string  `json:"trainMethod" form:"trainMethod" gorm:"comment:训练方式;size:20"` // SFT、DPO、CPT
	StoragePath   string  `json:"storagePath" form:"storagePath" gorm:"comment:存储路径;size:255"`
	Description   string  `json:"description" form:"description" gorm:"comment:描述;size:500"`
	UserId        *uint   `json:"userId" form:"userId" gorm:"comment:创建者ID"`
	LatestVersion string  `json:"latestVersion" form:"latestVersion" gorm:"comment:最新版本;size:20"`
	DataCount     int64   `json:"dataCount" form:"dataCount" gorm:"comment:数据量;default:0"`
	ImportStatus  string  `json:"importStatus" form:"importStatus" gorm:"comment:导入状态;size:20;default:'pending'"` // pending/success/failed
	PublishStatus bool    `json:"publishStatus" form:"publishStatus" gorm:"comment:发布状态;default:false"`
}

// TableName 数据集表名
func (Dataset) TableName() string {
	return "dataset"
}