package modeltraining

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modeltrainingModel "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
	modeltrainingReq "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DatasetVersionService struct{}

// CreateVersion 创建数据集版本
func (s *DatasetVersionService) CreateVersion(ctx context.Context, version *modeltrainingModel.DatasetVersion) error {
	// 检查数据集是否存在
	var dataset modeltrainingModel.Dataset
	if err := global.GVA_DB.Where("id = ?", version.DatasetId).First(&dataset).Error; err != nil {
		return fmt.Errorf("数据集不存在")
	}

	// 检查版本号是否已存在
	var count int64
	global.GVA_DB.Model(&modeltrainingModel.DatasetVersion{}).
		Where("dataset_id = ? AND version = ?", version.DatasetId, version.Version).
		Count(&count)
	if count > 0 {
		return fmt.Errorf("版本号已存在")
	}

	// 设置初始状态
	version.Status = "pending"

	err := global.GVA_DB.Create(version).Error
	if err != nil {
		return fmt.Errorf("创建版本失败: %v", err)
	}

	// 更新数据集的最新版本
	global.GVA_DB.Model(&dataset).Update("latest_version", version.Version)

	return nil
}

// DeleteVersion 删除数据集版本
func (s *DatasetVersionService) DeleteVersion(ctx context.Context, ID string) error {
	return global.GVA_DB.Delete(&modeltrainingModel.DatasetVersion{}, "id = ?", ID).Error
}

// GetVersionList 获取数据集版本列表
func (s *DatasetVersionService) GetVersionList(ctx context.Context, info modeltrainingReq.DatasetVersionSearch) ([]modeltrainingModel.DatasetVersion, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Model(&modeltrainingModel.DatasetVersion{}).Where("dataset_id = ?", info.DatasetId)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var versions []modeltrainingModel.DatasetVersion
	query := global.GVA_DB.Where("dataset_id = ?", info.DatasetId)
	if limit != 0 {
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Order("created_at DESC").Find(&versions).Error
	return versions, total, err
}

// GetVersion 根据ID获取版本
func (s *DatasetVersionService) GetVersion(ctx context.Context, ID string) (modeltrainingModel.DatasetVersion, error) {
	var version modeltrainingModel.DatasetVersion
	err := global.GVA_DB.Where("id = ?", ID).First(&version).Error
	if err != nil {
		return version, fmt.Errorf("版本不存在")
	}
	return version, nil
}

// UpdateVersionStatus 更新版本状态
func (s *DatasetVersionService) UpdateVersionStatus(ctx context.Context, ID uint, status string) error {
	return global.GVA_DB.Model(&modeltrainingModel.DatasetVersion{}).Where("id = ?", ID).Update("status", status).Error
}

// UpdateVersionSuccess 更新版本为成功状态并更新数据集统计
func (s *DatasetVersionService) UpdateVersionSuccess(ctx context.Context, ID uint, dataCount int64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 更新版本状态
		if err := tx.Model(&modeltrainingModel.DatasetVersion{}).Where("id = ?", ID).Updates(map[string]interface{}{
			"status":     "success",
			"data_count": dataCount,
		}).Error; err != nil {
			return err
		}

		// 获取版本信息
		var version modeltrainingModel.DatasetVersion
		if err := tx.Where("id = ?", ID).First(&version).Error; err != nil {
			return err
		}

		// 更新数据集的数据量和导入状态
		if err := tx.Model(&modeltrainingModel.Dataset{}).Where("id = ?", version.DatasetId).Updates(map[string]interface{}{
			"data_count":   dataCount,
			"import_status": "success",
		}).Error; err != nil {
			global.GVA_LOG.Error("更新数据集统计失败", zap.Error(err))
			return err
		}

		return nil
	})
}