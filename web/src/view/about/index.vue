<template>
  <div class="about-container">
    <!-- 头部横幅 -->
    <div class="about-banner">
      <div class="banner-content">
        <div class="logo-wrapper">
          <svg viewBox="0 0 640 512" class="logo-svg">
            <path d="M349.9 236.3h-66.1v-59.4h66.1v59.4zm0-204.3h-66.1v60.7h66.1V32zm78.2 144.8H362v59.4h66.1v-59.4zm-156.3-72.1h-66.1v60.1h66.1v-60.1zm78.1 0h-66.1v60.1h66.1v-60.1zm276.8 100c-14.4-9.7-47.6-13.2-73.1-8.4-3.3-24-16.7-44.9-41.1-63.7l-14-9.3-9.3 14c-18.4 27.8-23.4 73.6-3.7 103.8-8.7 4.7-25.8 11.1-48.4 10.7H2.4c-8.7 50.8 5.8 116.8 44 162.1 37.1 43.9 92.7 66.2 165.4 66.2 157.4 0 273.9-72.5 328.4-204.2 21.4.4 67.6.1 91.3-45.2 1.5-2.5 6.6-13.2 8.5-17.1l-13.3-8.9zm-511.1-27.9h-66v59.4h66.1v-59.4zm78.1 0h-66.1v59.4h66.1v-59.4zm78.1 0h-66.1v59.4h66.1v-59.4zm-78.1-72.1h-66.1v60.1h66.1v-60.1z"/>
          </svg>
        </div>
        <h1 class="banner-title">Docker GPU Manager</h1>
        <p class="banner-subtitle">容器化 GPU 算力资源管理平台</p>
        <p class="banner-desc">高效、灵活、智能的 GPU 资源调度解决方案</p>
      </div>
      <div class="banner-decoration"></div>
    </div>

    <!-- 主要内容 -->
    <div class="about-content">
      <!-- 平台特性 -->
      <div class="section">
        <h2 class="section-title">平台特性</h2>
        <div class="features-grid">
          <div class="feature-card">
            <div class="feature-icon feature-icon-1">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-7 14l-5-5 1.41-1.41L12 14.17l4.59-4.59L18 11l-6 6z"/>
              </svg>
            </div>
            <h3>实时监控</h3>
            <p>GPU 使用率、显存占用、温度功耗等关键指标实时监控</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon feature-icon-2">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M21 16V8a2 2 0 0 0-1-1.73l-7-4a2 2 0 0 0-2 0l-7 4A2 2 0 0 0 3 8v8a2 2 0 0 0 1 1.73l7 4a2 2 0 0 0 2 0l7-4A2 2 0 0 0 21 16z"/>
              </svg>
            </div>
            <h3>容器编排</h3>
            <p>基于 Docker 的容器管理，一键部署 AI 训练环境</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon feature-icon-3">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
            <h3>智能调度</h3>
            <p>自动分配空闲 GPU 资源，最大化算力利用效率</p>
          </div>

          <div class="feature-card">
            <div class="feature-icon feature-icon-4">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M12 3L1 9l11 6 9-4.91V17h2V9M5 13.18v4L12 21l7-3.82v-4L12 17l-7-3.82z"/>
              </svg>
            </div>
            <h3>模型训练</h3>
            <p>支持 SFT、DPO、CPT 等多种训练方式的可视化管理</p>
          </div>
        </div>
      </div>

      <!-- 系统信息 -->
      <div class="section">
        <h2 class="section-title">系统信息</h2>
        <div class="info-cards">
          <div class="info-card">
            <div class="info-header">
              <span class="info-label">技术栈</span>
            </div>
            <div class="info-list">
              <div class="info-item">
                <span class="info-key">前端框架</span>
                <span class="info-value">Vue 3 + Element Plus</span>
              </div>
              <div class="info-item">
                <span class="info-key">后端框架</span>
                <span class="info-value">Gin (Golang)</span>
              </div>
              <div class="info-item">
                <span class="info-key">容器技术</span>
                <span class="info-value">Docker Engine</span>
              </div>
              <div class="info-item">
                <span class="info-key">GPU 支持</span>
                <span class="info-value">NVIDIA CUDA / ROCm</span>
              </div>
            </div>
          </div>

          <div class="info-card">
            <div class="info-header">
              <span class="info-label">核心功能</span>
            </div>
            <div class="info-list">
              <div class="info-item">
                <span class="info-key">计算节点</span>
                <span class="info-value">多节点集群管理</span>
              </div>
              <div class="info-item">
                <span class="info-key">实例管理</span>
                <span class="info-value">GPU 实例生命周期</span>
              </div>
              <div class="info-item">
                <span class="info-key">镜像仓库</span>
                <span class="info-value">自定义镜像管理</span>
              </div>
              <div class="info-item">
                <span class="info-key">数据集</span>
                <span class="info-value">训练数据版本控制</span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 快速链接 -->
      <div class="section">
        <h2 class="section-title">相关链接</h2>
        <div class="links-grid">
          <a href="https://www.docker.com/" target="_blank" class="link-card">
            <svg viewBox="0 0 640 512" class="link-icon">
              <path d="M349.9 236.3h-66.1v-59.4h66.1v59.4zm0-204.3h-66.1v60.7h66.1V32zm78.2 144.8H362v59.4h66.1v-59.4zm-156.3-72.1h-66.1v60.1h66.1v-60.1zm78.1 0h-66.1v60.1h66.1v-60.1zm276.8 100c-14.4-9.7-47.6-13.2-73.1-8.4-3.3-24-16.7-44.9-41.1-63.7l-14-9.3-9.3 14c-18.4 27.8-23.4 73.6-3.7 103.8-8.7 4.7-25.8 11.1-48.4 10.7H2.4c-8.7 50.8 5.8 116.8 44 162.1 37.1 43.9 92.7 66.2 165.4 66.2 157.4 0 273.9-72.5 328.4-204.2 21.4.4 67.6.1 91.3-45.2 1.5-2.5 6.6-13.2 8.5-17.1l-13.3-8.9zm-511.1-27.9h-66v59.4h66.1v-59.4zm78.1 0h-66.1v59.4h66.1v-59.4zm78.1 0h-66.1v59.4h66.1v-59.4zm-78.1-72.1h-66.1v60.1h66.1v-60.1z" fill="currentColor"/>
            </svg>
            <span>Docker 官网</span>
          </a>
          <a href="https://www.nvidia.com/en-us/data-center/gpu-accelerated-applications/" target="_blank" class="link-card">
            <svg viewBox="0 0 24 24" class="link-icon">
              <path d="M4 4h16v16H4V4m2 2v12h12V6H6z" fill="currentColor"/>
            </svg>
            <span>NVIDIA GPU 应用</span>
          </a>
          <a href="https://vuejs.org/" target="_blank" class="link-card">
            <svg viewBox="0 0 128 128" class="link-icon">
              <path d="M78.7 64.8L60.6 32.4 42.5 64.8H78.7zM80.4 64.8H115.2L78.4 0 60.4 32.4 80.4 64.8zM38.4 64.8L0 128H40.8L60.8 96 38.4 64.8zM89.2 128H128L78.4 64.8 89.2 128z" fill="currentColor"/>
            </svg>
            <span>Vue.js</span>
          </a>
          <a href="https://gin-gonic.com/" target="_blank" class="link-card">
            <svg viewBox="0 0 24 24" class="link-icon">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.95-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" fill="currentColor"/>
            </svg>
            <span>Gin Framework</span>
          </a>
        </div>
      </div>

      <!-- 版权信息 -->
      <div class="copyright">
        <p>© 2024 Docker GPU Manager. All rights reserved.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
defineOptions({
  name: 'About'
})
</script>

<style scoped>
.about-container {
  min-height: calc(100vh - 120px);
  background: #f5f7fa;
}

/* 头部横幅 */
.about-banner {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
  padding: 60px 40px;
  position: relative;
  overflow: hidden;
}

.banner-content {
  max-width: 1100px;
  margin: 0 auto;
  text-align: center;
  position: relative;
  z-index: 1;
}

.logo-wrapper {
  width: 100px;
  height: 100px;
  margin: 0 auto 24px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.4);
}

.logo-svg {
  width: 60px;
  height: 60px;
  color: #fff;
}

.banner-title {
  font-size: 36px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 12px;
  letter-spacing: 1px;
}

.banner-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.9);
  margin: 0 0 8px;
}

.banner-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.banner-decoration {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.03'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  pointer-events: none;
}

/* 主要内容 */
.about-content {
  max-width: 1100px;
  margin: 0 auto;
  padding: 40px 20px;
}

.section {
  margin-bottom: 40px;
}

.section-title {
  font-size: 20px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 20px;
  padding-left: 12px;
  border-left: 4px solid #667eea;
}

/* 特性卡片 */
.features-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
}

.feature-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.feature-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.feature-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.feature-icon svg {
  width: 28px;
  height: 28px;
  color: #fff;
}

.feature-icon-1 { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.feature-icon-2 { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.feature-icon-3 { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.feature-icon-4 { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }

.feature-card h3 {
  font-size: 16px;
  font-weight: 600;
  color: #1a1a2e;
  margin: 0 0 8px;
}

.feature-card p {
  font-size: 13px;
  color: #86909c;
  margin: 0;
  line-height: 1.5;
}

/* 信息卡片 */
.info-cards {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20px;
}

.info-card {
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.info-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 14px 20px;
}

.info-label {
  font-size: 15px;
  font-weight: 600;
  color: #fff;
}

.info-list {
  padding: 8px 0;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  border-bottom: 1px solid #f0f2f5;
}

.info-item:last-child {
  border-bottom: none;
}

.info-key {
  font-size: 14px;
  color: #86909c;
}

.info-value {
  font-size: 14px;
  font-weight: 500;
  color: #1a1a2e;
}

/* 链接卡片 */
.links-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.link-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 24px 16px;
  background: #fff;
  border-radius: 12px;
  text-decoration: none;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

.link-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.link-icon {
  width: 40px;
  height: 40px;
  color: #667eea;
}

.link-card span {
  font-size: 13px;
  color: #1a1a2e;
  font-weight: 500;
}

/* 版权信息 */
.copyright {
  text-align: center;
  padding: 40px 0 20px;
  border-top: 1px solid #e4e7ed;
  margin-top: 20px;
}

.copyright p {
  font-size: 13px;
  color: #86909c;
  margin: 0;
}

/* 响应式 */
@media (max-width: 1000px) {
  .features-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .links-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .about-banner {
    padding: 40px 20px;
  }

  .banner-title {
    font-size: 28px;
  }

  .info-cards {
    grid-template-columns: 1fr;
  }

  .features-grid,
  .links-grid {
    grid-template-columns: 1fr;
  }
}
</style>