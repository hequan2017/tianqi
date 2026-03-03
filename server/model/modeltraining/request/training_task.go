package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

// TrainingTaskSearch 训练任务搜索请求
type TrainingTaskSearch struct {
	CreatedAtRange []time.Time `json:"createdAtRange" form:"createdAtRange[]"`
	Name           *string     `json:"name" form:"name"`
	TaskId         *string     `json:"taskId" form:"taskId"`
	Status         *string     `json:"status" form:"status"`
	TrainMethod    *string     `json:"trainMethod" form:"trainMethod"`
	UserId         *uint       `json:"userId" form:"userId"`
	BaseModel      *string     `json:"baseModel" form:"baseModel"`
	request.PageInfo
}

// CreateTrainingTaskReq 创建训练任务请求
type CreateTrainingTaskReq struct {
	Name            string  `json:"name" form:"name" binding:"required"`
	BaseModel       string  `json:"baseModel" form:"baseModel" binding:"required"`
	TrainMethod     string  `json:"trainMethod" form:"trainMethod" binding:"required"` // SFT、DPO、CPT
	TrainType       string  `json:"trainType" form:"trainType" binding:"required"`     // efficient/full
	TrainDatasetId  *uint   `json:"trainDatasetId" form:"trainDatasetId"`
	TrainVersionId  *uint   `json:"trainVersionId" form:"trainVersionId"`
	ValDatasetId    *uint   `json:"valDatasetId" form:"valDatasetId"`
	ValVersionId    *uint   `json:"valVersionId" form:"valVersionId"`
	ValSplitRatio   float64 `json:"valSplitRatio" form:"valSplitRatio"`
	OutputCount     int     `json:"outputCount" form:"outputCount"`
	ModelName       string  `json:"modelName" form:"modelName" binding:"required"`
	CheckpointInt   int     `json:"checkpointInterval" form:"checkpointInterval"`
	CheckpointUnit  string  `json:"checkpointUnit" form:"checkpointUnit"`
	Remark          string  `json:"remark" form:"remark"`
	// 训练参数
	BatchSize       int     `json:"batchSize" form:"batchSize"`
	LearningRate    float64 `json:"learningRate" form:"learningRate"`
	NEpochs         int     `json:"nEpochs" form:"nEpochs"`
	EvalSteps       int     `json:"evalSteps" form:"evalSteps"`
	LoraAlpha       int     `json:"loraAlpha" form:"loraAlpha"`
	LoraDropout     float64 `json:"loraDropout" form:"loraDropout"`
	LoraRank        int     `json:"loraRank" form:"loraRank"`
	LrSchedulerType string  `json:"lrSchedulerType" form:"lrSchedulerType"`
	MaxLength       int     `json:"maxLength" form:"maxLength"`
	WarmupRatio     float64 `json:"warmupRatio" form:"warmupRatio"`
	WeightDecay     float64 `json:"weightDecay" form:"weightDecay"`
}

// UpdateTrainingTaskReq 更新训练任务请求
type UpdateTrainingTaskReq struct {
	ID              uint    `json:"id" form:"id" binding:"required"`
	Name            string  `json:"name" form:"name"`
	ModelName       string  `json:"modelName" form:"modelName"`
	CheckpointInt   int     `json:"checkpointInterval" form:"checkpointInterval"`
	CheckpointUnit  string  `json:"checkpointUnit" form:"checkpointUnit"`
	Remark          string  `json:"remark" form:"remark"`
}

// StartTaskReq 启动训练任务请求
type StartTaskReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// StopTaskReq 停止训练任务请求
type StopTaskReq struct {
	ID uint `json:"id" form:"id" binding:"required"`
}

// TaskLogReq 获取训练日志请求
type TaskLogReq struct {
	ID   uint   `json:"id" form:"id" binding:"required"`
	Tail string `json:"tail" form:"tail"`
}

// ChatMessage 对话消息
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionReq 模型对话测试请求
type ChatCompletionReq struct {
	ID       uint          `json:"id" form:"id" binding:"required"`
	Model    string        `json:"model" form:"model" binding:"required"` // base 或 lora
	Messages []ChatMessage `json:"messages" form:"messages" binding:"required"`
}