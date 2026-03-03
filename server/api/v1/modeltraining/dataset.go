package modeltraining

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	modeltrainingModel "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining"
	modeltrainingReq "github.com/flipped-aurora/gin-vue-admin/server/model/modeltraining/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DatasetApi struct{}

// CreateDataset 创建数据集
// @Tags Dataset
// @Summary 创建数据集
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingModel.Dataset true "创建数据集"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /modeltraining/dataset/createDataset [post]
func (api *DatasetApi) CreateDataset(c *gin.Context) {
	ctx := c.Request.Context()

	var dataset modeltrainingModel.Dataset
	err := c.ShouldBindJSON(&dataset)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 自动设置当前用户ID
	userID := utils.GetUserID(c)
	if userID > 0 {
		dataset.UserId = &userID
	}

	err = datasetService.CreateDataset(ctx, &dataset)
	if err != nil {
		global.GVA_LOG.Error("创建数据集失败!", zap.Error(err))
		response.FailWithMessage("创建数据集失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(dataset, "创建成功", c)
}

// DeleteDataset 删除数据集
// @Tags Dataset
// @Summary 删除数据集
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "数据集ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /modeltraining/dataset/deleteDataset [delete]
func (api *DatasetApi) DeleteDataset(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("数据集ID不能为空", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err := datasetService.DeleteDataset(ctx, ID, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("删除数据集失败!", zap.Error(err))
		response.FailWithMessage("删除数据集失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// DeleteDatasetByIds 批量删除数据集
// @Tags Dataset
// @Summary 批量删除数据集
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param IDs query []string true "数据集ID列表"
// @Success 200 {object} response.Response{msg=string} "批量删除成功"
// @Router /modeltraining/dataset/deleteDatasetByIds [delete]
func (api *DatasetApi) DeleteDatasetByIds(c *gin.Context) {
	ctx := c.Request.Context()

	IDs := c.QueryArray("IDs[]")
	if len(IDs) == 0 {
		response.FailWithMessage("请选择要删除的数据集", c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err := datasetService.DeleteDatasetByIds(ctx, IDs, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("批量删除数据集失败!", zap.Error(err))
		response.FailWithMessage("批量删除数据集失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("批量删除成功", c)
}

// UpdateDataset 更新数据集
// @Tags Dataset
// @Summary 更新数据集
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingModel.Dataset true "更新数据集"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /modeltraining/dataset/updateDataset [put]
func (api *DatasetApi) UpdateDataset(c *gin.Context) {
	ctx := c.Request.Context()

	var dataset modeltrainingModel.Dataset
	err := c.ShouldBindJSON(&dataset)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = datasetService.UpdateDataset(ctx, dataset)
	if err != nil {
		global.GVA_LOG.Error("更新数据集失败!", zap.Error(err))
		response.FailWithMessage("更新数据集失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// FindDataset 查询数据集详情
// @Tags Dataset
// @Summary 查询数据集详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "数据集ID"
// @Success 200 {object} response.Response{data=modeltrainingModel.Dataset,msg=string} "查询成功"
// @Router /modeltraining/dataset/findDataset [get]
func (api *DatasetApi) FindDataset(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("数据集ID不能为空", c)
		return
	}

	result, err := datasetService.GetDataset(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("查询数据集失败!", zap.Error(err))
		response.FailWithMessage("查询数据集失败: "+err.Error(), c)
		return
	}
	response.OkWithData(result, c)
}

// GetDatasetList 分页获取数据集列表
// @Tags Dataset
// @Summary 分页获取数据集列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data query modeltrainingReq.DatasetSearch true "分页获取数据集列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /modeltraining/dataset/getDatasetList [get]
func (api *DatasetApi) GetDatasetList(c *gin.Context) {
	ctx := c.Request.Context()

	var search modeltrainingReq.DatasetSearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	list, total, err := datasetService.GetDatasetList(ctx, search, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("获取数据集列表失败!", zap.Error(err))
		response.FailWithMessage("获取数据集列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// GetDatasetDataSource 获取数据集数据源
// @Tags Dataset
// @Summary 获取数据集数据源
// @Accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{data=object,msg=string} "查询成功"
// @Router /modeltraining/dataset/getDatasetDataSource [get]
func (api *DatasetApi) GetDatasetDataSource(c *gin.Context) {
	ctx := c.Request.Context()

	dataSource, err := datasetService.GetDatasetDataSource(ctx)
	if err != nil {
		global.GVA_LOG.Error("查询数据源失败!", zap.Error(err))
		response.FailWithMessage("查询数据源失败: "+err.Error(), c)
		return
	}
	response.OkWithData(dataSource, c)
}

// CreateVersion 创建数据集版本
// @Tags Dataset
// @Summary 创建数据集版本
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body modeltrainingModel.DatasetVersion true "创建数据集版本"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /modeltraining/dataset/createVersion [post]
func (api *DatasetApi) CreateVersion(c *gin.Context) {
	ctx := c.Request.Context()

	var version modeltrainingModel.DatasetVersion
	err := c.ShouldBindJSON(&version)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = datasetVersionService.CreateVersion(ctx, &version)
	if err != nil {
		global.GVA_LOG.Error("创建版本失败!", zap.Error(err))
		response.FailWithMessage("创建版本失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// GetVersionList 获取数据集版本列表
// @Tags Dataset
// @Summary 获取数据集版本列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param datasetId query int true "数据集ID"
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /modeltraining/dataset/getVersionList [get]
func (api *DatasetApi) GetVersionList(c *gin.Context) {
	ctx := c.Request.Context()

	var search modeltrainingReq.DatasetVersionSearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := datasetVersionService.GetVersionList(ctx, search)
	if err != nil {
		global.GVA_LOG.Error("获取版本列表失败!", zap.Error(err))
		response.FailWithMessage("获取版本列表失败: "+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     search.Page,
		PageSize: search.PageSize,
	}, "获取成功", c)
}

// PublishDataset 发布数据集
// @Tags Dataset
// @Summary 发布数据集
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query uint true "数据集ID"
// @Success 200 {object} response.Response{msg=string} "发布成功"
// @Router /modeltraining/dataset/publishDataset [post]
func (api *DatasetApi) PublishDataset(c *gin.Context) {
	ctx := c.Request.Context()

	var req modeltrainingReq.PublishDatasetReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	authorityId := utils.GetUserAuthorityId(c)
	isAdmin := authorityId == 888

	err = datasetService.PublishDataset(ctx, req.ID, userID, isAdmin)
	if err != nil {
		global.GVA_LOG.Error("发布数据集失败!", zap.Error(err))
		response.FailWithMessage("发布数据集失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("发布成功", c)
}

// UploadDatasetFile 上传数据集文件
// @Tags Dataset
// @Summary 上传数据集文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param file formData file true "数据集文件"
// @Param datasetId formData int false "数据集ID（可选）"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "上传成功"
// @Router /modeltraining/dataset/uploadFile [post]
func (api *DatasetApi) UploadDatasetFile(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jsonl": true, ".xls": true, ".xlsx": true}
	if !allowedExts[ext] {
		response.FailWithMessage("不支持的文件格式，仅支持 jsonl、xls、xlsx", c)
		return
	}

	// 验证文件大小 (最大200MB)
	if header.Size > 200*1024*1024 {
		response.FailWithMessage("文件大小不能超过200MB", c)
		return
	}

	// 构建存储路径: uploads/file/dataset/
	uploadPath := filepath.Join("uploads", "file", "dataset")

	// 确保目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		global.GVA_LOG.Error("创建目录失败!", zap.Error(err))
		response.FailWithMessage("创建目录失败", c)
		return
	}

	// 生成唯一文件名: 时间戳_原文件名
	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s_%s", timestamp, header.Filename)
	filePath := filepath.Join(uploadPath, newFileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		global.GVA_LOG.Error("创建文件失败!", zap.Error(err))
		response.FailWithMessage("创建文件失败", c)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		global.GVA_LOG.Error("保存文件失败!", zap.Error(err))
		response.FailWithMessage("保存文件失败", c)
		return
	}

	// 如果有数据集ID，自动创建一个新版本（每次上传一个版本）
	datasetIdStr := c.DefaultPostForm("datasetId", "")
	if datasetIdStr != "" {
		datasetId, _ := strconv.ParseUint(datasetIdStr, 10, 64)
		if datasetId > 0 {
			versionNo, versionErr := getNextDatasetVersion(uint(datasetId))
			if versionErr != nil {
				global.GVA_LOG.Error("获取下一个版本号失败", zap.Error(versionErr))
				response.FailWithMessage("获取版本号失败", c)
				return
			}

			version := modeltrainingModel.DatasetVersion{
				DatasetId:   uint(datasetId),
				Version:     versionNo,
				Status:      "success",
				FileName:    header.Filename,
				FilePath:    filePath,
				StoragePath: filePath,
				FileSize:    header.Size,
				DataCount:   0,
			}
			if err := global.GVA_DB.Create(&version).Error; err != nil {
				global.GVA_LOG.Error("创建数据集版本失败!", zap.Error(err))
				response.FailWithMessage("创建数据集版本失败", c)
				return
			}

			// 更新数据集的最新版本和导入状态
			if err := global.GVA_DB.Model(&modeltrainingModel.Dataset{}).Where("id = ?", datasetId).Updates(map[string]interface{}{
				"storage_path":   filePath,
				"latest_version": versionNo,
				"import_status":  "success",
			}).Error; err != nil {
				global.GVA_LOG.Error("更新数据集状态失败!", zap.Error(err))
				response.FailWithMessage("更新数据集状态失败", c)
				return
			}
		}
	}

	response.OkWithDetailed(map[string]interface{}{
		"path": filePath,
		"name": header.Filename,
		"size": header.Size,
	}, "上传成功", c)
}

func getNextDatasetVersion(datasetID uint) (string, error) {
	var versions []modeltrainingModel.DatasetVersion
	if err := global.GVA_DB.Where("dataset_id = ?", datasetID).Find(&versions).Error; err != nil {
		return "", err
	}
	maxVer := 0
	for _, v := range versions {
		ver := strings.TrimSpace(v.Version)
		ver = strings.TrimPrefix(strings.ToUpper(ver), "V")
		if n, err := strconv.Atoi(ver); err == nil && n > maxVer {
			maxVer = n
		}
	}
	return fmt.Sprintf("V%d", maxVer+1), nil
}

// saveUploadedFile 保存上传的文件
func saveUploadedFile(c *gin.Context, file *multipart.FileHeader, uploadPath string) (string, error) {
	// 确保目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s_%s%s", timestamp, file.Filename[:len(file.Filename)-len(ext)], ext)
	filePath := filepath.Join(uploadPath, newFileName)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		return "", fmt.Errorf("保存文件失败: %v", err)
	}

	return filePath, nil
}

// UploadVersionFile 上传版本文件
// @Tags Dataset
// @Summary 上传版本文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce application/json
// @Param file formData file true "版本文件"
// @Param versionId formData int true "版本ID"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "上传成功"
// @Router /modeltraining/dataset/uploadVersionFile [post]
func (api *DatasetApi) UploadVersionFile(c *gin.Context) {
	// 获取版本ID
	versionIdStr := c.DefaultPostForm("versionId", "0")
	if versionIdStr == "0" {
		response.FailWithMessage("版本ID不能为空", c)
		return
	}

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("接收文件失败!", zap.Error(err))
		response.FailWithMessage("接收文件失败", c)
		return
	}
	defer file.Close()

	// 验证文件类型
	ext := filepath.Ext(header.Filename)
	allowedExts := map[string]bool{".jsonl": true, ".xls": true, ".xlsx": true}
	if !allowedExts[ext] {
		response.FailWithMessage("不支持的文件格式，仅支持 jsonl、xls、xlsx", c)
		return
	}

	// 验证文件大小 (最大200MB)
	if header.Size > 200*1024*1024 {
		response.FailWithMessage("文件大小不能超过200MB", c)
		return
	}

	// 构建存储路径: uploads/file/dataset/version/{versionId}/
	uploadPath := filepath.Join("uploads", "file", "dataset", "version", versionIdStr)

	// 确保目录存在
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		global.GVA_LOG.Error("创建目录失败!", zap.Error(err))
		response.FailWithMessage("创建目录失败", c)
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	newFileName := fmt.Sprintf("%s_%s", timestamp, header.Filename)
	filePath := filepath.Join(uploadPath, newFileName)

	// 创建目标文件
	dst, err := os.Create(filePath)
	if err != nil {
		global.GVA_LOG.Error("创建文件失败!", zap.Error(err))
		response.FailWithMessage("创建文件失败", c)
		return
	}
	defer dst.Close()

	// 复制文件内容
	if _, err := io.Copy(dst, file); err != nil {
		global.GVA_LOG.Error("保存文件失败!", zap.Error(err))
		response.FailWithMessage("保存文件失败", c)
		return
	}

	// 更新版本记录
	versionId, _ := strconv.ParseUint(versionIdStr, 10, 64)
	if versionId > 0 {
		err = global.GVA_DB.Model(&modeltrainingModel.DatasetVersion{}).Where("id = ?", versionId).Updates(map[string]interface{}{
			"file_name": header.Filename,
			"file_path": filePath,
			"file_size": header.Size,
			"status":    "success",
		}).Error
		if err != nil {
			global.GVA_LOG.Error("更新版本记录失败!", zap.Error(err))
			response.FailWithMessage("更新版本记录失败", c)
			return
		}
	}

	response.OkWithDetailed(map[string]interface{}{
		"path":      filePath,
		"name":      header.Filename,
		"size":      header.Size,
		"versionId": versionId,
	}, "上传成功", c)
}

// DeleteVersion 删除数据集版本
// @Tags Dataset
// @Summary 删除数据集版本
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param ID query string true "版本ID"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /modeltraining/dataset/deleteVersion [delete]
func (api *DatasetApi) DeleteVersion(c *gin.Context) {
	ctx := c.Request.Context()

	ID := c.Query("ID")
	if ID == "" {
		response.FailWithMessage("版本ID不能为空", c)
		return
	}

	err := datasetVersionService.DeleteVersion(ctx, ID)
	if err != nil {
		global.GVA_LOG.Error("删除版本失败!", zap.Error(err))
		response.FailWithMessage("删除版本失败: "+err.Error(), c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
