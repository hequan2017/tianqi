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

type ModelTestHistoryApi struct{}

// CreateTestHistory 创建模型测试历史
// @Tags ModelTestHistory
// @Summary 创建模型测试历史
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingReq.CreateModelTestReq true "创建模型测试历史"
// @Success 200 {object} response.Response{data=modeltraining.ModelTestHistory,msg=string} "创建成功"
// @Router /modeltraining/modelTest/createTestHistory [post]
func (api *ModelTestHistoryApi) CreateTestHistory(c *gin.Context) {
	ctx := c.Request.Context()

	var req modeltrainingReq.CreateModelTestReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	err = modelTestHistoryService.CreateTestHistory(ctx, &req, userIDPtr)
	if err != nil {
		global.GVA_LOG.Error("创建模型测试历史失败!", zap.Error(err))
		response.FailWithMessage("创建模型测试历史失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("保存成功", c)
}

// DeleteTestHistory 删除模型测试历史
// @Tags ModelTestHistory
// @Summary 删除模型测试历史
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "测试历史ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /modeltraining/modelTest/deleteTestHistory [delete]
func (api *ModelTestHistoryApi) DeleteTestHistory(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("测试历史ID不能为空", c)
		return
	}

	id, err := strconv.ParseUint(ID, 10, 64)
	if err != nil {
		response.FailWithMessage("测试历史ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	err = modelTestHistoryService.DeleteTestHistory(ctx, uint(id), userIDPtr, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("删除模型测试历史失败!", zap.Error(err))
		response.FailWithMessage("删除模型测试历史失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetTestHistoryList 获取模型测试历史列表
// @Tags ModelTestHistory
// @Summary 获取模型测试历史列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query modeltrainingReq.ModelTestHistorySearch true "获取模型测试历史列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /modeltraining/modelTest/getTestHistoryList [get]
func (api *ModelTestHistoryApi) GetTestHistoryList(c *gin.Context) {
	ctx := c.Request.Context()

	var search modeltrainingReq.ModelTestHistorySearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if search.TaskId == nil {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	list, total, err := modelTestHistoryService.GetTestHistoryList(ctx, search, userIDPtr, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("获取模型测试历史列表失败!", zap.Error(err))
		response.FailWithMessage("获取模型测试历史列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// ClearTestHistory 清空指定任务的测试历史
// @Tags ModelTestHistory
// @Summary 清空指定任务的测试历史
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param taskId query string true "训练任务ID"
// @Success 200 {object} response.Response{msg=string} "清空成功"
// @Router /modeltraining/modelTest/clearTestHistory [delete]
func (api *ModelTestHistoryApi) ClearTestHistory(c *gin.Context) {
	ctx := c.Request.Context()

	taskIdStr := c.Query("taskId")
	if taskIdStr == "" {
		response.FailWithMessage("训练任务ID不能为空", c)
		return
	}

	taskId, err := strconv.ParseUint(taskIdStr, 10, 64)
	if err != nil {
		response.FailWithMessage("训练任务ID格式错误", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	var userIDPtr *uint
	if userID != 0 {
		userIDPtr = &userID
	}

	err = modelTestHistoryService.ClearTestHistoryByTaskId(ctx, uint(taskId), userIDPtr, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("清空模型测试历史失败!", zap.Error(err))
		response.FailWithMessage("清空模型测试历史失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("清空成功", c)
}