// 训练任务管理
package modeltraining

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// TrainingTask 训练任务
type TrainingTask struct {
	global.GVA_MODEL
	Name           string     `json:"name" form:"name" gorm:"comment:任务名称;size:100;not null"`
	TaskId         string     `json:"taskId" form:"taskId" gorm:"comment:任务ID;size:50;uniqueIndex"` // 自动生成唯一ID
	UserId         *uint      `json:"userId" form:"userId" gorm:"comment:创建者ID;index"`
	BaseModel      string     `json:"baseModel" form:"baseModel" gorm:"comment:基础模型;size:100"`                  // Qwen3-1.7B
	TrainMethod    string     `json:"trainMethod" form:"trainMethod" gorm:"comment:训练方式;size:20"`               // SFT、DPO、CPT
	TrainType      string     `json:"trainType" form:"trainType" gorm:"comment:训练类型;size:20"`                   // efficient/full
	Status         string     `json:"status" form:"status" gorm:"comment:训练状态;size:20;default:'pending';index"` // pending/running/completed/failed
	TrainDatasetId *uint      `json:"trainDatasetId" form:"trainDatasetId" gorm:"comment:训练集ID"`
	TrainVersionId *uint      `json:"trainVersionId" form:"trainVersionId" gorm:"comment:训练集版本ID"`
	ValDatasetId   *uint      `json:"valDatasetId" form:"valDatasetId" gorm:"comment:验证集ID"`
	ValVersionId   *uint      `json:"valVersionId" form:"valVersionId" gorm:"comment:验证集版本ID"`
	ValSplitRatio  float64    `json:"valSplitRatio" form:"valSplitRatio" gorm:"comment:验证集切分比例"` // 0.1 = 10%
	OutputCount    int        `json:"outputCount" form:"outputCount" gorm:"comment:产出数量上限;default:5"`
	ModelName      string     `json:"modelName" form:"modelName" gorm:"comment:输出模型名称;size:100"`
	CheckpointInt  int        `json:"checkpointInterval" form:"checkpointInterval" gorm:"comment:Checkpoint保存间隔"`
	CheckpointUnit string     `json:"checkpointUnit" form:"checkpointUnit" gorm:"comment:Checkpoint间隔单位;size:10"` // epoch/step
	Progress       int        `json:"progress" form:"progress" gorm:"comment:训练进度(百分比);default:0"`
	StartTime      *time.Time `json:"startTime" form:"startTime" gorm:"comment:开始时间"`
	EndTime        *time.Time `json:"endTime" form:"endTime" gorm:"comment:结束时间"`
	NodeId         *uint      `json:"nodeId" form:"nodeId" gorm:"comment:执行节点ID;index"`
	InstanceId     *uint      `json:"instanceId" form:"instanceId" gorm:"comment:实例ID;index"`
	HostPort       *int       `json:"hostPort" form:"hostPort" gorm:"comment:训练容器端口"`
	ContainerId    *string    `json:"containerId" form:"containerId" gorm:"comment:训练容器ID;size:128"`
	ContainerName  *string    `json:"containerName" form:"containerName" gorm:"comment:训练容器名称;size:255"`
	CheckpointPath *string    `json:"checkpointPath" form:"checkpointPath" gorm:"comment:训练产出Checkpoint路径;size:512"`
	Remark         *string    `json:"remark" form:"remark" gorm:"comment:备注信息;size:1000"`
}

// TableName 训练任务表名
func (TrainingTask) TableName() string {
	return "training_task"
}
