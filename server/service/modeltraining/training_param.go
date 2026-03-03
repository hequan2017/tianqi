package modeltraining

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modeltrainingModel "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
)

type TrainingParamService struct{}

// GetParamByTaskId 根据任务ID获取训练参数
func (s *TrainingParamService) GetParamByTaskId(ctx context.Context, taskId uint) (modeltrainingModel.TrainingParam, error) {
	var param modeltrainingModel.TrainingParam
	err := global.GVA_DB.Where("task_id = ?", taskId).First(&param).Error
	if err != nil {
		return param, fmt.Errorf("训练参数不存在")
	}
	return param, nil
}

// UpdateParam 更新训练参数
func (s *TrainingParamService) UpdateParam(ctx context.Context, param *modeltrainingModel.TrainingParam) error {
	return global.GVA_DB.Model(&modeltrainingModel.TrainingParam{}).
		Where("task_id = ?", param.TaskId).
		Updates(param).Error
}