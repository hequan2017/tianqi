<template>
  <div class="create-task-container">
    <div class="breadcrumb">
      <span class="breadcrumb-item" @click="goBack">训练任务</span>
      <span class="breadcrumb-separator">/</span>
      <span class="breadcrumb-current">创建训练任务</span>
    </div>

    <h1 class="page-title">创建训练任务</h1>
    <p class="page-subtitle">面向 LoRA SFT 流程，先填写基础信息，再按需调整训练参数。</p>

    <div class="form-content">
      <el-form ref="formRef" :model="formData" label-width="120px">
        <div class="form-section">
          <div class="section-title">基础配置</div>

          <el-form-item label="任务名称" required>
            <el-input v-model="formData.name" placeholder="请输入任务名称" maxlength="50" show-word-limit clearable />
          </el-form-item>
          <p class="field-desc">用于任务列表和日志检索，建议与业务目标保持一致。</p>

          <el-form-item label="基础模型" required>
            <el-select v-model="formData.baseModel" filterable placeholder="请选择基础模型" style="width: 100%">
              <el-option v-for="item in dataSource.baseModel" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <p class="field-desc">训练将基于该模型进行 LoRA 微调，建议按显存选择模型规模。</p>

          <el-form-item label="输出模型名" required>
            <el-input v-model="formData.modelName" placeholder="请输入输出模型名" clearable />
          </el-form-item>
          <p class="field-desc">作为训练产物标识，默认值为 <code>swift-robot</code>。</p>

          <el-form-item label="输出数量">
            <el-input-number v-model="formData.outputCount" :min="1" :max="20" />
          </el-form-item>
          <p class="field-desc">单次任务预计输出模型版本数量，通常保持 <code>1</code>。</p>

          <el-form-item label="检查点间隔">
            <div class="inline-group">
              <el-input-number v-model="formData.checkpointInterval" :min="1" :max="1000" />
              <el-select v-model="formData.checkpointUnit" style="width: 120px">
                <el-option label="epoch" value="epoch" />
                <el-option label="step" value="step" />
              </el-select>
            </div>
          </el-form-item>
          <p class="field-desc">控制中间权重保存频率，避免意外中断导致训练成果丢失。</p>

          <el-form-item label="备注">
            <el-input v-model="formData.remark" type="textarea" :rows="3" placeholder="可选" />
          </el-form-item>
          <p class="field-desc">记录实验说明、数据版本或超参策略，便于复盘。</p>
        </div>

        <div class="form-section">
          <div class="section-title">数据集配置</div>

          <el-form-item label="训练数据集">
            <div class="dataset-row">
              <el-select
                v-model="formData.trainDatasetId"
                filterable
                clearable
                :loading="loadingDatasets"
                placeholder="选择训练数据集"
                style="width: 45%"
                @change="onTrainDatasetChange"
              >
                <el-option v-for="item in datasetOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select
                v-model="formData.trainVersionId"
                filterable
                clearable
                :disabled="!formData.trainDatasetId"
                :loading="loadingTrainVersions"
                placeholder="选择训练集版本"
                style="width: 45%"
              >
                <el-option v-for="item in trainVersionOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </div>
          </el-form-item>
          <p class="field-desc">先选择训练数据集，再选择该数据集下的具体版本。</p>

          <el-form-item label="验证数据集">
            <div class="dataset-row">
              <el-select
                v-model="formData.valDatasetId"
                filterable
                clearable
                :loading="loadingDatasets"
                placeholder="选择验证数据集（可选）"
                style="width: 45%"
                @change="onValDatasetChange"
              >
                <el-option v-for="item in datasetOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
              <el-select
                v-model="formData.valVersionId"
                filterable
                clearable
                :disabled="!formData.valDatasetId"
                :loading="loadingValVersions"
                placeholder="选择验证集版本"
                style="width: 45%"
              >
                <el-option v-for="item in valVersionOptions" :key="item.value" :label="item.label" :value="item.value" />
              </el-select>
            </div>
          </el-form-item>
          <p class="field-desc">验证数据集可不选；若选择了数据集，需要同步选择版本。</p>
        </div>

        <div class="form-section">
          <div class="section-title">可调训练参数</div>
          <div class="param-grid">
            <div class="param-item">
              <el-form-item label="Batch Size"><el-input-number v-model="formData.batchSize" :min="1" :max="1024" /></el-form-item>
              <p class="field-desc">每张卡每步训练样本数。值越大吞吐越高，但显存占用越大。<span class="range-note">范围：1 ~ 1024</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="Learning Rate"><el-input-number v-model="formData.learningRate" :min="0.0000001" :max="1" :step="0.00001" :precision="7" /></el-form-item>
              <p class="field-desc">学习率，决定参数更新步长。过大易震荡，过小收敛慢。<span class="range-note">范围：1e-7 ~ 1</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="Epochs"><el-input-number v-model="formData.nEpochs" :min="1" :max="200" /></el-form-item>
              <p class="field-desc">完整遍历训练集次数。默认 <code>1</code> 适合快速迭代验证。<span class="range-note">范围：1 ~ 200</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="Eval Steps"><el-input-number v-model="formData.evalSteps" :min="1" :max="10000" /></el-form-item>
              <p class="field-desc">每隔多少步进行一次评估，数值越小评估越频繁。<span class="range-note">范围：1 ~ 10000</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="LoRA Alpha"><el-input-number v-model="formData.loraAlpha" :min="1" :max="1024" /></el-form-item>
              <p class="field-desc">LoRA 缩放系数，配合 <code>LoRA Rank</code> 决定更新幅度。<span class="range-note">范围：1 ~ 1024</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="LoRA Dropout"><el-input-number v-model="formData.loraDropout" :min="0" :max="1" :step="0.01" /></el-form-item>
              <p class="field-desc">LoRA 分支 Dropout 比例，用于缓解过拟合。<span class="range-note">范围：0 ~ 1</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="LoRA Rank"><el-input-number v-model="formData.loraRank" :min="1" :max="256" /></el-form-item>
              <p class="field-desc">LoRA 低秩矩阵秩。值越大表达能力越强，训练成本也更高。<span class="range-note">范围：1 ~ 256</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="LR Scheduler">
                <el-select v-model="formData.lrSchedulerType" style="width: 100%">
                  <el-option label="linear" value="linear" />
                  <el-option label="cosine" value="cosine" />
                  <el-option label="constant" value="constant" />
                  <el-option label="constant_with_warmup" value="constant_with_warmup" />
                </el-select>
              </el-form-item>
              <p class="field-desc">学习率调度策略，默认 <code>cosine</code> 在中后期通常更平滑。<span class="range-note">可选：linear / cosine / constant / constant_with_warmup</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="Max Length"><el-input-number v-model="formData.maxLength" :min="128" :max="131072" /></el-form-item>
              <p class="field-desc">单样本最大 token 长度，超过会被截断。越大越耗显存。<span class="range-note">范围：128 ~ 131072</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="Warmup Ratio"><el-input-number v-model="formData.warmupRatio" :min="0" :max="1" :step="0.01" /></el-form-item>
              <p class="field-desc">预热阶段占比。预热可减少训练初期不稳定。<span class="range-note">范围：0 ~ 1</span></p>
            </div>
            <div class="param-item">
              <el-form-item label="Weight Decay"><el-input-number v-model="formData.weightDecay" :min="0" :max="1" :step="0.01" /></el-form-item>
              <p class="field-desc">权重衰减系数（L2），提升泛化能力，防止过拟合。<span class="range-note">范围：0 ~ 1</span></p>
            </div>
          </div>
        </div>

        <div class="form-section fixed-section">
          <div class="section-title">固定参数说明</div>
          <div class="fixed-grid">
            <div class="fixed-item"><span class="k">torch_dtype</span><span class="v">bfloat16，降低显存占用并保持训练稳定性。</span></div>
            <div class="fixed-item"><span class="k">target_modules</span><span class="v">all-linear，对线性层注入 LoRA。</span></div>
            <div class="fixed-item"><span class="k">gradient_accumulation_steps</span><span class="v">16，配合 batch size=1 提升等效批量。</span></div>
            <div class="fixed-item"><span class="k">save_steps / save_total_limit</span><span class="v">50 / 2，控制检查点频率与数量。</span></div>
            <div class="fixed-item"><span class="k">logging_steps</span><span class="v">5，日志刷新更及时。</span></div>
            <div class="fixed-item"><span class="k">output_dir</span><span class="v">output，训练结果目录。</span></div>
            <div class="fixed-item"><span class="k">system</span><span class="v"><code>You are a helpful assistant.</code>，系统指令模板。</span></div>
            <div class="fixed-item"><span class="k">dataloader_num_workers</span><span class="v">4，数据加载并行线程数。</span></div>
            <div class="fixed-item"><span class="k">model_author</span><span class="v">swift，模型作者标识。</span></div>
          </div>
        </div>
      </el-form>

      <div class="form-actions">
        <el-button @click="resetParams">恢复默认参数</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitForm">开始训练</el-button>
        <el-button @click="goBack">取消</el-button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { createTask, getTaskDataSource, getDefaultParams } from '@/api/modeltraining/trainingTask'
import { getDatasetList, getVersionList } from '@/api/modeltraining/dataset'
import { ElMessage } from 'element-plus'
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({
  name: 'CreateTrainingTask'
})

const router = useRouter()
const formRef = ref()
const submitLoading = ref(false)
const loadingDatasets = ref(false)
const loadingTrainVersions = ref(false)
const loadingValVersions = ref(false)
const datasetOptions = ref([])
const trainVersionOptions = ref([])
const valVersionOptions = ref([])

const dataSource = ref({
  baseModel: [
    { label: 'Qwen/Qwen3-1.7B', value: 'Qwen/Qwen3-1.7B' },
    { label: 'Qwen/Qwen3.5-0.8B', value: 'Qwen/Qwen3.5-0.8B' },
    { label: 'Qwen/Qwen2.5-7B-Instruct', value: 'Qwen/Qwen2.5-7B-Instruct' },
    { label: 'Qwen/Qwen3-7B', value: 'Qwen/Qwen3-7B' },
    { label: 'Qwen/Qwen3-14B', value: 'Qwen/Qwen3-14B' },
    { label: 'Llama3-8B', value: 'Llama3-8B' }
  ]
})

const formData = reactive({
  name: '',
  baseModel: 'Qwen/Qwen3-1.7B',
  trainMethod: 'SFT',
  trainType: 'efficient',
  modelName: 'swift-robot',
  outputCount: 1,
  checkpointInterval: 1,
  checkpointUnit: 'epoch',
  remark: '',
  batchSize: 1,
  learningRate: 0.0001,
  nEpochs: 1,
  evalSteps: 50,
  loraAlpha: 32,
  loraDropout: 0.05,
  loraRank: 8,
  lrSchedulerType: 'cosine',
  maxLength: 2048,
  warmupRatio: 0.05,
  weightDecay: 0.01,
  trainDatasetId: null,
  trainVersionId: null,
  valDatasetId: null,
  valVersionId: null,
  valSplitRatio: 10
})

const initialParamSnapshot = ref({})

const getVersionLabel = (row) => {
  return row.version || row.name || `v${row.ID}`
}

const loadDatasetOptions = async () => {
  loadingDatasets.value = true
  try {
    const res = await getDatasetList({
      page: 1,
      pageSize: 200,
      publishStatus: true,
      importStatus: 'success'
    })
    if (res.code === 0) {
      datasetOptions.value = (res.data?.list || []).map((d) => ({
        label: d.name,
        value: d.ID,
        raw: d
      }))
    }
  } finally {
    loadingDatasets.value = false
  }
}

const loadVersionsByDataset = async (datasetId, target) => {
  if (!datasetId) {
    if (target === 'train') trainVersionOptions.value = []
    if (target === 'val') valVersionOptions.value = []
    return
  }
  if (target === 'train') loadingTrainVersions.value = true
  if (target === 'val') loadingValVersions.value = true
  try {
    const res = await getVersionList({ datasetId, page: 1, pageSize: 200 })
    if (res.code === 0) {
      const list = (res.data?.list || []).filter((v) => v.status === undefined || v.status === 'success')
      const options = list.map((v) => ({ label: getVersionLabel(v), value: v.ID, raw: v }))
      if (target === 'train') trainVersionOptions.value = options
      if (target === 'val') valVersionOptions.value = options
    }
  } finally {
    if (target === 'train') loadingTrainVersions.value = false
    if (target === 'val') loadingValVersions.value = false
  }
}

const onTrainDatasetChange = async (datasetId) => {
  formData.trainVersionId = null
  await loadVersionsByDataset(datasetId, 'train')
}

const onValDatasetChange = async (datasetId) => {
  formData.valVersionId = null
  await loadVersionsByDataset(datasetId, 'val')
}

const loadDataSource = async () => {
  const res = await getTaskDataSource()
  if (res.code === 0) {
    dataSource.value = { ...dataSource.value, ...res.data }
  }
}

const loadDefaultParams = async () => {
  const res = await getDefaultParams()
  if (res.code === 0 && res.data) {
    Object.keys(res.data).forEach((k) => {
      if (Object.prototype.hasOwnProperty.call(formData, k)) {
        formData[k] = res.data[k]
      }
    })
  }
  initialParamSnapshot.value = {
    batchSize: formData.batchSize,
    learningRate: formData.learningRate,
    nEpochs: formData.nEpochs,
    evalSteps: formData.evalSteps,
    loraAlpha: formData.loraAlpha,
    loraDropout: formData.loraDropout,
    loraRank: formData.loraRank,
    lrSchedulerType: formData.lrSchedulerType,
    maxLength: formData.maxLength,
    warmupRatio: formData.warmupRatio,
    weightDecay: formData.weightDecay
  }
}

const resetParams = () => {
  Object.assign(formData, initialParamSnapshot.value)
  ElMessage.success('已恢复默认参数')
}

const goBack = () => {
  router.push({ name: 'trainingTask' })
}

const submitForm = async () => {
  if (!formData.name?.trim()) {
    ElMessage.warning('请输入任务名称')
    return
  }
  if (!formData.baseModel) {
    ElMessage.warning('请选择基础模型')
    return
  }
  if (!formData.modelName?.trim()) {
    ElMessage.warning('请输入输出模型名')
    return
  }
  if ((formData.trainDatasetId && !formData.trainVersionId) || (!formData.trainDatasetId && formData.trainVersionId)) {
    ElMessage.warning('训练数据集与版本需要同时选择')
    return
  }
  if ((formData.valDatasetId && !formData.valVersionId) || (!formData.valDatasetId && formData.valVersionId)) {
    ElMessage.warning('验证数据集与版本需要同时选择')
    return
  }
  const rangeRules = [
    [formData.batchSize >= 1 && formData.batchSize <= 1024, 'Batch Size 范围应为 1~1024'],
    [formData.learningRate >= 0.0000001 && formData.learningRate <= 1, 'Learning Rate 范围应为 1e-7~1'],
    [formData.nEpochs >= 1 && formData.nEpochs <= 200, 'Epochs 范围应为 1~200'],
    [formData.evalSteps >= 1 && formData.evalSteps <= 10000, 'Eval Steps 范围应为 1~10000'],
    [formData.loraAlpha >= 1 && formData.loraAlpha <= 1024, 'LoRA Alpha 范围应为 1~1024'],
    [formData.loraDropout >= 0 && formData.loraDropout <= 1, 'LoRA Dropout 范围应为 0~1'],
    [formData.loraRank >= 1 && formData.loraRank <= 256, 'LoRA Rank 范围应为 1~256'],
    [formData.maxLength >= 128 && formData.maxLength <= 131072, 'Max Length 范围应为 128~131072'],
    [formData.warmupRatio >= 0 && formData.warmupRatio <= 1, 'Warmup Ratio 范围应为 0~1'],
    [formData.weightDecay >= 0 && formData.weightDecay <= 1, 'Weight Decay 范围应为 0~1']
  ]
  const firstInvalid = rangeRules.find((rule) => !rule[0])
  if (firstInvalid) {
    ElMessage.warning(firstInvalid[1])
    return
  }

  submitLoading.value = true
  try {
    const payload = {
      ...formData,
      learningRate: Number(formData.learningRate) || 0.0001
    }
    const res = await createTask(payload)
    if (res.code === 0) {
      ElMessage.success('创建成功，已开始训练')
      router.push({ name: 'trainingTask' })
    }
  } finally {
    submitLoading.value = false
  }
}

onMounted(async () => {
  await loadDatasetOptions()
  await loadDataSource()
  await loadDefaultParams()
})
</script>

<style scoped>
.create-task-container {
  padding: 24px;
  background-color: #f3f6fb;
  min-height: 100vh;
}

.breadcrumb {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  font-size: 14px;
}

.breadcrumb-item {
  color: #2563eb;
  cursor: pointer;
}

.breadcrumb-separator {
  margin: 0 8px;
  color: #94a3b8;
}

.breadcrumb-current {
  color: #475569;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  margin-bottom: 8px;
  color: #0f172a;
}

.page-subtitle {
  margin: 0 0 20px;
  color: #64748b;
  font-size: 14px;
}

.form-content {
  background:
    linear-gradient(135deg, rgba(37, 99, 235, 0.06), rgba(16, 185, 129, 0.04)),
    #fff;
  border: 1px solid #dbe5f0;
  border-radius: 14px;
  padding: 24px;
  box-shadow: 0 8px 30px rgba(15, 23, 42, 0.06);
}

.form-section {
  padding: 18px 16px;
  border: 1px solid #e2e8f0;
  border-radius: 12px;
  background: #fff;
  margin-bottom: 16px;
}

.section-title {
  font-size: 16px;
  font-weight: 700;
  margin-bottom: 10px;
  color: #0f172a;
}

.inline-group {
  display: flex;
  gap: 12px;
  align-items: center;
}

.dataset-row {
  display: flex;
  gap: 12px;
  width: 100%;
}

.param-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(260px, 1fr));
  gap: 8px 18px;
}

.param-item :deep(.el-form-item) {
  margin-bottom: 6px;
}

.field-desc {
  margin: 0 0 12px 120px;
  font-size: 12px;
  line-height: 1.5;
  color: #64748b;
}

.range-note {
  margin-left: 6px;
  color: #2563eb;
}

.fixed-section {
  background: linear-gradient(180deg, #f8fafc, #ffffff);
}

.fixed-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(280px, 1fr));
  gap: 8px 14px;
}

.fixed-item {
  display: flex;
  gap: 8px;
  align-items: flex-start;
  font-size: 13px;
  color: #334155;
}

.fixed-item .k {
  min-width: 180px;
  font-family: Monaco, Menlo, monospace;
  color: #0f172a;
  background: #e2e8f0;
  border-radius: 6px;
  padding: 2px 6px;
}

.fixed-item .v {
  line-height: 1.5;
}

.form-actions {
  margin-top: 24px;
  display: flex;
  gap: 10px;
  justify-content: flex-end;
}

@media (max-width: 900px) {
  .dataset-row {
    flex-direction: column;
  }

  .dataset-row .el-select {
    width: 100% !important;
  }

  .param-grid {
    grid-template-columns: 1fr;
  }

  .fixed-grid {
    grid-template-columns: 1fr;
  }

  .field-desc {
    margin-left: 0;
  }
}
</style>
