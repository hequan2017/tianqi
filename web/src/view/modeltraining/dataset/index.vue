<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="elSearchFormRef" :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
        <el-form-item label="创建日期" prop="createdAtRange">
          <el-date-picker
            v-model="searchInfo.createdAtRange"
            class="!w-380px"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

        <el-form-item label="数据集名称" prop="name">
          <el-input v-model="searchInfo.name" placeholder="请输入数据集名称" clearable />
        </el-form-item>

        <el-form-item label="数据集类型" prop="type">
          <el-select v-model="searchInfo.type" filterable placeholder="请选择类型" clearable>
            <el-option v-for="item in dataSource.type" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="训练方式" prop="trainMethod">
          <el-select v-model="searchInfo.trainMethod" filterable placeholder="请选择训练方式" clearable>
            <el-option v-for="item in dataSource.trainMethod" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="导入状态" prop="importStatus">
          <el-select v-model="searchInfo.importStatus" filterable placeholder="请选择状态" clearable>
            <el-option v-for="item in dataSource.importStatus" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>

        <el-form-item label="发布状态" prop="publishStatus">
          <el-select v-model="searchInfo.publishStatus" filterable placeholder="请选择发布状态" clearable>
            <el-option label="已发布" :value="true" />
            <el-option label="未发布" :value="false" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openDialog()">新增</el-button>
        <el-button icon="delete" :disabled="!multipleSelection.length" @click="onDelete">删除</el-button>
      </div>

      <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />

        <el-table-column sortable align="left" label="创建日期" prop="CreatedAt" width="180">
          <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>

        <el-table-column align="left" label="数据集名称" prop="name" min-width="150" />

        <el-table-column align="left" label="类型" prop="type" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.type === 'training' ? 'primary' : 'success'">
              {{ scope.row.type === 'training' ? '训练集' : '验证集' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="left" label="训练方式" prop="trainMethod" width="100" />

        <el-table-column align="left" label="数据量" prop="dataCount" width="100" />

        <el-table-column align="left" label="最新版本" prop="latestVersion" width="100" />

        <el-table-column align="left" label="导入状态" prop="importStatus" width="100">
          <template #default="scope">
            <el-tag :type="getImportStatusType(scope.row.importStatus)">
              {{ getImportStatusText(scope.row.importStatus) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="left" label="发布状态" prop="publishStatus" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.publishStatus ? 'success' : 'info'">
              {{ scope.row.publishStatus ? '已发布' : '未发布' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column align="left" label="创建者" prop="userName" width="120" />

        <el-table-column align="left" label="操作" fixed="right" min-width="180">
          <template #default="scope">
            <el-button type="primary" link icon="view" @click="openVersionDialog(scope.row)">版本管理</el-button>
            <el-button type="primary" link icon="edit" @click="updateDatasetFunc(scope.row)">编辑</el-button>
            <el-button type="danger" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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

    <!-- 新增/编辑数据集弹窗 - 一比一复刻参考页面 -->
    <el-drawer
      destroy-on-close
      :size="appStore.drawerSize"
      v-model="dialogFormVisible"
      :show-close="false"
      class="dataset-drawer"
    >
      <template #header>
        <div class="drawer-header">
          <span class="drawer-title">{{ type === 'create' ? '新增数据集' : '编辑数据集' }}</span>
          <div class="drawer-actions">
            <el-button :loading="btnLoading" type="primary" @click="enterDialog">确认</el-button>
            <el-button @click="closeDialog">取消</el-button>
          </div>
        </div>
      </template>

      <div class="dataset-form-container">
        <!-- 数据集名称 -->
        <div class="form-item">
          <div class="form-label">
            <span>数据集名称</span>
            <span class="required">*</span>
          </div>
          <div class="form-input-wrapper">
            <el-input
              v-model="formData.name"
              placeholder="请输入数据集名称"
              maxlength="50"
              show-word-limit
              clearable
            />
          </div>
        </div>

        <!-- 数据集类型 -->
        <div class="form-item">
          <div class="form-label">数据集类型</div>
          <div class="form-card-group">
            <div
              class="type-card"
              :class="{ active: formData.type === 'training' }"
              @click="formData.type = 'training'"
            >
              <div class="card-title">训练集</div>
              <div class="card-desc">模型训练的数据集，训练任务提交时可切分验证集</div>
            </div>
            <div
              class="type-card"
              :class="{ active: formData.type === 'evaluation' }"
              @click="formData.type = 'evaluation'"
            >
              <div class="card-title">评测集</div>
              <div class="card-desc">模型评测的数据集，需要满足数据集规范</div>
            </div>
          </div>
        </div>

        <!-- 数据格式/训练场景 -->
        <div class="form-item">
          <div class="form-label">数据格式</div>
          <div class="form-section">
            <div class="section-label">训练场景</div>
            <div class="tag-group">
              <div
                v-for="item in trainScenes"
                :key="item.value"
                class="tag-item"
                :class="{ active: formData.format === item.value }"
                @click="formData.format = item.value"
              >
                {{ item.label }}
              </div>
            </div>
          </div>

          <!-- 训练方式 -->
          <div class="form-section">
            <div class="section-label">训练方式</div>
            <div class="method-btn-group">
              <el-button
                v-for="item in trainMethods"
                :key="item.value"
                :type="formData.trainMethod === item.value ? 'primary' : 'default'"
                @click="formData.trainMethod = item.value"
              >
                {{ item.label }}
              </el-button>
            </div>
            <div class="method-desc">{{ getTrainMethodDesc(formData.trainMethod) }}</div>
            <div class="data-format-link" @click="showFormatDialog = true">
              <el-icon><Document /></el-icon>
              <span>数据格式说明</span>
            </div>
          </div>

          <!-- 样例下载 -->
          <div class="form-section">
            <div class="section-label">样例下载</div>
            <div class="download-group">
              <div class="download-item" @click="downloadSample('jsonl')">
                <el-icon><Download /></el-icon>
                <span>Jsonl格式</span>
              </div>
              <div class="download-item" @click="downloadSample('excel')">
                <el-icon><Download /></el-icon>
                <span>Excel格式</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 导入方式 -->
        <div class="form-item">
          <div class="form-label">导入方式</div>
          <div class="import-method">
            <el-icon><Upload /></el-icon>
            <span>本地上传</span>
          </div>
        </div>

        <!-- 文件上传 -->
        <div class="form-item">
          <el-upload
            ref="uploadRef"
            class="upload-area"
            drag
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :file-list="fileList"
            :limit="10"
            accept=".jsonl,.xls,.xlsx"
          >
            <div class="upload-content">
              <el-icon class="upload-icon"><UploadFilled /></el-icon>
              <div class="upload-text">点击或将文件拖拽到这里上传 ({{ fileList.length }}/10)</div>
              <div class="upload-tip">支持扩展名：jsonl、xls、xlsx，文件最大200MB。点击确认后上传，每个文件自动创建一个版本。</div>
            </div>
          </el-upload>
        </div>

        <!-- 是否发布 -->
        <div class="form-item">
          <div class="form-label">发布状态</div>
          <div class="publish-switch-wrapper">
            <el-switch
              v-model="formData.publishStatus"
              active-text="发布"
              inactive-text="不发布"
            />
            <span class="publish-tip">{{ formData.publishStatus ? '发布后数据集将公开可见' : '数据集默认不发布，可稍后在列表中发布' }}</span>
          </div>
        </div>
      </div>
    </el-drawer>

    <!-- 版本管理弹窗 -->
    <el-drawer destroy-on-close size="70%" v-model="versionDialogVisible" title="版本管理">
      <div class="version-manager">
        <!-- 版本列表 -->
        <div class="version-list">
          <div class="version-header">
            <span class="version-title">版本列表</span>
          </div>

          <div class="version-items">
            <div
              v-for="version in versionList"
              :key="version.ID"
              class="version-item"
              :class="{ active: selectedVersion?.ID === version.ID }"
              @click="selectVersion(version)"
            >
              <div class="version-info">
                <span class="version-name">{{ version.version }}</span>
                <el-tag size="small" :type="getVersionStatusType(version.status)">
                  {{ getVersionStatusText(version.status) }}
                </el-tag>
              </div>
              <div class="version-meta">
                <span v-if="version.dataCount">{{ version.dataCount }} 条数据</span>
                <span v-if="version.fileSize">{{ formatFileSize(version.fileSize) }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 版本详情 -->
        <div class="version-detail" v-if="selectedVersion">
          <div class="detail-header">
            <span class="detail-title">{{ selectedVersion.version }} - 版本详情</span>
            <el-button type="danger" size="small" icon="Delete" @click="deleteVersionRow(selectedVersion)">删除版本</el-button>
          </div>

          <el-descriptions :column="2" border class="detail-info">
            <el-descriptions-item label="版本号">{{ selectedVersion.version }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getVersionStatusType(selectedVersion.status)">
                {{ getVersionStatusText(selectedVersion.status) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="数据量">{{ selectedVersion.dataCount || 0 }} 条</el-descriptions-item>
            <el-descriptions-item label="文件大小">{{ formatFileSize(selectedVersion.fileSize) }}</el-descriptions-item>
            <el-descriptions-item label="创建时间" :span="2">{{ formatDate(selectedVersion.CreatedAt) }}</el-descriptions-item>
            <el-descriptions-item label="版本说明" :span="2">{{ selectedVersion.description || '暂无说明' }}</el-descriptions-item>
          </el-descriptions>

          <!-- 文件上传区域 -->
          <div class="file-section">
            <div class="section-header">
              <span class="section-title">数据文件</span>
            </div>

            <div v-if="selectedVersion.fileName" class="file-info-card">
              <div class="file-icon">
                <el-icon size="32"><Document /></el-icon>
              </div>
              <div class="file-details">
                <div class="file-name">{{ selectedVersion.fileName }}</div>
                <div class="file-meta">
                  <span>{{ formatFileSize(selectedVersion.fileSize) }}</span>
                  <span>{{ selectedVersion.filePath }}</span>
                </div>
              </div>
              <div class="file-actions">
                <el-button type="primary" size="small" @click="downloadVersionFile(selectedVersion)">下载</el-button>
                <el-button type="danger" size="small" @click="deleteVersionFile(selectedVersion)">删除文件</el-button>
              </div>
            </div>

            <div v-else class="upload-section">
              <el-upload
                class="version-upload"
                drag
                :action="versionUploadAction"
                :headers="uploadHeaders"
                :data="{ versionId: selectedVersion.ID }"
                :show-file-list="false"
                :on-success="handleVersionUploadSuccess"
                :on-error="handleVersionUploadError"
                accept=".jsonl,.xls,.xlsx"
              >
                <div class="upload-content">
                  <el-icon class="upload-icon"><UploadFilled /></el-icon>
                  <div class="upload-text">点击或拖拽文件到此上传</div>
                  <div class="upload-tip">支持 jsonl、xls、xlsx 格式，最大 200MB</div>
                </div>
              </el-upload>
            </div>
          </div>
        </div>

        <div v-else class="version-empty">
          <el-empty description="请选择一个版本查看详情" />
        </div>
      </div>
    </el-drawer>

    <!-- 数据格式说明弹窗 -->
    <el-dialog v-model="showFormatDialog" title="数据格式说明" width="700px">
      <div class="format-content">
        <h4>SFT 数据格式</h4>
        <p>一种多轮对话的训练数据，支持文本生成类模型的微调训练。</p>
        <pre class="code-block">
{
  "messages": [
    {"role": "system", "content": "你是一个有帮助的助手。"},
    {"role": "user", "content": "用户问题"},
    {"role": "assistant", "content": "助手回答"}
  ]
}</pre>

        <h4>DPO 数据格式</h4>
        <p>偏好对齐训练数据，包含chosen和rejected两个回答。</p>
        <pre class="code-block">
{
  "chosen": [{"role": "user", "content": "问题"}, {"role": "assistant", "content": "优质回答"}],
  "rejected": [{"role": "user", "content": "问题"}, {"role": "assistant", "content": "较差回答"}]
}</pre>

        <h4>CPT 数据格式</h4>
        <p>持续预训练数据，用于领域知识注入。</p>
        <pre class="code-block">
{
  "text": "领域知识文本内容..."
}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import {
  createDataset,
  deleteDataset,
  deleteDatasetByIds,
  updateDataset,
  findDataset,
  getDatasetList,
  getDatasetDataSource,
  getVersionList,
  publishDataset,
  uploadDatasetFile
} from '@/api/modeltraining/dataset'

import { formatDate } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted, computed } from 'vue'
import { useAppStore } from '@/pinia'
import { useUserStore } from '@/pinia/modules/user'
import { Document, Download, Upload, UploadFilled } from '@element-plus/icons-vue'

defineOptions({
  name: 'Dataset'
})

const appStore = useAppStore()
const userStore = useUserStore()
const btnLoading = ref(false)
const uploadRef = ref()

// 上传配置（版本管理区仍使用 action 上传）
const uploadHeaders = computed(() => ({ 'x-token': userStore.token }))

// 版本上传配置
const versionUploadAction = computed(() => {
  return import.meta.env.VITE_BASE_API + '/modeltraining/dataset/uploadVersionFile'
})

// 选中的版本
const selectedVersion = ref(null)

// 训练场景选项
const trainScenes = ref([
  { label: '文本生成', value: '文本生成' },
  { label: '图片理解', value: '图片理解' },
  { label: '图生视频(首帧)', value: '图生视频(首帧)' },
  { label: '图生视频(首尾帧)', value: '图生视频(首尾帧)' }
])

// 训练方式选项
const trainMethods = ref([
  { label: 'SFT', value: 'SFT' },
  { label: 'DPO', value: 'DPO' },
  { label: 'CPT', value: 'CPT' }
])

// 数据源
const dataSource = ref({
  type: [
    { label: '训练集', value: 'training' },
    { label: '验证集', value: 'evaluation' }
  ],
  trainMethod: [
    { label: 'SFT', value: 'SFT' },
    { label: 'DPO', value: 'DPO' },
    { label: 'CPT', value: 'CPT' }
  ],
  importStatus: [
    { label: '待导入', value: 'pending' },
    { label: '导入成功', value: 'success' },
    { label: '导入失败', value: 'failed' }
  ]
})

// 表单数据
const formData = ref({
  name: '',
  type: 'training',
  format: '文本生成',
  trainMethod: 'SFT',
  storagePath: '',
  description: '',
  publishStatus: false
})

// 文件列表
const fileList = ref([])

// 数据格式说明弹窗
const showFormatDialog = ref(false)

// 验证规则
const rule = reactive({
  name: [{ required: true, message: '请输入数据集名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择数据集类型', trigger: 'change' }],
  trainMethod: [{ required: true, message: '请选择训练方式', trigger: 'change' }]
})

const elFormRef = ref()
const elSearchFormRef = ref()

// 表格控制
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})

// 多选
const multipleSelection = ref([])
const handleSelectionChange = (val) => {
  multipleSelection.value = val
}

// 弹窗控制
const type = ref('')
const dialogFormVisible = ref(false)
const versionDialogVisible = ref(false)
const currentDataset = ref(null)
const versionList = ref([])

// 获取表格数据
const getTableData = async () => {
  const table = await getDatasetList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

// 获取数据源
const getDataSource = async () => {
  const res = await getDatasetDataSource()
  if (res.code === 0) {
    dataSource.value = { ...dataSource.value, ...res.data }
  }
}

// 搜索
const onSubmit = () => {
  page.value = 1
  getTableData()
}

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 打开弹窗
const openDialog = () => {
  type.value = 'create'
  formData.value = {
    name: '',
    type: 'training',
    format: '文本生成',
    trainMethod: 'SFT',
    storagePath: '',
    description: '',
    publishStatus: false
  }
  fileList.value = []
  dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
  dialogFormVisible.value = false
  fileList.value = []
}

// 确定提交
const enterDialog = async () => {
  if (!formData.value.name) {
    ElMessage.warning('请输入数据集名称')
    return
  }
  if (!formData.value.trainMethod) {
    ElMessage.warning('请选择训练方式')
    return
  }

  btnLoading.value = true
  let res
  let datasetId = 0
  if (type.value === 'create') {
    res = await createDataset(formData.value)
    datasetId = res?.data?.ID || 0
  } else {
    res = await updateDataset(formData.value)
    datasetId = formData.value.ID || 0
  }
  btnLoading.value = false
  if (res.code === 0) {
    if (fileList.value.length > 0 && datasetId > 0) {
      const { successCount, failCount } = await uploadFilesAsVersions(datasetId)
      if (failCount > 0) {
        ElMessage.warning(`数据集保存成功，文件上传成功 ${successCount} 个，失败 ${failCount} 个`)
      } else {
        ElMessage.success(`操作成功，已上传 ${successCount} 个文件并创建版本`)
      }
    } else {
      ElMessage.success('操作成功')
    }
    closeDialog()
    getTableData()
  }
}

// 编辑
const updateDatasetFunc = async (row) => {
  const res = await findDataset({ ID: row.ID })
  if (res.code === 0) {
    formData.value = res.data
    type.value = 'update'
    dialogFormVisible.value = true
  }
}

// 删除
const deleteRow = (row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const res = await deleteDataset({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === 1 && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 批量删除
const onDelete = async () => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    const IDs = multipleSelection.value.map(item => item.ID)
    const res = await deleteDatasetByIds({ IDs })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      if (tableData.value.length === IDs.length && page.value > 1) {
        page.value--
      }
      getTableData()
    }
  })
}

// 发布数据集
const handlePublish = async (row) => {
  ElMessageBox.confirm('确定要发布此数据集吗？发布后将公开可见。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'info'
  }).then(async () => {
    const res = await publishDataset({ ID: row.ID })
    if (res.code === 0) {
      ElMessage.success('发布成功')
      getTableData()
    }
  })
}

// 打开版本管理弹窗
const openVersionDialog = async (row) => {
  currentDataset.value = row
  versionDialogVisible.value = true
  await loadVersionList(row.ID)
}

// 加载版本列表
const loadVersionList = async (datasetId) => {
  const res = await getVersionList({ datasetId, page: 1, pageSize: 100 })
  if (res.code === 0) {
    versionList.value = res.data.list || []
  }
}

// 删除版本
const deleteVersionRow = (row) => {
  ElMessageBox.confirm('确定要删除此版本吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    // TODO: 调用删除版本API
    ElMessage.success('删除成功')
    await loadVersionList(currentDataset.value.ID)
  })
}

// 文件上传处理
const handleFileChange = (file, list) => {
  fileList.value = list
}

const handleFileRemove = (file, list) => {
  fileList.value = list
}

const uploadFilesAsVersions = async (datasetId) => {
  let successCount = 0
  let failCount = 0
  for (const item of fileList.value) {
    const raw = item.raw
    if (!raw) continue
    const form = new FormData()
    form.append('file', raw)
    form.append('datasetId', String(datasetId))
    try {
      const res = await uploadDatasetFile(form)
      if (res.code === 0) {
        successCount++
      } else {
        failCount++
      }
    } catch (e) {
      failCount++
    }
  }
  return { successCount, failCount }
}

// 下载样例
const downloadSample = (type) => {
  // TODO: 实现样例下载
  ElMessage.info(`下载${type}格式样例`)
}

// 获取训练方式描述
const getTrainMethodDesc = (method) => {
  const descMap = {
    SFT: '一种多轮对话的训练数据，支持文本生成类模型的微调训练',
    DPO: '偏好对齐训练，需要提供chosen和rejected两种回答',
    CPT: '持续预训练，用于领域知识注入'
  }
  return descMap[method] || ''
}

// 辅助函数
const getImportStatusType = (status) => {
  const map = { pending: 'warning', success: 'success', failed: 'danger' }
  return map[status] || 'info'
}

const getImportStatusText = (status) => {
  const map = { pending: '待导入', success: '导入成功', failed: '导入失败' }
  return map[status] || status
}

const getVersionStatusType = (status) => {
  const map = { pending: 'warning', uploading: 'primary', success: 'success', failed: 'danger' }
  return map[status] || 'info'
}

const getVersionStatusText = (status) => {
  const map = { pending: '待处理', uploading: '上传中', success: '成功', failed: '失败' }
  return map[status] || status
}

const formatFileSize = (bytes) => {
  if (!bytes || bytes === 0) return '-'
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB'
  if (bytes < 1024 * 1024 * 1024) return (bytes / 1024 / 1024).toFixed(2) + ' MB'
  return (bytes / 1024 / 1024 / 1024).toFixed(2) + ' GB'
}

onMounted(() => {
  getTableData()
  getDataSource()
})
</script>

<style scoped>
.dataset-drawer :deep(.el-drawer__header) {
  margin-bottom: 0;
  padding: 16px 24px;
  border-bottom: 1px solid #e8e8e8;
}

.drawer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.drawer-title {
  font-size: 18px;
  font-weight: 600;
  color: #1a1a1a;
}

.drawer-actions {
  display: flex;
  gap: 8px;
}

.dataset-form-container {
  padding: 24px;
}

.form-item {
  margin-bottom: 24px;
}

.form-label {
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 8px;
}

.form-label .required {
  color: #f56c6c;
  margin-left: 4px;
}

.form-input-wrapper {
  width: 100%;
}

/* 卡片选择样式 */
.form-card-group {
  display: flex;
  gap: 16px;
  margin-top: 8px;
}

.type-card {
  flex: 1;
  padding: 16px;
  border: 1px solid #e8e8e8;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
}

.type-card:hover {
  border-color: #409eff;
}

.type-card.active {
  border-color: #409eff;
  background-color: #ecf5ff;
}

.card-title {
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.card-desc {
  font-size: 12px;
  color: #86909c;
  line-height: 1.5;
}

/* 表单区块样式 */
.form-section {
  margin-top: 16px;
}

.section-label {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 8px;
}

/* 标签选择样式 */
.tag-group {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-item {
  padding: 6px 16px;
  border: 1px solid #e8e8e8;
  border-radius: 4px;
  font-size: 14px;
  color: #4e5969;
  cursor: pointer;
  transition: all 0.3s;
}

.tag-item:hover {
  border-color: #409eff;
  color: #409eff;
}

.tag-item.active {
  border-color: #409eff;
  background-color: #409eff;
  color: #fff;
}

/* 训练方式按钮组 */
.method-btn-group {
  display: flex;
  gap: 8px;
}

.method-desc {
  margin-top: 8px;
  font-size: 12px;
  color: #86909c;
  line-height: 1.5;
}

.data-format-link {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-top: 8px;
  font-size: 13px;
  color: #409eff;
  cursor: pointer;
}

.data-format-link:hover {
  color: #66b1ff;
}

/* 样例下载样式 */
.download-group {
  display: flex;
  gap: 16px;
}

.download-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #409eff;
  cursor: pointer;
}

.download-item:hover {
  color: #66b1ff;
}

/* 存储位置样式 */
.storage-info {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background-color: #f7f8fa;
  border-radius: 4px;
  font-size: 14px;
  color: #4e5969;
}

/* 导入方式样式 */
.import-notice {
  display: flex;
  gap: 12px;
  padding: 16px;
  background-color: #fffbe6;
  border-radius: 8px;
  margin-bottom: 12px;
}

.notice-icon {
  color: #faad14;
  font-size: 20px;
}

.notice-title {
  font-size: 14px;
  font-weight: 500;
  color: #1a1a1a;
  margin-bottom: 4px;
}

.notice-desc {
  font-size: 12px;
  color: #86909c;
  line-height: 1.6;
}

.import-method {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 16px;
  background-color: #f7f8fa;
  border-radius: 4px;
  font-size: 14px;
  color: #4e5969;
}

/* 上传区域样式 */
.upload-area {
  width: 100%;
}

.upload-area :deep(.el-upload-dragger) {
  width: 100%;
  height: auto;
  min-height: 160px;
  border: 1px dashed #c9cdd4;
  border-radius: 8px;
  background-color: #f7f8fa;
}

.upload-area :deep(.el-upload-dragger:hover) {
  border-color: #409eff;
}

.upload-content {
  padding: 32px;
  text-align: center;
}

.upload-icon {
  font-size: 48px;
  color: #c9cdd4;
}

.upload-text {
  margin-top: 16px;
  font-size: 14px;
  color: #4e5969;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #86909c;
}

/* 数据格式说明弹窗样式 */
.format-content {
  line-height: 1.8;
}

.format-content h4 {
  margin: 16px 0 8px;
  color: #1a1a1a;
}

.format-content h4:first-child {
  margin-top: 0;
}

.format-content p {
  color: #4e5969;
  font-size: 14px;
}

.code-block {
  background-color: #f7f8fa;
  padding: 12px 16px;
  border-radius: 6px;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  overflow-x: auto;
  color: #4e5969;
}

/* 发布开关样式 */
.publish-switch-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
}

.publish-tip {
  font-size: 12px;
  color: #86909c;
}
</style>
