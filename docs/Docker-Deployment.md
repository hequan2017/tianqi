# 🐳 Docker 部署

本页面详细介绍使用 Docker 和 Kubernetes 部署天启算力管理平台的方法。

## Docker Compose 部署

### 目录结构

```
deploy/docker-compose/
└── docker-compose.yaml
```

### 快速启动

```bash
cd deploy/docker-compose

# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f server
docker-compose logs -f web
```

### 默认端口

| 服务 | 端口 | 说明 |
|------|------|------|
| MySQL | 13306 | 数据库 |
| Redis | 16379 | 缓存 |
| 后端服务 | 8888 | API 服务 |
| 前端服务 | 8080 | Web 界面 |

### 修改配置

编辑 `docker-compose.yaml` 修改配置：

```yaml
version: '3.8'

services:
  mysql:
    image: mysql:8.0
    ports:
      - "13306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: gva
    volumes:
      - mysql_data:/var/lib/mysql

  redis:
    image: redis:7-alpine
    ports:
      - "16379:6379"
    volumes:
      - redis_data:/data

  server:
    image: your-registry/docker-gpu-manage-server:latest
    ports:
      - "8888:8888"
      - "2026:2026"  # SSH 跳板机端口
    depends_on:
      - mysql
      - redis
    volumes:
      - ./config.yaml:/app/config.yaml
      - server_logs:/app/log

  web:
    image: your-registry/docker-gpu-manage-web:latest
    ports:
      - "8080:80"
    depends_on:
      - server

volumes:
  mysql_data:
  redis_data:
  server_logs:
```

### 停止服务

```bash
# 停止服务（保留数据）
docker-compose down

# 停止服务并删除数据卷（危险！会清空数据库）
docker-compose down -v
```

---

## 手动构建镜像

### 构建后端镜像

```bash
cd server

# 构建镜像
docker build -t docker-gpu-manage-server:latest .

# 查看 Dockerfile
cat Dockerfile
```

后端 Dockerfile 示例：

```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/config.docker.yaml ./config.yaml

EXPOSE 8888 2026
CMD ["./server"]
```

### 构建前端镜像

```bash
cd web

# 构建镜像
docker build -t docker-gpu-manage-web:latest .

# 查看 Dockerfile
cat Dockerfile
```

前端 Dockerfile 示例：

```dockerfile
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

---

## Kubernetes 部署

### 目录结构

```
deploy/kubernetes/
├── server/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
└── web/
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    └── configmap.yaml
```

### 部署步骤

#### 1. 创建命名空间

```bash
kubectl create namespace docker-gpu-manage
```

#### 2. 部署数据库（可选）

如果使用外部数据库，跳过此步骤。

```bash
# 部署 MySQL
kubectl apply -f mysql/ -n docker-gpu-manage

# 部署 Redis
kubectl apply -f redis/ -n docker-gpu-manage
```

#### 3. 部署后端服务

```bash
cd deploy/kubernetes

# 创建配置
kubectl apply -f server/configmap.yaml -n docker-gpu-manage

# 部署服务
kubectl apply -f server/deployment.yaml -n docker-gpu-manage
kubectl apply -f server/service.yaml -n docker-gpu-manage
```

#### 4. 部署前端服务

```bash
# 部署前端
kubectl apply -f web/ -n docker-gpu-manage
```

#### 5. 配置 Ingress（可选）

```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: docker-gpu-manage
  namespace: docker-gpu-manage
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "100m"
spec:
  ingressClassName: nginx
  rules:
    - host: gpu.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: web
                port:
                  number: 80
          - path: /api
            pathType: Prefix
            backend:
              service:
                name: server
                port:
                  number: 8888
```

### 查看部署状态

```bash
# 查看所有资源
kubectl get all -n docker-gpu-manage

# 查看 Pod 状态
kubectl get pods -n docker-gpu-manage

# 查看日志
kubectl logs -f deployment/server -n docker-gpu-manage

# 进入 Pod
kubectl exec -it <pod-name> -n docker-gpu-manage -- /bin/sh
```

### 配置示例

#### Server Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server
  namespace: docker-gpu-manage
spec:
  replicas: 1
  selector:
    matchLabels:
      app: server
  template:
    metadata:
      labels:
        app: server
    spec:
      containers:
        - name: server
          image: your-registry/docker-gpu-manage-server:latest
          ports:
            - containerPort: 8888
            - containerPort: 2026
          volumeMounts:
            - name: config
              mountPath: /app/config.yaml
              subPath: config.yaml
          resources:
            requests:
              memory: "256Mi"
              cpu: "200m"
            limits:
              memory: "512Mi"
              cpu: "500m"
      volumes:
        - name: config
          configMap:
            name: server-config
```

#### Server Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: server
  namespace: docker-gpu-manage
spec:
  selector:
    app: server
  ports:
    - name: http
      port: 8888
      targetPort: 8888
    - name: ssh
      port: 2026
      targetPort: 2026
  type: ClusterIP
```

---

## 高可用部署

### 多副本部署

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: server
spec:
  replicas: 3  # 多副本
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
```

### 健康检查

```yaml
spec:
  containers:
    - name: server
      livenessProbe:
        httpGet:
          path: /health
          port: 8888
        initialDelaySeconds: 30
        periodSeconds: 10
      readinessProbe:
        httpGet:
          path: /health
          port: 8888
        initialDelaySeconds: 5
        periodSeconds: 5
```

### 资源限制

```yaml
resources:
  requests:
    memory: "256Mi"
    cpu: "200m"
  limits:
    memory: "1Gi"
    cpu: "1000m"
```

---

## 注意事项

1. **SSH 跳板机端口**：如需使用 SSH 跳板机功能，确保端口 2026 可访问
2. **数据持久化**：生产环境务必配置数据卷持久化
3. **密钥管理**：使用 Kubernetes Secret 管理敏感配置
4. **网络策略**：根据需要配置 NetworkPolicy 限制访问
5. **日志收集**：建议配置集中日志收集（如 ELK、Loki）

