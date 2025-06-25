<template>
  <div class="vhd-install">
    <el-row :gutter="20">
      <!-- 左侧配置面板 -->
      <el-col :span="14">
        <!-- 服务器选择 -->
        <el-card class="config-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Connection /></el-icon>
              <span>服务器选择</span>
              <el-button 
                type="primary" 
                size="small" 
                @click="refreshServers"
                :loading="downloadStore.isLoadingServers"
              >
                <el-icon><Refresh /></el-icon>
                刷新
              </el-button>
            </div>
          </template>
          
          <el-form label-position="top">
            <el-form-item label="服务器">
              <el-select 
                v-model="selectedServerId" 
                placeholder="请选择服务器"
                style="width: 100%"
                @change="onServerChange"
                :loading="downloadStore.isLoadingServers"
              >
                <el-option
                  v-for="server in downloadStore.servers"
                  :key="server.id"
                  :label="server.name"
                  :value="server.id"
                >
                  <span style="float: left">{{ server.name }}</span>
                  <span style="float: right; color: var(--el-text-color-secondary)">
                    {{ server.type }}
                  </span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- VHD选择和筛选 -->
        <el-card class="config-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Download /></el-icon>
              <span>VHD选择</span>
              <el-tag v-if="downloadStore.filteredVHDs.length" type="info" size="small">
                {{ downloadStore.filteredVHDs.length }} 个可用
              </el-tag>
            </div>
          </template>
          
          <el-form label-position="top">
            <!-- 筛选条件 -->
            <el-row :gutter="16">
              <el-col :span="8">
                <el-form-item label="系统类型">
                  <el-select 
                    v-model="downloadStore.filters.systemType" 
                    @change="downloadStore.applyFilters"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="type in systemTypeOptions"
                      :key="type"
                      :label="type"
                      :value="type"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="语言">
                  <el-select 
                    v-model="downloadStore.filters.language" 
                    @change="downloadStore.applyFilters"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="lang in languageOptions"
                      :key="lang"
                      :label="lang"
                      :value="lang"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="8">
                <el-form-item label="启动模式">
                  <el-select 
                    v-model="downloadStore.filters.bootMode" 
                    @change="downloadStore.applyFilters"
                    style="width: 100%"
                  >
                    <el-option
                      v-for="mode in bootModeOptions"
                      :key="mode"
                      :label="mode"
                      :value="mode"
                    />
                  </el-select>
                </el-form-item>
              </el-col>
            </el-row>

            <!-- VHD版本选择 -->
            <el-form-item label="VHD版本">
              <el-select 
                v-model="selectedVHDId" 
                placeholder="请先选择服务器"
                style="width: 100%"
                @change="onVHDChange"
                :loading="downloadStore.isLoadingVHDs"
                filterable
              >
                <el-option
                  v-for="vhd in downloadStore.filteredVHDs"
                  :key="vhd.filename"
                  :label="vhd.displayName"
                  :value="vhd.filename"
                >
                  <div class="vhd-option">
                    <span class="vhd-name">{{ vhd.displayName }}</span>
                    <div class="vhd-details">
                      <el-tag size="small" type="info">{{ vhd.system }}</el-tag>
                      <el-tag size="small" type="success">{{ vhd.language }}</el-tag>
                      <el-tag size="small" :type="vhd.bootMode === 'UEFI' ? 'warning' : 'primary'">
                        {{ vhd.bootMode }}
                      </el-tag>
                    </div>
                  </div>
                </el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 高级选项 -->
        <el-card class="config-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Setting /></el-icon>
              <span>高级选项</span>
            </div>
          </template>
          
          <el-form label-position="top">
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item label="自定义密码">
                  <el-input 
                    v-model="advancedOptions.password"
                    type="password"
                    placeholder="留空使用默认密码"
                    show-password
                  />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="SSH端口">
                  <el-input-number 
                    v-model="advancedOptions.sshPort"
                    :min="1"
                    :max="65535"
                    placeholder="22"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item>
              <el-checkbox v-model="advancedOptions.enableSSH">
                启用SSH服务
              </el-checkbox>
            </el-form-item>
            
            <el-form-item>
              <el-checkbox v-model="advancedOptions.autoReboot">
                安装完成后自动重启
              </el-checkbox>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>

      <!-- 右侧操作和信息面板 -->
      <el-col :span="10">
        <!-- 选中的VHD信息 -->
        <el-card v-if="downloadStore.selectedVHD" class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>VHD信息</span>
            </div>
          </template>
          
          <div class="vhd-info">
            <div class="info-item">
              <span class="label">名称：</span>
              <span class="value">{{ downloadStore.selectedVHD.displayName }}</span>
            </div>
            <div class="info-item">
              <span class="label">文件名：</span>
              <span class="value">{{ downloadStore.selectedVHD.filename }}</span>
            </div>
            <div class="info-item">
              <span class="label">系统：</span>
              <el-tag type="primary">{{ downloadStore.selectedVHD.system }}</el-tag>
            </div>
            <div class="info-item">
              <span class="label">版本：</span>
              <el-tag type="success">{{ downloadStore.selectedVHD.version }}</el-tag>
            </div>
            <div class="info-item">
              <span class="label">语言：</span>
              <el-tag type="info">{{ downloadStore.selectedVHD.language }}</el-tag>
            </div>
            <div class="info-item">
              <span class="label">启动模式：</span>
              <el-tag :type="downloadStore.selectedVHD.bootMode === 'UEFI' ? 'warning' : 'primary'">
                {{ downloadStore.selectedVHD.bootMode }}
              </el-tag>
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
            <el-button 
              type="primary"
              size="large"
              @click="startDownload"
              :disabled="!downloadStore.selectedVHD"
              :loading="isDownloading"
              style="width: 100%"
            >
              <el-icon><Download /></el-icon>
              {{ isDownloading ? '下载中...' : '开始下载' }}
            </el-button>
            
            <el-button 
              type="success"
              size="large"
              @click="startInstall"
              :disabled="!canInstall"
              style="width: 100%; margin-top: 10px"
            >
              <el-icon><Upload /></el-icon>
              开始安装
            </el-button>
          </div>
        </el-card>

        <!-- 下载进度 -->
        <el-card v-if="downloadStore.activeDownloads.length" class="progress-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Timer /></el-icon>
              <span>下载进度</span>
            </div>
          </template>
          
          <div v-for="download in downloadStore.activeDownloads" :key="download.id" class="download-item">
            <div class="download-header">
              <span class="download-name">{{ download.displayName }}</span>
              <el-button 
                type="danger" 
                size="small" 
                @click="cancelDownload(download.id)"
                circle
              >
                <el-icon><Close /></el-icon>
              </el-button>
            </div>
            <el-progress 
              :percentage="download.progress" 
              :status="download.status === 'failed' ? 'exception' : ''"
            />
            <div class="download-info">
              <span class="speed">{{ download.speed }}</span>
              <span class="time">{{ formatTime(download.startTime) }}</span>
            </div>
          </div>
        </el-card>

        <!-- 下载历史 -->
        <el-card v-if="downloadStore.downloadHistory.length" class="history-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Clock /></el-icon>
              <span>下载历史</span>
              <el-button 
                type="danger" 
                size="small" 
                @click="clearHistory"
              >
                清空
              </el-button>
            </div>
          </template>
          
          <div class="history-list">
            <div 
              v-for="item in downloadStore.downloadHistory.slice(0, 5)" 
              :key="item.id" 
              class="history-item"
            >
              <div class="history-header">
                <span class="history-name">{{ item.displayName }}</span>
                <el-tag 
                  :type="getStatusType(item.status)" 
                  size="small"
                >
                  {{ getStatusText(item.status) }}
                </el-tag>
              </div>
              <div class="history-time">
                {{ formatTime(item.startTime) }}
              </div>
              <div v-if="item.status === 'failed'" class="history-actions">
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="retryDownload(item)"
                >
                  重试
                </el-button>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useDownloadStore } from '../stores/download'
import { useAppStore } from '../stores/app'

const downloadStore = useDownloadStore()
const appStore = useAppStore()

// 响应式数据
const selectedServerId = ref('')
const selectedVHDId = ref('')
const isDownloading = ref(false)

// 高级选项
const advancedOptions = ref({
  password: '',
  sshPort: 22,
  enableSSH: false,
  autoReboot: false
})

// 计算属性
const systemTypeOptions = computed(() => downloadStore.getSystemTypeOptions())
const languageOptions = computed(() => downloadStore.getLanguageOptions())
const bootModeOptions = computed(() => downloadStore.getBootModeOptions())

const canInstall = computed(() => {
  // 检查是否有已下载的文件可以安装
  return downloadStore.downloadHistory.some(item => item.status === 'completed')
})

// 监听器
watch(() => downloadStore.selectedServer, (server) => {
  if (server) {
    selectedServerId.value = server.id
  }
})

watch(() => downloadStore.selectedVHD, (vhd) => {
  if (vhd) {
    selectedVHDId.value = vhd.filename
  }
})

// 方法
const refreshServers = async () => {
  try {
    await downloadStore.loadServers()
    ElMessage.success('服务器列表刷新成功')
  } catch (error) {
    ElMessage.error('刷新服务器列表失败: ' + error.message)
  }
}

const onServerChange = (serverId) => {
  const server = downloadStore.servers.find(s => s.id === serverId)
  if (server) {
    downloadStore.selectServer(server)
  }
}

const onVHDChange = (filename) => {
  const vhd = downloadStore.filteredVHDs.find(v => v.filename === filename)
  if (vhd) {
    downloadStore.selectVHD(vhd)
  }
}

const startDownload = async () => {
  if (!downloadStore.selectedVHD) {
    ElMessage.warning('请先选择要下载的VHD')
    return
  }

  try {
    isDownloading.value = true
    appStore.addLog('info', `开始下载: ${downloadStore.selectedVHD.displayName}`)
    
    await downloadStore.startDownload()
    
    ElMessage.success('下载完成')
    appStore.addLog('info', '下载完成')
  } catch (error) {
    ElMessage.error('下载失败: ' + error.message)
    appStore.addLog('error', `下载失败: ${error.message}`)
  } finally {
    isDownloading.value = false
  }
}

const startInstall = async () => {
  const result = await ElMessageBox.confirm(
    '确定要开始安装吗？这将重装您的系统。',
    '确认安装',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)

  if (!result) return

  try {
    appStore.showProgress('正在准备安装...', 0)
    
    // 这里调用安装逻辑
    // await installVHD(downloadStore.selectedVHD)
    
    ElMessage.success('安装已开始，请注意查看系统状态')
    appStore.addLog('info', '安装已开始')
  } catch (error) {
    ElMessage.error('安装失败: ' + error.message)
    appStore.addLog('error', `安装失败: ${error.message}`)
  } finally {
    appStore.hideProgress()
  }
}

const cancelDownload = (downloadId) => {
  downloadStore.cancelDownload(downloadId)
  ElMessage.info('下载已取消')
}

const clearHistory = () => {
  downloadStore.clearDownloadHistory()
  ElMessage.success('下载历史已清空')
}

const retryDownload = async (historyItem) => {
  try {
    await downloadStore.retryDownload(historyItem)
    ElMessage.success('重新下载已开始')
  } catch (error) {
    ElMessage.error('重新下载失败: ' + error.message)
  }
}

const getStatusType = (status) => {
  const types = {
    completed: 'success',
    failed: 'danger',
    downloading: 'primary',
    cancelled: 'info'
  }
  return types[status] || 'info'
}

const getStatusText = (status) => {
  const texts = {
    completed: '已完成',
    failed: '失败',
    downloading: '下载中',
    cancelled: '已取消'
  }
  return texts[status] || '未知'
}

const formatTime = (timestamp) => {
  return new Date(timestamp).toLocaleString()
}

// 生命周期
onMounted(async () => {
  appStore.addLog('info', 'VHD重装页面已加载')
  
  try {
    await downloadStore.init()
  } catch (error) {
    appStore.addLog('error', `初始化失败: ${error.message}`)
  }
})
</script>

<style scoped>
.vhd-install {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}

.config-card,
.info-card,
.action-card,
.progress-card,
.history-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
}

.card-header .el-button {
  margin-left: auto;
}

.vhd-option {
  width: 100%;
}

.vhd-name {
  font-weight: 500;
  display: block;
  margin-bottom: 4px;
}

.vhd-details {
  display: flex;
  gap: 4px;
}

.vhd-info .info-item {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
}

.vhd-info .label {
  min-width: 80px;
  font-weight: 500;
  color: var(--el-text-color-secondary);
}

.vhd-info .value {
  flex: 1;
}

.action-buttons {
  text-align: center;
}

.download-item {
  margin-bottom: 15px;
  padding: 10px;
  border: 1px solid var(--el-border-color);
  border-radius: 4px;
}

.download-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.download-name {
  font-weight: 500;
  font-size: 14px;
}

.download-info {
  display: flex;
  justify-content: space-between;
  margin-top: 8px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.history-list {
  max-height: 300px;
  overflow-y: auto;
}

.history-item {
  padding: 8px;
  border-bottom: 1px solid var(--el-border-color-lighter);
}

.history-item:last-child {
  border-bottom: none;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.history-name {
  font-weight: 500;
  font-size: 14px;
}

.history-time {
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.history-actions {
  margin-top: 8px;
  text-align: right;
}
</style>