package modeltraining

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/computenode"
	"github.com/flipped-aurora/gin-vue-admin/server/model/imageregistry"
	instanceModel "github.com/flipped-aurora/gin-vue-admin/server/model/instance"
	modeltrainingModel "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
	modeltrainingReq "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/product"
	instanceSvc "github.com/flipped-aurora/gin-vue-admin/server/service/instance"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TrainingTaskWithUser 包含用户信息的训练任务结构体
type TrainingTaskWithUser struct {
	modeltrainingModel.TrainingTask
	UserName         string `json:"userName"`         // 创建用户名
	TrainDatasetName string `json:"trainDatasetName"` // 训练集名称
	ValDatasetName   string `json:"valDatasetName"`   // 验证集名称
}

type TrainingTaskService struct{}

const (
	trainingNodeName       = "chengdu"
	trainingImage          = "modelscope-registry.cn-hangzhou.cr.aliyuncs.com/modelscope-repo/modelscope:ubuntu22.04-cuda12.8.1-py311-torch2.9.0-vllm0.13.0-modelscope1.33.0-swift3.12.5"
	trainingPortBase       = 8001
	trainingImageName      = "modelscope-training-auto"
	trainingSpecNamePrefix = "training-auto-spec-"
)

var (
	trainingDockerService   = &instanceSvc.DockerService{}
	trainingPortMu          sync.Mutex
	trainingInstanceService = &instanceSvc.InstanceService{}
)

// CreateTask 创建训练任务
func (s *TrainingTaskService) CreateTask(ctx context.Context, req *modeltrainingReq.CreateTrainingTaskReq, userID uint) (*modeltrainingModel.TrainingTask, error) {
	// 生成唯一任务ID
	taskId := generateTaskId()

	task := &modeltrainingModel.TrainingTask{
		Name:           req.Name,
		TaskId:         taskId,
		UserId:         &userID,
		BaseModel:      req.BaseModel,
		TrainMethod:    req.TrainMethod,
		TrainType:      req.TrainType,
		Status:         "pending",
		TrainDatasetId: req.TrainDatasetId,
		TrainVersionId: req.TrainVersionId,
		ValDatasetId:   req.ValDatasetId,
		ValVersionId:   req.ValVersionId,
		ValSplitRatio:  req.ValSplitRatio,
		OutputCount:    req.OutputCount,
		ModelName:      req.ModelName,
		CheckpointInt:  req.CheckpointInt,
		CheckpointUnit: req.CheckpointUnit,
		Progress:       0,
	}

	if req.Remark != "" {
		task.Remark = &req.Remark
	}

	// 设置默认值
	if task.OutputCount == 0 {
		task.OutputCount = 5
	}
	if task.CheckpointInt == 0 {
		task.CheckpointInt = 1
	}
	if task.CheckpointUnit == "" {
		task.CheckpointUnit = "epoch"
	}

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建任务
		if err := tx.Create(task).Error; err != nil {
			return fmt.Errorf("创建训练任务失败: %v", err)
		}

		// 创建训练参数
		param := &modeltrainingModel.TrainingParam{
			TaskId:          task.ID,
			BatchSize:       req.BatchSize,
			LearningRate:    req.LearningRate,
			NEpochs:         req.NEpochs,
			EvalSteps:       req.EvalSteps,
			LoraAlpha:       req.LoraAlpha,
			LoraDropout:     req.LoraDropout,
			LoraRank:        req.LoraRank,
			LrSchedulerType: req.LrSchedulerType,
			MaxLength:       req.MaxLength,
			WarmupRatio:     req.WarmupRatio,
			WeightDecay:     req.WeightDecay,
		}

		// 设置默认参数
		if param.BatchSize == 0 {
			param.BatchSize = 4
		}
		if param.LearningRate == 0 {
			param.LearningRate = 0.0001
		}
		if param.NEpochs == 0 {
			param.NEpochs = 3
		}
		if param.LoraAlpha == 0 {
			param.LoraAlpha = 16
		}
		if param.LoraRank == 0 {
			param.LoraRank = 8
		}
		if param.MaxLength == 0 {
			param.MaxLength = 2048
		}
		if param.LrSchedulerType == "" {
			param.LrSchedulerType = "cosine"
		}

		if err := tx.Create(param).Error; err != nil {
			return fmt.Errorf("创建训练参数失败: %v", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 创建任务后自动启动训练流程：创建容器并执行微调命令
	if err := s.launchTrainingTask(ctx, task.ID); err != nil {
		return nil, err
	}

	return task, nil
}

// generateTaskId 生成唯一任务ID
func generateTaskId() string {
	return fmt.Sprintf("train-%s", uuid.New().String()[:8])
}

// DeleteTask 删除训练任务
func (s *TrainingTaskService) DeleteTask(ctx context.Context, ID string, userID uint, isAdmin bool) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.Where("id = ?", ID).First(&task).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	// 权限检查
	if !isAdmin && (task.UserId == nil || *task.UserId != userID) {
		return fmt.Errorf("无权删除此训练任务")
	}

	// 检查任务状态
	if task.Status == "running" || task.Status == "serving" || task.Status == "completed" {
		return fmt.Errorf("正在运行的任务无法删除，请先停止任务")
	}

	// 清理关联实例（会同时清理容器）
	if task.InstanceId != nil && *task.InstanceId > 0 {
		if delErr := trainingInstanceService.DeleteInstance(ctx, strconv.FormatUint(uint64(*task.InstanceId), 10), 0, true); delErr != nil {
			global.GVA_LOG.Warn("删除训练任务关联实例失败", zap.Uint("taskID", task.ID), zap.Uint("instanceID", *task.InstanceId), zap.Error(delErr))
		}
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 删除训练参数
		if err := tx.Where("task_id = ?", ID).Delete(&modeltrainingModel.TrainingParam{}).Error; err != nil {
			global.GVA_LOG.Warn("删除训练参数失败", zap.Error(err))
		}

		// 删除任务
		return tx.Delete(&task).Error
	})
}

// DeleteTaskByIds 批量删除训练任务
func (s *TrainingTaskService) DeleteTaskByIds(ctx context.Context, IDs []string, userID uint, isAdmin bool) error {
	for _, id := range IDs {
		if err := s.DeleteTask(ctx, id, userID, isAdmin); err != nil {
			if err.Error() == "无权删除此训练任务" || err.Error() == "正在运行的任务无法删除，请先停止任务" {
				return err
			}
			global.GVA_LOG.Warn("删除训练任务失败", zap.String("id", id), zap.Error(err))
		}
	}
	return nil
}

// UpdateTask 更新训练任务
func (s *TrainingTaskService) UpdateTask(ctx context.Context, req *modeltrainingReq.UpdateTrainingTaskReq) error {
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.ModelName != "" {
		updates["model_name"] = req.ModelName
	}
	if req.CheckpointInt > 0 {
		updates["checkpoint_interval"] = req.CheckpointInt
	}
	if req.CheckpointUnit != "" {
		updates["checkpoint_unit"] = req.CheckpointUnit
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	return global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).Where("id = ?", req.ID).Updates(updates).Error
}

// GetTask 根据ID获取训练任务
func (s *TrainingTaskService) GetTask(ctx context.Context, ID string) (TrainingTaskWithUser, error) {
	var result TrainingTaskWithUser
	err := global.GVA_DB.Table("training_task").
		Select(`training_task.*,
			sys_users.username as user_name,
			d1.name as train_dataset_name,
			d2.name as val_dataset_name`).
		Joins("LEFT JOIN sys_users ON training_task.user_id = sys_users.id").
		Joins("LEFT JOIN dataset d1 ON training_task.train_dataset_id = d1.id").
		Joins("LEFT JOIN dataset d2 ON training_task.val_dataset_id = d2.id").
		Where("training_task.id = ? AND training_task.deleted_at IS NULL", ID).
		First(&result).Error
	if err != nil {
		return result, fmt.Errorf("获取训练任务失败: %v", err)
	}
	s.reconcileServingTaskStatus(ctx, &result.TrainingTask)
	return result, nil
}

// GetTaskList 分页获取训练任务列表
func (s *TrainingTaskService) GetTaskList(ctx context.Context, info modeltrainingReq.TrainingTaskSearch, userID uint, isAdmin bool) ([]TrainingTaskWithUser, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := global.GVA_DB.Table("training_task").
		Select(`training_task.*,
			sys_users.username as user_name,
			d1.name as train_dataset_name,
			d2.name as val_dataset_name`).
		Joins("LEFT JOIN sys_users ON training_task.user_id = sys_users.id").
		Joins("LEFT JOIN dataset d1 ON training_task.train_dataset_id = d1.id").
		Joins("LEFT JOIN dataset d2 ON training_task.val_dataset_id = d2.id")

	var tasks []TrainingTaskWithUser

	// 时间范围筛选
	if len(info.CreatedAtRange) == 2 {
		db = db.Where("training_task.created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}

	// 权限控制
	if !isAdmin {
		db = db.Where("training_task.user_id = ?", userID)
	}

	// 条件筛选
	if info.Name != nil && *info.Name != "" {
		db = db.Where("training_task.name LIKE ?", "%"+*info.Name+"%")
	}
	if info.TaskId != nil && *info.TaskId != "" {
		db = db.Where("training_task.task_id = ?", *info.TaskId)
	}
	if info.Status != nil && *info.Status != "" {
		db = db.Where("training_task.status = ?", *info.Status)
	}
	if info.TrainMethod != nil && *info.TrainMethod != "" {
		db = db.Where("training_task.train_method = ?", *info.TrainMethod)
	}
	if info.BaseModel != nil && *info.BaseModel != "" {
		db = db.Where("training_task.base_model LIKE ?", "%"+*info.BaseModel+"%")
	}
	if info.UserId != nil {
		db = db.Where("training_task.user_id = ?", *info.UserId)
	}

	// 统计总数
	var total int64
	countDB := global.GVA_DB.Model(&modeltrainingModel.TrainingTask{})
	if len(info.CreatedAtRange) == 2 {
		countDB = countDB.Where("created_at BETWEEN ? AND ?", info.CreatedAtRange[0], info.CreatedAtRange[1])
	}
	if !isAdmin {
		countDB = countDB.Where("user_id = ?", userID)
	}
	if info.Name != nil && *info.Name != "" {
		countDB = countDB.Where("name LIKE ?", "%"+*info.Name+"%")
	}
	if info.TaskId != nil && *info.TaskId != "" {
		countDB = countDB.Where("task_id = ?", *info.TaskId)
	}
	if info.Status != nil && *info.Status != "" {
		countDB = countDB.Where("status = ?", *info.Status)
	}
	if info.TrainMethod != nil && *info.TrainMethod != "" {
		countDB = countDB.Where("train_method = ?", *info.TrainMethod)
	}
	if info.BaseModel != nil && *info.BaseModel != "" {
		countDB = countDB.Where("base_model LIKE ?", "%"+*info.BaseModel+"%")
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

	err := db.Order("training_task.created_at DESC").Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}
	for i := range tasks {
		s.reconcileServingTaskStatus(ctx, &tasks[i].TrainingTask)
	}
	return tasks, total, nil
}

// StartTask 启动训练任务
func (s *TrainingTaskService) StartTask(ctx context.Context, ID uint, userID uint, isAdmin bool) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.Where("id = ?", ID).First(&task).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	// 权限检查
	if !isAdmin && (task.UserId == nil || *task.UserId != userID) {
		return fmt.Errorf("无权操作此训练任务")
	}

	// 检查状态
	if task.Status == "running" {
		return fmt.Errorf("任务已在运行中")
	}

	return s.launchTrainingTask(ctx, task.ID)
}

// StopTask 停止训练任务
func (s *TrainingTaskService) StopTask(ctx context.Context, ID uint, userID uint, isAdmin bool) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.Where("id = ?", ID).First(&task).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	// 权限检查
	if !isAdmin && (task.UserId == nil || *task.UserId != userID) {
		return fmt.Errorf("无权操作此训练任务")
	}

	// 检查状态
	if task.Status != "running" && task.Status != "serving" {
		return fmt.Errorf("任务未在运行中")
	}

	// 尝试停止对应容器
	if task.ContainerId != nil && *task.ContainerId != "" && task.NodeId != nil {
		var node computenode.ComputeNode
		if err := global.GVA_DB.Where("id = ?", *task.NodeId).First(&node).Error; err == nil {
			cli, err := trainingDockerService.CreateDockerClient(&node)
			if err == nil {
				defer cli.Close()
				timeout := 10
				if stopErr := cli.ContainerStop(ctx, *task.ContainerId, container.StopOptions{Timeout: &timeout}); stopErr != nil {
					global.GVA_LOG.Warn("停止训练容器失败", zap.Uint("taskID", task.ID), zap.String("containerID", *task.ContainerId), zap.Error(stopErr))
				}
			}
		}
	}

	_ = s.appendTaskLog(task.ID, "任务被手动停止")

	// 更新状态和结束时间
	now := time.Now()
	return global.GVA_DB.Model(&task).Updates(map[string]interface{}{
		"status":   "stopped",
		"end_time": now,
	}).Error
}

// GetTaskLogs 获取训练日志
func (s *TrainingTaskService) GetTaskLogs(ctx context.Context, ID uint, tail string) (string, error) {
	tailLines := 100
	if tail != "" {
		if n, err := strconv.Atoi(tail); err == nil && n > 0 {
			tailLines = n
		}
	}

	// 只读取容器日志（实时）
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.Where("id = ?", ID).First(&task).Error; err != nil {
		return "", fmt.Errorf("训练任务不存在")
	}
	if task.ContainerId != nil && *task.ContainerId != "" && task.NodeId != nil {
		var node computenode.ComputeNode
		if err := global.GVA_DB.Where("id = ?", *task.NodeId).First(&node).Error; err == nil {
			dockerLogs, err := trainingDockerService.GetContainerLogs(ctx, &node, *task.ContainerId, strconv.Itoa(tailLines))
			if err == nil && strings.TrimSpace(dockerLogs) != "" {
				return dockerLogs, nil
			}
		}
	}
	return "暂无日志（任务可能尚未启动或启动失败，请查看任务状态）", nil
}

// GetTaskDataSource 获取数据源
func (s *TrainingTaskService) GetTaskDataSource(ctx context.Context) (map[string][]map[string]any, error) {
	res := make(map[string][]map[string]any)

	// 训练状态
	statuses := []map[string]any{
		{"label": "待执行", "value": "pending"},
		{"label": "运行中", "value": "running"},
		{"label": "服务中", "value": "serving"},
		{"label": "已完成", "value": "completed"},
		{"label": "失败", "value": "failed"},
		{"label": "已停止", "value": "stopped"},
	}
	res["status"] = statuses

	// 训练方式
	trainMethods := []map[string]any{
		{"label": "SFT", "value": "SFT"},
		{"label": "DPO", "value": "DPO"},
		{"label": "CPT", "value": "CPT"},
	}
	res["trainMethod"] = trainMethods

	// 训练类型
	trainTypes := []map[string]any{
		{"label": "高效训练", "value": "efficient"},
		{"label": "全量训练", "value": "full"},
	}
	res["trainType"] = trainTypes

	// 基础模型列表
	baseModels := []map[string]any{
		{"label": "Qwen/Qwen3-1.7B", "value": "Qwen/Qwen3-1.7B"},
		{"label": "Qwen/Qwen3.5-0.8B", "value": "Qwen/Qwen3.5-0.8B"},
		{"label": "Qwen/Qwen2.5-7B-Instruct", "value": "Qwen/Qwen2.5-7B-Instruct"},
		{"label": "Qwen/Qwen3-7B", "value": "Qwen/Qwen3-7B"},
		{"label": "Qwen/Qwen3-14B", "value": "Qwen/Qwen3-14B"},
		{"label": "Llama3-8B", "value": "Llama3-8B"},
	}
	res["baseModel"] = baseModels

	// 已发布的数据集列表
	datasets := make([]map[string]any, 0)
	global.GVA_DB.Table("dataset").
		Where("deleted_at IS NULL AND publish_status = ? AND import_status = ?", true, "success").
		Select("name as label, id as value").Scan(&datasets)
	res["datasets"] = datasets

	// 用户列表
	users := make([]map[string]any, 0)
	global.GVA_DB.Table("sys_users").Where("deleted_at IS NULL").Select("username as label, id as value").Scan(&users)
	res["userId"] = users

	return res, nil
}

// GetDefaultParams 获取默认训练参数
func (s *TrainingTaskService) GetDefaultParams(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"batchSize":       1,
		"learningRate":    0.0001,
		"nEpochs":         1,
		"evalSteps":       50,
		"loraAlpha":       32,
		"loraDropout":     0.05,
		"loraRank":        8,
		"lrSchedulerType": "cosine",
		"maxLength":       2048,
		"warmupRatio":     0.05,
		"weightDecay":     0.01,
	}
}

func (s *TrainingTaskService) launchTrainingTask(ctx context.Context, taskID uint) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.Where("id = ?", taskID).First(&task).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	var param modeltrainingModel.TrainingParam
	if err := global.GVA_DB.Where("task_id = ?", taskID).First(&param).Error; err != nil {
		return fmt.Errorf("训练参数不存在: %v", err)
	}

	node, err := s.getTrainingNode()
	if err != nil {
		_ = s.failTask(task.ID, fmt.Sprintf("获取算力节点失败: %v", err))
		return err
	}

	trainingPortMu.Lock()
	hostPort, err := s.allocateTrainingPort()
	trainingPortMu.Unlock()
	if err != nil {
		_ = s.failTask(task.ID, fmt.Sprintf("分配端口失败: %v", err))
		return err
	}

	imageID, specID, err := s.ensureInstanceResources(ctx, node)
	if err != nil {
		_ = s.failTask(task.ID, fmt.Sprintf("准备实例资源失败: %v", err))
		return err
	}

	instanceName := sanitizeContainerName(fmt.Sprintf("modelscope-%s", task.TaskId))
	userID := int64(0)
	if task.UserId != nil {
		userID = int64(*task.UserId)
	}
	nodeID := int64(node.ID)
	inst := &instanceModel.Instance{
		ImageId:    &imageID,
		SpecId:     &specID,
		UserId:     &userID,
		NodeId:     &nodeID,
		Name:       &instanceName,
		HostPort:   &hostPort,
		StartupCmd: ptrString("tail -f /dev/null"),
	}
	if err := trainingInstanceService.CreateInstance(ctx, inst); err != nil {
		_ = s.failTask(task.ID, fmt.Sprintf("通过实例管理创建训练容器失败: %v", err))
		return err
	}
	if inst.ContainerId == nil || *inst.ContainerId == "" {
		_ = s.failTask(task.ID, "实例已创建但容器ID为空")
		return fmt.Errorf("实例已创建但容器ID为空")
	}

	now := time.Now()
	task.ContainerId = inst.ContainerId
	task.ContainerName = inst.ContainerName
	task.NodeId = &node.ID
	task.InstanceId = &inst.ID
	task.HostPort = &hostPort
	task.Status = "running"
	task.StartTime = &now
	task.Progress = 1
	task.EndTime = nil

	if err := global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).Where("id = ?", task.ID).Updates(map[string]any{
		"container_id":   task.ContainerId,
		"container_name": task.ContainerName,
		"node_id":        node.ID,
		"instance_id":    inst.ID,
		"host_port":      hostPort,
		"status":         task.Status,
		"start_time":     task.StartTime,
		"progress":       task.Progress,
		"end_time":       task.EndTime,
	}).Error; err != nil {
		return fmt.Errorf("更新任务状态失败: %v", err)
	}

	containerName := safeString(task.ContainerName)
	containerID := safeString(task.ContainerId)
	_ = s.appendTaskLog(task.ID, fmt.Sprintf("任务启动: 节点=%s 容器=%s(%s) 端口=%d 镜像=%s", safeString(node.Name), containerName, containerID, hostPort, trainingImage))
	_ = s.appendTaskLog(task.ID, fmt.Sprintf("实例管理已创建实例: name=%s imageId=%d specId=%d nodeId=%d", instanceName, imageID, specID, node.ID))
	_ = s.appendTaskLog(task.ID, "开始执行训练命令...")

	go s.runSwiftTraining(task.ID, *node, containerID, task, param)

	return nil
}

func (s *TrainingTaskService) runSwiftTraining(taskID uint, node computenode.ComputeNode, containerID string, task modeltrainingModel.TrainingTask, param modeltrainingModel.TrainingParam) {
	ctx := context.Background()
	cli, err := trainingDockerService.CreateDockerClient(&node)
	if err != nil {
		_ = s.failTask(taskID, fmt.Sprintf("创建Docker客户端失败: %v", err))
		return
	}
	defer cli.Close()

	baseCommand := buildSwiftSFTCommand(task, param)
	command := wrapTrainingCommandForDockerLogs(baseCommand, taskID)
	_ = s.appendTaskLog(taskID, baseCommand)

	execResp, err := cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd: []string{"bash", "-lc", command},
	})
	if err != nil {
		_ = s.failTask(taskID, fmt.Sprintf("创建训练进程失败: %v", err))
		return
	}

	if err := cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{}); err != nil {
		_ = s.failTask(taskID, fmt.Sprintf("启动训练进程失败: %v", err))
		return
	}

	// 用于记录上次更新的进度，避免频繁更新数据库
	lastProgress := 0
	// 进度更新最小间隔（百分比）
	minProgressDelta := 1

	// 轮询执行状态，同时解析训练进度
	for {
		// 1. 检查执行状态
		info, err := cli.ContainerExecInspect(ctx, execResp.ID)
		if err != nil {
			_ = s.failTask(taskID, fmt.Sprintf("获取训练进程状态失败: %v", err))
			return
		}

		// 2. 如果训练进程还在运行，解析日志获取进度
		if info.Running {
			// 从容器日志解析进度
			progress := s.parseTrainingProgressFromLogs(ctx, &node, containerID, param.NEpochs)

			// 如果进度变化超过阈值，更新数据库
			if progress > 0 && progress-lastProgress >= minProgressDelta {
				lastProgress = progress
				_ = global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).Where("id = ?", taskID).Update("progress", progress).Error
				global.GVA_LOG.Debug("训练进度更新",
					zap.Uint("taskID", taskID),
					zap.Int("progress", progress))
			}

			time.Sleep(3 * time.Second)
			continue
		}

		// 3. 训练进程已结束
		now := time.Now()
		if info.ExitCode == 0 {
			// 训练成功完成
			checkpointPath, cpErr := s.extractLastModelCheckpoint(ctx, node, containerID)
			if cpErr != nil {
				_ = s.failTask(taskID, fmt.Sprintf("训练完成但提取 last_model_checkpoint 失败: %v", cpErr))
				s.stopTrainingContainer(context.Background(), &node, containerID)
				return
			}
			_ = s.appendTaskLog(taskID, fmt.Sprintf("识别到 last_model_checkpoint: %s", checkpointPath))
			_ = s.appendTaskLog(taskID, "训练任务执行完成，可通过「启动服务」API 启动 vLLM 推理服务")

			// 训练完成，保存 checkpoint 路径，状态设为 completed（不自动启动 vLLM）
			_ = global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).Where("id = ?", taskID).Updates(map[string]any{
				"status":          "completed",
				"progress":        100,
				"end_time":        now,
				"checkpoint_path": checkpointPath,
			}).Error
		} else {
			_ = s.failTask(taskID, fmt.Sprintf("训练任务执行失败，退出码=%d", info.ExitCode))
			s.stopTrainingContainer(context.Background(), &node, containerID)
		}
		return
	}
}

// parseTrainingProgressFromLogs 从容器日志解析训练进度
// 返回值：0-100 的进度百分比
func (s *TrainingTaskService) parseTrainingProgressFromLogs(ctx context.Context, node *computenode.ComputeNode, containerID string, totalEpochs int) int {
	// 获取最近的训练日志
	logs, err := trainingDockerService.GetContainerLogs(ctx, node, containerID, "500")
	if err != nil || logs == "" {
		global.GVA_LOG.Debug("获取训练日志失败或为空", zap.Error(err))
		return 0
	}

	// Swift 训练日志进度解析
	// 常见格式:
	// 1. [INFO:swift] Train: 0%| | 0/100 [00:00<?, ?it/s]
	// 2. [INFO:swift] Train: 50%|████ | 50/100 [01:23<01:23, 1.00it/s]
	// 3. [INFO:swift] epoch 1/3
	// 4. [INFO:swift] step: 100/1000
	// 5. {'loss': 0.5, 'epoch': 1.0}

	// 优先解析 Swift 的进度条格式: Train: 50%|████ | 50/100
	if match := regexp.MustCompile(`Train:\s*(\d+)%`).FindStringSubmatch(logs); len(match) > 1 {
		if percent, err := strconv.Atoi(match[1]); err == nil && percent >= 0 && percent <= 100 {
			global.GVA_LOG.Debug("解析到进度条格式", zap.Int("percent", percent))
			return percent
		}
	}

	// 解析 epoch 进度: epoch 1/3 或 [INFO:swift] epoch 1/3
	epochMatches := regexp.MustCompile(`(?i)epoch[:\s]*(\d+)/(\d+)`).FindAllStringSubmatch(logs, -1)
	if len(epochMatches) > 0 {
		lastMatch := epochMatches[len(epochMatches)-1]
		if len(lastMatch) >= 3 {
			currentEpoch, err1 := strconv.Atoi(lastMatch[1])
			totalEpoch := totalEpochs
			if err1 != nil || totalEpoch <= 0 {
				if totalFromLog, err2 := strconv.Atoi(lastMatch[2]); err2 == nil && totalFromLog > 0 {
					totalEpoch = totalFromLog
				}
			}
			if totalEpoch > 0 && currentEpoch > 0 {
				percent := (currentEpoch * 100) / totalEpoch
				if percent > 100 {
					percent = 100
				}
				global.GVA_LOG.Debug("解析到epoch进度", zap.Int("currentEpoch", currentEpoch), zap.Int("totalEpoch", totalEpoch), zap.Int("percent", percent))
				return percent
			}
		}
	}

	// 解析 step 进度: step: 100/1000 或 step=100/1000
	stepMatches := regexp.MustCompile(`(?i)step[:\s]*(\d+)/(\d+)`).FindAllStringSubmatch(logs, -1)
	if len(stepMatches) > 0 {
		lastMatch := stepMatches[len(stepMatches)-1]
		if len(lastMatch) >= 3 {
			currentStep, err1 := strconv.Atoi(lastMatch[1])
			totalSteps, err2 := strconv.Atoi(lastMatch[2])
			if err1 == nil && err2 == nil && totalSteps > 0 && currentStep > 0 {
				percent := (currentStep * 100) / totalSteps
				if percent > 100 {
					percent = 100
				}
				global.GVA_LOG.Debug("解析到step进度", zap.Int("currentStep", currentStep), zap.Int("totalSteps", totalSteps), zap.Int("percent", percent))
				return percent
			}
		}
	}

	// 解析 JSON 格式的进度信息: {'epoch': 1.0} 或 {"epoch": 1.0}
	jsonEpochMatches := regexp.MustCompile(`['"]epoch['"]\s*:\s*([0-9.]+)`).FindAllStringSubmatch(logs, -1)
	if len(jsonEpochMatches) > 0 && totalEpochs > 0 {
		lastMatch := jsonEpochMatches[len(jsonEpochMatches)-1]
		if len(lastMatch) >= 2 {
			if epochFloat, err := strconv.ParseFloat(lastMatch[1], 64); err == nil && epochFloat > 0 {
				percent := int(epochFloat * 100 / float64(totalEpochs))
				if percent > 100 {
					percent = 100
				}
				global.GVA_LOG.Debug("解析到JSON epoch进度", zap.Float64("epochFloat", epochFloat), zap.Int("totalEpochs", totalEpochs), zap.Int("percent", percent))
				return percent
			}
		}
	}

	// 解析 loss 出现次数作为估算
	lossMatches := regexp.MustCompile(`(?i)loss['"]?\s*[:=]\s*([0-9.]+)`).FindAllStringSubmatch(logs, -1)
	if len(lossMatches) > 0 && totalEpochs > 0 {
		lossCount := len(lossMatches)
		// 假设每个 epoch 产生约 10 个 loss 日志
		estimatedEpoch := lossCount / 10
		if estimatedEpoch > 0 {
			percent := (estimatedEpoch * 100) / totalEpochs
			if percent > 95 {
				percent = 95
			}
			if percent < 1 {
				percent = 1
			}
			global.GVA_LOG.Debug("基于loss次数估算进度", zap.Int("lossCount", lossCount), zap.Int("estimatedEpoch", estimatedEpoch), zap.Int("percent", percent))
			return percent
		}
	}

	global.GVA_LOG.Debug("未能解析训练进度", zap.String("logsPreview", truncateString(logs, 500)))
	return 0
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func (s *TrainingTaskService) getTrainingNode() (*computenode.ComputeNode, error) {
	var node computenode.ComputeNode
	if err := global.GVA_DB.Where("name = ? AND deleted_at IS NULL", trainingNodeName).First(&node).Error; err != nil {
		return nil, fmt.Errorf("未找到算力节点 %s", trainingNodeName)
	}
	return &node, nil
}

func (s *TrainingTaskService) allocateTrainingPort() (int, error) {
	var maxPort int
	err := global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).
		Where("host_port >= ?", trainingPortBase).
		Select("COALESCE(MAX(host_port), 0)").
		Scan(&maxPort).Error
	if err != nil {
		return 0, err
	}
	if maxPort < trainingPortBase {
		return trainingPortBase, nil
	}
	return maxPort + 1, nil
}

func (s *TrainingTaskService) ensureInstanceResources(ctx context.Context, node *computenode.ComputeNode) (int64, int64, error) {
	imageID, err := s.ensureTrainingImage(ctx)
	if err != nil {
		return 0, 0, err
	}
	specID, err := s.ensureTrainingSpec(ctx, node)
	if err != nil {
		return 0, 0, err
	}
	return imageID, specID, nil
}

func (s *TrainingTaskService) ensureTrainingImage(ctx context.Context) (int64, error) {
	var image imageregistry.ImageRegistry
	if err := global.GVA_DB.Where("address = ? AND deleted_at IS NULL", trainingImage).First(&image).Error; err == nil {
		return int64(image.ID), nil
	}

	name := trainingImageName
	addr := trainingImage
	desc := "训练任务自动创建镜像"
	source := "modeltraining-auto"
	isOnShelf := true
	supportMemorySplit := false
	newImage := imageregistry.ImageRegistry{
		Name:               &name,
		Address:            &addr,
		Description:        &desc,
		Source:             &source,
		IsOnShelf:          &isOnShelf,
		SupportMemorySplit: &supportMemorySplit,
	}
	if err := global.GVA_DB.Create(&newImage).Error; err != nil {
		return 0, err
	}
	return int64(newImage.ID), nil
}

func (s *TrainingTaskService) ensureTrainingSpec(ctx context.Context, node *computenode.ComputeNode) (int64, error) {
	specName := trainingSpecNamePrefix + strconv.FormatUint(uint64(node.ID), 10)
	var spec product.ProductSpec
	if err := global.GVA_DB.Where("name = ? AND deleted_at IS NULL", specName).First(&spec).Error; err == nil {
		// 训练任务固定使用单卡，避免节点GPU登记数与实际设备不一致导致 unknown device
		if spec.GpuCount == nil || *spec.GpuCount != 1 {
			one := int64(1)
			_ = global.GVA_DB.Model(&product.ProductSpec{}).Where("id = ?", spec.ID).Update("gpu_count", one).Error
			spec.GpuCount = &one
		}
		return int64(spec.ID), nil
	}

	// 训练容器固定申请1张卡，训练命令中用 CUDA_VISIBLE_DEVICES=0
	var gpuCount int64 = 1
	var gpuModel string = "auto"
	if node.GpuName != nil && strings.TrimSpace(*node.GpuName) != "" {
		gpuModel = *node.GpuName
	}
	cpuCores := int64(4)
	memGB := int64(16)
	systemDiskGB := int64(50)
	dataDiskGB := int64(100)
	price := float64(0)
	isOnShelf := true
	supportMemorySplit := false
	memoryCap := int64(0)
	if node.MemoryCapacity != nil && *node.MemoryCapacity > 0 {
		memoryCap = *node.MemoryCapacity
	}

	newSpec := product.ProductSpec{
		Name:               &specName,
		GpuModel:           &gpuModel,
		GpuCount:           &gpuCount,
		MemoryCapacity:     &memoryCap,
		CpuCores:           &cpuCores,
		MemoryGb:           &memGB,
		SystemDiskGb:       &systemDiskGB,
		DataDiskGb:         &dataDiskGB,
		PricePerHour:       &price,
		IsOnShelf:          &isOnShelf,
		SupportMemorySplit: &supportMemorySplit,
	}
	if err := global.GVA_DB.Create(&newSpec).Error; err != nil {
		return 0, err
	}
	return int64(newSpec.ID), nil
}

func (s *TrainingTaskService) taskLogPath(taskID uint) string {
	logDir := strings.TrimSpace(global.GVA_CONFIG.Zap.Director)
	if logDir == "" {
		logDir = "log"
	}
	return filepath.Join(logDir, "training", fmt.Sprintf("task-%d.log", taskID))
}

func (s *TrainingTaskService) appendTaskLog(taskID uint, content string) error {
	if err := os.MkdirAll(filepath.Dir(s.taskLogPath(taskID)), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(s.taskLogPath(taskID), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "[%s] %s\n", time.Now().Format("2006-01-02 15:04:05"), content)
	return err
}

func (s *TrainingTaskService) failTask(taskID uint, msg string) error {
	_ = s.appendTaskLog(taskID, msg)
	now := time.Now()
	return global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).Where("id = ?", taskID).Updates(map[string]any{
		"status":   "failed",
		"end_time": now,
	}).Error
}

func readLastNLines(path string, n int) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "暂无日志", nil
		}
		return "", err
	}
	defer f.Close()

	if n <= 0 {
		n = 100
	}
	lines := make([]string, 0, n)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	if len(lines) == 0 {
		return "暂无日志", nil
	}
	return strings.Join(lines, "\n"), nil
}

func buildSwiftSFTCommand(task modeltrainingModel.TrainingTask, param modeltrainingModel.TrainingParam) string {
	model := resolveVllmBaseModel(task.BaseModel)
	learningRate := param.LearningRate
	if learningRate <= 0 {
		learningRate = 1e-4
	}
	nEpochs := param.NEpochs
	if nEpochs <= 0 {
		nEpochs = 1
	}
	loraRank := param.LoraRank
	if loraRank <= 0 {
		loraRank = 8
	}
	loraAlpha := param.LoraAlpha
	if loraAlpha <= 0 {
		loraAlpha = 32
	}
	evalSteps := param.EvalSteps
	if evalSteps <= 0 {
		evalSteps = 50
	}
	maxLength := param.MaxLength
	if maxLength <= 0 {
		maxLength = 2048
	}
	warmupRatio := param.WarmupRatio
	if warmupRatio <= 0 {
		warmupRatio = 0.05
	}
	modelName := task.ModelName
	if strings.TrimSpace(modelName) == "" {
		modelName = "swift-robot"
	}

	return fmt.Sprintf(`CUDA_VISIBLE_DEVICES=0 swift sft \
    --model %s \
    --train_type lora \
    --dataset 'AI-ModelScope/alpaca-gpt4-data-zh#100' \
    --torch_dtype bfloat16 \
    --num_train_epochs %d \
    --per_device_train_batch_size 1 \
    --per_device_eval_batch_size 1 \
    --learning_rate %g \
    --lora_rank %d \
    --lora_alpha %d \
    --target_modules all-linear \
    --gradient_accumulation_steps 16 \
    --eval_steps %d \
    --save_steps 50 \
    --save_total_limit 2 \
    --logging_steps 5 \
    --max_length %d \
    --output_dir output \
    --system 'You are a helpful assistant.' \
    --warmup_ratio %.2f \
    --dataloader_num_workers 4 \
    --model_author swift \
    --model_name %s`, shQuote(model), nEpochs, learningRate, loraRank, loraAlpha, evalSteps, maxLength, warmupRatio, shQuote(modelName))
}

func sanitizeContainerName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, " ", "-")
	if name == "" {
		return fmt.Sprintf("modelscope-%d", time.Now().Unix())
	}
	return name
}

func safeString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func shQuote(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "swift-robot"
	}
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

func wrapTrainingCommandForDockerLogs(base string, taskID uint) string {
	logFile := fmt.Sprintf("/tmp/training-task-%d.log", taskID)
	// 输出同时落到容器日志与容器内文件，便于 docker logs 查看
	return fmt.Sprintf("(%s) 2>&1 | tee -a %s > /proc/1/fd/1", base, logFile)
}

func (s *TrainingTaskService) extractLastModelCheckpoint(ctx context.Context, node computenode.ComputeNode, containerID string) (string, error) {
	logs, err := trainingDockerService.GetContainerLogs(ctx, &node, containerID, "20000")
	if err != nil {
		return "", fmt.Errorf("获取训练容器日志失败: %w", err)
	}

	// 匹配多种格式:
	// 1. last_model_checkpoint: /output/v0-20260303-111956/checkpoint-7
	// 2. [INFO:swift] last_model_checkpoint: /output/checkpoint-xxx
	// 3. last_model_checkpoint:/output/checkpoint (无空格)
	re := regexp.MustCompile(`last_model_checkpoint[:\s]+([^\s"'\n\r]+)`)
	matches := re.FindAllStringSubmatch(logs, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("未在日志中找到 last_model_checkpoint")
	}
	checkpointPath := strings.TrimSpace(matches[len(matches)-1][1])
	if checkpointPath == "" {
		return "", fmt.Errorf("last_model_checkpoint 为空")
	}

	// 记录找到的 checkpoint 路径
	global.GVA_LOG.Info("训练完成，找到 checkpoint 路径",
		zap.String("checkpointPath", checkpointPath))

	return checkpointPath, nil
}

func (s *TrainingTaskService) startVllmServerInContainer(ctx context.Context, node computenode.ComputeNode, containerID string, task modeltrainingModel.TrainingTask, checkpointPath string) error {
	cli, err := trainingDockerService.CreateDockerClient(&node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %w", err)
	}
	defer cli.Close()

	baseModel := resolveVllmBaseModel(task.BaseModel)
	vllmCmd := fmt.Sprintf(
		"nohup python -m vllm.entrypoints.openai.api_server \\\n"+
			"    --model %s \\\n"+
			"    --served-model-name base \\\n"+
			"    --enable-lora \\\n"+
			"    --lora-modules lora=%s \\\n"+
			"    --max-lora-rank 64 \\\n"+
			"    --gpu-memory-utilization 0.8 \\\n"+
			"    --dtype float16 \\\n"+
			"    --max-model-len 4096 \\\n"+
			"    --port 8000 >/tmp/vllm-api-server.log 2>&1 &",
		shQuote(baseModel), shQuote(checkpointPath),
	)

	_ = s.appendTaskLog(task.ID, fmt.Sprintf("启动 vLLM 命令: %s", vllmCmd))

	execResp, err := cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd: []string{"bash", "-lc", vllmCmd},
	})
	if err != nil {
		return fmt.Errorf("创建 vLLM 进程失败: %w", err)
	}
	if err := cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{}); err != nil {
		return fmt.Errorf("启动 vLLM 进程失败: %w", err)
	}

	return nil
}

func (s *TrainingTaskService) waitVllmReady(ctx context.Context, node computenode.ComputeNode, containerID string, timeout time.Duration, taskID uint) error {
	deadline := time.Now().Add(timeout)
	checkCmd := `python -c "import urllib.request; urllib.request.urlopen('http://127.0.0.1:8000/v1/models', timeout=3).read()"`
	var lastErr error
	var lastLogLine int

	for time.Now().Before(deadline) {
		// 获取最新的容器日志并追加到任务日志
		logs, err := trainingDockerService.GetContainerLogs(ctx, &node, containerID, "100")
		if err == nil && logs != "" {
			lines := strings.Split(logs, "\n")
			if len(lines) > lastLogLine {
				newLines := lines[lastLogLine:]
				for _, line := range newLines {
					line = strings.TrimSpace(line)
					if line != "" {
						_ = s.appendTaskLog(taskID, line)
					}
				}
				lastLogLine = len(lines)
			}
		}

		err = s.execInContainerAndWait(ctx, node, containerID, checkCmd)
		if err == nil {
			return nil
		}
		lastErr = err
		time.Sleep(3 * time.Second)
	}

	if lastErr != nil {
		return fmt.Errorf("vLLM 在 %s 内未就绪: %w", timeout.String(), lastErr)
	}
	return fmt.Errorf("vLLM 在 %s 内未就绪", timeout.String())
}

func (s *TrainingTaskService) execInContainerAndWait(ctx context.Context, node computenode.ComputeNode, containerID string, cmd string) error {
	cli, err := trainingDockerService.CreateDockerClient(&node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %w", err)
	}
	defer cli.Close()

	execResp, err := cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd: []string{"bash", "-lc", cmd},
	})
	if err != nil {
		return fmt.Errorf("创建容器执行命令失败: %w", err)
	}

	if err := cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{}); err != nil {
		return fmt.Errorf("启动容器执行命令失败: %w", err)
	}

	for {
		info, err := cli.ContainerExecInspect(ctx, execResp.ID)
		if err != nil {
			return fmt.Errorf("获取容器执行状态失败: %w", err)
		}
		if !info.Running {
			if info.ExitCode != 0 {
				return fmt.Errorf("容器命令退出码=%d", info.ExitCode)
			}
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
}

func (s *TrainingTaskService) reconcileServingTaskStatus(ctx context.Context, task *modeltrainingModel.TrainingTask) {
	if task == nil || task.Status != "serving" {
		return
	}
	if task.NodeId == nil || task.ContainerId == nil || strings.TrimSpace(*task.ContainerId) == "" {
		s.markServingTaskFailed(task, "serving 状态缺少容器信息")
		return
	}

	var node computenode.ComputeNode
	if err := global.GVA_DB.Where("id = ?", *task.NodeId).First(&node).Error; err != nil {
		s.markServingTaskFailed(task, fmt.Sprintf("读取节点失败: %v", err))
		return
	}

	cli, err := trainingDockerService.CreateDockerClient(&node)
	if err != nil {
		s.markServingTaskFailed(task, fmt.Sprintf("创建 docker 客户端失败: %v", err))
		return
	}
	defer cli.Close()

	containerID := strings.TrimSpace(*task.ContainerId)
	inspect, err := cli.ContainerInspect(ctx, containerID)
	if err != nil {
		s.markServingTaskFailed(task, fmt.Sprintf("检查容器状态失败: %v", err))
		return
	}
	if inspect.State == nil || !inspect.State.Running {
		s.markServingTaskFailed(task, "vLLM 容器已停止")
	}
}

func (s *TrainingTaskService) markServingTaskFailed(task *modeltrainingModel.TrainingTask, reason string) {
	if task == nil || task.Status != "serving" {
		return
	}
	now := time.Now()
	task.Status = "failed"
	task.EndTime = &now
	_ = global.GVA_DB.Model(&modeltrainingModel.TrainingTask{}).
		Where("id = ? AND status = ?", task.ID, "serving").
		Updates(map[string]any{
			"status":   "failed",
			"end_time": now,
		}).Error
	_ = s.appendTaskLog(task.ID, fmt.Sprintf("检测到服务异常退出，任务状态自动更新为 failed: %s", reason))
}

func resolveVllmBaseModel(baseModel string) string {
	baseModel = strings.TrimSpace(baseModel)
	if baseModel == "" {
		return "Qwen/Qwen3-1.7B"
	}
	if strings.Contains(baseModel, "/") {
		return baseModel
	}
	if strings.HasPrefix(baseModel, "Qwen") {
		return "Qwen/" + baseModel
	}
	return baseModel
}

// StartService 启动 vLLM 推理服务
func (s *TrainingTaskService) StartService(ctx context.Context, taskID uint, userID uint, isAdmin bool) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	// 权限检查
	if !isAdmin && (task.UserId == nil || *task.UserId != userID) {
		return fmt.Errorf("无权操作此训练任务")
	}

	// 状态检查：只有 completed 状态才能启动服务
	if task.Status != "completed" {
		return fmt.Errorf("只有已完成的任务才能启动服务，当前状态: %s", task.Status)
	}

	if task.CheckpointPath == nil || *task.CheckpointPath == "" {
		return fmt.Errorf("任务缺少 checkpoint 路径，无法启动服务")
	}

	if task.ContainerId == nil || *task.ContainerId == "" {
		return fmt.Errorf("任务容器不存在")
	}

	if task.NodeId == nil {
		return fmt.Errorf("任务节点不存在")
	}

	node, err := s.getTrainingNode()
	if err != nil {
		return fmt.Errorf("获取算力节点失败: %v", err)
	}

	_ = s.appendTaskLog(taskID, "正在启动 vLLM 推理服务...")

	if err := s.startVllmServerInContainer(ctx, *node, *task.ContainerId, task, *task.CheckpointPath); err != nil {
		_ = s.appendTaskLog(taskID, fmt.Sprintf("启动 vLLM 失败: %v", err))
		return fmt.Errorf("启动 vLLM 失败: %v", err)
	}

	if err := s.waitVllmReady(ctx, *node, *task.ContainerId, 120*time.Second, taskID); err != nil {
		_ = s.appendTaskLog(taskID, fmt.Sprintf("vLLM 健康检查失败: %v", err))
		return fmt.Errorf("vLLM 健康检查失败: %v", err)
	}

	_ = s.appendTaskLog(taskID, "vLLM 推理服务已启动，可通过模型测试功能进行对话")

	return global.GVA_DB.Model(&task).Updates(map[string]interface{}{
		"status": "serving",
	}).Error
}

// StopService 停止 vLLM 推理服务
func (s *TrainingTaskService) StopService(ctx context.Context, taskID uint, userID uint, isAdmin bool) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	// 权限检查
	if !isAdmin && (task.UserId == nil || *task.UserId != userID) {
		return fmt.Errorf("无权操作此训练任务")
	}

	// 状态检查：只有 serving 状态才能停止服务
	if task.Status != "serving" {
		return fmt.Errorf("只有服务中的任务才能停止服务，当前状态: %s", task.Status)
	}

	if task.ContainerId == nil || *task.ContainerId == "" {
		return fmt.Errorf("任务容器不存在")
	}

	node, err := s.getTrainingNode()
	if err != nil {
		return fmt.Errorf("获取算力节点失败: %v", err)
	}

	_ = s.appendTaskLog(taskID, "正在停止 vLLM 推理服务...")

	// 在容器内停止 vLLM 进程
	cli, err := trainingDockerService.CreateDockerClient(node)
	if err != nil {
		return fmt.Errorf("创建 Docker 客户端失败: %v", err)
	}
	defer cli.Close()

	// 查找并杀死 vLLM 进程
	killCmd := "pkill -f 'vllm.entrypoints.openai.api_server' || true"
	execResp, err := cli.ContainerExecCreate(ctx, *task.ContainerId, container.ExecOptions{
		Cmd: []string{"bash", "-lc", killCmd},
	})
	if err == nil {
		_ = cli.ContainerExecStart(ctx, execResp.ID, container.ExecStartOptions{})
	}

	_ = s.appendTaskLog(taskID, "vLLM 推理服务已停止")

	return global.GVA_DB.Model(&task).Updates(map[string]interface{}{
		"status": "completed",
	}).Error
}

// MarkCompleted 手动标记训练完成
func (s *TrainingTaskService) MarkCompleted(ctx context.Context, taskID uint, userID uint, isAdmin bool) error {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("训练任务不存在")
	}

	// 权限检查
	if !isAdmin && (task.UserId == nil || *task.UserId != userID) {
		return fmt.Errorf("无权操作此训练任务")
	}

	// 状态检查：只有 running 状态才能手动标记完成
	if task.Status != "running" {
		return fmt.Errorf("只有训练中的任务才能标记完成，当前状态: %s", task.Status)
	}

	if task.ContainerId == nil || *task.ContainerId == "" {
		return fmt.Errorf("任务容器不存在")
	}

	node, err := s.getTrainingNode()
	if err != nil {
		return fmt.Errorf("获取算力节点失败: %v", err)
	}

	// 尝试从容器日志中提取 checkpoint
	checkpointPath, err := s.extractLastModelCheckpoint(ctx, *node, *task.ContainerId)
	if err != nil {
		_ = s.appendTaskLog(taskID, fmt.Sprintf("手动标记完成，但未能提取 checkpoint: %v", err))
	} else {
		_ = s.appendTaskLog(taskID, fmt.Sprintf("手动标记完成，识别到 checkpoint: %s", checkpointPath))
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":   "completed",
		"progress": 100,
		"end_time": now,
	}
	if checkpointPath != "" {
		updates["checkpoint_path"] = checkpointPath
	}

	return global.GVA_DB.Model(&task).Updates(updates).Error
}

// ChatCompletion 调用 vLLM API 进行对话测试
func (s *TrainingTaskService) ChatCompletion(ctx context.Context, taskID uint, model string, messages []modeltrainingReq.ChatMessage) (map[string]interface{}, error) {
	var task modeltrainingModel.TrainingTask
	if err := global.GVA_DB.First(&task, taskID).Error; err != nil {
		return nil, fmt.Errorf("训练任务不存在")
	}

	if task.Status != "serving" {
		return nil, fmt.Errorf("任务未处于服务状态，当前状态: %s", task.Status)
	}

	if task.HostPort == nil || *task.HostPort == 0 {
		return nil, fmt.Errorf("任务服务端口未配置")
	}

	// 构建 vLLM API 请求
	vllmMessages := make([]map[string]string, len(messages))
	for i, msg := range messages {
		vllmMessages[i] = map[string]string{
			"role":    msg.Role,
			"content": msg.Content,
		}
	}

	reqBody := map[string]interface{}{
		"model":    model, // base 或 lora
		"messages": vllmMessages,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("构建请求失败: %w", err)
	}

	// 获取节点信息
	node, err := s.getTrainingNode()
	if err != nil {
		return nil, fmt.Errorf("获取训练节点失败: %w", err)
	}

	// 通过 Docker exec 在容器内调用 vLLM API
	cli, err := trainingDockerService.CreateDockerClient(node)
	if err != nil {
		return nil, fmt.Errorf("创建 Docker 客户端失败: %w", err)
	}
	defer cli.Close()

	containerID := *task.ContainerId
	vllmURL := "http://127.0.0.1:8000/v1/chat/completions"

	// 使用 heredoc 方式传递 JSON，避免 shell 转义问题
	curlCmd := fmt.Sprintf(`curl -s -S --max-time 30 -X POST '%s' -H 'Content-Type: application/json' -d @- << 'EOF'
%s
EOF`, vllmURL, string(reqBytes))

	global.GVA_LOG.Debug("执行模型测试命令", zap.String("cmd", curlCmd))

	execResp, err := cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd:          []string{"bash", "-lc", curlCmd},
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("创建执行命令失败: %w", err)
	}

	// 获取执行输出
	output, err := cli.ContainerExecAttach(ctx, execResp.ID, container.ExecStartOptions{})
	if err != nil {
		return nil, fmt.Errorf("执行命令失败: %w", err)
	}
	defer output.Close()

	// 读取输出
	outputBytes, err := io.ReadAll(output.Reader)
	if err != nil {
		return nil, fmt.Errorf("读取输出失败: %w", err)
	}

	result := string(outputBytes)
	// 去除 Docker 多路复用协议头 (8字节头 + 数据)
	// Docker 使用 [stream_type(1字节)][0(3字节)][size(4字节)] 格式
	result = stripDockerHeader(result)
	result = stripANSI(result)
	result = strings.TrimSpace(result)

	// 如果结果为空，可能是服务未启动
	if result == "" {
		return nil, fmt.Errorf("vLLM 服务无响应，请确认服务已启动")
	}

	global.GVA_LOG.Debug("模型测试原始响应", zap.String("result", result))

	// 解析 JSON 响应
	var vllmResp map[string]interface{}
	if err := json.Unmarshal([]byte(result), &vllmResp); err != nil {
		// 检查是否是 curl 错误
		if strings.Contains(result, "curl:") || strings.Contains(result, "Could not connect") {
			return nil, fmt.Errorf("vLLM 服务连接失败: %s", result)
		}
		return nil, fmt.Errorf("解析响应失败: %w, 原始响应: %s", err, result)
	}

	// 检查 vLLM 返回的错误
	if errMsg, ok := vllmResp["error"]; ok {
		return nil, fmt.Errorf("vLLM 返回错误: %v", errMsg)
	}

	return vllmResp, nil
}

func (s *TrainingTaskService) stopTrainingContainer(ctx context.Context, node *computenode.ComputeNode, containerID string) {
	if node == nil || containerID == "" {
		return
	}
	cli, err := trainingDockerService.CreateDockerClient(node)
	if err != nil {
		global.GVA_LOG.Warn("停止训练容器失败：创建Docker客户端失败", zap.String("containerID", containerID), zap.Error(err))
		return
	}
	defer cli.Close()
	timeout := 10
	if err := cli.ContainerStop(ctx, containerID, container.StopOptions{Timeout: &timeout}); err != nil {
		global.GVA_LOG.Warn("停止训练容器失败", zap.String("containerID", containerID), zap.Error(err))
	}
}

func ptrString(v string) *string {
	return &v
}

// stripANSI 去除 ANSI 控制字符
func stripANSI(s string) string {
	// 匹配 ANSI 转义序列
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b\][^\x07]*\x07|\x1b[()][AB012]`)
	return re.ReplaceAllString(s, "")
}

// stripDockerHeader 去除 Docker 多路复用协议头
// Docker 使用 8 字节头: [stream_type(1)][0(3)][size(4)] + 数据
func stripDockerHeader(s string) string {
	// 直接查找 JSON 起始位置
	jsonStart := strings.Index(s, "{")
	if jsonStart > 0 {
		return s[jsonStart:]
	}
	// 如果不是 JSON，尝试去除第一个 8 字节头
	if len(s) > 8 {
		// 检查是否是 Docker 头格式
		firstByte := s[0]
		if firstByte == 1 || firstByte == 2 { // stdout 或 stderr
			// 读取 size (big-endian uint32)
			if len(s) >= 8 {
				size := int(s[4])<<24 | int(s[5])<<16 | int(s[6])<<8 | int(s[7])
				if len(s) >= 8+size {
					return s[8 : 8+size]
				}
			}
		}
	}
	return s
}
