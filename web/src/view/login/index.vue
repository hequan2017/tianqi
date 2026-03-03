<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>

    <div class="login-wrapper">
      <!-- 左侧品牌区域 -->
      <div class="brand-section">
        <div class="brand-content">
          <!-- Logo -->
          <div class="logo-container">
            <svg viewBox="0 0 640 512" class="logo-icon">
              <path d="M349.9 236.3h-66.1v-59.4h66.1v59.4zm0-204.3h-66.1v60.7h66.1V32zm78.2 144.8H362v59.4h66.1v-59.4zm-156.3-72.1h-66.1v60.1h66.1v-60.1zm78.1 0h-66.1v60.1h66.1v-60.1zm276.8 100c-14.4-9.7-47.6-13.2-73.1-8.4-3.3-24-16.7-44.9-41.1-63.7l-14-9.3-9.3 14c-18.4 27.8-23.4 73.6-3.7 103.8-8.7 4.7-25.8 11.1-48.4 10.7H2.4c-8.7 50.8 5.8 116.8 44 162.1 37.1 43.9 92.7 66.2 165.4 66.2 157.4 0 273.9-72.5 328.4-204.2 21.4.4 67.6.1 91.3-45.2 1.5-2.5 6.6-13.2 8.5-17.1l-13.3-8.9zm-511.1-27.9h-66v59.4h66.1v-59.4zm78.1 0h-66.1v59.4h66.1v-59.4zm78.1 0h-66.1v59.4h66.1v-59.4zm-78.1-72.1h-66.1v60.1h66.1v-60.1z"/>
            </svg>
          </div>

          <h1 class="brand-title">Docker GPU Manager</h1>
          <p class="brand-subtitle">容器化 GPU 算力资源管理平台</p>

          <!-- 特性列表 -->
          <div class="features-list">
            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14l-5-5 1.41-1.41L12 14.17l4.59-4.59L18 11l-6 6z"/>
                </svg>
              </div>
              <div class="feature-content">
                <h3>实时监控</h3>
                <p>GPU 使用率、显存、温度等关键指标实时展示</p>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
                </svg>
              </div>
              <div class="feature-content">
                <h3>容器编排</h3>
                <p>基于 Docker 的一键部署，快速构建 AI 训练环境</p>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
                </svg>
              </div>
              <div class="feature-content">
                <h3>智能调度</h3>
                <p>自动分配空闲 GPU 资源，最大化算力利用率</p>
              </div>
            </div>

            <div class="feature-item">
              <div class="feature-icon">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 3L1 9l11 6 9-4.91V17h2V9M5 13.18v4L12 21l7-3.82v-4L12 17l-7-3.82z"/>
                </svg>
              </div>
              <div class="feature-content">
                <h3>模型训练</h3>
                <p>支持 SFT、DPO、CPT 等多种训练方式的可视化管理</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录区域 -->
      <div class="login-section">
        <div class="login-box">
          <div class="login-header">
            <h2>欢迎登录</h2>
            <p>请输入您的账号信息</p>
          </div>

          <el-form
            ref="loginForm"
            :model="loginFormData"
            :rules="rules"
            :validate-on-rule-change="false"
            @keyup.enter="submitForm"
            class="login-form"
          >
            <el-form-item prop="username">
              <el-input
                v-model="loginFormData.username"
                size="large"
                placeholder="请输入用户名"
                :prefix-icon="User"
              />
            </el-form-item>

            <el-form-item prop="password">
              <el-input
                v-model="loginFormData.password"
                show-password
                size="large"
                type="password"
                placeholder="请输入密码"
                :prefix-icon="Lock"
              />
            </el-form-item>

            <el-form-item v-if="loginFormData.openCaptcha" prop="captcha">
              <div class="captcha-row">
                <el-input
                  v-model="loginFormData.captcha"
                  placeholder="请输入验证码"
                  size="large"
                  class="captcha-input"
                />
                <div class="captcha-img" @click="loginVerify()">
                  <img v-if="picPath" :src="picPath" alt="验证码" />
                </div>
              </div>
            </el-form-item>

            <el-form-item>
              <el-button
                type="primary"
                size="large"
                class="login-btn"
                :loading="loading"
                @click="submitForm"
              >
                登 录
              </el-button>
            </el-form-item>

            <div class="login-footer">
              <el-button type="primary" link @click="checkInit">
                首次使用？前往初始化
              </el-button>
            </div>
          </el-form>
        </div>

        <!-- 底部信息 -->
        <div class="bottom-info">
          <span>© 2024 Docker GPU Manager. All rights reserved.</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { captcha } from '@/api/user'
import { checkDB } from '@/api/initdb'
import { reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/pinia/modules/user'
import { User, Lock } from '@element-plus/icons-vue'

defineOptions({
  name: 'Login'
})

const router = useRouter()
const loading = ref(false)
const captchaRequiredLength = ref(6)

// 验证函数
const checkUsername = (rule, value, callback) => {
  if (value.length < 5) {
    return callback(new Error('请输入正确的用户名'))
  } else {
    callback()
  }
}
const checkPassword = (rule, value, callback) => {
  if (value.length < 6) {
    return callback(new Error('请输入正确的密码'))
  } else {
    callback()
  }
}
const checkCaptcha = (rule, value, callback) => {
  if (!loginFormData.openCaptcha) {
    return callback()
  }
  const sanitizedValue = (value || '').replace(/\s+/g, '')
  if (!sanitizedValue) {
    return callback(new Error('请输入验证码'))
  }
  if (!/^\d+$/.test(sanitizedValue)) {
    return callback(new Error('验证码须为数字'))
  }
  if (sanitizedValue.length < captchaRequiredLength.value) {
    return callback(
      new Error(`请输入至少${captchaRequiredLength.value}位数字验证码`)
    )
  }
  if (sanitizedValue !== value) {
    loginFormData.captcha = sanitizedValue
  }
  callback()
}

// 获取验证码
const loginVerify = async () => {
  const ele = await captcha()
  captchaRequiredLength.value = Number(ele.data?.captchaLength) || 0
  picPath.value = ele.data?.picPath
  loginFormData.captchaId = ele.data?.captchaId
  loginFormData.openCaptcha = ele.data?.openCaptcha
}
loginVerify()

// 登录相关操作
const loginForm = ref(null)
const picPath = ref('')
const loginFormData = reactive({
  username: 'admin',
  password: '',
  captcha: '',
  captchaId: '',
  openCaptcha: false
})
const rules = reactive({
  username: [{ validator: checkUsername, trigger: 'blur' }],
  password: [{ validator: checkPassword, trigger: 'blur' }],
  captcha: [{ validator: checkCaptcha, trigger: 'blur' }]
})

const userStore = useUserStore()
const login = async () => {
  return await userStore.LoginIn(loginFormData)
}
const submitForm = () => {
  loginForm.value.validate(async (v) => {
    if (!v) {
      ElMessage({
        type: 'error',
        message: '请正确填写登录信息',
        showClose: true
      })
      return false
    }

    loading.value = true
    try {
      const flag = await login()
      if (!flag) {
        await loginVerify()
        return false
      }
      return true
    } finally {
      loading.value = false
    }
  })
}

// 跳转初始化
const checkInit = async () => {
  const res = await checkDB()
  if (res.code === 0) {
    if (res.data?.needInit) {
      userStore.NeedInit()
      await router.push({ name: 'Init' })
    } else {
      ElMessage({
        type: 'info',
        message: '初始化已完成，请直接登录'
      })
    }
  }
}
</script>

<style scoped>
.login-container {
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰圆 */
.bg-decoration {
  position: absolute;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.circle-1 {
  width: 600px;
  height: 600px;
  top: -200px;
  left: -200px;
}

.circle-2 {
  width: 400px;
  height: 400px;
  bottom: -100px;
  right: -100px;
}

.circle-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  right: 10%;
  background: rgba(255, 255, 255, 0.05);
}

.login-wrapper {
  display: flex;
  width: 90%;
  max-width: 1100px;
  height: 600px;
  background: #fff;
  border-radius: 20px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
  overflow: hidden;
  position: relative;
  z-index: 1;
}

/* 左侧品牌区域 */
.brand-section {
  flex: 1;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  padding: 50px 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
}

.brand-section::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.03'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  pointer-events: none;
}

.brand-content {
  text-align: center;
  position: relative;
  z-index: 1;
}

.logo-container {
  width: 100px;
  height: 100px;
  margin: 0 auto 30px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.logo-icon {
  width: 60px;
  height: 60px;
  color: #fff;
}

.brand-title {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 10px;
  letter-spacing: 1px;
}

.brand-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  margin: 0 0 40px;
}

.features-list {
  text-align: left;
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 16px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateX(5px);
}

.feature-icon {
  width: 44px;
  height: 44px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.feature-icon svg {
  width: 24px;
  height: 24px;
  color: #fff;
}

.feature-content h3 {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 4px;
}

.feature-content p {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
  line-height: 1.4;
}

/* 右侧登录区域 */
.login-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 50px 60px;
  background: #fff;
}

.login-box {
  width: 100%;
  max-width: 360px;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.login-header h2 {
  font-size: 28px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 10px;
}

.login-header p {
  font-size: 14px;
  color: #86909c;
  margin: 0;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 10px;
  box-shadow: 0 0 0 1px #e4e7ed inset;
  padding: 4px 15px;
  height: 48px;
}

.login-form :deep(.el-input__wrapper:hover) {
  box-shadow: 0 0 0 1px #c0c4cc inset;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  box-shadow: 0 0 0 2px #667eea inset;
}

.login-form :deep(.el-input__inner) {
  height: 40px;
  font-size: 15px;
}

.login-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.captcha-row {
  display: flex;
  gap: 12px;
  width: 100%;
}

.captcha-input {
  flex: 1;
}

.captcha-img {
  width: 120px;
  height: 48px;
  border-radius: 10px;
  overflow: hidden;
  cursor: pointer;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.captcha-img img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.login-btn {
  width: 100%;
  height: 48px;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 500;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  transition: all 0.3s ease;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.login-footer {
  text-align: center;
  margin-top: 20px;
}

.bottom-info {
  position: absolute;
  bottom: 20px;
  text-align: center;
  width: 100%;
}

.bottom-info span {
  font-size: 12px;
  color: #c9cdd4;
}

/* 响应式 */
@media (max-width: 900px) {
  .login-wrapper {
    flex-direction: column;
    height: auto;
    max-height: 90vh;
  }

  .brand-section {
    padding: 40px 30px;
  }

  .brand-title {
    font-size: 24px;
  }

  .features-list {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .feature-item {
    padding: 12px;
  }

  .feature-content h3 {
    font-size: 13px;
  }

  .feature-content p {
    font-size: 11px;
  }

  .login-section {
    padding: 40px 30px;
  }
}

@media (max-width: 600px) {
  .login-container {
    padding: 20px;
  }

  .login-wrapper {
    border-radius: 16px;
  }

  .features-list {
    grid-template-columns: 1fr;
  }

  .feature-item:hover {
    transform: none;
  }
}
</style>