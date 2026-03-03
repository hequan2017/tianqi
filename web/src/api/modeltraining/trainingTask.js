import service from '@/utils/request'

/**
 * 训练任务管理 API
 */

// 创建训练任务
export const createTask = (data) => {
  return service({
    url: '/modeltraining/trainingTask/createTask',
    method: 'post',
    data
  })
}

// 删除训练任务
export const deleteTask = (params) => {
  return service({
    url: '/modeltraining/trainingTask/deleteTask',
    method: 'delete',
    params
  })
}

// 批量删除训练任务
export const deleteTaskByIds = (params) => {
  return service({
    url: '/modeltraining/trainingTask/deleteTaskByIds',
    method: 'delete',
    params
  })
}

// 更新训练任务
export const updateTask = (data) => {
  return service({
    url: '/modeltraining/trainingTask/updateTask',
    method: 'put',
    data
  })
}

// 查询训练任务详情
export const findTask = (params) => {
  return service({
    url: '/modeltraining/trainingTask/findTask',
    method: 'get',
    params
  })
}

// 分页获取训练任务列表
export const getTaskList = (params) => {
  return service({
    url: '/modeltraining/trainingTask/getTaskList',
    method: 'get',
    params
  })
}

// 启动训练任务
export const startTask = (params) => {
  return service({
    url: '/modeltraining/trainingTask/startTask',
    method: 'post',
    params
  })
}

// 停止训练任务
export const stopTask = (params) => {
  return service({
    url: '/modeltraining/trainingTask/stopTask',
    method: 'post',
    params
  })
}

// 获取训练日志
export const getTaskLogs = (params) => {
  return service({
    url: '/modeltraining/trainingTask/getTaskLogs',
    method: 'get',
    params
  })
}

// 获取训练任务数据源
export const getTaskDataSource = () => {
  return service({
    url: '/modeltraining/trainingTask/getTaskDataSource',
    method: 'get'
  })
}

// 获取默认训练参数
export const getDefaultParams = () => {
  return service({
    url: '/modeltraining/trainingTask/getDefaultParams',
    method: 'get'
  })
}

// 模型对话测试
export const chatCompletion = (data) => {
  return service({
    url: '/modeltraining/trainingTask/chatCompletion',
    method: 'post',
    data
  })
}

// 启动推理服务
export const startService = (params) => {
  return service({
    url: '/modeltraining/trainingTask/startService',
    method: 'post',
    params
  })
}

// 停止推理服务
export const stopService = (params) => {
  return service({
    url: '/modeltraining/trainingTask/stopService',
    method: 'post',
    params
  })
}

// 手动标记训练完成
export const markCompleted = (params) => {
  return service({
    url: '/modeltraining/trainingTask/markCompleted',
    method: 'post',
    params
  })
}

// ========== 模型测试历史 API ==========

// 创建测试历史
export const createTestHistory = (data) => {
  return service({
    url: '/modeltraining/modelTest/createTestHistory',
    method: 'post',
    data
  })
}

// 删除测试历史
export const deleteTestHistory = (params) => {
  return service({
    url: '/modeltraining/modelTest/deleteTestHistory',
    method: 'delete',
    params
  })
}

// 获取测试历史列表
export const getTestHistoryList = (params) => {
  return service({
    url: '/modeltraining/modelTest/getTestHistoryList',
    method: 'get',
    params
  })
}

// 清空测试历史
export const clearTestHistory = (params) => {
  return service({
    url: '/modeltraining/modelTest/clearTestHistory',
    method: 'delete',
    params
  })
}