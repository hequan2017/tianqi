import service from '@/utils/request'

/**
 * 数据集管理 API
 */

// 创建数据集
export const createDataset = (data) => {
  return service({
    url: '/modeltraining/dataset/createDataset',
    method: 'post',
    data
  })
}

// 删除数据集
export const deleteDataset = (params) => {
  return service({
    url: '/modeltraining/dataset/deleteDataset',
    method: 'delete',
    params
  })
}

// 批量删除数据集
export const deleteDatasetByIds = (params) => {
  return service({
    url: '/modeltraining/dataset/deleteDatasetByIds',
    method: 'delete',
    params
  })
}

// 更新数据集
export const updateDataset = (data) => {
  return service({
    url: '/modeltraining/dataset/updateDataset',
    method: 'put',
    data
  })
}

// 查询数据集详情
export const findDataset = (params) => {
  return service({
    url: '/modeltraining/dataset/findDataset',
    method: 'get',
    params
  })
}

// 分页获取数据集列表
export const getDatasetList = (params) => {
  return service({
    url: '/modeltraining/dataset/getDatasetList',
    method: 'get',
    params
  })
}

// 获取数据集数据源
export const getDatasetDataSource = () => {
  return service({
    url: '/modeltraining/dataset/getDatasetDataSource',
    method: 'get'
  })
}

// 创建数据集版本
export const createVersion = (data) => {
  return service({
    url: '/modeltraining/dataset/createVersion',
    method: 'post',
    data
  })
}

// 获取数据集版本列表
export const getVersionList = (params) => {
  return service({
    url: '/modeltraining/dataset/getVersionList',
    method: 'get',
    params
  })
}

// 发布数据集
export const publishDataset = (params) => {
  return service({
    url: '/modeltraining/dataset/publishDataset',
    method: 'post',
    params
  })
}

// 上传数据集文件
export const uploadDatasetFile = (formData, onProgress) => {
  return service({
    url: '/modeltraining/dataset/uploadFile',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: onProgress
  })
}

// 上传版本文件
export const uploadVersionFile = (formData, onProgress) => {
  return service({
    url: '/modeltraining/dataset/uploadVersionFile',
    method: 'post',
    data: formData,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    onUploadProgress: onProgress
  })
}

// 删除版本
export const deleteVersion = (params) => {
  return service({
    url: '/modeltraining/dataset/deleteVersion',
    method: 'delete',
    params
  })
}