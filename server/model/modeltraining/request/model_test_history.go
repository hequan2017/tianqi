package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// ModelTestHistorySearch 模型测试历史搜索请求
type ModelTestHistorySearch struct {
	TaskId *uint `json:"taskId" form:"taskId" binding:"required"`
	request.PageInfo
}

// CreateModelTestReq 创建模型测试请求
type CreateModelTestReq struct {
	TaskId     uint   `json:"taskId" form:"taskId" binding:"required"`
	Question   string `json:"question" form:"question" binding:"required"`
	BaseAnswer string `json:"baseAnswer" form:"baseAnswer"`
	LoraAnswer string `json:"loraAnswer" form:"loraAnswer"`
}