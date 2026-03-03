package dashboard

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/computenode"
	dashboardRes "github.com/flipped-aurora/gin-vue-admin/server/model/dashboard/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/imageregistry"
	"github.com/flipped-aurora/gin-vue-admin/server/model/instance"
	"github.com/flipped-aurora/gin-vue-admin/server/model/product"
)

type DashboardService struct{}

// GetDashboardStats 获取仪表盘统计数据
func (d *DashboardService) GetDashboardStats() (res dashboardRes.DashboardStats, err error) {
	// 获取实例统计
	res.InstanceStats, err = d.getInstanceStats()
	if err != nil {
		return
	}

	// 获取产品规格统计
	res.ProductStats, err = d.getProductStats()
	if err != nil {
		return
	}

	// 获取算力节点统计
	res.NodeStats, err = d.getNodeStats()
	if err != nil {
		return
	}

	// 获取镜像库统计
	res.ImageStats, err = d.getImageStats()
	if err != nil {
		return
	}

	// 获取最近实例列表
	res.RecentInstance, err = d.getRecentInstance()
	if err != nil {
		return
	}

	return
}

// getInstanceStats 获取实例统计
func (d *DashboardService) getInstanceStats() (stats dashboardRes.InstanceStats, err error) {
	db := global.GVA_DB.Model(&instance.Instance{})

	// 总数
	if err = db.Count(&stats.Total).Error; err != nil {
		return
	}

	// 运行中
	if err = db.Where("container_status = ?", "running").Count(&stats.Running).Error; err != nil {
		return
	}

	// 已停止
	if err = db.Where("container_status = ?", "exited").Count(&stats.Stopped).Error; err != nil {
		return
	}

	// 其他状态
	stats.OtherStatus = stats.Total - stats.Running - stats.Stopped

	return
}

// getProductStats 获取产品规格统计
func (d *DashboardService) getProductStats() (stats dashboardRes.ProductStats, err error) {
	db := global.GVA_DB.Model(&product.ProductSpec{})

	// 总数
	if err = db.Count(&stats.Total).Error; err != nil {
		return
	}

	// 已上架
	if err = db.Where("is_on_shelf = ?", true).Count(&stats.OnShelf).Error; err != nil {
		return
	}

	// 已下架
	stats.OffShelf = stats.Total - stats.OnShelf

	return
}

// getNodeStats 获取算力节点统计
func (d *DashboardService) getNodeStats() (stats dashboardRes.NodeStats, err error) {
	db := global.GVA_DB.Model(&computenode.ComputeNode{})

	// 总数
	if err = db.Count(&stats.Total).Error; err != nil {
		return
	}

	// 已上架
	if err = db.Where("is_on_shelf = ?", true).Count(&stats.OnShelf).Error; err != nil {
		return
	}

	// 已下架
	stats.OffShelf = stats.Total - stats.OnShelf

	// 在线(Docker已连接)
	if err = db.Where("docker_status = ?", "connected").Count(&stats.Online).Error; err != nil {
		return
	}

	// 离线
	stats.Offline = stats.Total - stats.Online

	// GPU总数和显存总量
	var nodes []computenode.ComputeNode
	if err = db.Select("gpu_count, memory_capacity").Find(&nodes).Error; err != nil {
		return
	}

	for _, node := range nodes {
		if node.GpuCount != nil {
			stats.TotalGpu += *node.GpuCount
		}
		if node.MemoryCapacity != nil {
			stats.TotalMemory += *node.MemoryCapacity
		}
	}

	return
}

// getImageStats 获取镜像库统计
func (d *DashboardService) getImageStats() (stats dashboardRes.ImageStats, err error) {
	db := global.GVA_DB.Model(&imageregistry.ImageRegistry{})

	// 总数
	if err = db.Count(&stats.Total).Error; err != nil {
		return
	}

	// 已上架
	if err = db.Where("is_on_shelf = ?", true).Count(&stats.OnShelf).Error; err != nil {
		return
	}

	// 已下架
	stats.OffShelf = stats.Total - stats.OnShelf

	return
}

// getRecentInstance 获取最近实例列表
func (d *DashboardService) getRecentInstance() (list []dashboardRes.InstanceInfo, err error) {
	var instances []instance.Instance
	if err = global.GVA_DB.Model(&instance.Instance{}).
		Order("created_at DESC").
		Limit(5).
		Find(&instances).Error; err != nil {
		return
	}

	for _, ins := range instances {
		info := dashboardRes.InstanceInfo{
			ID:        ins.ID,
			CreatedAt: ins.CreatedAt.Format("2006-01-02 15:04"),
		}
		if ins.Name != nil {
			info.Name = *ins.Name
		}
		if ins.ContainerStatus != nil {
			info.ContainerStatus = *ins.ContainerStatus
		}
		if ins.CpuUsagePercent != nil {
			info.CpuUsage = *ins.CpuUsagePercent
		}
		if ins.MemoryUsagePercent != nil {
			info.MemoryUsage = *ins.MemoryUsagePercent
		}
		if ins.GpuMemoryUsageRate != nil {
			info.GpuUsage = *ins.GpuMemoryUsageRate
		}
		list = append(list, info)
	}

	return
}