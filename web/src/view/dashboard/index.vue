<template>
  <div class="tech-dashboard">
    <!-- 顶部统计卡片 -->
    <div class="stats-grid">
      <!-- 实例统计 -->
      <div class="tech-card instance-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <div class="card-header">
            <div class="icon-wrapper blue">
              <el-icon :size="28"><Briefcase /></el-icon>
            </div>
            <div class="card-title">实例总数</div>
          </div>
          <div class="card-value">{{ animatedStats.instanceTotal }}</div>
          <div class="card-footer">
            <div class="stat-item running">
              <span class="dot"></span>
              <span>运行中 {{ stats.instanceStats.running }}</span>
            </div>
            <div class="stat-item stopped">
              <span class="dot"></span>
              <span>已停止 {{ stats.instanceStats.stopped }}</span>
            </div>
          </div>
        </div>
        <div class="card-corner tl"></div>
        <div class="card-corner tr"></div>
        <div class="card-corner bl"></div>
        <div class="card-corner br"></div>
      </div>

      <!-- 产品规格 -->
      <div class="tech-card product-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <div class="card-header">
            <div class="icon-wrapper purple">
              <el-icon :size="28"><Baseball /></el-icon>
            </div>
            <div class="card-title">产品规格</div>
          </div>
          <div class="card-value">{{ animatedStats.productTotal }}</div>
          <div class="card-footer">
            <div class="stat-item success">
              <span class="dot"></span>
              <span>已上架 {{ stats.productStats.onShelf }}</span>
            </div>
            <div class="stat-item muted">
              <span class="dot"></span>
              <span>已下架 {{ stats.productStats.offShelf }}</span>
            </div>
          </div>
        </div>
        <div class="card-corner tl"></div>
        <div class="card-corner tr"></div>
        <div class="card-corner bl"></div>
        <div class="card-corner br"></div>
      </div>

      <!-- 算力节点 -->
      <div class="tech-card node-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <div class="card-header">
            <div class="icon-wrapper green">
              <el-icon :size="28"><Monitor /></el-icon>
            </div>
            <div class="card-title">算力节点</div>
          </div>
          <div class="card-value">{{ animatedStats.nodeTotal }}</div>
          <div class="card-footer">
            <div class="stat-item online">
              <span class="dot pulse"></span>
              <span>在线 {{ stats.nodeStats.online }}</span>
            </div>
            <div class="stat-item offline">
              <span class="dot"></span>
              <span>离线 {{ stats.nodeStats.offline }}</span>
            </div>
          </div>
        </div>
        <div class="card-corner tl"></div>
        <div class="card-corner tr"></div>
        <div class="card-corner bl"></div>
        <div class="card-corner br"></div>
      </div>

      <!-- 镜像库 -->
      <div class="tech-card image-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <div class="card-header">
            <div class="icon-wrapper orange">
              <el-icon :size="28"><Box /></el-icon>
            </div>
            <div class="card-title">镜像库</div>
          </div>
          <div class="card-value">{{ animatedStats.imageTotal }}</div>
          <div class="card-footer">
            <div class="stat-item success">
              <span class="dot"></span>
              <span>已上架 {{ stats.imageStats.onShelf }}</span>
            </div>
            <div class="stat-item muted">
              <span class="dot"></span>
              <span>已下架 {{ stats.imageStats.offShelf }}</span>
            </div>
          </div>
        </div>
        <div class="card-corner tl"></div>
        <div class="card-corner tr"></div>
        <div class="card-corner bl"></div>
        <div class="card-corner br"></div>
      </div>
    </div>

    <!-- 中间区域 -->
    <div class="middle-section">
      <!-- GPU资源面板 -->
      <div class="tech-panel gpu-panel">
        <div class="panel-header">
          <div class="panel-title">
            <span class="title-icon">◈</span>
            GPU资源概览
          </div>
        </div>
        <div class="panel-body">
          <div class="gpu-stats">
            <div class="gpu-stat-item">
              <div class="gpu-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="3" y="6" width="18" height="12" rx="2"/>
                  <line x1="7" y1="10" x2="7" y2="14"/>
                  <line x1="11" y1="10" x2="11" y2="14"/>
                  <line x1="15" y1="10" x2="15" y2="14"/>
                  <line x1="19" y1="10" x2="19" y2="14"/>
                </svg>
              </div>
              <div class="gpu-info">
                <div class="gpu-label">GPU总数</div>
                <div class="gpu-value">
                  <span class="value-number">{{ stats.nodeStats.totalGpu }}</span>
                  <span class="value-unit">块</span>
                </div>
              </div>
              <div class="gpu-bar">
                <div class="bar-fill" :style="{ width: gpuUsagePercent + '%' }"></div>
              </div>
            </div>
            <div class="gpu-stat-item">
              <div class="gpu-icon memory">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                  <rect x="4" y="4" width="16" height="16" rx="2"/>
                  <rect x="8" y="8" width="8" height="8"/>
                </svg>
              </div>
              <div class="gpu-info">
                <div class="gpu-label">显存总量</div>
                <div class="gpu-value">
                  <span class="value-number">{{ stats.nodeStats.totalMemory }}</span>
                  <span class="value-unit">GB</span>
                </div>
              </div>
              <div class="gpu-bar memory">
                <div class="bar-fill" :style="{ width: memoryUsagePercent + '%' }"></div>
              </div>
            </div>
          </div>
        </div>
        <div class="panel-border"></div>
      </div>

      <!-- 实例状态分布 -->
      <div class="tech-panel status-panel">
        <div class="panel-header">
          <div class="panel-title">
            <span class="title-icon">◈</span>
            实例状态分布
          </div>
        </div>
        <div class="panel-body">
          <div class="status-chart" v-if="stats.instanceStats.total > 0">
            <div class="chart-ring">
              <svg viewBox="0 0 120 120">
                <defs>
                  <linearGradient id="runningGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#00d4aa"/>
                    <stop offset="100%" style="stop-color:#00ff88"/>
                  </linearGradient>
                  <linearGradient id="stoppedGrad" x1="0%" y1="0%" x2="100%" y2="100%">
                    <stop offset="0%" style="stop-color:#ff4757"/>
                    <stop offset="100%" style="stop-color:#ff6b7a"/>
                  </linearGradient>
                </defs>
                <circle cx="60" cy="60" r="50" fill="none" stroke="#1a2a3a" stroke-width="12"/>
                <circle cx="60" cy="60" r="50" fill="none" stroke="url(#runningGrad)" stroke-width="12"
                  :stroke-dasharray="runningArc" :stroke-dashoffset="-stoppedArc"
                  stroke-linecap="round" class="ring-segment"/>
                <circle cx="60" cy="60" r="50" fill="none" stroke="url(#stoppedGrad)" stroke-width="12"
                  :stroke-dasharray="stoppedArc" :stroke-dashoffset="0"
                  stroke-linecap="round" class="ring-segment"/>
              </svg>
              <div class="chart-center">
                <div class="center-value">{{ stats.instanceStats.total }}</div>
                <div class="center-label">总数</div>
              </div>
            </div>
            <div class="chart-legend">
              <div class="legend-item">
                <span class="legend-dot running"></span>
                <span class="legend-label">运行中</span>
                <span class="legend-value">{{ stats.instanceStats.running }}</span>
              </div>
              <div class="legend-item">
                <span class="legend-dot stopped"></span>
                <span class="legend-label">已停止</span>
                <span class="legend-value">{{ stats.instanceStats.stopped }}</span>
              </div>
              <div class="legend-item">
                <span class="legend-dot other"></span>
                <span class="legend-label">其他</span>
                <span class="legend-value">{{ stats.instanceStats.otherStatus }}</span>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <el-icon :size="48"><Document /></el-icon>
            <span>暂无实例数据</span>
          </div>
        </div>
        <div class="panel-border"></div>
      </div>
    </div>

    <!-- 最近实例列表 -->
    <div class="tech-panel table-panel">
      <div class="panel-header">
        <div class="panel-title">
          <span class="title-icon">◈</span>
          最近创建的实例
        </div>
        <div class="panel-action" @click="goToInstance">
          <span>查看全部</span>
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>
      <div class="panel-body">
        <div class="tech-table" v-if="stats.recentInstance.length > 0">
          <div class="table-header">
            <div class="th name">实例名称</div>
            <div class="th status">状态</div>
            <div class="th resources">资源使用</div>
            <div class="th time">创建时间</div>
          </div>
          <div class="table-body">
            <div class="table-row" v-for="(item, index) in stats.recentInstance" :key="index">
              <div class="td name">
                <span class="instance-name">{{ item.name }}</span>
              </div>
              <div class="td status">
                <span class="status-badge" :class="item.containerStatus">
                  {{ getStatusText(item.containerStatus) }}
                </span>
              </div>
              <div class="td resources">
                <div class="resource-bars">
                  <div class="resource-item">
                    <span class="resource-label">CPU</span>
                    <div class="mini-bar">
                      <div class="mini-fill cpu" :style="{ width: (item.cpuUsage || 0) + '%' }"></div>
                    </div>
                    <span class="resource-value">{{ (item.cpuUsage || 0).toFixed(1) }}%</span>
                  </div>
                  <div class="resource-item">
                    <span class="resource-label">MEM</span>
                    <div class="mini-bar">
                      <div class="mini-bar mem" :style="{ width: (item.memoryUsage || 0) + '%' }"></div>
                    </div>
                    <span class="resource-value">{{ (item.memoryUsage || 0).toFixed(1) }}%</span>
                  </div>
                  <div class="resource-item">
                    <span class="resource-label">GPU</span>
                    <div class="mini-bar">
                      <div class="mini-bar gpu" :style="{ width: (item.gpuUsage || 0) + '%' }"></div>
                    </div>
                    <span class="resource-value">{{ (item.gpuUsage || 0).toFixed(1) }}%</span>
                  </div>
                </div>
              </div>
              <div class="td time">{{ item.createdAt }}</div>
            </div>
          </div>
        </div>
        <div v-else class="empty-state">
          <el-icon :size="48"><Folder /></el-icon>
          <span>暂无实例数据</span>
        </div>
      </div>
      <div class="panel-border"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { getDashboardStats } from '@/api/dashboard/dashboard'
import {
  Briefcase,
  Baseball,
  Monitor,
  Box,
  ArrowRight,
  Document,
  Folder
} from '@element-plus/icons-vue'

defineOptions({
  name: 'Dashboard'
})

const router = useRouter()

const stats = ref({
  instanceStats: { total: 0, running: 0, stopped: 0, otherStatus: 0 },
  productStats: { total: 0, onShelf: 0, offShelf: 0 },
  nodeStats: { total: 0, online: 0, offline: 0, totalGpu: 0, totalMemory: 0 },
  imageStats: { total: 0, onShelf: 0, offShelf: 0 },
  recentInstance: []
})

// 动画数字
const animatedStats = reactive({
  instanceTotal: 0,
  productTotal: 0,
  nodeTotal: 0,
  imageTotal: 0
})

// GPU使用率模拟
const gpuUsagePercent = computed(() => Math.min(75, Math.random() * 100))
const memoryUsagePercent = computed(() => Math.min(60, Math.random() * 100))

// 圆环图计算
const circumference = 2 * Math.PI * 50
const runningArc = computed(() => {
  const total = stats.value.instanceStats.total || 1
  return (stats.value.instanceStats.running / total) * circumference
})
const stoppedArc = computed(() => {
  const total = stats.value.instanceStats.total || 1
  return (stats.value.instanceStats.stopped / total) * circumference
})

// 数字动画
const animateNumber = (target, value) => {
  const duration = 1000
  const start = animatedStats[target]
  const startTime = performance.now()

  const animate = (currentTime) => {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)
    const easeOut = 1 - Math.pow(1 - progress, 3)
    animatedStats[target] = Math.floor(start + (value - start) * easeOut)

    if (progress < 1) {
      requestAnimationFrame(animate)
    }
  }
  requestAnimationFrame(animate)
}

// 获取状态文本
const getStatusText = (status) => {
  const map = {
    running: '运行中',
    exited: '已停止',
    created: '已创建',
    paused: '已暂停',
    restarting: '重启中'
  }
  return map[status] || status || '未知'
}

// 跳转
const goToInstance = () => router.push('/instance')

// 获取数据
const fetchDashboardStats = async () => {
  try {
    const res = await getDashboardStats()
    if (res.code === 0) {
      stats.value = res.data
      // 触发数字动画
      animateNumber('instanceTotal', res.data.instanceStats.total)
      animateNumber('productTotal', res.data.productStats.total)
      animateNumber('nodeTotal', res.data.nodeStats.total)
      animateNumber('imageTotal', res.data.imageStats.total)
    }
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
  }
}

onMounted(() => {
  fetchDashboardStats()
})
</script>

<style lang="scss" scoped>
.tech-dashboard {
  padding: 20px;
}

// 统计卡片网格
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  margin-bottom: 20px;

  @media (max-width: 1200px) {
    grid-template-columns: repeat(2, 1fr);
  }
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

// 科技风格卡片
.tech-card {
  position: relative;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

  &:hover {
    border-color: #409eff;
    transform: translateY(-4px);
    box-shadow: 0 10px 40px rgba(64, 158, 255, 0.15);

    .card-glow {
      opacity: 1;
    }
  }

  .card-glow {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, transparent, #409eff, transparent);
    opacity: 0.5;
    transition: opacity 0.3s;
  }

  .card-content {
    padding: 24px;
  }

  .card-header {
    display: flex;
    align-items: center;
    margin-bottom: 16px;

    .icon-wrapper {
      width: 48px;
      height: 48px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 12px;

      &.blue { background: rgba(64, 158, 255, 0.1); color: #409eff; }
      &.purple { background: rgba(167, 139, 250, 0.1); color: #a78bfa; }
      &.green { background: rgba(0, 212, 170, 0.1); color: #00d4aa; }
      &.orange { background: rgba(255, 150, 50, 0.1); color: #ff9632; }
    }

    .card-title {
      font-size: 14px;
      color: #606266;
      font-weight: 500;
    }
  }

  .card-value {
    font-size: 42px;
    font-weight: 700;
    color: #303133;
    margin-bottom: 16px;
    font-family: 'Orbitron', 'Courier New', monospace;
  }

  .card-footer {
    display: flex;
    gap: 20px;

    .stat-item {
      display: flex;
      align-items: center;
      font-size: 13px;
      color: #909399;

      .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        margin-right: 8px;

        &.pulse {
          animation: pulse 2s infinite;
        }
      }

      &.running .dot { background: #00d4aa; }
      &.stopped .dot { background: #ff4757; }
      &.online .dot { background: #00d4aa; animation: pulse 2s infinite; }
      &.offline .dot { background: #909399; }
      &.success .dot { background: #00d4aa; }
      &.muted .dot { background: #c0c4cc; }
    }
  }

  // 角标装饰
  .card-corner {
    position: absolute;
    width: 12px;
    height: 12px;
    border: 1px solid #409eff;
    opacity: 0.3;

    &.tl { top: 4px; left: 4px; border-right: none; border-bottom: none; }
    &.tr { top: 4px; right: 4px; border-left: none; border-bottom: none; }
    &.bl { bottom: 4px; left: 4px; border-right: none; border-top: none; }
    &.br { bottom: 4px; right: 4px; border-left: none; border-top: none; }
  }
}

// 中间区域
.middle-section {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 20px;
  margin-bottom: 20px;

  @media (max-width: 1024px) {
    grid-template-columns: 1fr;
  }
}

// 科技面板
.tech-panel {
  position: relative;
  background: #fff;
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.04);

  .panel-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16px 20px;
    border-bottom: 1px solid #ebeef5;
  }

  .panel-title {
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    display: flex;
    align-items: center;

    .title-icon {
      color: #409eff;
      margin-right: 10px;
    }
  }

  .panel-action {
    display: flex;
    align-items: center;
    gap: 6px;
    color: #409eff;
    font-size: 14px;
    cursor: pointer;
    transition: all 0.2s;

    &:hover {
      color: #66b1ff;
    }
  }

  .panel-body {
    padding: 20px;
  }

  .panel-border {
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 2px;
    background: linear-gradient(90deg, transparent, rgba(64, 158, 255, 0.3), transparent);
  }
}

// GPU统计
.gpu-stats {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .gpu-stat-item {
    display: flex;
    align-items: center;
    gap: 16px;

    .gpu-icon {
      width: 50px;
      height: 50px;
      background: rgba(64, 158, 255, 0.1);
      border-radius: 10px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #409eff;

      svg {
        width: 28px;
        height: 28px;
      }

      &.memory {
        background: rgba(255, 150, 50, 0.1);
        color: #ff9632;
      }
    }

    .gpu-info {
      flex: 1;

      .gpu-label {
        font-size: 13px;
        color: #909399;
        margin-bottom: 4px;
      }

      .gpu-value {
        display: flex;
        align-items: baseline;
        gap: 6px;

        .value-number {
          font-size: 28px;
          font-weight: 700;
          color: #303133;
          font-family: 'Orbitron', monospace;
        }

        .value-unit {
          font-size: 14px;
          color: #909399;
        }
      }
    }

    .gpu-bar {
      width: 120px;
      height: 6px;
      background: #f0f2f5;
      border-radius: 3px;
      overflow: hidden;

      .bar-fill {
        height: 100%;
        background: linear-gradient(90deg, #409eff, #67c23a);
        border-radius: 3px;
        transition: width 1s ease;
      }

      &.memory .bar-fill {
        background: linear-gradient(90deg, #ff9632, #ffcc00);
      }
    }
  }
}

// 状态图表
.status-chart {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 40px;

  .chart-ring {
    position: relative;
    width: 150px;
    height: 150px;

    svg {
      width: 100%;
      height: 100%;
      transform: rotate(-90deg);
    }

    .ring-segment {
      transition: stroke-dasharray 1s ease;
    }

    .chart-center {
      position: absolute;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      text-align: center;

      .center-value {
        font-size: 32px;
        font-weight: 700;
        color: #303133;
        font-family: 'Orbitron', monospace;
      }

      .center-label {
        font-size: 12px;
        color: #909399;
      }
    }
  }

  .chart-legend {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .legend-item {
      display: flex;
      align-items: center;
      gap: 10px;

      .legend-dot {
        width: 12px;
        height: 12px;
        border-radius: 3px;

        &.running { background: linear-gradient(135deg, #00d4aa, #00ff88); }
        &.stopped { background: linear-gradient(135deg, #ff4757, #ff6b7a); }
        &.other { background: #c0c4cc; }
      }

      .legend-label {
        color: #606266;
        font-size: 14px;
        width: 60px;
      }

      .legend-value {
        color: #303133;
        font-weight: 600;
        font-size: 16px;
      }
    }
  }
}

// 科技表格
.tech-table {
  .table-header {
    display: grid;
    grid-template-columns: 1fr 100px 1fr 120px;
    gap: 16px;
    padding: 12px 16px;
    background: #f5f7fa;
    border-radius: 8px;
    margin-bottom: 8px;

    .th {
      font-size: 13px;
      color: #909399;
      font-weight: 500;

      &.name { padding-left: 8px; }
    }
  }

  .table-body {
    .table-row {
      display: grid;
      grid-template-columns: 1fr 100px 1fr 120px;
      gap: 16px;
      padding: 16px;
      border-radius: 8px;
      transition: all 0.2s;

      &:hover {
        background: #f5f7fa;
      }

      .td {
        display: flex;
        align-items: center;

        &.name {
          .instance-name {
            color: #303133;
            font-weight: 500;
          }
        }

        &.status {
          .status-badge {
            padding: 4px 12px;
            border-radius: 20px;
            font-size: 12px;
            font-weight: 500;

            &.running {
              background: rgba(0, 212, 170, 0.1);
              color: #00d4aa;
              border: 1px solid rgba(0, 212, 170, 0.3);
            }

            &.exited, &.stopped {
              background: rgba(255, 71, 87, 0.1);
              color: #ff4757;
              border: 1px solid rgba(255, 71, 87, 0.3);
            }

            &.created, &.paused, &.restarting {
              background: rgba(144, 147, 153, 0.1);
              color: #909399;
              border: 1px solid rgba(144, 147, 153, 0.3);
            }
          }
        }

        &.resources {
          .resource-bars {
            display: flex;
            gap: 16px;

            .resource-item {
              display: flex;
              align-items: center;
              gap: 6px;

              .resource-label {
                font-size: 11px;
                color: #909399;
                width: 28px;
              }

              .mini-bar {
                width: 50px;
                height: 4px;
                background: #e4e7ed;
                border-radius: 2px;
                overflow: hidden;

                .mini-fill {
                  height: 100%;
                  border-radius: 2px;
                  transition: width 0.5s ease;

                  &.cpu { background: #409eff; }
                  &.mem { background: #a78bfa; }
                  &.gpu { background: #00d4aa; }
                }
              }

              .resource-value {
                font-size: 11px;
                color: #606266;
                width: 40px;
              }
            }
          }
        }

        &.time {
          color: #909399;
          font-size: 13px;
        }
      }
    }
  }
}

// 空状态
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: #c0c4cc;

  .el-icon {
    margin-bottom: 16px;
    opacity: 0.5;
  }

  span {
    font-size: 14px;
  }
}

// 脉冲动画
@keyframes pulse {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.6; transform: scale(1.2); }
}

// 表格面板响应式
.table-panel {
  .tech-table {
    .table-header, .table-row {
      @media (max-width: 900px) {
        grid-template-columns: 1fr 80px;
        gap: 12px;

        .th.resources, .td.resources,
        .th.time, .td.time {
          display: none;
        }
      }
    }
  }
}
</style>