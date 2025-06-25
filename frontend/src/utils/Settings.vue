<template>
  <div class="settings">
    <el-row :gutter="20">
      <el-col :span="16">
        <!-- 基本设置 -->
        <el-card class="setting-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Setting /></el-icon>
              <span>基本设置</span>
            </div>
          </template>
          
          <el-form label-position="left" label-width="120px">
            <el-form-item label="主题">
              <el-radio-group v-model="appStore.theme" @change="onThemeChange">
                <el-radio label="light">浅色主题</el-radio>
                <el-radio label="dark">深色主题</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="语言">
              <el-select v-model="language" style="width: 200px">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" disabled />
              </el-select>
            </el-form-item>
            
            <el-form-item label="自动刷新">
              <el-switch v-model="autoRefresh" />
              <span class="setting-desc">自动刷新服务器列表和备份历史</span>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 下载设置 -->
        <el-card class="setting-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Download /></el-icon>
              <span>下载设置</span>
            </div>
          </template>
          
          <el-form label-position="left" label-width="120px">
            <el-form-item label="并发下载数">
              <el-input-number 
                v-model="downloadConfig.maxConcurrent"
                :min="1"
                :max="5"
                style="width: 150px"
              />
              <span class="setting-desc">同时进行的最大下载数量</span>
            </el-form-item>
            
            <el-form-item label="重试次数">
              <el-input-number 
                v-model="downloadConfig.retryCount"
                :min="0"
                :max="10"
                style="width: 150px"
              />
              <span class="setting-desc">下载失败时的重试次数</span>
            </el-form-item>
            
            <el-form-item label="超时时间">
              <el-input-number 
                v-model="timeoutMinutes"
                :min="1"
                :max="60"
                style="width: 150px"
              />
              <span class="setting-desc">分钟，单个文件下载超时时间</span>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 备份设置 -->
        <el-card class="setting-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><FolderOpened /></el-icon>
              <span>备份设置</span>
            </div>
          </template>
          
          <el-form label-position="left" label-width="120px">
            <el-form-item label="默认启用压缩">
              <el-switch v-model="driverStore.backupConfig.enableCompression" />
              <span class="setting-desc">新建备份时默认启用7z压缩</span>
            </el-form-item>
            
            <el-form-item label="自动清理">
              <el-switch v-model="driverStore.backupConfig.autoCleanup" />
              <span class="setting-desc">自动删除超出数量限制的旧备份</span>
            </el-form-item>
            
            <el-form-item v-if="driverStore.backupConfig.autoCleanup" label="最大备份数">
              <el-input-number 
                v-model="driverStore.backupConfig.maxBackups"
                :min="1"
                :max="50"
                style="width: 150px"
              />
              <span class="setting-desc">保留的最大备份数量</span>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 日志设置 -->
        <el-card class="setting-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Document /></el-icon>
              <span>日志设置</span>
            </div>
          </template>
          
          <el-form label-position="left" label-width="120px">
            <el-form-item label="日志级别">
              <el-select v-model="logLevel" style="width: 200px">
                <el-option label="调试" value="debug" />
                <el-option label="信息" value="info" />
                <el-option label="警告" value="warning" />
                <el-option label="错误" value="error" />
              </el-select>
            </el-form-item>
            
            <el-form-item label="最大日志数">
              <el-input-number 
                v-model="maxLogs"
                :min="100"
                :max="10000"
                :step="100"
                style="width: 150px"
              />
              <span class="setting-desc">内存中保留的最大日志条数</span>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <el-col :span="8">
        <!-- 应用信息 -->
        <el-card class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>应用信息</span>
            </div>
          </template>
          
          <div class="app-info">
            <div class="info-item">
              <span class="label">应用名称:</span>
              <span class="value">SystemReinstaller</span>
            </div>
            <div class="info-item">
              <span class="label">版本:</span>
              <span class="value">v2.0.0</span>
            </div>
            <div class="info-item">
              <span class="label">构建时间:</span>
              <span class="value">{{ buildTime }}</span>
            </div>
            <div class="info-item">
              <span class="label">运行环境:</span>
              <span class="value">{{ runtime }}</span>
            </div>
          </div>
        </el-card>

        <!-- 系统信息 -->
        <el-card class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Monitor /></el-icon>
              <span>系统信息</span>
            </div>
          </template>
          
          <div class="system-info">
            <div class="info-item">
              <span class="label">计算机名:</span>
              <span class="value">{{ systemInfo.computerName || '获取中...' }}</span>
            </div>
            <div class="info-item">
              <span class="label">操作系统:</span>
              <span class="value">{{ systemInfo.osVersion || '获取中...' }}</span>
            </div>
            <div class="info-item">
              <span class="label">系统架构:</span>
              <span class="value">{{ systemInfo.arch || '获取中...' }}</span>
            </div>
            <div class="info-item">
              <span class="label">启动时间:</span>
              <span class="value">{{ systemInfo.timestamp || '获取中...' }}</span>
            </div>
          </div>
        </el-card>

        <!-- 操作按钮 -->
        <el-card class="action-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Operation /></el-icon>
              <span>操作</span>
            </div>
          </template>
          
          <div class="action-buttons">
            <el-button type="primary" @click="saveSettings" style="width: 100%">
              <el-icon><Check /></el-icon>
              保存设置
            </el-button>
            
            <el-button @click="resetSettings" style="width: 100%">
              <el-icon><RefreshLeft /></el-icon>
              重置设置
            </el-button>
            
            <el-button @click="exportSettings" style="width: 100%">
              <el-icon><Download /></el-icon>
              导出设置
            </el-button>
            
            <el-button @click="importSettings" style="width: 100%">
              <el-icon><Upload /></el-icon>
              导入设置
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAppStore } from '../stores/app'
import { useDriverStore } from '../stores/driver'
import { GetSystemInfo } from '../utils/wails'

const appStore = useAppStore()
const driverStore = useDriverStore()

// 响应式数据
const language = ref('zh-CN')
const autoRefresh = ref(true)
const logLevel = ref('info')
const maxLogs = ref(1000)

// 下载配置
const downloadConfig = ref({
  maxConcurrent: 3,
  retryCount: 3,
  timeout: 300000
})

// 计算属性
const timeoutMinutes = computed({
  get: () => downloadConfig.value.timeout / 60000,
  set: (value) => {
    downloadConfig.value.timeout = value * 60000
  }
})

const systemInfo = computed(() => appStore.systemInfo)

const buildTime = ref(new Date().toLocaleString())
const runtime = ref('Wails + Vue 3')

// 方法
const onThemeChange = (theme) => {
  appStore.setTheme(theme)
  document.documentElement.classList.toggle('dark', theme === 'dark')
}

const saveSettings = () => {
  try {
    // 保存所有设置到localStorage
    const settings = {
      language: language.value,
      autoRefresh: autoRefresh.value,
      logLevel: logLevel.value,
      maxLogs: maxLogs.value,
      downloadConfig: downloadConfig.value,
      backupConfig: driverStore.backupConfig
    }
    
    localStorage.setItem('app-settings', JSON.stringify(settings))
    
    ElMessage.success('设置已保存')
    appStore.addLog('info', '设置已保存')
  } catch (error) {
    ElMessage.error('保存设置失败: ' + error.message)
  }
}

const resetSettings = async () => {
  const result = await ElMessageBox.confirm(
    '确定要重置所有设置到默认值吗？',
    '确认重置',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)

  if (!result) return

  try {
    // 重置所有设置
    language.value = 'zh-CN'
    autoRefresh.value = true
    logLevel.value = 'info'
    maxLogs.value = 1000
    downloadConfig.value = {
      maxConcurrent: 3,
      retryCount: 3,
      timeout: 300000
    }
    
    // 重置主题
    appStore.setTheme('light')
    
    // 清除localStorage
    localStorage.removeItem('app-settings')
    
    ElMessage.success('设置已重置')
    appStore.addLog('info', '设置已重置为默认值')
  } catch (error) {
    ElMessage.error('重置设置失败: ' + error.message)
  }
}

const exportSettings = () => {
  try {
    const settings = {
      language: language.value,
      autoRefresh: autoRefresh.value,
      logLevel: logLevel.value,
      maxLogs: maxLogs.value,
      downloadConfig: downloadConfig.value,
      backupConfig: driverStore.backupConfig,
      theme: appStore.theme,
      exportTime: new Date().toISOString()
    }
    
    const blob = new Blob([JSON.stringify(settings, null, 2)], { 
      type: 'application/json;charset=utf-8' 
    })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `system-reinstaller-settings-${new Date().toISOString().slice(0, 10)}.json`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
    
    ElMessage.success('设置已导出')
    appStore.addLog('info', '设置已导出到文件')
  } catch (error) {
    ElMessage.error('导出设置失败: ' + error.message)
  }
}

const importSettings = () => {
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = '.json'
  input.onchange = (event) => {
    const file = event.target.files[0]
    if (!file) return
    
    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const settings = JSON.parse(e.target.result)
        
        // 验证设置格式
        if (typeof settings !== 'object') {
          throw new Error('无效的设置文件格式')
        }
        
        // 应用设置
        if (settings.language) language.value = settings.language
        if (typeof settings.autoRefresh === 'boolean') autoRefresh.value = settings.autoRefresh
        if (settings.logLevel) logLevel.value = settings.logLevel
        if (settings.maxLogs) maxLogs.value = settings.maxLogs
        if (settings.downloadConfig) Object.assign(downloadConfig.value, settings.downloadConfig)
        if (settings.backupConfig) Object.assign(driverStore.backupConfig, settings.backupConfig)
        if (settings.theme) appStore.setTheme(settings.theme)
        
        ElMessage.success('设置已导入')
        appStore.addLog('info', '设置已从文件导入')
      } catch (error) {
        ElMessage.error('导入设置失败: ' + error.message)
      }
    }
    reader.readAsText(file)
  }
  input.click()
}

// 生命周期
onMounted(async () => {
  appStore.addLog('info', '设置页面已加载')
  
  try {
    // 获取系统信息
    const info = await GetSystemInfo()
    appStore.updateSystemInfo(info)
  } catch (error) {
    appStore.addLog('error', `获取系统信息失败: ${error.message}`)
  }
  
  // 加载保存的设置
  const savedSettings = localStorage.getItem('app-settings')
  if (savedSettings) {
    try {
      const settings = JSON.parse(savedSettings)
      if (settings.language) language.value = settings.language
      if (typeof settings.autoRefresh === 'boolean') autoRefresh.value = settings.autoRefresh
      if (settings.logLevel) logLevel.value = settings.logLevel
      if (settings.maxLogs) maxLogs.value = settings.maxLogs
      if (settings.downloadConfig) Object.assign(downloadConfig.value, settings.downloadConfig)
    } catch (error) {
      appStore.addLog('error', `加载设置失败: ${error.message}`)
    }
  }
})
</script>

<style scoped>
.settings {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

.setting-card,
.info-card,
.action-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.setting-desc {
  margin-left: 12px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.app-info,
.system-info {
  font-size: 14px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.info-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
}

.label {
  font-weight: 500;
  color: var(--el-text-color-secondary);
}

.value {
  color: var(--el-text-color-primary);
}

.action-buttons {
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>