// 模型测试历史
package modeltraining

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ModelTestHistory 模型测试历史
type ModelTestHistory struct {
	global.GVA_MODEL
	TaskId      uint       `json:"taskId" form:"taskId" gorm:"comment:训练任务ID;not null;index"`
	UserId      *uint      `json:"userId" form:"userId" gorm:"comment:创建者ID;index"`
	Question    string     `json:"question" form:"question" gorm:"comment:测试问题;type:text;not null"`
	BaseAnswer  string     `json:"baseAnswer" form:"baseAnswer" gorm:"comment:基础模型回复;type:text"`
	LoraAnswer  string     `json:"loraAnswer" form:"loraAnswer" gorm:"comment:LoRA模型回复;type:text"`
	TestTime    *time.Time `json:"testTime" form:"testTime" gorm:"comment:测试时间"`
}

// TableName 模型测试历史表名
func (ModelTestHistory) TableName() string {
	return "model_test_history"
}