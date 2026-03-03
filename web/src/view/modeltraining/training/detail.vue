<template>
  <div class="task-detail-container">
    <div class="breadcrumb">
      <span class="breadcrumb-item" @click="goBack">训练任务</span>
      <span class="breadcrumb-separator">/</span>
      <span class="breadcrumb-current">任务详情</span>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      <span>加载中...</span>
    </div>

    <template v-else>
      <div class="detail-card">
        <div class="card-header">
          <div class="title-section">
            <h2 class="task-title">{{ taskDetail.name }}</h2>
            <el-tag :type="getStatusType(taskDetail.status)" size="large">
              {{ getStatusText(taskDetail.status) }}
            </el-tag>
          </div>
          <div class="action-section">
            <el-button
              v-if="taskDetail.status === 'running'"
              type="success"
              @click="handleMarkCompleted"
            >确认完成</el-button>
            <el-button
              v-if="taskDetail.status === 'running' || taskDetail.status === 'serving'"
              type="warning"
              @click="handleStop"
            >停止任务</el-button>
            <el-button
              v-if="canStart(taskDetail.status)"
              type="primary"
              @click="handleStart"
            >启动任务</el-button>
            <el-button
              v-if="taskDetail.status === 'completed'"
              type="success"
              @click="handleStartService"
            >启动服务</el-button>
            <el-button
              v-if="taskDetail.status === 'serving'"
              type="warning"
              @click="handleStopService"
            >停止服务</el-button>
            <el-button
              v-if="taskDetail.instanceId"
              type="primary"
              plain
              @click="goToInstanceDetail"
            >实例详情</el-button>
            <el-button @click="viewLogs">查看日志</el-button>
          </div>
        </div>

        <!-- 基本信息 -->
        <div class="section-title">
          <el-icon><Document /></el-icon>
          <span>基本信息</span>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">任务ID</span>
            <span class="info-value">
              <code class="code-text">{{ taskDetail.taskId }}</code>
              <el-icon class="copy-btn" @click="copyTaskId"><DocumentCopy /></el-icon>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">状态</span>
            <span class="info-value">
              <el-tag :type="getStatusType(taskDetail.status)" size="small">
                {{ getStatusText(taskDetail.status) }}
              </el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">训练进度</span>
            <span class="info-value">
              <template v-if="taskDetail.status === 'running' || taskDetail.status === 'serving'">
                <el-progress :percentage="taskDetail.progress || 0" :stroke-width="6" style="width: 120px" />
                <span class="progress-percent">{{ taskDetail.progress || 0 }}%</span>
              </template>
              <template v-else-if="taskDetail.status === 'completed'">
                <el-tag type="success" size="small">已完成 100%</el-tag>
              </template>
              <template v-else>-</template>
            </span>
          </div>
        </div>

        <!-- 模型配置 -->
        <div class="section-title">
          <el-icon><Cpu /></el-icon>
          <span>模型配置</span>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">基础模型</span>
            <span class="info-value">
              <el-tag type="info" size="small">{{ taskDetail.baseModel || '-' }}</el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">训练方式</span>
            <span class="info-value">{{ getTrainMethodText(taskDetail.trainMethod) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">训练类型</span>
            <span class="info-value">
              <el-tag :type="taskDetail.trainType === 'efficient' ? 'success' : 'primary'" size="small">
                {{ taskDetail.trainType === 'efficient' ? '高效训练' : '全参训练' }}
              </el-tag>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">输出模型</span>
            <span class="info-value">{{ taskDetail.modelName || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">产出数量上限</span>
            <span class="info-value">{{ taskDetail.outputCount || 5 }} 个</span>
          </div>
          <div class="info-item">
            <span class="info-label">Checkpoint间隔</span>
            <span class="info-value">
              <template v-if="taskDetail.checkpointInterval">
                每 {{ taskDetail.checkpointInterval }} {{ taskDetail.checkpointUnit === 'epoch' ? '轮' : '步' }} 保存
              </template>
              <template v-else>-</template>
            </span>
          </div>
        </div>

        <!-- 数据集配置 -->
        <div class="section-title">
          <el-icon><FolderOpened /></el-icon>
          <span>数据集配置</span>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">训练集ID</span>
            <span class="info-value">{{ taskDetail.trainDatasetId || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">训练集版本ID</span>
            <span class="info-value">{{ taskDetail.trainVersionId || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">验证集ID</span>
            <span class="info-value">{{ taskDetail.valDatasetId || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">验证集版本ID</span>
            <span class="info-value">{{ taskDetail.valVersionId || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">验证集切分比例</span>
            <span class="info-value">
              <template v-if="taskDetail.valSplitRatio">
                <el-tag type="warning" size="small">{{ (taskDetail.valSplitRatio * 100).toFixed(0) }}%</el-tag>
              </template>
              <template v-else>-</template>
            </span>
          </div>
        </div>

        <!-- 执行信息 -->
        <div class="section-title">
          <el-icon><Monitor /></el-icon>
          <span>执行信息</span>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">执行节点ID</span>
            <span class="info-value">{{ taskDetail.nodeId || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">实例ID</span>
            <span class="info-value">
              <template v-if="taskDetail.instanceId">
                <el-link type="primary" @click="goToInstanceDetail">{{ taskDetail.instanceId }}</el-link>
              </template>
              <template v-else>-</template>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">训练端口</span>
            <span class="info-value">
              <template v-if="taskDetail.hostPort">
                <el-tag type="info" size="small">{{ taskDetail.hostPort }}</el-tag>
              </template>
              <template v-else>-</template>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">容器ID</span>
            <span class="info-value">
              <template v-if="taskDetail.containerId">
                <code class="code-text code-short">{{ taskDetail.containerId }}</code>
                <el-icon class="copy-btn" @click="copyContainerId"><DocumentCopy /></el-icon>
              </template>
              <template v-else>-</template>
            </span>
          </div>
          <div class="info-item">
            <span class="info-label">容器名称</span>
            <span class="info-value">{{ taskDetail.containerName || '-' }}</span>
          </div>
          <div class="info-item" v-if="taskDetail.checkpointPath">
            <span class="info-label">Checkpoint路径</span>
            <span class="info-value path-value">
              <el-tooltip :content="taskDetail.checkpointPath" placement="top">
                <span class="path-text">{{ taskDetail.checkpointPath }}</span>
              </el-tooltip>
              <el-icon class="copy-btn" @click="copyCheckpointPath"><DocumentCopy /></el-icon>
            </span>
          </div>
        </div>

        <!-- 时间信息 -->
        <div class="section-title">
          <el-icon><Clock /></el-icon>
          <span>时间信息</span>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">创建时间</span>
            <span class="info-value">{{ formatDate(taskDetail.CreatedAt) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">开始时间</span>
            <span class="info-value">{{ taskDetail.startTime ? formatDate(taskDetail.startTime) : '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">结束时间</span>
            <span class="info-value">{{ taskDetail.endTime ? formatDate(taskDetail.endTime) : '-' }}</span>
          </div>
          <div class="info-item" v-if="taskDetail.startTime">
            <span class="info-label">运行时长</span>
            <span class="info-value">{{ getDuration() }}</span>
          </div>
        </div>
      </div>

      <div class="detail-card" v-if="taskDetail.remark">
        <h3 class="card-title">备注信息</h3>
        <div class="remark-content">{{ taskDetail.remark }}</div>
      </div>

          </template>

    <el-dialog v-model="logDialogVisible" title="训练日志" width="900px">
      <div class="log-header">
        <el-button size="small" @click="refreshLogs">刷新日志</el-button>
        <el-select v-model="logTail" size="small" style="width: 120px; margin-left: 12px" @change="refreshLogs">
          <el-option label="最近100行" :value="100" />
          <el-option label="最近500行" :value="500" />
          <el-option label="最近1000行" :value="1000" />
        </el-select>
      </div>
      <div class="log-container" ref="logContainerRef">
        <pre>{{ logContent }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { findTask, stopTask, startTask, getTaskLogs, startService, stopService, markCompleted } from '@/api/modeltraining/trainingTask'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { DocumentCopy, Loading, InfoFilled, Document, Cpu, FolderOpened, Monitor, Clock } from '@element-plus/icons-vue'

defineOptions({
  name: 'TrainingTaskDetail'
})

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const taskDetail = ref({})

const logDialogVisible = ref(false)
const logContent = ref('')
const logTail = ref(500)
const logContainerRef = ref(null)
let detailTimer = null
let logTimer = null

const canStart = (status) => ['pending', 'failed', 'stopped'].includes(status)

const getTaskDetail = async () => {
  loading.value = true
  try {
    const res = await findTask({ ID: route.params.id })
    if (res.code === 0) {
      taskDetail.value = res.data?.task || res.data || {}
      if (taskDetail.value.status === 'running' || taskDetail.value.status === 'serving') {
        startAutoRefresh()
      } else {
        stopAutoRefresh()
      }
    }
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  router.push({ name: 'trainingTask' })
}

const goToInstanceDetail = () => {
  if (!taskDetail.value.instanceId) {
    ElMessage.warning('当前任务未关联实例')
    return
  }
  router.push({ name: 'Instance', query: { instanceId: String(taskDetail.value.instanceId) } })
}

const copyTaskId = () => {
  navigator.clipboard.writeText(taskDetail.value.taskId || '')
  ElMessage.success('已复制任务ID')
}

const handleStop = async () => {
  ElMessageBox.confirm('确定要停止该训练任务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await stopTask({ ID: taskDetail.value.ID })
    if (res.code === 0) {
      ElMessage.success('已停止')
      getTaskDetail()
    }
  })
}

const handleStart = async () => {
  ElMessageBox.confirm('确定要启动该训练任务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    const res = await startTask({ ID: taskDetail.value.ID })
    if (res.code === 0) {
      ElMessage.success('已启动')
      getTaskDetail()
    }
  })
}

const handleStartService = async () => {
  ElMessageBox.confirm('确定要启动推理服务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    const res = await startService({ ID: taskDetail.value.ID })
    if (res.code === 0) {
      ElMessage.success('推理服务已启动')
      getTaskDetail()
    }
  })
}

const handleStopService = async () => {
  ElMessageBox.confirm('确定要停止推理服务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await stopService({ ID: taskDetail.value.ID })
    if (res.code === 0) {
      ElMessage.success('推理服务已停止')
      getTaskDetail()
    }
  })
}

const handleMarkCompleted = async () => {
  ElMessageBox.confirm('确定要标记训练已完成吗？\n\n这将停止训练进程并尝试提取 checkpoint。', '确认完成', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    const res = await markCompleted({ ID: taskDetail.value.ID })
    if (res.code === 0) {
      ElMessage.success('已标记为训练完成')
      getTaskDetail()
    }
  })
}

const viewLogs = async () => {
  logDialogVisible.value = true
  await refreshLogs()
}

const refreshLogs = async () => {
  const res = await getTaskLogs({ ID: taskDetail.value.ID, tail: logTail.value })
  if (res.code === 0) {
    logContent.value = res.data || '暂无日志'
    await nextTick()
    scrollLogsToBottom()
  }
}

const startAutoRefresh = () => {
  if (detailTimer) return
  detailTimer = setInterval(async () => {
    await getTaskDetail()
  }, 5000)
}

const stopAutoRefresh = () => {
  if (!detailTimer) return
  clearInterval(detailTimer)
  detailTimer = null
}

const startLogAutoRefresh = () => {
  if (logTimer) return
  logTimer = setInterval(async () => {
    if (logDialogVisible.value) await refreshLogs()
  }, 2000)
}

const stopLogAutoRefresh = () => {
  if (!logTimer) return
  clearInterval(logTimer)
  logTimer = null
}

const scrollLogsToBottom = () => {
  const el = logContainerRef.value
  if (el) {
    el.scrollTop = el.scrollHeight
  }
}

const getStatusType = (status) => {
  const map = {
    pending: 'info',
    running: 'primary',
    serving: 'success',
    completed: 'success',
    failed: 'danger',
    stopped: 'warning'
  }
  return map[status] || 'info'
}

const getStatusText = (status) => {
  const map = {
    pending: '待执行',
    running: '运行中',
    serving: '服务中',
    completed: '已完成',
    failed: '失败',
    stopped: '已停止'
  }
  return map[status] || status
}

const getTrainMethodText = (method) => {
  const map = {
    SFT: '监督微调 (SFT)',
    DPO: '直接偏好优化 (DPO)',
    CPT: '持续预训练 (CPT)'
  }
  return map[method] || method || '-'
}

const getDuration = () => {
  if (!taskDetail.value.startTime) return '-'
  const start = new Date(taskDetail.value.startTime)
  const end = taskDetail.value.endTime ? new Date(taskDetail.value.endTime) : new Date()
  const diff = end - start

  const hours = Math.floor(diff / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)

  if (hours > 0) {
    return `${hours}小时 ${minutes}分钟`
  } else if (minutes > 0) {
    return `${minutes}分钟 ${seconds}秒`
  }
  return `${seconds}秒`
}

const copyContainerId = () => {
  navigator.clipboard.writeText(taskDetail.value.containerId || '')
  ElMessage.success('已复制容器ID')
}

const copyCheckpointPath = () => {
  navigator.clipboard.writeText(taskDetail.value.checkpointPath || '')
  ElMessage.success('已复制Checkpoint路径')
}

onMounted(() => {
  getTaskDetail()
})

onBeforeUnmount(() => {
  stopAutoRefresh()
  stopLogAutoRefresh()
})

watch(logDialogVisible, async (visible) => {
  if (visible) {
    await refreshLogs()
    startLogAutoRefresh()
  } else {
    stopLogAutoRefresh()
  }
})
</script>

<style scoped>
.task-detail-container {
  height: 100%;
  max-height: calc(100vh - 120px);
  overflow-y: auto;
  padding: 24px;
  background-color: #f5f7fa;
  box-sizing: border-box;
}

.breadcrumb {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  font-size: 14px;
}

.breadcrumb-item {
  color: #409eff;
  cursor: pointer;
}

.breadcrumb-item:hover {
  color: #66b1ff;
}

.breadcrumb-separator {
  margin: 0 8px;
  color: #c0c4cc;
}

.breadcrumb-current {
  color: #606266;
}

.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 0;
  color: #86909c;
  gap: 12px;
}

.detail-card {
  background-color: #fff;
  border-radius: 8px;
  padding: 24px;
  margin-bottom: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 24px;
  padding-bottom: 16px;
  border-bottom: 1px solid #e8e8e8;
}

.title-section {
  display: flex;
  align-items: center;
  gap: 12px;
}

.task-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0;
}

.action-section {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 16px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  margin: 24px 0 12px;
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
  font-size: 14px;
  font-weight: 600;
  color: #1a1a1a;
}

.section-title:first-of-type {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
}

.section-title .el-icon {
  color: #409eff;
}

.code-text {
  background-color: #f5f7fa;
  padding: 2px 8px;
  border-radius: 4px;
  font-family: Monaco, Menlo, monospace;
  font-size: 13px;
  color: #606266;
}

.code-text.code-short {
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: middle;
}

.path-value {
  max-width: 100%;
  overflow: hidden;
}

.path-text {
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  display: inline-block;
  vertical-align: middle;
  color: #606266;
  cursor: pointer;
}

.path-text:hover {
  color: #409eff;
}

.progress-percent {
  font-weight: 600;
  color: #409eff;
  font-size: 13px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px 24px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  color: #86909c;
}

.info-value {
  font-size: 14px;
  color: #1a1a1a;
  display: flex;
  align-items: center;
  gap: 8px;
}

.copy-btn {
  cursor: pointer;
  color: #86909c;
}

.copy-btn:hover {
  color: #409eff;
}

.progress-section {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #e8e8e8;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  font-size: 14px;
  color: #1a1a1a;
}

.progress-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.progress-text {
  font-weight: 600;
  color: #409eff;
  font-size: 18px;
}

.progress-tag {
  display: flex;
  align-items: center;
  gap: 4px;
}

.progress-label {
  font-weight: 600;
}

.progress-hint {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  font-size: 12px;
  color: #86909c;
}

.remark-content {
  font-size: 14px;
  color: #4e5969;
  line-height: 1.6;
  white-space: pre-wrap;
}

.log-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.log-container {
  background-color: #1a1a1a;
  color: #52c41a;
  padding: 16px;
  border-radius: 8px;
  font-family: Monaco, Menlo, monospace;
  font-size: 13px;
  overflow: auto;
  max-height: 500px;
}

.log-container pre {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
}

@media (max-width: 1200px) {
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .task-detail-container {
    max-height: calc(100vh - 96px);
    padding: 16px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .card-header {
    flex-direction: column;
    gap: 16px;
  }
}
</style>
