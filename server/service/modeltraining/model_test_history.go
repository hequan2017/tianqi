package modeltraining

import (
	"context"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	modeltraining "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
	modeltrainingReq "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ModelTestHistoryService struct{}

// CreateTestHistory 创建模型测试历史
func (s *ModelTestHistoryService) CreateTestHistory(ctx context.Context, req *modeltrainingReq.CreateModelTestReq, userID *uint) error {
	now := time.Now()
	history := modeltraining.ModelTestHistory{
		TaskId:     req.TaskId,
		UserId:     userID,
		Question:   req.Question,
		BaseAnswer: req.BaseAnswer,
		LoraAnswer: req.LoraAnswer,
		TestTime:   &now,
	}

	if err := global.GVA_DB.WithContext(ctx).Create(&history).Error; err != nil {
		global.GVA_LOG.Error("创建模型测试历史失败!", zap.Error(err))
		return err
	}
	return nil
}

// DeleteTestHistory 删除模型测试历史
func (s *ModelTestHistoryService) DeleteTestHistory(ctx context.Context, id uint, userID *uint, isAdmin bool) error {
	query := global.GVA_DB.WithContext(ctx).Where("id = ?", id)
	if !isAdmin && userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	result := query.Delete(&modeltraining.ModelTestHistory{})
	if result.Error != nil {
		global.GVA_LOG.Error("删除模型测试历史失败!", zap.Error(result.Error))
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteTestHistoryByTaskId 删除指定任务的所有测试历史
func (s *ModelTestHistoryService) DeleteTestHistoryByTaskId(ctx context.Context, taskId uint) error {
	if err := global.GVA_DB.WithContext(ctx).Where("task_id = ?", taskId).Delete(&modeltraining.ModelTestHistory{}).Error; err != nil {
		global.GVA_LOG.Error("删除模型测试历史失败!", zap.Error(err))
		return err
	}
	return nil
}

// GetTestHistoryList 获取模型测试历史列表
func (s *ModelTestHistoryService) GetTestHistoryList(ctx context.Context, req modeltrainingReq.ModelTestHistorySearch, userID *uint, isAdmin bool) (list []modeltraining.ModelTestHistory, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.GVA_DB.WithContext(ctx).Model(&modeltraining.ModelTestHistory{})

	if req.TaskId != nil {
		db = db.Where("task_id = ?", *req.TaskId)
	}

	if !isAdmin && userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	if err = db.Count(&total).Error; err != nil {
		global.GVA_LOG.Error("获取模型测试历史总数失败!", zap.Error(err))
		return
	}

	if err = db.Order("test_time DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		global.GVA_LOG.Error("获取模型测试历史列表失败!", zap.Error(err))
		return
	}

	return
}

// GetTestHistoryDetail 获取模型测试历史详情
func (s *ModelTestHistoryService) GetTestHistoryDetail(ctx context.Context, id uint) (modeltraining.ModelTestHistory, error) {
	var history modeltraining.ModelTestHistory
	if err := global.GVA_DB.WithContext(ctx).Where("id = ?", id).First(&history).Error; err != nil {
		global.GVA_LOG.Error("获取模型测试历史详情失败!", zap.Error(err))
		return history, err
	}
	return history, nil
}

// ClearTestHistoryByTaskId 清空指定任务的测试历史
func (s *ModelTestHistoryService) ClearTestHistoryByTaskId(ctx context.Context, taskId uint, userID *uint, isAdmin bool) error {
	query := global.GVA_DB.WithContext(ctx).Where("task_id = ?", taskId)
	if !isAdmin && userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if err := query.Delete(&modeltraining.ModelTestHistory{}).Error; err != nil {
		global.GVA_LOG.Error("清空模型测试历史失败!", zap.Error(err))
		return err
	}
	return nil
}

// GetTestHistoryPage 分页获取模型测试历史
func (s *ModelTestHistoryService) GetTestHistoryPage(ctx context.Context, info request.PageInfo, taskId *uint, userID *uint, isAdmin bool) (list []modeltraining.ModelTestHistory, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.WithContext(ctx).Model(&modeltraining.ModelTestHistory{})

	if taskId != nil {
		db = db.Where("task_id = ?", *taskId)
	}

	if !isAdmin && userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	if err = db.Count(&total).Error; err != nil {
		return
	}

	if err = db.Order("test_time DESC").Limit(limit).Offset(offset).Find(&list).Error; err != nil {
		return
	}

	return
}