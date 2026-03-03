package modeltraining

import (
	"context"
	"fmt"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	modeltrainingModel "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
	modeltrainingReq "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining/request"
	"go.uber.org/zap"
)

// DatasetWithUser 包含用户信息的数据集结构体
type DatasetWithUser struct {
	modeltrainingModel.Dataset
	UserName string `json:"userName"` // 创建用户名
}

type DatasetService struct{}

// CreateDataset 创建数据集
func (s *DatasetService) CreateDataset(ctx context.Context, dataset *modeltrainingModel.Dataset) error {
	// 设置初始状态
	dataset.ImportStatus = "pending"
	// 如果前端没有设置发布状态，默认为不发布
	// PublishStatus 已经通过请求传入，无需额外设置

	err := global.GVA_DB.Create(dataset).Error
	if err != nil {
		global.GVA_LOG.Error("创建数据集失败", zap.Error(err))
		return fmt.Errorf("创建数据集失败: %v", err)
	}
	return nil
}

// DeleteDataset 删除数据集
func (s *DatasetService) DeleteDataset(ctx context.Context, ID string, userID uint, isAdmin bool) error {
	var dataset modeltrainingModel.Dataset
	if err := global.GVA_DB.Where("id = ?", ID).First(&dataset).Error; err != nil {
		return fmt.Errorf("数据集不存在: %v", err)
	}

	// 权限检查：普通用户只能删除自己创建的数据集
	if !isAdmin && (dataset.UserId == nil || *dataset.UserId != userID) {
		return fmt.Errorf("无权删除此数据集")
	}

	// 删除关联的版本
	if err := global.GVA_DB.Where("dataset_id = ?", ID).Delete(&modeltrainingModel.DatasetVersion{}).Error; err != nil {
		global.GVA_LOG.Warn("删除数据集版本失败", zap.Error(err))
	}

	// 删除数据集
	return global.GVA_DB.Delete(&dataset).Error
}

// DeleteDatasetByIds 批量删除数据集
func (s *DatasetService) DeleteDatasetByIds(ctx context.Context, IDs []string, userID uint, isAdmin bool) error {
	for _, id := range IDs {
		if err := s.DeleteDataset(ctx, id, userID, isAdmin); err != nil {
			global.GVA_LOG.Warn("删除数据集失败", zap.String("id", id), zap.Error(err))
			if err.Error() == "无权删除此数据集" {
				return err
			}
		}
	}
	return nil
}

// UpdateDataset 更新数据集
func (s *DatasetService) UpdateDataset(ctx context.Context, dataset modeltrainingModel.Dataset) error {
	return global.GVA_DB.Model(&modeltrainingModel.Dataset{}).Where("id = ?", dataset.ID).Updates(&dataset).Error
}

// GetDataset 根据ID获取数据集
func (s *DatasetService) GetDataset(ctx context.Context, ID string) (DatasetWithUser, error) {
	var result DatasetWithUser
	err := global.GVA_DB.Table("dataset").
		Select("dataset.*, sys_users.username as user_name").
		Joins("LEFT JOIN sys_users ON dataset.user_id = sys_users.id").
		Where("dataset.id = ? AND dataset.deleted_at IS NULL", ID).
		First(&result).Error
	if err != nil {
		return result, fmt.Errorf("获取数据集失败: %v", err)
	}
	return result, nil
}

// GetDatasetList 分页获取数据集列表
func (s *DatasetService) GetDatasetList(ctx context.Context, info modeltrainingReq.DatasetSearch, userID uint, isAdmin bool) ([]DatasetWithUser, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Table("dataset").
		Select("dataset.*, sys_users.username as user_name").
		Joins("LEFT JOIN sys_users ON dataset.user_id = sys_users.id")

	var datasets []DatasetWithUser

	// 时间范围筛选
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("dataset.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	// 权限控制：普通用户只能看到自己创建的数据集
	if !isAdmin {
		db = db.Where("dataset.user_id = ?", userID)
	}

	// 条件筛选
	if info.Name != nil && *info.Name != "" {
		db = db.Where("dataset.name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Type != nil && *info.Type != "" {
		db = db.Where("dataset.type = ?", *info.Type)
	}
	if info.TrainMethod != nil && *info.TrainMethod != "" {
		db = db.Where("dataset.train_method = ?", *info.TrainMethod)
	}
	if info.ImportStatus != nil && *info.ImportStatus != "" {
		db = db.Where("dataset.import_status = ?", *info.ImportStatus)
	}
	if info.PublishStatus != nil {
		db = db.Where("dataset.publish_status = ?", *info.PublishStatus)
	}
	if info.UserId != nil {
		db = db.Where("dataset.user_id = ?", *info.UserId)
	}

	// 统计总数
	var total int64
	countDB := global.GVA_DB.Model(&modeltrainingModel.Dataset{})
	if len(info.CreatedAtRange) == 2 {
		countDB = countDB.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if !isAdmin {
		countDB = countDB.Where("user_id = ?", userID)
	}
	if info.Name != nil && *info.Name != "" {
		countDB = countDB.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.Type != nil && *info.Type != "" {
		countDB = countDB.Where("type = ?", *info.Type)
	}
	if info.TrainMethod != nil && *info.TrainMethod != "" {
		countDB = countDB.Where("train_method = ?", *info.TrainMethod)
	}
	if info.ImportStatus != nil && *info.ImportStatus != "" {
		countDB = countDB.Where("import_status = ?", *info.ImportStatus)
	}
	if info.PublishStatus != nil {
		countDB = countDB.Where("publish_status = ?", *info.PublishStatus)
	}
	if info.UserId != nil {
		countDB = countDB.Where("user_id = ?", *info.UserId)
	}
	if err := countDB.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err := db.Order("dataset.created_at DESC").Find(&datasets).Error
	return datasets, total, err
}

// GetDatasetDataSource 获取数据源
func (s *DatasetService) GetDatasetDataSource(ctx context.Context) (map[string][]map[string]any, error) {
	res := make(map[string][]map[string]any)

	// 数据集类型
	datasetTypes := []map[string]any{
		{"label": "训练集", "value": "training"},
		{"label": "验证集", "value": "evaluation"},
	}
	res["type"] = datasetTypes

	// 训练方式
	trainMethods := []map[string]any{
		{"label": "SFT", "value": "SFT"},
		{"label": "DPO", "value": "DPO"},
		{"label": "CPT", "value": "CPT"},
	}
	res["trainMethod"] = trainMethods

	// 导入状态
	importStatuses := []map[string]any{
		{"label": "待导入", "value": "pending"},
		{"label": "导入成功", "value": "success"},
		{"label": "导入失败", "value": "failed"},
	}
	res["importStatus"] = importStatuses

	// 用户列表
	users := make([]map[string]any, 0)
	global.GVA_DB.Table("sys_users").Where("deleted_at IS NULL").Select("username as label, id as value").Scan(&users)
	res["userId"] = users

	return res, nil
}

// PublishDataset 发布数据集
func (s *DatasetService) PublishDataset(ctx context.Context, ID uint, userID uint, isAdmin bool) error {
	var dataset modeltrainingModel.Dataset
	if err := global.GVA_DB.Where("id = ?", ID).First(&dataset).Error; err != nil {
		return fmt.Errorf("数据集不存在")
	}

	// 权限检查
	if !isAdmin && (dataset.UserId == nil || *dataset.UserId != userID) {
		return fmt.Errorf("无权发布此数据集")
	}

	// 检查导入状态
	if dataset.ImportStatus != "success" {
		return fmt.Errorf("数据集导入未成功，无法发布")
	}

	return global.GVA_DB.Model(&dataset).Update("publish_status", true).Error
}