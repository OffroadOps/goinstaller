<template>
  <div id="app" class="app-container">
    <!-- 标题栏 -->
    <div class="title-bar">
      <div class="title-content">
        <el-icon class="title-icon"><Monitor /></el-icon>
        <span class="title-text">SystemReinstaller v2.0 - 系统重装助手</span>
      </div>
      <div class="title-actions">
        <el-button type="info" size="small" @click="toggleTheme">
          <el-icon><Moon v-if="isDark" /><Sunny v-else /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 主内容区域 -->
    <div class="main-content">
      <!-- 选项卡导航 -->
      <div class="tab-nav">
        <el-tabs v-model="activeTab" type="card" class="main-tabs" @tab-change="onTabChange">
          <el-tab-pane label="VHD重装" name="vhd">
            <template #label>
              <span class="tab-label">
                <el-icon><Download /></el-icon>
                VHD重装
              </span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="ISO安装" name="iso">
            <template #label>
              <span class="tab-label">
                <el-icon><CdRom /></el-icon>
                ISO安装
              </span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="Linux安装" name="linux">
            <template #label>
              <span class="tab-label">
                <el-icon><Monitor /></el-icon>
                Linux安装
              </span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="驱动管理" name="driver">
            <template #label>
              <span class="tab-label">
                <el-icon><Setting /></el-icon>
                驱动管理
              </span>
            </template>
          </el-tab-pane>
          <el-tab-pane label="系统设置" name="settings">
            <template #label>
              <span class="tab-label">
                <el-icon><Tools /></el-icon>
                系统设置
              </span>
            </template>
          </el-tab-pane>
        </el-tabs>
      </div>

      <!-- 内容区域 -->
      <div class="content-area">
        <div class="content-wrapper">
          <!-- 主要内容 -->
          <div class="main-panel">
            <router-view v-slot="{ Component }">
              <transition name="fade" mode="out-in">
                <component :is="Component" v-if="Component" />
                <div v-else class="loading-placeholder">
                  <el-icon class="is-loading"><Loading /></el-icon>
                  <span>加载中...</span>
                </div>
              </transition>
            </router-view>
          </div>

          <!-- 日志面板 -->
          <div class="log-panel">
            <div class="log-header">
              <span class="log-title">
                <el-icon><Document /></el-icon>
                系统日志
              </span>
              <div class="log-actions">
                <el-button size="small" @click="clearLog">
                  <el-icon><Delete /></el-icon>
                  清空
                </el-button>
                <el-button size="small" @click="exportLog">
                  <el-icon><Download /></el-icon>
                  导出
                </el-button>
              </div>
            </div>
            <div class="log-content">
              <el-scrollbar height="100%">
                <div class="log-text" ref="logContainer">
                  <div v-for="(log, index) in logs" :key="index" class="log-line" :class="log.type">
                    <span class="log-time">{{ log.time }}</span>
                    <span class="log-level">{{ log.level }}</span>
                    <span class="log-message">{{ log.message }}</span>
                  </div>
                </div>
              </el-scrollbar>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 进度条 -->
    <div v-if="showProgress" class="progress-overlay">
      <div class="progress-content">
        <el-progress 
          :percentage="progressValue" 
          :status="progressStatus"
          :stroke-width="6"
        />
        <div class="progress-text">{{ progressText }}</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAppStore } from './stores/app'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()

// 响应式数据
const activeTab = ref('vhd')
const isDark = ref(false)
const logs = ref([])
const showProgress = ref(false)
const progressValue = ref(0)
const progressStatus = ref('')
const progressText = ref('')

// 监听路由变化
watch(() => route.name, (newRoute) => {
  if (newRoute) {
    activeTab.value = newRoute
  }
}, { immediate: true })

// 选项卡变化处理
const onTabChange = (tabName) => {
  router.push({ name: tabName }).catch(() => {
    // 忽略导航重复错误
  })
}

// 主题切换
const toggleTheme = () => {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// 日志操作
const clearLog = () => {
  logs.value = []
  appStore.clearLogs()
}

const exportLog = () => {
  appStore.exportLogs()
}

// 添加日志
const addLog = (level, message) => {
  const log = {
    time: new Date().toLocaleTimeString(),
    level: level.toUpperCase(),
    message,
    type: level.toLowerCase()
  }
  logs.value.push(log)
  
  // 自动滚动到底部
  nextTick(() => {
    const container = document.querySelector('.log-text')
    if (container) {
      container.scrollTop = container.scrollHeight
    }
  })
}

// 组件挂载
onMounted(() => {
  console.log('App mounted')
  
  // 初始化主题
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
  
  // 初始化路由
  if (!route.name) {
    router.push({ name: 'vhd' }).catch(() => {})
  } else {
    activeTab.value = route.name
  }
  
  // 初始化store
  appStore.init()
  
  // 添加欢迎日志
  addLog('info', '系统重装助手已启动')
  addLog('info', '欢迎使用 SystemReinstaller v2.0')
})

// 暴露方法给全局使用
window.addLog = addLog
window.showProgress = (text, value, status) => {
  showProgress.value = true
  progressText.value = text
  progressValue.value = value
  progressStatus.value = status
}
window.hideProgress = () => {
  showProgress.value = false
}
</script>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color-page);
}

.title-bar {
  height: 50px;
  background: var(--el-bg-color);
  border-bottom: 1px solid var(--el-border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  -webkit-app-region: drag;
}

.title-content {
  display: flex;
  align-items: center;
  gap: 10px;
}

.title-icon {
  font-size: 20px;
  color: var(--el-color-primary);
}

.title-text {
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.title-actions {
  -webkit-app-region: no-drag;
}

.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.tab-nav {
  background: var(--el-bg-color);
  padding: 0 20px;
  border-bottom: 1px solid var(--el-border-color);
}

.main-tabs {
  --el-tabs-header-height: 50px;
}

.tab-label {
  display: flex;
  align-items: center;
  gap: 6px;
}

.content-area {
  flex: 1;
  overflow: hidden;
}

.content-wrapper {
  height: 100%;
  display: flex;
}

.main-panel {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

.loading-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 200px;
  color: var(--el-text-color-secondary);
}

.log-panel {
  width: 400px;
  border-left: 1px solid var(--el-border-color);
  background: var(--el-bg-color);
  display: flex;
  flex-direction: column;
}

.log-header {
  height: 50px;
  padding: 0 15px;
  border-bottom: 1px solid var(--el-border-color);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.log-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.log-actions {
  display: flex;
  gap: 8px;
}

.log-content {
  flex: 1;
  overflow: hidden;
}

.log-text {
  padding: 10px;
  font-family: 'Consolas', 'Monaco', monospace;
  font-size: 12px;
  line-height: 1.4;
}

.log-line {
  margin-bottom: 2px;
  display: flex;
  gap: 8px;
}

.log-time {
  color: var(--el-text-color-secondary);
  min-width: 80px;
}

.log-level {
  min-width: 50px;
  font-weight: bold;
}

.log-level.info {
  color: var(--el-color-primary);
}

.log-level.warning {
  color: var(--el-color-warning);
}

.log-level.error {
  color: var(--el-color-danger);
}

.log-level.debug {
  color: var(--el-color-info);
}

.log-message {
  flex: 1;
  color: var(--el-text-color-primary);
}

.progress-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.progress-content {
  background: var(--el-bg-color);
  padding: 30px;
  border-radius: 8px;
  min-width: 300px;
  text-align: center;
}

.progress-text {
  margin-top: 15px;
  color: var(--el-text-color-primary);
  font-size: 14px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>