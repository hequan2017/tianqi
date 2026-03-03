package modeltraining

import (
	"strconv"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	modeltrainingReq "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TrainingTaskApi struct{}

// CreateTask 创建训练任务
// @Tags TrainingTask
// @Summary 创建训练任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingReq.CreateTrainingTaskReq true "创建训练任务"
// @Success 200 {object} response.Response{data=modeltrainingModel.TrainingTask,msg=string} "创建成功"
// @Router /modeltraining/trainingTask/createTask [post]
func (api *TrainingTaskApi) CreateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req modeltrainingReq.CreateTrainingTaskReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	if userID == 0 {
		response.FailWithMessage("用户未登录", c)
		return
	}

	task, err := trainingTaskService.CreateTask(ctx, &req, userID)
	if err != nil {
		global.GVA_LOG.Error("创建训练任务失败!", zap.Error(err))
		response.FailWithMessage("创建训练任务失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(task, "创建成功", c)
}

// DeleteTask 删除训练任务
// @Tags TrainingTask
// @Summary 删除训练任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /modeltraining/trainingTask/deleteTask [delete]
func (api *TrainingTaskApi) DeleteTask(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err := trainingTaskService.DeleteTask(ctx, ID, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("删除训练任务失败!", zap.Error(err))
		response.FailWithMessage("删除训练任务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteTaskByIds 批量删除训练任务
// @Tags TrainingTask
// @Summary 批量删除训练任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param IDs query []string true "训练任务ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /modeltraining/trainingTask/deleteTaskByIds [delete]
func (api *TrainingTaskApi) DeleteTaskByIds(c *gin.Context) {
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	if len(IDs) == 0 {
		response.FailWithMessage("请选择要删除的训练任务", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err := trainingTaskService.DeleteTaskByIds(ctx, IDs, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("批量删除训练任务失败!", zap.Error(err))
		response.FailWithMessage("批量删除训练任务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateTask 更新训练任务
// @Tags TrainingTask
// @Summary 更新训练任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingReq.UpdateTrainingTaskReq true "更新训练任务"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /modeltraining/trainingTask/updateTask [put]
func (api *TrainingTaskApi) UpdateTask(c *gin.Context) {
	ctx := c.Request.Context()

	var req modeltrainingReq.UpdateTrainingTaskReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = trainingTaskService.UpdateTask(ctx, &req)
	if err != nil {
		global.GVA_LOG.Error("更新训练任务失败!", zap.Error(err))
		response.FailWithMessage("更新训练任务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindTask 查询训练任务详情
// @Tags TrainingTask
// @Summary 查询训练任务详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "训练任务ID"
// @Success 200 {object} response.Response{data=modeltrainingModel.TrainingTask,msg=string} "查询成功"
// @Router /modeltraining/trainingTask/findTask [get]
func (api *TrainingTaskApi) FindTask(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	result, err := trainingTaskService.GetTask(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询训练任务失败!", zap.Error(err))
		response.FailWithMessage("查询训练任务失败: "+err.Error(), c)
		return
	}

	// 获取训练参数
	taskId, _ := strconv.ParseUint(ID, 10, 64)
	param, _ := trainingParamService.GetParamByTaskId(ctx, uint(taskId))

	response.OkWithDetailed(gin.H{
		"task":  result,
		"param": param,
	}, "查询成功", c)
}

// GetTaskList 分页获取训练任务列表
// @Tags TrainingTask
// @Summary 分页获取训练任务列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query modeltrainingReq.TrainingTaskSearch true "分页获取训练任务列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /modeltraining/trainingTask/getTaskList [get]
func (api *TrainingTaskApi) GetTaskList(c *gin.Context) {
	ctx := c.Request.Context()

	var search modeltrainingReq.TrainingTaskSearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	list, total, err := trainingTaskService.GetTaskList(ctx, search, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("获取训练任务列表失败!", zap.Error(err))
		response.FailWithMessage("获取训练任务列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// StartTask 启动训练任务
// @Tags TrainingTask
// @Summary 启动训练任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "启动成功"
// @Router /modeltraining/trainingTask/startTask [post]
func (api *TrainingTaskApi) StartTask(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Query("ID")
	if IDStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err = trainingTaskService.StartTask(ctx, uint(ID), userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("启动训练任务失败!", zap.Error(err))
		response.FailWithMessage("启动训练任务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("启动成功", c)
}

// StopTask 停止训练任务
// @Tags TrainingTask
// @Summary 停止训练任务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "停止成功"
// @Router /modeltraining/trainingTask/stopTask [post]
func (api *TrainingTaskApi) StopTask(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Query("ID")
	if IDStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err = trainingTaskService.StopTask(ctx, uint(ID), userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("停止训练任务失败!", zap.Error(err))
		response.FailWithMessage("停止训练任务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("停止成功", c)
}

// GetTaskLogs 获取训练日志
// @Tags TrainingTask
// @Summary 获取训练日志
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "训练任务ID"
// @Param tail query string false "日志行数"
// @Success 200 {object} response.Response{data=string,msg=string} "获取成功"
// @Router /modeltraining/trainingTask/getTaskLogs [get]
func (api *TrainingTaskApi) GetTaskLogs(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Query("ID")
	if IDStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	tail := c.DefaultQuery("tail", "100")

	logs, err := trainingTaskService.GetTaskLogs(ctx, uint(ID), tail)
	if err != nil {
		global.GVA_LOG.Error("获取训练日志失败!", zap.Error(err))
		response.FailWithMessage("获取训练日志失败: "+err.Error(), c)
		return
	}
	response.OkWithData(logs, c)
}

// GetTaskDataSource 获取训练任务数据源
// @Tags TrainingTask
// @Summary 获取训练任务数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /modeltraining/trainingTask/getTaskDataSource [get]
func (api *TrainingTaskApi) GetTaskDataSource(c *gin.Context) {
	ctx := c.Request.Context()

	dataSource, err := trainingTaskService.GetTaskDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询数据源失败!", zap.Error(err))
		response.FailWithMessage("查询数据源失败: "+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// GetDefaultParams 获取默认训练参数
// @Tags TrainingTask
// @Summary 获取默认训练参数
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "获取成功"
// @Router /modeltraining/trainingTask/getDefaultParams [get]
func (api *TrainingTaskApi) GetDefaultParams(c *gin.Context) {
	ctx := c.Request.Context()

	params := trainingTaskService.GetDefaultParams(ctx)
	response.OkWithData(params, c)
}

// ChatCompletion 模型对话测试
// @Tags TrainingTask
// @Summary 模型对话测试
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingReq.ChatCompletionReq true "对话请求"
// @Success 200 {object} response.Response{data=object,msg=string} "请求成功"
// @Router /modeltraining/trainingTask/chatCompletion [post]
func (api *TrainingTaskApi) ChatCompletion(c *gin.Context) {
	ctx := c.Request.Context()

	var req modeltrainingReq.ChatCompletionReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if req.ID == 0 {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	result, err := trainingTaskService.ChatCompletion(ctx, req.ID, req.Model, req.Messages)
	if err != nil {
		global.GVA_LOG.Error("模型对话测试失败!", zap.Error(err))
		response.FailWithMessage("模型对话测试失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// StartService 启动推理服务
// @Tags TrainingTask
// @Summary 启动 vLLM 推理服务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "启动成功"
// @Router /modeltraining/trainingTask/startService [post]
func (api *TrainingTaskApi) StartService(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Query("ID")
	if IDStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err = trainingTaskService.StartService(ctx, uint(ID), userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("启动推理服务失败!", zap.Error(err))
		response.FailWithMessage("启动推理服务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("推理服务已启动", c)
}

// StopService 停止推理服务
// @Tags TrainingTask
// @Summary 停止 vLLM 推理服务
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "停止成功"
// @Router /modeltraining/trainingTask/stopService [post]
func (api *TrainingTaskApi) StopService(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Query("ID")
	if IDStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err = trainingTaskService.StopService(ctx, uint(ID), userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("停止推理服务失败!", zap.Error(err))
		response.FailWithMessage("停止推理服务失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("推理服务已停止", c)
}

// MarkCompleted 手动标记训练完成
// @Tags TrainingTask
// @Summary 手动标记训练完成
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "标记成功"
// @Router /modeltraining/trainingTask/markCompleted [post]
func (api *TrainingTaskApi) MarkCompleted(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Query("ID")
	if IDStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	ID, err := strconv.ParseUint(IDStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err = trainingTaskService.MarkCompleted(ctx, uint(ID), userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("标记训练完成失败!", zap.Error(err))
		response.FailWithMessage("标记训练完成失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("已标记为训练完成", c)
}