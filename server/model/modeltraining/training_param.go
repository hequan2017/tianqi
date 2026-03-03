// 训练参数配置
package modeltraining

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TrainingParam 训练参数配置
type TrainingParam struct {
	global.GVA_MODEL
	TaskId          uint    `json:"taskId" form:"taskId" gorm:"comment:任务ID;not null;uniqueIndex"`
	BatchSize       int     `json:"batchSize" form:"batchSize" gorm:"comment:批次大小"`
	LearningRate    float64 `json:"learningRate" form:"learningRate" gorm:"comment:学习率"`
	NEpochs         int     `json:"nEpochs" form:"nEpochs" gorm:"comment:训练轮数"`
	EvalSteps       int     `json:"evalSteps" form:"evalSteps" gorm:"comment:验证步数"`
	LoraAlpha       int     `json:"loraAlpha" form:"loraAlpha" gorm:"comment:LoRa缩放系数"`
	LoraDropout     float64 `json:"loraDropout" form:"loraDropout" gorm:"comment:LoRa Dropout"`
	LoraRank        int     `json:"loraRank" form:"loraRank" gorm:"comment:LoRa秩值"`
	LrSchedulerType string  `json:"lrSchedulerType" form:"lrSchedulerType" gorm:"comment:学习率调整策略;size:50"`
	MaxLength       int     `json:"maxLength" form:"maxLength" gorm:"comment:序列长度"`
	WarmupRatio     float64 `json:"warmupRatio" form:"warmupRatio" gorm:"comment:学习率预热比例"`
	WeightDecay     float64 `json:"weightDecay" form:"weightDecay" gorm:"comment:权重衰减"`
}

// TableName 训练参数表名
func (TrainingParam) TableName() string {
	return "training_param"
}