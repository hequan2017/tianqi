import re

# 读取文件
with open(r'D:\devops\docker-gpu-manage\web\src\view\modeltraining\training\detail.vue', 'r', encoding='utf-8') as f:
    content = f.read()

# 找到模型测试区域的开始和结束
start_marker = '<!-- 模型测试区域 - 仅在服务状态时显示 -->'
end_marker = '</el-collapse-item>\n        </el-collapse>\n      </div>\n    </template>'

start_idx = content.find(start_marker)
if start_idx == -1:
    print("Start marker not found!")
    exit(1)

# 找到这个区域的结束
end_idx = content.find('</div>\n    </template>', start_idx)
if end_idx == -1:
    print("End marker not found!")
    exit(1)

# 新的内容
new_content = '''<!-- 模型测试区域 - 仅在服务状态时显示 -->
      <div class="detail-card model-test-section" v-if="taskDetail.status === 'serving'">
        <el-collapse v-model="activeTestPanel" class="test-collapse">
          <el-collapse-item name="test">
            <template #title>
              <h3 class="card-title" style="margin: 0">模型测试</h3>
            </template>

            <!-- 结果对比区域 - 默认显示 -->
            <div class="result-section">
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
                      <div class="result-content">{{ baseTestResult }}</div>
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
                      <div class="result-content">{{ loraTestResult }}</div>
                    </div>
                    <div v-else class="result-empty">暂无结果，请输入问题进行测试</div>
                  </div>
                </div>
              </div>
            </div>

            <!-- 公共输入区域 -->
            <div class="test-input-area">
              <el-input
                v-model="testInput"
                type="textarea"
                :rows="3"
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
                <span class="test-hint">将同时发送给基础模型和训练后的 LoRA 模型</span>
              </div>
            </div>
          </el-collapse-item>
        </el-collapse>
      </div>
    </template>'''

# 找到正确的结束位置
# 从 start_idx 开始查找 </div>\n    </template>
temp_idx = start_idx
lines = content[start_idx:].split('\n')
line_count = 0
found = False
for i, line in enumerate(lines):
    if '</div>' in line and i + 1 < len(lines) and '    </template>' in lines[i + 1]:
        end_pos = start_idx + sum(len(l) + 1 for l in lines[:i + 2])
        # 确保不是太早结束
        if 'test-input-area' in content[start_idx:end_pos]:
            found = True
            break

if not found:
    # 使用另一种方式查找
    end_pattern = '</div>\n      </div>\n    </template>'
    end_idx = content.find(end_pattern, start_idx)
    if end_idx != -1:
        end_idx += len(end_pattern)
    else:
        end_idx = content.find('</div>\n    </template>', start_idx + 100)
        if end_idx != -1:
            end_idx += len('</div>\n    </template>')

print(f"Start: {start_idx}, End: {end_idx}")
print(f"Replacing {end_idx - start_idx} characters")

# 替换内容
new_file_content = content[:start_idx] + new_content + content[end_idx:]

# 写回文件
with open(r'D:\devops\docker-gpu-manage\web\src\view\modeltraining\training\detail.vue', 'w', encoding='utf-8') as f:
    f.write(new_file_content)

print("File modified successfully!")