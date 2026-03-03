package dashboard

// DashboardStats 仪表盘统计响应
type DashboardStats struct {
	InstanceStats  InstanceStats  `json:"instanceStats"`  // 实例统计
	ProductStats   ProductStats   `json:"productStats"`   // 产品规格统计
	NodeStats      NodeStats      `json:"nodeStats"`      // 算力节点统计
	ImageStats     ImageStats     `json:"imageStats"`     // 镜像库统计
	RecentInstance []InstanceInfo `json:"recentInstance"` // 最近实例列表
}

// InstanceStats 实例统计
type InstanceStats struct {
	Total       int64 `json:"total"`       // 总数
	Running     int64 `json:"running"`     // 运行中
	Stopped     int64 `json:"stopped"`     // 已停止
	OtherStatus int64 `json:"otherStatus"` // 其他状态
}

// ProductStats 产品规格统计
type ProductStats struct {
	Total    int64 `json:"total"`    // 总数
	OnShelf  int64 `json:"onShelf"`  // 已上架
	OffShelf int64 `json:"offShelf"` // 已下架
}

// NodeStats 算力节点统计
type NodeStats struct {
	Total       int64 `json:"total"`       // 总数
	OnShelf     int64 `json:"onShelf"`     // 已上架
	OffShelf    int64 `json:"offShelf"`    // 已下架
	Online      int64 `json:"online"`      // 在线(Docker已连接)
	Offline     int64 `json:"offline"`     // 离线(Docker未连接)
	TotalGpu    int64 `json:"totalGpu"`    // GPU总数
	TotalMemory int64 `json:"totalMemory"` // 显存总量(GB)
}

// ImageStats 镜像库统计
type ImageStats struct {
	Total    int64 `json:"total"`    // 总数
	OnShelf  int64 `json:"onShelf"`  // 已上架
	OffShelf int64 `json:"offShelf"` // 已下架
}

// InstanceInfo 实例简要信息
type InstanceInfo struct {
	ID              uint    `json:"ID"`
	Name            string  `json:"name"`
	ContainerStatus string  `json:"containerStatus"`
	CpuUsage        float64 `json:"cpuUsage"`
	MemoryUsage     float64 `json:"memoryUsage"`
	GpuUsage        float64 `json:"gpuUsage"`
	CreatedAt       string  `json:"createdAt"`
}