<template>
  <div class="model-test-container">
    <div class="breadcrumb">
      <span class="breadcrumb-item" @click="goBack">训练任务</span>
      <span class="breadcrumb-separator">/</span>
      <span class="breadcrumb-current">模型测试</span>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      <span>加载中...</span>
    </div>

    <template v-else>
      <!-- 任务信息卡片 -->
      <div class="detail-card task-info-card">
        <div class="card-header">
          <div class="title-section">
            <h2 class="task-title">{{ taskDetail.name }}</h2>
            <el-tag :type="getStatusType(taskDetail.status)" size="large">
              {{ getStatusText(taskDetail.status) }}
            </el-tag>
          </div>
          <div class="action-section">
            <el-button type="primary" plain @click="goToDetail">任务详情</el-button>
          </div>
        </div>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">基础模型</span>
            <span class="info-value">{{ taskDetail.baseModel }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">训练方式</span>
            <span class="info-value">{{ taskDetail.trainMethod }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">输出模型</span>
            <span class="info-value">{{ taskDetail.modelName || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">训练端口</span>
            <span class="info-value">{{ taskDetail.hostPort || '-' }}</span>
          </div>
        </div>
      </div>

      <!-- 非服务状态提示 -->
      <div v-if="taskDetail.status !== 'serving'" class="detail-card not-serving-card">
        <el-empty description="当前任务未处于服务状态">
          <template #image>
            <el-icon :size="64" color="#c0c4cc"><WarningFilled /></el-icon>
          </template>
          <p class="hint-text">请先在任务详情页启动推理服务后再进行模型测试</p>
          <el-button type="primary" @click="goToDetail">前往任务详情</el-button>
        </el-empty>
      </div>

      <!-- 模型测试区域 - 仅在服务状态时显示 -->
      <template v-else>
        <!-- 结果对比区域 -->
        <div class="detail-card result-card">
          <div class="result-header">
            <el-icon><DocumentCopy /></el-icon>
            <span>模型回复对比</span>
            <el-tag v-if="baseTestResult || loraTestResult" type="success" size="small">已生成</el-tag>
            <el-tag v-else type="info" size="small">等待测试</el-tag>
          </div>
          <div class="test-results-grid">
            <!-- 基础模型结果 -->
            <div class="test-panel result-panel">
              <div class="test-panel-header">
                <span class="test-panel-title">基础模型 (base)</span>
                <el-tag type="info" size="small">原始模型</el-tag>
              </div>
              <div class="test-panel-body">
                <div v-if="baseTestLoading" class="result-loading">
                  <el-icon class="is-loading"><Loading /></el-icon>
                  <span>正在生成回复...</span>
                </div>
                <div v-else-if="baseTestResult" class="test-result">
                  <!-- 思考过程（如果有） -->
                  <template v-if="parseThinkContent(baseTestResult).think">
                    <div class="think-section" @click="toggleBaseThink">
                      <div class="think-header">
                        <el-icon><Cpu /></el-icon>
                        <span>思考过程</span>
                        <el-icon :class="['toggle-icon', { expanded: baseThinkExpanded }]">
                          <ArrowDown />
                        </el-icon>
                      </div>
                      <el-collapse-transition>
                        <div v-show="baseThinkExpanded" class="think-content">
                          {{ parseThinkContent(baseTestResult).think }}
                        </div>
                      </el-collapse-transition>
                    </div>
                  </template>
                  <!-- 正常回复 -->
                  <div class="result-main">
                    <div class="result-label" v-if="parseThinkContent(baseTestResult).think">回复内容</div>
                    <div class="result-text">{{ parseThinkContent(baseTestResult).content }}</div>
                  </div>
                </div>
                <div v-else class="result-empty">暂无结果，请输入问题进行测试</div>
              </div>
            </div>
            <!-- LoRA 模型结果 -->
            <div class="test-panel result-panel">
              <div class="test-panel-header">
                <span class="test-panel-title">训练模型 (lora)</span>
                <el-tag type="success" size="small">微调后</el-tag>
              </div>
              <div class="test-panel-body">
                <div v-if="loraTestLoading" class="result-loading">
                  <el-icon class="is-loading"><Loading /></el-icon>
                  <span>正在生成回复...</span>
                </div>
                <div v-else-if="loraTestResult" class="test-result">
                  <!-- 思考过程（如果有） -->
                  <template v-if="parseThinkContent(loraTestResult).think">
                    <div class="think-section" @click="toggleLoraThink">
                      <div class="think-header">
                        <el-icon><Cpu /></el-icon>
                        <span>思考过程</span>
                        <el-icon :class="['toggle-icon', { expanded: loraThinkExpanded }]">
                          <ArrowDown />
                        </el-icon>
                      </div>
                      <el-collapse-transition>
                        <div v-show="loraThinkExpanded" class="think-content">
                          {{ parseThinkContent(loraTestResult).think }}
                        </div>
                      </el-collapse-transition>
                    </div>
                  </template>
                  <!-- 正常回复 -->
                  <div class="result-main">
                    <div class="result-label" v-if="parseThinkContent(loraTestResult).think">回复内容</div>
                    <div class="result-text">{{ parseThinkContent(loraTestResult).content }}</div>
                  </div>
                </div>
                <div v-else class="result-empty">暂无结果，请输入问题进行测试</div>
              </div>
            </div>
          </div>
        </div>

        <!-- 输入区域 -->
        <div class="detail-card input-card">
          <h3 class="card-title">测试输入</h3>
          <div class="test-input-area">
            <el-input
              v-model="testInput"
              type="textarea"
              :rows="4"
              placeholder="输入测试问题，如：你是谁"
              :disabled="testLoading"
              class="test-input"
            />
            <div class="test-actions">
              <el-button
                type="primary"
                :loading="testLoading"
                @click="testBothModels"
                size="large"
              >
                <el-icon v-if="!testLoading"><Promotion /></el-icon>
                同时测试两个模型
              </el-button>
              <div class="action-hints">
                <span class="test-hint">将同时发送给基础模型和训练后的 LoRA 模型</span>
                <el-button
                  v-if="baseTestResult || loraTestResult"
                  type="default"
                  @click="clearResults"
                >
                  清空结果
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <!-- 历史记录 -->
        <div class="detail-card history-card" v-if="testHistory.length > 0">
          <div class="history-header">
            <h3 class="card-title">测试历史</h3>
            <el-button type="danger" link @click="clearHistory">
              <el-icon><Delete /></el-icon>
              清空历史
            </el-button>
          </div>
          <div class="history-list">
            <div v-for="item in testHistory" :key="item.ID" class="history-item">
              <div class="history-meta">
                <div class="history-question">
                  <el-icon><User /></el-icon>
                  <span>{{ item.question }}</span>
                </div>
                <div class="history-time">
                  <el-icon><Clock /></el-icon>
                  <span>{{ formatDate(item.testTime || item.CreatedAt) }}</span>
                </div>
              </div>
              <div class="history-answers">
                <div class="answer-item">
                  <div class="answer-header">
                    <el-tag type="info" size="small">base</el-tag>
                    <el-button type="danger" link size="small" @click="deleteHistoryItem(item)">
                      <el-icon><Delete /></el-icon>
                    </el-button>
                  </div>
                  <span class="answer-text">{{ item.baseAnswer }}</span>
                </div>
                <div class="answer-item">
                  <div class="answer-header">
                    <el-tag type="success" size="small">lora</el-tag>
                  </div>
                  <span class="answer-text">{{ item.loraAnswer }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </template>
    </template>
  </div>
</template>

<script setup>
import { findTask, chatCompletion, createTestHistory, getTestHistoryList, clearTestHistory, deleteTestHistory } from '@/api/modeltraining/trainingTask'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { DocumentCopy, Loading, Promotion, User, WarningFilled, ArrowDown, Cpu, Delete, Clock } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/format'

defineOptions({
  name: 'ModelTest'
})

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const taskDetail = ref({})

// 模型测试相关
const testInput = ref('')
const testLoading = ref(false)
const baseTestResult = ref('')
const baseTestLoading = ref(false)
const loraTestResult = ref('')
const loraTestLoading = ref(false)
const testHistory = ref([])

// 结果展开状态
const baseThinkExpanded = ref(false)
const loraThinkExpanded = ref(false)

const getTaskDetail = async () => {
  loading.value = true
  try {
    const res = await findTask({ ID: route.params.id })
    if (res.code === 0) {
      taskDetail.value = res.data?.task || res.data || {}
      // 加载测试历史
      if (taskDetail.value.status === 'serving') {
        await loadTestHistory()
      }
    }
  } finally {
    loading.value = false
  }
}

// 加载测试历史
const loadTestHistory = async () => {
  try {
    const res = await getTestHistoryList({
      taskId: route.params.id,
      page: 1,
      pageSize: 20
    })
    if (res.code === 0) {
      testHistory.value = res.data?.list || []
    }
  } catch (err) {
    console.error('加载测试历史失败:', err)
  }
}

const goBack = () => {
  router.push({ name: 'trainingTask' })
}

const goToDetail = () => {
  router.push({ name: 'TrainingTaskDetail', params: { id: route.params.id } })
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

// 解析思考过程和回复内容
const parseThinkContent = (text) => {
  if (!text) return { think: '', content: '' }
  
  // 匹配 <think>...</think> 标签
  const thinkMatch = text.match(/<think>([\s\S]*?)<\/think>/)
  
  if (thinkMatch) {
    const think = thinkMatch[1].trim()
    // 移除 think 标签后的内容作为回复内容
    const content = text.replace(/<think>[\s\S]*?<\/think>/, '').trim()
    return { think, content }
  }
  
  return { think: '', content: text }
}

// 切换思考过程展开状态
const toggleBaseThink = () => {
  baseThinkExpanded.value = !baseThinkExpanded.value
}

const toggleLoraThink = () => {
  loraThinkExpanded.value = !loraThinkExpanded.value
}

// 模型测试 - 同时测试两个模型
const testBothModels = async () => {
  if (!testInput.value || !testInput.value.trim()) {
    ElMessage.warning('请输入测试问题')
    return
  }

  testLoading.value = true
  baseTestLoading.value = true
  loraTestLoading.value = true
  baseTestResult.value = ''
  loraTestResult.value = ''
  baseThinkExpanded.value = false
  loraThinkExpanded.value = false

  const question = testInput.value.trim()

  try {
    // 并行请求两个模型
    const [baseRes, loraRes] = await Promise.all([
      chatCompletion({
        id: taskDetail.value.ID,
        model: 'base',
        messages: [{ role: 'user', content: question }]
      }),
      chatCompletion({
        id: taskDetail.value.ID,
        model: 'lora',
        messages: [{ role: 'user', content: question }]
      })
    ])

    // 处理基础模型结果
    let baseAnswer = ''
    if (baseRes.code === 0 && baseRes.data) {
      baseAnswer = baseRes.data.choices?.[0]?.message?.content || '无响应内容'
      baseTestResult.value = baseAnswer
    } else {
      baseAnswer = '请求失败: ' + (baseRes.msg || '未知错误')
      baseTestResult.value = baseAnswer
    }

    // 处理 LoRA 模型结果
    let loraAnswer = ''
    if (loraRes.code === 0 && loraRes.data) {
      loraAnswer = loraRes.data.choices?.[0]?.message?.content || '无响应内容'
      loraTestResult.value = loraAnswer
    } else {
      loraAnswer = '请求失败: ' + (loraRes.msg || '未知错误')
      loraTestResult.value = loraAnswer
    }

    // 保存到数据库
    try {
      await createTestHistory({
        taskId: taskDetail.value.ID,
        question,
        baseAnswer,
        loraAnswer
      })
      // 重新加载历史记录
      await loadTestHistory()
    } catch (err) {
      console.error('保存测试历史失败:', err)
    }
  } catch (err) {
    ElMessage.error('请求失败: ' + (err.message || err))
  } finally {
    testLoading.value = false
    baseTestLoading.value = false
    loraTestLoading.value = false
  }
}

const clearResults = () => {
  baseTestResult.value = ''
  loraTestResult.value = ''
  testInput.value = ''
}

const clearHistory = async () => {
  ElMessageBox.confirm('确定要清空所有测试历史吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    try {
      const res = await clearTestHistory({ taskId: route.params.id })
      if (res.code === 0) {
        ElMessage.success('已清空')
        testHistory.value = []
      }
    } catch (err) {
      ElMessage.error('清空失败: ' + (err.message || err))
    }
  }).catch(() => {})
}

// 删除单条历史记录
const deleteHistoryItem = async (item) => {
  try {
    const res = await deleteTestHistory({ ID: item.ID })
    if (res.code === 0) {
      ElMessage.success('已删除')
      await loadTestHistory()
    }
  } catch (err) {
    ElMessage.error('删除失败: ' + (err.message || err))
  }
}

onMounted(() => {
  getTaskDetail()
})
</script>

<style scoped>
.model-test-container {
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
  margin-bottom: 20px;
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

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
  margin: 0 0 16px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
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
}

/* 非服务状态提示 */
.not-serving-card {
  text-align: center;
  padding: 60px 24px;
}

.hint-text {
  color: #86909c;
  margin-bottom: 20px;
}

/* 结果区域 */
.result-card {
  padding-bottom: 20px;
}

.result-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 20px;
  font-size: 16px;
  font-weight: 600;
  color: #1a1a1a;
}

.test-results-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.test-panel {
  border: 1px solid #e4e7ed;
  border-radius: 12px;
  overflow: hidden;
  transition: box-shadow 0.3s;
  background: #fff;
}

.test-panel:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

.result-panel {
  min-height: 200px;
  max-height: 400px;
}

.test-panel-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 18px;
  background: linear-gradient(135deg, #fafbfc 0%, #f4f6f8 100%);
  border-bottom: 1px solid #e4e7ed;
}

.test-panel-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.test-panel-body {
  padding: 18px;
  min-height: 150px;
  max-height: 320px;
  overflow-y: auto;
}

.result-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 40px;
  color: #409eff;
}

.result-empty {
  color: #909399;
  text-align: center;
  padding: 40px;
  font-size: 14px;
}

.test-result {
  background: linear-gradient(135deg, #f8f9fb 0%, #ffffff 100%);
  border-radius: 8px;
  overflow: hidden;
}

/* 思考过程样式 */
.think-section {
  margin-bottom: 12px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  overflow: hidden;
}

.think-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: #f5f7fa;
  cursor: pointer;
  font-size: 13px;
  color: #606266;
  user-select: none;
}

.think-header:hover {
  background: #eef0f3;
}

.think-header .el-icon:first-child {
  color: #909399;
}

.think-header span {
  flex: 1;
  font-weight: 500;
}

.think-content {
  padding: 12px 14px;
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  background: #fafbfc;
  border-top: 1px solid #e4e7ed;
}

/* 回复内容样式 */
.result-main {
  padding: 12px 0;
}

.result-label {
  font-size: 12px;
  color: #86909c;
  margin-bottom: 8px;
}

.result-text {
  font-size: 14px;
  color: #303133;
  line-height: 1.8;
  white-space: pre-wrap;
  word-break: break-word;
}

/* 输入区域 */
.input-card .card-title {
  margin-bottom: 16px;
}

.test-input-area .test-input {
  margin-bottom: 16px;
}

.test-input-area :deep(.el-textarea__inner) {
  border-radius: 10px;
  border: 1px solid #dcdfe6;
  padding: 14px 16px;
  font-size: 14px;
  line-height: 1.7;
  resize: vertical;
  min-height: 120px;
}

.test-input-area :deep(.el-textarea__inner:focus) {
  border-color: #409eff;
  box-shadow: 0 0 0 3px rgba(64, 158, 255, 0.1);
}

.test-actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.test-actions .el-button {
  border-radius: 8px;
  padding: 12px 24px;
  font-weight: 500;
}

.action-hints {
  display: flex;
  align-items: center;
  gap: 16px;
}

.test-hint {
  color: #86909c;
  font-size: 13px;
}

/* 历史记录 */
.history-card {
  padding-bottom: 20px;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.history-header .card-title {
  margin: 0;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.history-item {
  border: 1px solid #e4e7ed;
  border-radius: 8px;
  padding: 16px;
  background: #fafbfc;
}

.history-meta {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 12px;
}

.history-question {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 14px;
  color: #303133;
  font-weight: 500;
  flex: 1;
}

.history-question .el-icon {
  margin-top: 2px;
  color: #409eff;
}

.history-time {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #86909c;
  flex-shrink: 0;
}

.history-answers {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px;
}

.answer-item {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px;
  background: #fff;
  border-radius: 6px;
}

.answer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.answer-text {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
}

@media (max-width: 1200px) {
  .info-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .test-results-grid,
  .history-answers {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .model-test-container {
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

  .test-actions {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>