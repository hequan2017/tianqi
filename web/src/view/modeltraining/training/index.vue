<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="search-form" @keyup.enter="onSubmit">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="searchInfo.name" placeholder="请输入任务名称" clearable style="width: 220px" />
        </el-form-item>

        <el-form-item label="训练状态" prop="status">
          <el-select v-model="searchInfo.status" placeholder="请选择状态" clearable style="width: 140px">
            <el-option v-for="item in dataSource.status" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="训练方式" prop="trainMethod">
          <el-select v-model="searchInfo.trainMethod" placeholder="请选择方式" clearable style="width: 120px">
            <el-option v-for="item in dataSource.trainMethod" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="基础模型" prop="baseModel">
          <el-select v-model="searchInfo.baseModel" placeholder="请选择模型" filterable clearable style="width: 180px">
            <el-option v-for="item in dataSource.baseModel" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="创建时间" prop="createdAtRange">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            style="width: 360px"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="Search" @click="onSubmit">查询</el-button>
          <el-button icon="Refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="Plus" @click="goCreateTask">创建训练任务</el-button>
        <el-button icon="Delete" :disabled="!multipleSelection.length" @click="onDelete">批量删除</el-button>
      </div>

      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="50" />

        <el-table-column align="left" label="任务名称" min-width="180">
          <template #default="scope">
            <span class="task-name" @click="viewDetail(scope.row)">{{ scope.row.name }}</span>
          </template>
        </el-table-column>

        <el-table-column align="center" label="任务ID" width="180">
          <template #default="scope">
            <div class="task-id-cell">
              <span class="task-id">{{ scope.row.taskId }}</span>
              <el-tooltip content="复制" placement="top">
                <el-icon class="copy-icon" @click.stop="copyTaskId(scope.row.taskId)"><DocumentCopy /></el-icon>
              </el-tooltip>
            </div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="基础模型" prop="baseModel" width="170" />

        <el-table-column align="center" label="训练方式" width="140">
          <template #default="scope">
            <el-tag size="small" effect="plain">{{ scope.row.trainMethod }}</el-tag>
            <el-tag v-if="scope.row.trainType" size="small" type="info" effect="plain" style="margin-left: 4px">
              {{ scope.row.trainType === 'efficient' ? '高效' : '全参' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="center" label="状态" width="110">
          <template #default="scope">
            <div class="status-cell">
              <el-icon :class="['status-icon', 'status-' + scope.row.status]">
                <component :is="getStatusIcon(scope.row.status)" />
              </el-icon>
              <span>{{ getStatusText(scope.row.status) }}</span>
            </div>
          </template>
        </el-table-column>

        <el-table-column align="center" label="进度" width="140">
          <template #default="scope">
            <template v-if="scope.row.status === 'running'">
              <el-progress :percentage="scope.row.progress || 0" :stroke-width="6" style="width: 100px" />
            </template>
            <template v-else-if="scope.row.status === 'serving' || scope.row.status === 'completed'">
              <span class="progress-done">100%</span>
            </template>
            <template v-else>
              <span class="progress-none">-</span>
            </template>
          </template>
        </el-table-column>

        <el-table-column sortable align="center" label="创建时间" prop="CreatedAt" width="170">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="center" label="操作" fixed="right" width="240">
          <template #default="scope">
            <div class="action-btns">
              <el-button type="primary" link size="small" @click="viewDetail(scope.row)">详情</el-button>
              <el-button
                v-if="scope.row.status === 'serving'"
                type="success"
                link
                size="small"
                @click="goModelTest(scope.row)"
              >测试</el-button>
              <el-button
                v-if="scope.row.status === 'running' || scope.row.status === 'serving'"
                type="warning"
                link
                size="small"
                @click="handleStop(scope.row)"
              >停止</el-button>
              <el-button type="info" link size="small" @click="viewLogs(scope.row)">日志</el-button>
              <el-button type="danger" link size="small" @click="deleteRow(scope.row)">删除</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          layout="total, sizes, prev, pager, next, jumper"
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="logDialogVisible" title="训练日志" width="900px" destroy-on-close>
      <div class="log-header">
        <span class="log-title">任务：{{ currentTask?.name }}</span>
        <el-button size="small" icon="Refresh" @click="refreshLogs(true)">刷新</el-button>
      </div>
      <div class="log-container" ref="logContainerRef">
        <pre>{{ logContent }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { deleteTask, deleteTaskByIds, getTaskList, stopTask, getTaskLogs, getTaskDataSource } from '@/api/modeltraining/trainingTask'
import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted, onBeforeUnmount, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { DocumentCopy, Loading, CircleCheck, CircleClose, VideoPause } from '@element-plus/icons-vue'

defineOptions({
  name: 'TrainingTask'
})

const router = useRouter()
const elSearchFormRef = ref()

const dataSource = ref({
  status: [
    { label: '待执行', value: 'pending' },
    { label: '运行中', value: 'running' },
    { label: '服务中', value: 'serving' },
    { label: '已完成', value: 'completed' },
    { label: '失败', value: 'failed' },
    { label: '已停止', value: 'stopped' }
  ],
  trainMethod: [
    { label: 'SFT', value: 'SFT' },
    { label: 'DPO', value: 'DPO' },
    { label: 'CPT', value: 'CPT' }
  ],
  baseModel: [
    { label: 'Qwen3-1.7B', value: 'Qwen3-1.7B' },
    { label: 'Qwen3-7B', value: 'Qwen3-7B' },
    { label: 'Qwen3-14B', value: 'Qwen3-14B' },
    { label: 'Llama3-8B', value: 'Llama3-8B' }
  ]
})

const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

const logDialogVisible = ref(false)
const logContent = ref('')
const currentTask = ref(null)
const logContainerRef = ref(null)
let logTimer = null

const getTableData = async () => {
  const table = await getTaskList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

const getDataSource = async () => {
  const res = await getTaskDataSource()
  if (res.code === 0) {
    dataSource.value = { ...dataSource.value, ...res.data }
  }
}

const onSubmit = () => {
  page.value = 1
  getTableData()
}

const onReset = () => {
  searchInfo.value = {}
  page.value = 1
  getTableData()
}

const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

const goCreateTask = () => {
  router.push({ name: 'createTrainingTask' })
}

const viewDetail = (row) => {
  router.push({ name: 'TrainingTaskDetail', params: { id: row.ID } })
}

const goModelTest = (row) => {
  router.push({ name: 'ModelTest', params: { id: row.ID } })
}

const copyTaskId = (taskId) => {
  navigator.clipboard.writeText(taskId)
  ElMessage.success('已复制任务ID')
}

const handleStop = async (row) => {
  ElMessageBox.confirm('确定要停止该训练任务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await stopTask({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('停止成功')
      getTableData()
    }
  })
}

const viewLogs = async (row) => {
  currentTask.value = row
  logDialogVisible.value = true
  await refreshLogs()
}

const refreshLogs = async (showToast = false) => {
  if (!currentTask.value) return
  const res = await getTaskLogs({ ID: currentTask.value.ID, tail: '500' })
  if (res.code === 0) {
    logContent.value = res.data || '暂无日志'
    await nextTick()
    scrollLogsToBottom()
    if (showToast) ElMessage.success('已刷新')
  }
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

const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除该训练任务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteTask({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

const onDelete = async () => {
  ElMessageBox.confirm('确定要删除选中的训练任务吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = multipleSelection.value.map((item) => item.ID)
    const res = await deleteTaskByIds({ IDs })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

const getStatusIcon = (status) => {
  const map = {
    pending: VideoPause,
    running: Loading,
    serving: CircleCheck,
    completed: CircleCheck,
    failed: CircleClose,
    stopped: VideoPause
  }
  return map[status] || VideoPause
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

onMounted(() => {
  getTableData()
  getDataSource()
})

onBeforeUnmount(() => {
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
.search-form {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.task-name {
  color: #409eff;
  cursor: pointer;
  font-weight: 500;
}

.task-name:hover {
  text-decoration: underline;
}

.task-id-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.task-id {
  font-family: Monaco, Menlo, monospace;
  font-size: 12px;
  color: #86909c;
}

.copy-icon {
  cursor: pointer;
  color: #c9cdd4;
  font-size: 14px;
}

.copy-icon:hover {
  color: #409eff;
}

.status-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
}

.status-icon {
  font-size: 14px;
}

.status-running {
  color: #409eff;
  animation: spin 1s linear infinite;
}

.status-serving,
.status-completed {
  color: #52c41a;
}

.status-failed {
  color: #f56c6c;
}

.status-pending,
.status-stopped {
  color: #faad14;
}

@keyframes spin {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.progress-done {
  color: #52c41a;
  font-weight: 500;
}

.progress-none {
  color: #c9cdd4;
}

.action-btns {
  display: flex;
  justify-content: center;
  gap: 4px;
}

.log-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.log-title {
  font-size: 14px;
  color: #4e5969;
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
</style>
