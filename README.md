[简体中文](README.md) | [English](README.en.md)

# 天启算力管理平台（TianQi GPU Manager）

**Docker GPU 算力资源管理平台** —— 企业级的 GPU 容器化资源管理和调度系统，帮助组织高效、安全地管理和分配 GPU 算力资源。

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vuedotjs&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-API%20v27-2496ED?logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-EPL--2.0-blue)

## 📖 项目介绍

天启算力管理平台是一套前后端分离的企业级 GPU 算力管理平台，基于 [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) 脚手架二次开发。平台将分散的 GPU 服务器统一纳管为算力节点，通过"镜像库 + 产品规格 + 算力节点"三者组合自动创建 GPU 容器实例，并提供从创建、连接、监控到销毁的全生命周期管理。

随着 AI / 深度学习的普及，GPU 算力稀缺且昂贵。传统管理方式普遍存在资源利用率低、多节点管理成本高、缺乏权限控制、运维靠人工等问题。本平台针对这些痛点提供：统一资源调度、HAMi 显存切分实现细粒度分配、基于 RBAC 的多用户隔离、SSH 跳板机与 Web 终端双通道访问，以及容器级实时监控和自动化巡检。

平台还集成了模型训练模块，支持数据集版本管理与 SFT/DPO/CPT 等训练方式，训练基于 ModelScope Swift 框架，完成后可自动拉起 vLLM 推理服务，形成"算力分配 → 训练 → 推理"的闭环。

适合：AI/ML 训练平台、科研计算平台、企业内部 GPU 算力池、教育机构教学实验环境等场景。

## ✨ 功能特性

- 📊 **可视化仪表盘**：实例、产品规格、算力节点、镜像库统计概览与最近实例一览
- 🖥️ **算力节点管理**：多 GPU 节点统一纳管，支持 Docker TLS 安全连接，定时自动检测节点 Docker 连通状态
- 📦 **镜像库管理**：多镜像源统一管理，支持标记镜像是否可用于显存切分、上架/下架
- 💰 **产品规格管理**：定义 GPU 型号/数量/显存/CPU/内存/磁盘与小时定价，支持无显卡（GPU=0）的 CPU 规格
- 🐳 **容器实例全生命周期**：创建、启动、停止、重启、删除、日志查看；按规格资源需求智能匹配节点；删除实例自动清理挂载的数据卷
- 📈 **资源监控**：CPU/内存使用率、网络 I/O、块设备 I/O、进程数，运行中实例每 5 秒自动刷新
- 💻 **Web 终端**：WebSocket + xterm，浏览器内直接操作容器（bash/sh，终端窗口自适应）
- 🔐 **SSH 跳板机**：内置 SSH 服务（默认 2026 端口），系统账号密码认证，登录后交互式选择并连接自己的容器，管理员可见全部容器
- ⚡ **HAMi 显存切分**：挂载 HAMi-core 至容器 `/libvgpu/build`，注入 `LD_PRELOAD`、`CUDA_DEVICE_MEMORY_LIMIT`、`CUDA_DEVICE_SM_LIMIT` 环境变量，将单卡显存切分为多个虚拟 GPU
- 🤖 **模型训练**：数据集管理与版本控制（jsonl/xls/xlsx，单文件最大 200MB，批量上传自动生成版本）；SFT/DPO/CPT 训练方式，高效训练（LoRA）与全量训练；基于 Swift 框架执行训练，完成后自动启动 vLLM 推理服务
- 👥 **RBAC 权限**：基于 Casbin 的角色权限控制，普通用户仅能操作自己创建的实例，管理员（authorityId=888）可管理全部资源
- ⏰ **定时任务**：gcron 每 30 秒合并巡检节点 Docker 状态与容器指标，每日自动清理数据库日志
- 🧩 **开发者友好**：Swagger 接口文档、内置 MCP Server（SSE）、Web 引导式数据库初始化

## 🛠 技术栈

**后端**（`server/go.mod`）：

| 分类 | 组件 |
|------|------|
| 语言/框架 | Go 1.23+ / Gin 1.10.0 |
| ORM/数据库 | GORM 1.25.12（MySQL / PostgreSQL / SQLite / MSSQL / Oracle） |
| 认证/权限 | golang-jwt v5.2.2 / Casbin v2.103.0 |
| 容器/终端 | Docker API v27.0.0 / gorilla/websocket 1.5.3 |
| 任务/日志/配置 | gcron（GoFrame v2.9.5）/ Zap 1.27.0 / Viper 1.19.0 |
| 缓存/AI | go-redis v9.7.0 / mcp-go v0.41.1 |

**前端**（`web/package.json`）：

| 分类 | 组件 |
|------|------|
| 框架 | Vue 3.5 / Vite |
| UI/图表 | Element Plus 2.10 / ECharts 5.5 / UnoCSS 66 |
| 状态/路由/请求 | Pinia 2.2 / Vue Router 4.4 / Axios 1.8 |
| 终端 | xterm 5.3 / @vueuse/core 11 |

**部署**：Docker、Docker Compose、Kubernetes（`deploy/kubernetes/`）；训练镜像基于 Ubuntu 22.04 + CUDA 12.8.1 + PyTorch 2.9.0 + vLLM 0.13.0 + Swift 3.12.5。

## 🚀 快速开始

### 环境要求

- 后端：Go 1.23+、MySQL 5.7+（或 PostgreSQL/SQLite 等）、Redis（可选）
- 前端：Node.js 20+，npm 或 pnpm
- 节点：各算力节点需安装 Docker（显存切分需额外制作 HAMi-core 镜像，见下文）

### 方式一：本地开发部署

```bash
git clone https://github.com/hequan2017/tianqi.git
cd tianqi
\mv server/config.yaml.bak server/config.yaml   # 初始化后端配置

# 启动后端（默认监听 8890）
cd server
go mod download
go run main.go

# 启动前端（默认监听 8080）
cd web
npm install        # 或 pnpm install
npm run dev        # 或 pnpm dev
```

访问 `http://localhost:8080`，系统会自动检测数据库状态并跳转到初始化页面，按提示填写数据库连接信息完成初始化。初始化完成后使用默认管理员账号登录：`admin` / `123456`（**首次登录请及时修改密码**）。

常用地址：

- 前端：http://localhost:8080
- 后端 API：http://localhost:8890
- Swagger：http://localhost:8890/swagger/index.html

在支持 MCP 的 AI IDE（如 Cursor）中可直接挂载内置 MCP Server：

```json
{
  "mcpServers": {
    "GVA Helper": {
      "url": "http://127.0.0.1:8890/sse"
    }
  }
}
```

### 方式二：Docker Compose 一键部署

```bash
cd deploy/docker-compose
docker-compose up -d      # 启动 web/server/mysql/redis
docker-compose ps         # 查看服务状态
docker-compose logs -f    # 查看日志
```

默认端口：前端 8080、后端 8888、MySQL 13306、Redis 16379。启动完成后访问 `http://localhost:8080` 完成 Web 引导初始化。

### 方式三：Kubernetes 部署

```bash
kubectl apply -f deploy/kubernetes/server/
kubectl apply -f deploy/kubernetes/web/
```

### 开启显存切分（可选）

1. 在 Docker Server 上参考 [HAMi-core](https://github.com/Project-HAMi/HAMi-core) 构建并制作镜像；
2. 将节点上实际的 HAMi-core 目录路径（如 `/root/HAMi-core/build`）填入"算力节点"的 HAMi-core 目录字段。

创建支持显存切分的容器时，平台会自动挂载该目录到容器 `/libvgpu/build` 并注入相关环境变量。

### 关键配置（server/config.yaml）

```yaml
system:
  db-type: mysql        # 数据库类型
  addr: 8890            # 后端监听端口
  use-redis: false      # 生产环境建议开启

jumpbox:
  enabled: true         # 是否启用 SSH 跳板机
  port: 2026            # SSH 监听端口

jwt:
  signing-key: your-key # 生产环境请修改
  expires-time: 7d
```

前端环境变量（`web/.env.development`）：`VITE_BASE_API`（请求前缀，默认 `/api`）、`VITE_BASE_PATH`、`VITE_SERVER_PORT`（后端端口）、`VITE_CLI_PORT`（前端端口）。

## 📁 目录结构

```
├── server/                 # 后端（Go / Gin）
│   ├── api/v1/             # API 控制器（computenode/imageregistry/product/instance/modeltraining...）
│   ├── service/            # 业务逻辑（jumpbox SSH跳板机、modeltraining 训练服务等）
│   ├── model/              # 数据模型
│   ├── router/             # 路由
│   ├── source/             # 初始化数据（菜单、API、训练模块 SQL）
│   └── config.yaml         # 配置文件（由 config.yaml.bak 复制）
├── web/                    # 前端（Vue3 / Vite）
│   └── src/view/           # 页面（dashboard/instance/modeltraining 等）
├── deploy/                 # docker-compose / kubernetes 部署配置
└── docs/                   # 项目文档与截图
```

## 📸 项目截图

![系统截图](docs/4.png)
![系统截图](docs/1.png)
![系统截图](docs/2.png)
![系统截图](docs/3.png)
![系统截图](docs/5.png)

## 📚 文档

- [开发文档](docs/DEVELOPMENT.md) —— 功能实现逻辑、调用流程、架构设计
- [快速开始](docs/Quick-Start.md) / [功能模块](docs/Features.md) / [配置说明](docs/Configuration.md)
- [Docker 部署](docs/Docker-Deployment.md) / [常见问题](docs/FAQ.md) / [API 参考](docs/API-Reference.md)
- Swagger 接口文档：启动后端后访问 http://localhost:8890/swagger/index.html

## 🔗 相关项目

- [gin-vue-admin](https://github.com/flipped-aurora/gin-vue-admin) —— 本项目所基于的后台管理脚手架
- [HAMi-core](https://github.com/Project-HAMi/HAMi-core) —— GPU 显存切分组件，用于制作 vGPU 镜像

## 📄 License

本项目基于 [EPL-2.0](LICENSE)（Eclipse Public License 2.0）开源。
