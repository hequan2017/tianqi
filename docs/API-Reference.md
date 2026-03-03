# 🔧 API 文档

本页面介绍天启算力管理平台的后端 API 接口。

## API 概述

- **基础路径**：`/api/v1` 或直接访问（视配置而定）
- **认证方式**：JWT Token（通过 `x-token` 请求头传递）
- **响应格式**：JSON
- **Swagger 文档**：`http://localhost:8888/swagger/index.html`

## 通用响应格式

### 成功响应

```json
{
  "code": 0,
  "data": { ... },
  "msg": "操作成功"
}
```

### 错误响应

```json
{
  "code": 7,
  "data": null,
  "msg": "错误信息"
}
```

## 认证相关

### 用户登录

```http
POST /base/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456",
  "captcha": "xxxx",
  "captchaId": "xxx"
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "user": {
      "ID": 1,
      "username": "admin",
      "nickName": "管理员",
      "authorityId": 888
    },
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "expiresAt": 1735516800
  },
  "msg": "登录成功"
}
```

### 获取验证码

```http
POST /base/captcha
```

---

## 镜像库管理

### 获取镜像库列表

```http
GET /imageregistry/getImageRegistryList
x-token: <your-token>

Query Parameters:
- page: 页码（默认 1）
- pageSize: 每页数量（默认 10）
- name: 名称搜索（可选）
- isListed: 是否上架（可选）
```

### 创建镜像库

```http
POST /imageregistry/createImageRegistry
x-token: <your-token>
Content-Type: application/json

{
  "name": "Ubuntu 22.04 CUDA",
  "address": "nvidia/cuda:12.0-runtime-ubuntu22.04",
  "description": "CUDA 12.0 运行时环境",
  "source": "Docker Hub",
  "supportMemorySplit": true,
  "isListed": true,
  "remark": ""
}
```

### 更新镜像库

```http
PUT /imageregistry/updateImageRegistry
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1,
  "name": "Ubuntu 22.04 CUDA",
  "address": "nvidia/cuda:12.0-runtime-ubuntu22.04",
  ...
}
```

### 删除镜像库

```http
DELETE /imageregistry/deleteImageRegistry
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1
}
```

---

## 算力节点管理

### 获取算力节点列表

```http
GET /computenode/getComputeNodeList
x-token: <your-token>

Query Parameters:
- page: 页码
- pageSize: 每页数量
- name: 名称搜索
- isListed: 是否上架
```

### 创建算力节点

```http
POST /computenode/createComputeNode
x-token: <your-token>
Content-Type: application/json

{
  "name": "GPU-Node-01",
  "region": "华东",
  "cpu": "Intel Xeon Gold 6248",
  "memory": "256GB",
  "systemDisk": "500GB",
  "dataDisk": "2TB",
  "publicIp": "1.2.3.4",
  "privateIp": "192.168.1.100",
  "sshPort": 22,
  "username": "root",
  "password": "******",
  "gpuName": "NVIDIA RTX 4090",
  "gpuCount": 8,
  "memoryCapacity": 24,
  "hamiCore": "/root/HAMi-core/build",
  "dockerAddress": "tcp://192.168.1.100:2376",
  "useTls": true,
  "caCert": "-----BEGIN CERTIFICATE-----...",
  "clientCert": "-----BEGIN CERTIFICATE-----...",
  "clientKey": "-----BEGIN RSA PRIVATE KEY-----...",
  "isListed": true
}
```

### 测试 Docker 连接

```http
POST /computenode/testDockerConnection
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1
}
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "status": "connected",
    "message": "Docker 连接成功"
  },
  "msg": "测试完成"
}
```

---

## 产品规格管理

### 获取产品规格列表

```http
GET /product/getProductSpecList
x-token: <your-token>

Query Parameters:
- page: 页码
- pageSize: 每页数量
- name: 名称搜索
- gpuModel: GPU 型号
- isListed: 是否上架
```

### 创建产品规格

```http
POST /product/createProductSpec
x-token: <your-token>
Content-Type: application/json

{
  "name": "4090-2卡-16核-64G",
  "gpuModel": "NVIDIA RTX 4090",
  "gpuCount": 2,
  "memoryCapacity": 48,
  "supportMemorySplit": false,
  "cpuCores": 16,
  "memoryGb": 64,
  "systemDiskGb": 100,
  "dataDiskGb": 500,
  "pricePerHour": 10.0,
  "isListed": true
}
```

---

## 实例管理

### 获取实例列表

```http
GET /instance/getInstanceList
x-token: <your-token>

Query Parameters:
- page: 页码
- pageSize: 每页数量
- name: 实例名称
- status: 容器状态（running/exited/creating/failed）
```

> 注：普通用户只能看到自己创建的实例，管理员可以看到所有实例。

### 创建实例

```http
POST /instance/createInstance
x-token: <your-token>
Content-Type: application/json

{
  "imageId": 1,
  "productSpecId": 1,
  "computeNodeId": 1,
  "name": "my-training-instance",
  "remark": "深度学习训练"
}
```

### 获取匹配的算力节点

根据产品规格查询可用的算力节点：

```http
POST /instance/getMatchedNodes
x-token: <your-token>
Content-Type: application/json

{
  "productSpecId": 1
}
```

### 启动实例

```http
POST /instance/startInstance
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1
}
```

### 停止实例

```http
POST /instance/stopInstance
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1
}
```

### 重启实例

```http
POST /instance/restartInstance
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1
}
```

### 删除实例

```http
DELETE /instance/deleteInstance
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1
}
```

> 注：删除实例会同时删除容器及其挂载的数据卷。

### 获取容器日志

```http
GET /instance/getContainerLogs
x-token: <your-token>

Query Parameters:
- id: 实例 ID
- tail: 返回行数（默认 100）
```

### 获取容器统计信息

```http
GET /instance/getContainerStats
x-token: <your-token>

Query Parameters:
- id: 实例 ID
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "cpuPercent": 25.5,
    "memoryUsage": 4294967296,
    "memoryLimit": 68719476736,
    "memoryPercent": 6.25,
    "networkRx": 1048576,
    "networkTx": 524288,
    "blockRead": 10485760,
    "blockWrite": 5242880,
    "pids": 42
  },
  "msg": "获取成功"
}
```

### Web 终端

```http
WebSocket /instance/terminal
x-token: <your-token>

Query Parameters:
- id: 实例 ID
```

---

## 数据集管理

### 获取数据集列表

```http
GET /modeltraining/dataset/getDatasetList
x-token: <your-token>

Query Parameters:
- page: 页码
- pageSize: 每页数量
- name: 数据集名称（模糊搜索）
- type: 数据集类型（training/evaluation）
- publishStatus: 发布状态
```

### 创建数据集

```http
POST /modeltraining/dataset/createDataset
x-token: <your-token>
Content-Type: application/json

{
  "name": "我的训练数据集",
  "type": "training",
  "format": "文本生成",
  "trainMethod": "SFT",
  "description": "用于模型微调的训练数据"
}
```

### 更新数据集

```http
PUT /modeltraining/dataset/updateDataset
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1,
  "name": "更新后的名称",
  ...
}
```

### 删除数据集

```http
DELETE /modeltraining/dataset/deleteDataset?ID=1
x-token: <your-token>
```

### 发布数据集

```http
POST /modeltraining/dataset/publishDataset?ID=1
x-token: <your-token>
```

### 上传数据集文件

```http
POST /modeltraining/dataset/uploadFile
x-token: <your-token>
Content-Type: multipart/form-data

Form Data:
- file: 数据集文件（支持 .jsonl, .xls, .xlsx，最大 200MB）
- datasetId: 数据集ID（可选，自动创建新版本）
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "path": "uploads/file/dataset/20250101120000_data.jsonl",
    "name": "data.jsonl",
    "size": 1048576
  },
  "msg": "上传成功"
}
```

### 获取数据集版本列表

```http
GET /modeltraining/dataset/getVersionList
x-token: <your-token>

Query Parameters:
- datasetId: 数据集ID
- page: 页码
- pageSize: 每页数量
```

### 删除数据集版本

```http
DELETE /modeltraining/dataset/deleteVersion?ID=1
x-token: <your-token>
```

---

## 训练任务管理

### 获取训练任务列表

```http
GET /modeltraining/trainingTask/getTaskList
x-token: <your-token>

Query Parameters:
- page: 页码
- pageSize: 每页数量
- name: 任务名称（模糊搜索）
- status: 任务状态（pending/running/completed/failed）
```

> 注：普通用户只能看到自己创建的任务，管理员可以看到所有任务。

### 创建训练任务

```http
POST /modeltraining/trainingTask/createTask
x-token: <your-token>
Content-Type: application/json

{
  "name": "我的训练任务",
  "baseModel": "Qwen3-1.7B",
  "trainMethod": "SFT",
  "trainType": "efficient",
  "trainDatasetId": 1,
  "trainVersionId": 1,
  "valSplitRatio": 0.1,
  "modelName": "my-fine-tuned-model",
  "outputCount": 5,
  "checkpointInterval": 500,
  "checkpointUnit": "step",
  "nodeId": 1,
  "remark": "备注信息",
  "trainingParam": {
    "batchSize": 4,
    "learningRate": 0.0001,
    "nEpochs": 3,
    "evalSteps": 100,
    "loraAlpha": 16,
    "loraDropout": 0.05,
    "loraRank": 8,
    "lrSchedulerType": "cosine",
    "maxLength": 2048,
    "warmupRatio": 0.1,
    "weightDecay": 0.01
  }
}
```

### 更新训练任务

```http
PUT /modeltraining/trainingTask/updateTask
x-token: <your-token>
Content-Type: application/json

{
  "ID": 1,
  "name": "更新后的任务名称",
  ...
}
```

### 删除训练任务

```http
DELETE /modeltraining/trainingTask/deleteTask?ID=1
x-token: <your-token>
```

### 查询训练任务详情

```http
GET /modeltraining/trainingTask/findTask?ID=1
x-token: <your-token>
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "task": {
      "ID": 1,
      "name": "我的训练任务",
      "taskId": "train_abc123",
      "status": "running",
      "progress": 45,
      ...
    },
    "param": {
      "batchSize": 4,
      "learningRate": 0.0001,
      ...
    }
  },
  "msg": "查询成功"
}
```

### 启动训练任务

```http
POST /modeltraining/trainingTask/startTask?ID=1
x-token: <your-token>
```

### 停止训练任务

```http
POST /modeltraining/trainingTask/stopTask?ID=1
x-token: <your-token>
```

### 获取训练日志

```http
GET /modeltraining/trainingTask/getTaskLogs
x-token: <your-token>

Query Parameters:
- ID: 训练任务ID
- tail: 返回行数（默认 100）
```

### 获取默认训练参数

```http
GET /modeltraining/trainingTask/getDefaultParams
x-token: <your-token>
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "batchSize": 4,
    "learningRate": 0.0001,
    "nEpochs": 3,
    "evalSteps": 100,
    "loraAlpha": 16,
    "loraDropout": 0.05,
    "loraRank": 8,
    "lrSchedulerType": "cosine",
    "maxLength": 2048,
    "warmupRatio": 0.1,
    "weightDecay": 0.01
  },
  "msg": "获取成功"
}
```

---

## 仪表盘

### 获取仪表盘统计数据

```http
GET /dashboard/stats
x-token: <your-token>
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "nodeCount": 5,
    "instanceCount": 20,
    "runningInstanceCount": 15,
    "datasetCount": 10,
    "trainingTaskCount": 8
  },
  "msg": "获取成功"
}
```

---

## SSH 跳板机

SSH 跳板机不通过 HTTP API 访问，而是通过 SSH 协议：

```bash
ssh -p 2026 username@server-ip
```

前端可通过以下 API 获取 SSH 连接信息：

```http
GET /instance/getSSHCommand
x-token: <your-token>

Query Parameters:
- id: 实例 ID
```

**响应**：

```json
{
  "code": 0,
  "data": {
    "command": "ssh -p 2026 admin@192.168.112.148"
  },
  "msg": "获取成功"
}
```

---

## 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 7 | 操作失败 |
| 401 | 未授权（Token 无效或过期） |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 500 | 服务器内部错误 |

---

## 完整 API 文档

完整的 API 文档请访问 Swagger UI：

```
http://localhost:8888/swagger/index.html
```

Swagger 文档提供：
- 所有 API 接口列表
- 请求/响应参数详情
- 在线测试功能

