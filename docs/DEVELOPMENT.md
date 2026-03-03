# Docker GPU 算力管理平台 - 开发说明文档

## 一、项目概述

### 1.1 项目定位

**Docker GPU 算力管理平台** 是一个企业级的 GPU 容器化资源管理和调度系统，基于 gin-vue-admin (GVA) 框架开发。平台提供从资源管理到容器实例全生命周期的完整解决方案，并集成模型训练功能。

### 1.2 技术栈

| 层级 | 技术选型 | 版本 |
|------|----------|------|
| **后端框架** | Go + Gin | 1.23+ / 1.10.0 |
| **ORM** | GORM | 1.25.12 |
| **认证授权** | JWT + Casbin | 5.2.2 / 2.103.0 |
| **前端框架** | Vue 3 + Composition API | 3.5.7 |
| **UI 组件库** | Element Plus | 2.10.2 |
| **构建工具** | Vite | 6.2.3 |
| **状态管理** | Pinia | 2.2.2 |
| **数据库** | MySQL (主) | 5.7+ |
| **容器技术** | Docker API | v27.0.0 |

---

## 二、架构设计

### 2.1 分层架构

```
┌─────────────────────────────────────────────────────┐
│                    Router 层                         │
│              (路由定义、中间件配置)                    │
├─────────────────────────────────────────────────────┤
│                    API 层                            │
│         (请求处理、参数校验、响应封装)                 │
├─────────────────────────────────────────────────────┤
│                   Service 层                         │
│            (业务逻辑、数据操作)                       │
├─────────────────────────────────────────────────────┤
│                    Model 层                          │
│          (数据模型、请求/响应结构体)                   │
└─────────────────────────────────────────────────────┘
```

### 2.2 模块间调用关系

```
前端(Vue) → API请求 → Router → API → Service → Model → Database
                              ↓
                         中间件(JWT/Casbin)
```

### 2.3 Enter.go 组管理模式

```go
// service/enter.go
var ServiceGroupApp = new(ServiceGroup)

type ServiceGroup struct {
    InstanceServiceGroup      instance.ServiceGroup
    ComputenodeServiceGroup   computenode.ServiceGroup
    ModeltrainingServiceGroup modeltraining.ServiceGroup
    // ...
}

// 调用方式
service.ServiceGroupApp.InstanceServiceGroup.InstanceService.CreateInstance(ctx, inst)
```

---

## 三、核心功能模块详解

### 3.1 实例管理模块 (Instance)

#### 3.1.1 数据模型

**文件**: `server/model/instance/instance.go`

```go
type Instance struct {
    global.GVA_MODEL
    ImageId            *int64   // 镜像ID
    SpecId             *int64   // 产品规格ID
    UserId             *int64   // 用户ID
    NodeId             *int64   // 算力节点ID
    ContainerId        *string  // Docker容器ID
    ContainerName      *string  // 容器名称
    Name               *string  // 实例名称
    ContainerStatus    *string  // 容器状态
    CpuUsagePercent    *float64 // CPU使用率
    MemoryUsagePercent *float64 // 内存使用率
    GpuMemoryUsageRate *float64 // GPU显存使用率
}
```

#### 3.1.2 创建实例流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   前端请求    │ ──→ │  API 层校验   │ ──→ │  Service 层  │
└──────────────┘     └──────────────┘     └──────────────┘
                                                 │
                                                 ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  创建数据库   │ ←── │ 智能匹配节点  │ ←── │ 获取镜像/规格 │
│    记录      │     │              │     │    信息      │
└──────────────┘     └──────────────┘     └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  更新数据库   │ ←── │ 启动容器     │ ←── │ 创建Docker   │
│  容器信息    │     │              │     │   容器       │
└──────────────┘     └──────────────┘     └──────────────┘
```

**核心代码实现** (`server/service/instance/instance.go`):

```go
func (s *InstanceService) CreateInstance(ctx context.Context, inst *Instance) error {
    // 1. 获取镜像信息
    var image ImageRegistry
    global.GVA_DB.Where("id = ?", *inst.ImageId).First(&image)

    // 2. 获取产品规格
    var spec ProductSpec
    global.GVA_DB.Where("id = ?", *inst.SpecId).First(&spec)

    // 3. 获取算力节点
    var node ComputeNode
    global.GVA_DB.Where("id = ?", *inst.NodeId).First(&node)

    // 4. 创建数据库记录
    inst.ContainerStatus = "creating"
    global.GVA_DB.Create(inst)

    // 5. 构建容器配置
    config := dockerService.BuildContainerConfig(&image, &spec, &node, containerName)

    // 6. 创建Docker容器
    containerID, err := dockerService.CreateContainer(ctx, &node, config)

    // 7. 更新实例记录
    global.GVA_DB.Model(inst).Updates(map[string]interface{}{
        "container_id":     containerID,
        "container_status": "running",
    })
}
```

#### 3.1.3 Docker 容器创建流程

**文件**: `server/service/instance/docker.go`

```go
func (d *DockerService) CreateContainer(ctx context.Context, node *ComputeNode, config *ContainerConfig) (string, error) {
    // 1. 创建 Docker 客户端（支持 TLS）
    cli, err := d.CreateDockerClient(node)

    // 2. 构建容器配置
    containerConfig := &container.Config{
        Image: config.Image,
        Env:   []string{}, // 环境变量
    }

    // 3. 构建主机配置
    hostConfig := &container.HostConfig{
        NanoCPUs: config.CPUCores * 1e9,        // CPU 限制
        Memory:   config.MemoryGB * 1024^3,     // 内存限制
    }

    // 4. GPU 配置
    if config.GPUCount > 0 {
        hostConfig.DeviceRequests = []container.DeviceRequest{{
            Driver: "nvidia",
            Count:  int(config.GPUCount),
        }}
    }

    // 5. 显存切分支持
    if config.SupportMemorySplit {
        // 挂载 HAMi 库
        hostConfig.Mounts = append(hostConfig.Mounts, mount.Mount{
            Source: hamiCorePath,
            Target: "/libvgpu/build",
        })
        // 注入环境变量
        containerConfig.Env = append(containerConfig.Env,
            "LD_PRELOAD=/libvgpu/build/libvgpu.so",
            fmt.Sprintf("CUDA_DEVICE_MEMORY_LIMIT=%dg", config.MemoryCapacity),
        )
    }

    // 6. 创建并启动容器
    resp, _ := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, name)
    cli.ContainerStart(ctx, resp.ID, container.StartOptions{})

    return resp.ID, nil
}
```

#### 3.1.4 智能节点匹配算法

**文件**: `server/service/instance/instance.go:GetAvailableNodes`

```
输入: 产品规格ID
输出: 可用节点列表

算法流程:
1. 获取产品规格需求 (GPU型号/数量、显存、CPU、内存、磁盘)
2. 获取所有上架节点
3. 统计每个节点已占用资源
   - 遍历该节点所有实例
   - 累加 GPU/CPU/内存/磁盘使用量
   - 按卡分配显存（支持切分/不支持切分两种模式）
4. 筛选满足需求的节点
   - GPU型号匹配
   - 资源充足（考虑显存切分）
5. 返回可用节点列表
```

**显存分配逻辑**:

```go
// 支持显存切分
if supportMemorySplit {
    // 可以分配到任意有剩余空间的卡
    for _, card := range cards {
        if card.remaining >= requiredMemoryPerCard {
            card.usage += requiredMemoryPerCard
            break
        }
    }
}

// 不支持显存切分
else {
    // 需要完全未使用的卡
    for _, card := range cards {
        if card.usage == 0 {
            card.usage = fullCapacity
            break
        }
    }
}
```

#### 3.1.5 资源监控实现

**文件**: `server/service/instance/docker.go:GetContainerStats`

```go
func (d *DockerService) GetContainerStats(ctx context.Context, node *ComputeNode, containerID string) (*ContainerStats, error) {
    // 方式1: 通过 docker stats 命令
    stats := d.getContainerStatsViaCLI(ctx, node, containerID)

    // 方式2: 通过 Docker API
    stream, _ := cli.ContainerStats(ctx, containerID, false)

    // 计算 CPU 使用率（归一化到 0-100%）
    cpuPercent := (cpuDelta / systemDelta) * numCPUs * 100
    normPercent := cpuPercent / assignedCPUs  // 按分配核数归一化

    // 获取 GPU 显存信息
    gpuSize, gpuRate := d.getGPUMemoryInfo(ctx, cli, containerID)
    // 通过 nvidia-smi 命令获取

    return &ContainerStats{
        CPUUsagePercent:    normPercent,
        MemoryUsagePercent: memPercent,
        GPUMemoryUsageRate: gpuRate,
    }
}
```

---

### 3.2 算力节点管理模块 (ComputeNode)

#### 3.2.1 数据模型

```go
type ComputeNode struct {
    global.GVA_MODEL
    Name           *string  // 节点名称
    PublicIp       *string  // 公网IP
    PrivateIp      *string  // 内网IP
    SshPort        *int64   // SSH端口
    GpuName        *string  // GPU型号
    GpuCount       *int64   // GPU数量
    MemoryCapacity *int64   // 显存容量
    DockerAddress  *string  // Docker连接地址
    UseTls         *bool    // 是否使用TLS
    CaCert         *string  // CA证书
    ClientCert     *string  // 客户端证书
    ClientKey      *string  // 客户端私钥
    DockerStatus   *string  // Docker连接状态
    HamiCore       *string  // HAMi-core目录
}
```

#### 3.2.2 Docker 连接测试

```go
func (d *DockerService) TestDockerConnection(ctx context.Context, node *ComputeNode) (bool, string) {
    // 1. 创建 Docker 客户端
    cli, err := d.CreateDockerClient(node)

    // 2. TLS 配置
    if useTLS {
        tlsConfig := d.createTLSConfig(caCert, clientCert, clientKey)
        httpClient := &http.Client{
            Transport: &http.Transport{TLSClientConfig: tlsConfig},
        }
    }

    // 3. Ping 测试
    _, err = cli.Ping(ctx)

    return err == nil, "连接成功"
}
```

#### 3.2.3 定时状态检查

**文件**: `server/service/instance/cron.go`

```go
// 每30秒执行一次
func CheckAllInstancesStatus() {
    // 1. 检查所有节点 Docker 连接状态
    nodes := getAllNodes()
    for _, node := range nodes {
        connected, _ := dockerService.TestDockerConnection(ctx, node)
        updateNodeDockerStatus(node.ID, connected)
    }

    // 2. 刷新所有容器状态和资源使用率
    instances := getAllInstances()
    for _, inst := range instances {
        status, _ := dockerService.GetContainerStatus(ctx, node, containerID)
        stats, _ := dockerService.GetContainerStats(ctx, node, containerID)

        updateInstanceStatus(inst.ID, status, stats)
    }
}
```

---

### 3.3 模型训练模块 (ModelTraining)

#### 3.3.1 数据模型

**训练任务** (`server/model/modeltraining/training_task.go`):

```go
type TrainingTask struct {
    global.GVA_MODEL
    Name           string     // 任务名称
    TaskId         string     // 唯一ID
    BaseModel      string     // 基础模型
    TrainMethod    string     // 训练方式 SFT/DPO/CPT
    TrainType      string     // 训练类型 efficient/full
    Status         string     // 状态 pending/running/serving/completed/failed/stopped
    Progress       int        // 进度百分比
    NodeId         *uint      // 执行节点
    InstanceId     *uint      // 实例ID
    ContainerId    *string    // 容器ID
    HostPort       *int       // 训练端口
}
```

**训练参数** (`server/model/modeltraining/training_param.go`):

```go
type TrainingParam struct {
    TaskId          uint    // 任务ID
    BatchSize       int     // 批次大小
    LearningRate    float64 // 学习率
    NEpochs         int     // 训练轮数
    LoraAlpha       int     // LoRA参数
    LoraRank        int     // LoRA秩
    MaxLength       int     // 序列长度
    LrSchedulerType string  // 学习率调度器
}
```

#### 3.3.2 创建训练任务流程

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  接收请求    │ ──→ │  生成TaskId  │ ──→ │  创建数据库  │
│              │     │ train-{uuid} │     │    记录      │
└──────────────┘     └──────────────┘     └──────────────┘
                                                 │
                                                 ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  分配端口    │ ←── │  获取训练节点 │ ←── │  创建训练参数 │
│  8001+       │     │  chengdu     │     │              │
└──────────────┘     └──────────────┘     └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  更新任务    │ ←── │  创建容器    │ ←── │  准备镜像/规格│
│  状态        │     │  (通过实例管理)│    │  (自动创建)   │
└──────────────┘     └──────────────┘     └──────────────┘
       │
       ▼
┌──────────────┐     ┌──────────────┐
│  异步执行    │ ──→ │  Swift 训练  │
│  训练命令    │     │  vLLM 服务   │
└──────────────┘     └──────────────┘
```

#### 3.3.3 训练执行核心逻辑

**文件**: `server/service/modeltraining/training_task.go`

```go
func (s *TrainingTaskService) launchTrainingTask(ctx context.Context, taskID uint) error {
    // 1. 获取任务和参数
    task, param := getTaskAndParam(taskID)

    // 2. 获取训练节点（固定 chengdu）
    node := getTrainingNode()  // 节点名称: chengdu

    // 3. 分配端口（从 8001 开始递增）
    hostPort := allocateTrainingPort()

    // 4. 确保镜像和规格存在（自动创建）
    imageID, specID := ensureInstanceResources(node)

    // 5. 创建训练容器（通过实例管理模块）
    inst := &Instance{
        ImageId:    imageID,
        SpecId:     specID,
        NodeId:     node.ID,
        HostPort:   hostPort,
        StartupCmd: "tail -f /dev/null",  // 保持容器运行
    }
    instanceService.CreateInstance(ctx, inst)

    // 6. 更新任务状态
    updateTaskStatus(taskID, "running", inst.ContainerId)

    // 7. 异步执行训练
    go runSwiftTraining(taskID, node, containerID, task, param)
}
```

#### 3.3.4 Swift 训练命令生成

```go
func buildSwiftSFTCommand(task TrainingTask, param TrainingParam) string {
    return fmt.Sprintf(`
CUDA_VISIBLE_DEVICES=0 swift sft \
    --model %s \
    --train_type lora \
    --dataset 'AI-ModelScope/alpaca-gpt4-data-zh#100' \
    --num_train_epochs %d \
    --per_device_train_batch_size 1 \
    --learning_rate %g \
    --lora_rank %d \
    --lora_alpha %d \
    --output_dir output \
    --model_name %s
`, model, nEpochs, learningRate, loraRank, loraAlpha, modelName)
}
```

#### 3.3.5 训练完成后启动 vLLM 服务

```go
func (s *TrainingTaskService) runSwiftTraining(...) {
    // 1. 执行训练命令
    execResp, _ := cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
        Cmd: []string{"bash", "-lc", trainingCommand},
    })
    cli.ContainerExecStart(ctx, execResp.ID, ...)

    // 2. 等待训练完成
    for {
        info, _ := cli.ContainerExecInspect(ctx, execResp.ID)
        if !info.Running {
            break
        }
        time.Sleep(3 * time.Second)
    }

    // 3. 提取 last_model_checkpoint
    checkpointPath := extractLastModelCheckpoint(containerID)

    // 4. 启动 vLLM 服务
    vllmCmd := fmt.Sprintf(`
nohup python -m vllm.entrypoints.openai.api_server \
    --model %s \
    --enable-lora \
    --lora-modules lora=%s \
    --port 8000 &
`, baseModel, checkpointPath)

    // 5. 更新状态为 serving
    updateTaskStatus(taskID, "serving")
}
```

---

### 3.4 SSH 跳板机模块 (Jumpbox)

#### 3.4.1 架构设计

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│   用户SSH    │ ──→ │  跳板机服务  │ ──→ │  Docker容器  │
│   客户端     │     │  (端口2026)  │     │  (目标实例)  │
└──────────────┘     └──────────────┘     └──────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │  用户认证    │
                     │  (系统账号)  │
                     └──────────────┘
```

#### 3.4.2 核心实现

**文件**: `server/service/jumpbox/jumpbox.go`

```go
func StartJumpboxServer() error {
    // 1. 配置 SSH 服务器
    config := &ssh.ServerConfig{
        PasswordCallback: PasswordAuth,  // 密码认证
    }

    // 2. 生成/加载主机密钥
    signer := getHostKey()
    config.AddHostKey(signer)

    // 3. 监听端口
    listener, _ := net.Listen("tcp", ":2026")

    // 4. 接受连接
    for {
        conn, _ := listener.Accept()
        go handleConnection(conn)
    }
}
```

#### 3.4.3 会话处理流程

```go
func HandleSession(newChannel ssh.NewChannel, userID, authorityID uint) {
    // 1. 接受会话通道
    channel, requests, _ := newChannel.Accept()

    // 2. 处理 PTY 请求（终端大小）
    for req := range requests {
        if req.Type == "pty-req" {
            // 解析终端大小
        }
        if req.Type == "shell" {
            // 启动交互式会话
        }
    }

    // 3. 显示用户容器列表
    instances := getUserInstances(userID)
    fmt.Fprintf(channel, "请选择要连接的容器:\n")
    for i, inst := range instances {
        fmt.Fprintf(channel, "%d. %s\n", i+1, inst.Name)
    }

    // 4. 等待用户选择
    choice := readInput(channel)

    // 5. 连接到选中的容器
    container := instances[choice]
    execConfig := container.ExecOptions{
        Cmd:          []string{"/bin/bash"},
        Tty:          true,
        AttachStdin:  true,
        AttachStdout: true,
        AttachStderr: true,
    }

    // 6. 执行 exec 并转发数据
    execResp, _ := cli.ContainerExecCreate(ctx, containerID, execConfig)
    execConn, _ := cli.ContainerExecAttach(ctx, execResp.ID, ...)

    // 双向数据转发
    go io.Copy(channel, execConn.Reader)
    go io.Copy(execConn, channel)
}
```

---

### 3.5 用户权限模块 (System)

#### 3.5.1 RBAC 权限模型

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│    用户      │ ←── │  用户-角色   │ ──→ │    角色      │
│  (sys_users) │     │  关联表      │     │(sys_authorities)│
└──────────────┘     └──────────────┘     └──────────────┘
                            │
                            ▼
                     ┌──────────────┐
                     │  角色-菜单   │
                     │  角色-API    │
                     │  Casbin规则  │
                     └──────────────┘
```

#### 3.5.2 权限验证流程

```go
// 1. JWT 中间件验证
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        claims, err := jwt.ParseToken(token)
        c.Set("userId", claims.UserID)
        c.Set("authorityId", claims.AuthorityId)
    }
}

// 2. Casbin 权限检查
func CasbinHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        userId := c.GetInt("userId")
        path := c.Request.URL.Path
        method := c.Request.Method

        ok, _ := casbinEnforcer.Enforce(authorityId, path, method)
        if !ok {
            c.AbortWithStatus(403)
        }
    }
}

// 3. 数据权限过滤
func GetInstanceList(userID uint, isAdmin bool) {
    db := global.GVA_DB.Model(&Instance{})
    if !isAdmin {
        db = db.Where("user_id = ?", userID)
    }
}
```

---

## 四、前端模块详解

### 4.1 API 封装规范

**文件**: `web/src/api/instance/instance.js`

```javascript
import service from '@/utils/request'

/**
 * 获取实例列表
 * @param {Object} data - 查询参数
 * @param {number} data.page - 页码
 * @param {number} data.pageSize - 每页数量
 * @returns {Promise<{code: number, data: {list: Array, total: number}}>}
 */
export const getInstanceList = (data) => {
  return service({
    url: '/instance/getInstanceList',
    method: 'post',
    data
  })
}

/**
 * 创建实例
 * @param {Object} data - 实例信息
 * @returns {Promise<{code: number, msg: string}>}
 */
export const createInstance = (data) => {
  return service({
    url: '/instance/createInstance',
    method: 'post',
    data
  })
}
```

### 4.2 组件开发规范

**文件**: `web/src/view/instance/instance/instance.vue`

```vue
<template>
  <div class="gva-table">
    <!-- 搜索表单 -->
    <el-form :model="searchForm">
      <el-form-item label="实例名称">
        <el-input v-model="searchForm.name" />
      </el-form-item>
    </el-form>

    <!-- 数据表格 -->
    <el-table :data="tableData" v-loading="loading">
      <el-table-column prop="name" label="实例名称" />
      <el-table-column prop="containerStatus" label="状态">
        <template #default="{ row }">
          <el-tag :type="statusTypeMap[row.containerStatus]">
            {{ row.containerStatus }}
          </el-tag>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getInstanceList } from '@/api/instance/instance'

// 响应式数据
const tableData = ref([])
const loading = ref(false)

// 获取列表
const getTableData = async () => {
  loading.value = true
  try {
    const res = await getInstanceList({
      page: pagination.value.page,
      pageSize: pagination.value.pageSize,
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      pagination.value.total = res.data.total
    }
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  getTableData()
})
</script>
```

---

## 五、定时任务

### 5.1 系统健康检查

**文件**: `server/initialize/timer.go`

```go
// 任务名称: system-health-check
// 执行频率: 每 30 秒
// 任务内容:
//   1. 检查所有节点 Docker 连接状态
//   2. 刷新所有容器状态和资源使用率

func RegisterHealthCheckCron() {
    gcron.Add(ctx, "*/30 * * * * *", func(ctx context.Context) {
        // 检查节点状态
        checkNodesDockerStatus()
        // 刷新实例状态
        refreshInstancesStatus()
    }, "system-health-check")
}
```

---

## 六、配置说明

### 6.1 后端配置 (`server/config.yaml`)

```yaml
# 系统配置
system:
  db-type: mysql
  addr: 8890
  use-redis: false

# 数据库配置
mysql:
  path: 127.0.0.1
  port: "3306"
  db-name: docker_gpu_manage
  username: root
  password: "123456"

# JWT 配置
jwt:
  signing-key: your-secret-key
  expires-time: "7d"

# SSH 跳板机配置
jumpbox:
  enabled: true
  port: 2026
  server-ip: "192.168.112.148"
```

### 6.2 前端配置

**文件**: `web/.env.development`

```env
VITE_BASE_API=/api
VITE_BASE_PATH=http://127.0.0.1
VITE_SERVER_PORT=8890
VITE_CLI_PORT=8080
```

---

## 七、数据库初始化

### 7.1 自动初始化

系统采用 GVA 框架的初始化机制，首次启动时会自动检测数据库状态并引导初始化。初始化数据以 Go 代码形式定义在 `server/source/` 目录：

```
server/source/
├── system/                          # 系统核心初始化
│   ├── menu.go                      # 菜单数据
│   ├── api.go                       # API数据
│   ├── casbin.go                    # 权限规则
│   ├── authority.go                 # 角色数据
│   ├── user.go                      # 用户数据
│   └── ...
├── example/                         # 示例模块
│   └── file_upload_download.go
├── modeltraining_init.sql           # 模型训练完整SQL
└── modeltraining_init_simple.sql    # 模型训练简化SQL
```

### 7.2 初始化顺序

系统按以下顺序执行初始化（通过 `initOrder` 常量控制）：

```go
const (
    initOrderSystem       = 1  // 系统基础表
    initOrderAuthority    = 2  // 角色表
    initOrderApiIgnore    = 3  // API忽略列表
    initOrderCasbin       = 4  // Casbin权限表
    initOrderApi          = 5  // API表
    initOrderMenu         = 6  // 菜单表
    // ...
)
```

### 7.3 手动SQL初始化

对于已存在的数据库，可以手动执行SQL初始化模型训练模块：

```bash
# MySQL 8.0+ 推荐使用完整版
mysql -u root -p docker_gpu_manage < server/source/modeltraining_init.sql

# 其他MySQL版本使用简化版（需替换占位符）
mysql -u root -p docker_gpu_manage < server/source/modeltraining_init_simple.sql
```

### 7.4 数据表结构

**数据集表 (dataset)：**

```sql
CREATE TABLE `dataset` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL COMMENT '数据集名称',
  `type` varchar(20) NOT NULL COMMENT '数据集类型(training/evaluation)',
  `format` varchar(50) DEFAULT NULL COMMENT '数据格式',
  `train_method` varchar(20) DEFAULT NULL COMMENT '训练方式(SFT/DPO/CPT)',
  `storage_path` varchar(255) DEFAULT NULL COMMENT '存储路径',
  `description` varchar(500) DEFAULT NULL COMMENT '描述',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者ID',
  `latest_version` varchar(20) DEFAULT NULL COMMENT '最新版本',
  `data_count` bigint DEFAULT 0 COMMENT '数据量',
  `import_status` varchar(20) DEFAULT 'pending' COMMENT '导入状态',
  `publish_status` tinyint(1) DEFAULT 0 COMMENT '发布状态',
  PRIMARY KEY (`id`),
  KEY `idx_dataset_deleted_at` (`deleted_at`),
  KEY `idx_dataset_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**训练任务表 (training_task)：**

```sql
CREATE TABLE `training_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `name` varchar(100) NOT NULL COMMENT '任务名称',
  `task_id` varchar(50) DEFAULT NULL COMMENT '任务ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者ID',
  `base_model` varchar(100) DEFAULT NULL COMMENT '基础模型',
  `train_method` varchar(20) DEFAULT NULL COMMENT '训练方式',
  `train_type` varchar(20) DEFAULT NULL COMMENT '训练类型',
  `status` varchar(20) DEFAULT 'pending' COMMENT '训练状态',
  `progress` int DEFAULT 0 COMMENT '训练进度(百分比)',
  `node_id` bigint unsigned DEFAULT NULL COMMENT '执行节点ID',
  `instance_id` bigint unsigned DEFAULT NULL COMMENT '实例ID',
  `host_port` int DEFAULT NULL COMMENT '训练容器端口',
  `container_id` varchar(128) DEFAULT NULL COMMENT '训练容器ID',
  `container_name` varchar(255) DEFAULT NULL COMMENT '训练容器名称',
  `checkpoint_path` varchar(512) DEFAULT NULL COMMENT '训练产出Checkpoint路径',
  -- ... 其他字段
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_training_task_task_id` (`task_id`),
  KEY `idx_training_task_status` (`status`),
  KEY `idx_training_task_node_id` (`node_id`),
  KEY `idx_training_task_instance_id` (`instance_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**模型测试历史表 (model_test_history)：**

```sql
CREATE TABLE `model_test_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `task_id` bigint unsigned NOT NULL COMMENT '训练任务ID',
  `user_id` bigint unsigned DEFAULT NULL COMMENT '创建者ID',
  `question` text NOT NULL COMMENT '测试问题',
  `base_answer` text COMMENT '基础模型回复',
  `lora_answer` text COMMENT 'LoRA模型回复',
  `test_time` datetime(3) DEFAULT NULL COMMENT '测试时间',
  PRIMARY KEY (`id`),
  KEY `idx_model_test_history_task_id` (`task_id`),
  KEY `idx_model_test_history_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### 7.5 API和权限初始化

模型训练模块包含以下API组：

| API分组 | API数量 | 说明 |
|--------|--------|------|
| 数据集管理 | 13个 | 数据集CRUD、版本管理、文件上传 |
| 训练任务 | 15个 | 任务CRUD、训练控制、推理服务 |
| 模型测试 | 4个 | 测试历史CRUD |

权限角色：
- **888 (超级管理员)**：拥有所有API权限
- **9528 (测试角色)**：拥有模型训练模块所有权限

---

## 八、部署指南

### 8.1 开发环境

```bash
# 后端
cd server
go mod download
go run main.go

# 前端
cd web
npm install
npm run dev
```

### 8.2 Docker 部署

```bash
cd deploy/docker-compose
docker-compose up -d
```

---

## 九、常见问题

### 9.1 Docker 连接失败

- 检查 Docker 服务状态
- 检查 TLS 证书配置
- 检查网络连通性

### 9.2 GPU 显存分配问题

- 确认镜像是否支持显存切分
- 检查 HAMi-core 目录配置
- 查看 CUDA_DEVICE_MEMORY_LIMIT 环境变量

---

## 十、版本历史

| 版本 | 日期 | 更新内容 |
|------|------|----------|
| v1.0.0 | - | 初始版本 |
| v1.1.0 | - | 新增模型训练模块 |
| v1.2.0 | - | 新增显存切分支持 |
| v1.3.0 | - | 完善训练任务功能（推理服务、对话测试等） |

**文档维护**：请随着项目迭代及时更新此文档。