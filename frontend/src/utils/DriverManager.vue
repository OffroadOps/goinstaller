<template>
  <div class="driver-manager">
    <el-row :gutter="20">
      <!-- 左侧备份设置 -->
      <el-col :span="16">
        <!-- 备份设置 -->
        <el-card class="config-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><FolderOpened /></el-icon>
              <span>备份设置</span>
            </div>
          </template>
          
          <el-form label-position="top">
            <el-form-item label="备份保存位置">
              <el-input 
                v-model="driverStore.backupPath" 
                placeholder="选择备份保存位置"
                readonly
              >
                <template #append>
                  <el-button @click="selectBackupPath">
                    <el-icon><FolderOpened /></el-icon>
                    浏览
                  </el-button>
                </template>
              </el-input>
            </el-form-item>
            
            <el-row :gutter="16">
              <el-col :span="12">
                <el-form-item>
                  <el-checkbox v-model="driverStore.backupConfig.enableCompression">
                    启用7z压缩
                  </el-checkbox>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item>
                  <el-checkbox v-model="driverStore.backupConfig.autoCleanup">
                    自动清理旧备份
                  </el-checkbox>
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item v-if="driverStore.backupConfig.autoCleanup" label="最大备份数量">
              <el-input-number 
                v-model="driverStore.backupConfig.maxBackups"
                :min="1"
                :max="50"
                style="width: 150px"
              />
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 操作按钮 -->
        <el-card class="action-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Operation /></el-icon>
              <span>驱动操作</span>
            </div>
          </template>
          
          <div class="action-buttons">
            <el-button 
              type="primary"
              size="large"
              @click="startBackup"
              :loading="driverStore.isBackingUp"
              icon="Download"
            >
              {{ driverStore.isBackingUp ? '备份中...' : '开始备份' }}
            </el-button>
            
            <el-button 
              type="success"
              size="large"
              @click="startRestore"
              :loading="driverStore.isRestoring"
              :disabled="!driverStore.selectedBackup"
              icon="Upload"
            >
              {{ driverStore.isRestoring ? '恢复中...' : '恢复驱动' }}
            </el-button>
            
            <el-button 
              type="danger"
              size="large"
              @click="deleteSelectedBackup"
              :disabled="!driverStore.selectedBackup"
              icon="Delete"
            >
              删除备份
            </el-button>
          </div>
        </el-card>

        <!-- 备份历史 -->
        <el-card class="history-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Clock /></el-icon>
              <span>备份历史</span>
              <div class="header-actions">
                <el-tag v-if="backupStats.total" type="info" size="small">
                  共 {{ backupStats.total }} 个备份
                </el-tag>
                <el-button 
                  type="primary" 
                  size="small" 
                  @click="refreshBackupHistory"
                  :loading="driverStore.isLoading"
                >
                  <el-icon><Refresh /></el-icon>
                  刷新
                </el-button>
              </div>
            </div>
          </template>
          
          <div class="backup-list">
            <div 
              v-if="!driverStore.backupRecords.length && !driverStore.isLoading" 
              class="empty-state"
            >
              <el-empty description="暂无备份记录" />
            </div>
            
            <div v-else>
              <div 
                v-for="record in driverStore.backupRecords" 
                :key="record.path"
                class="backup-item"
                :class="{ active: driverStore.selectedBackup?.path === record.path }"
                @click="selectBackup(record)"
              >
                <div class="backup-header">
                  <div class="backup-name">
                    <el-icon>
                      <el-icon v-if="record.isCompressed"><Archive /></el-icon>
                      <el-icon v-else><Folder /></el-icon>
                    </el-icon>
                    <span>{{ record.name }}</span>
                    <el-tag 
                      v-if="record.isCompressed" 
                      size="small" 
                      type="warning"
                    >
                      7z
                    </el-tag>
                  </div>
                  <div class="backup-actions">
                    <el-dropdown @command="handleBackupAction">
                      <el-button type="text" size="small">
                        <el-icon><MoreFilled /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :command="{ action: 'restore', record }">
                            <el-icon><Upload /></el-icon>
                            恢复
                          </el-dropdown-item>
                          <el-dropdown-item :command="{ action: 'delete', record }" divided>
                            <el-icon><Delete /></el-icon>
                            删除
                          </el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                </div>
                
                <div class="backup-info">
                  <div class="info-row">
                    <span class="label">创建时间:</span>
                    <span class="value">{{ driverStore.formatTime(record.time) }}</span>
                  </div>
                  <div class="info-row">
                    <span class="label">文件大小:</span>
                    <span class="value">{{ record.size }}</span>
                  </div>
                  <div v-if="record.driverCount" class="info-row">
                    <span class="label">驱动数量:</span>
                    <span class="value">{{ record.driverCount }} 个</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧信息面板 -->
      <el-col :span="8">
        <!-- 选中备份信息 -->
        <el-card v-if="driverStore.selectedBackup" class="info-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><InfoFilled /></el-icon>
              <span>备份详情</span>
            </div>
          </template>
          
          <div class="backup-details">
            <div class="detail-item">
              <span class="label">名称:</span>
              <span class="value">{{ driverStore.selectedBackup.name }}</span>
            </div>
            <div class="detail-item">
              <span class="label">路径:</span>
              <span class="value" :title="driverStore.selectedBackup.path">
                {{ driverStore.selectedBackup.path }}
              </span>
            </div>
            <div class="detail-item">
              <span class="label">类型:</span>
              <el-tag 
                :type="driverStore.selectedBackup.isCompressed ? 'warning' : 'primary'"
                size="small"
              >
                {{ driverStore.selectedBackup.isCompressed ? '压缩文件' : '文件夹' }}
              </el-tag>
            </div>
            <div class="detail-item">
              <span class="label">大小:</span>
              <span class="value">{{ driverStore.selectedBackup.size }}</span>
            </div>
            <div class="detail-item">
              <span class="label">创建时间:</span>
              <span class="value">{{ driverStore.formatTime(driverStore.selectedBackup.time) }}</span>
            </div>
            <div v-if="driverStore.selectedBackup.driverCount" class="detail-item">
              <span class="label">驱动数量:</span>
              <span class="value">{{ driverStore.selectedBackup.driverCount }} 个</span>
            </div>
          </div>
        </el-card>

        <!-- 系统驱动信息 -->
        <el-card class="drivers-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><Setting /></el-icon>
              <span>系统驱动</span>
              <el-button 
                type="primary" 
                size="small" 
                @click="loadSystemDrivers"
                :loading="driverStore.isLoading"
              >
                <el-icon><Refresh /></el-icon>
                扫描
              </el-button>
            </div>
          </template>
          
          <div v-if="!systemDrivers.length && !driverStore.isLoading" class="empty-state">
            <el-empty description="点击扫描获取系统驱动信息" />
          </div>
          
          <div v-else class="drivers-summary">
            <el-statistic 
              title="检测到的驱动" 
              :value="systemDrivers.length" 
              suffix="个"
            />
            
            <div class="driver-types">
              <div v-for="type in driverTypes" :key="type.name" class="type-item">
                <span class="type-name">{{ type.name }}:</span>
                <span class="type-count">{{ type.count }} 个</span>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 备份统计 -->
        <el-card class="stats-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <el-icon><DataAnalysis /></el-icon>
              <span>备份统计</span>
            </div>
          </template>
          
          <div class="stats-grid">
            <div class="stat-item">
              <div class="stat-value">{{ backupStats.total }}</div>
              <div class="stat-label">总备份数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ backupStats.compressed }}</div>
              <div class="stat-label">压缩备份</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">{{ backupStats.uncompressed }}</div>
              <div class="stat-label">文件夹备份</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useDriverStore } from '../stores/driver'
import { useAppStore } from '../stores/app'

const driverStore = useDriverStore()
const appStore = useAppStore()

// 响应式数据
const systemDrivers = ref([])

// 计算属性
const backupStats = computed(() => driverStore.getBackupStats())

const driverTypes = computed(() => {
  const types = {}
  systemDrivers.value.forEach(driver => {
    const type = driver.driverClass || '未知'
    types[type] = (types[type] || 0) + 1
  })
  
  return Object.entries(types).map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5) // 只显示前5种类型
})

// 方法
const selectBackupPath = async () => {
  // 这里需要调用文件选择对话框
  // 在真实环境中会调用Wails的文件选择API
  try {
    ElMessage.info('文件选择功能需要在Wails环境中运行')
    // const path = await SelectDirectory()
    // driverStore.setBackupPath(path)
  } catch (error) {
    ElMessage.error('选择路径失败: ' + error.message)
  }
}

const startBackup = async () => {
  const result = await ElMessageBox.confirm(
    '确定要备份系统驱动吗？此操作需要管理员权限，可能需要较长时间。',
    '确认备份',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)

  if (!result) return

  try {
    appStore.addLog('info', '开始备份系统驱动')
    appStore.showProgress('正在备份驱动...', 0)
    
    const result = await driverStore.backupSystemDrivers()
    
    ElMessage.success(`备份完成！共备份 ${result.driverCount} 个驱动`)
    appStore.addLog('info', `备份完成，共 ${result.driverCount} 个驱动`)
    
    // 刷新备份历史
    await refreshBackupHistory()
  } catch (error) {
    ElMessage.error('备份失败: ' + error.message)
    appStore.addLog('error', `备份失败: ${error.message}`)
  } finally {
    appStore.hideProgress()
  }
}

const startRestore = async () => {
  if (!driverStore.selectedBackup) {
    ElMessage.warning('请先选择要恢复的备份')
    return
  }

  const result = await ElMessageBox.confirm(
    `确定要从备份 "${driverStore.selectedBackup.name}" 恢复驱动吗？\n\n注意：此操作需要管理员权限。`,
    '确认恢复',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).catch(() => false)

  if (!result) return

  try {
    appStore.addLog('info', `开始从备份恢复驱动: ${driverStore.selectedBackup.name}`)
    appStore.showProgress('正在恢复驱动...', 0)
    
    await driverStore.restoreDriversFromBackup(driverStore.selectedBackup.path)
    
    ElMessage.success('驱动恢复完成！建议重启系统以确保驱动正常工作。')
    appStore.addLog('info', '驱动恢复完成')
  } catch (error) {
    ElMessage.error('恢复失败: ' + error.message)
    appStore.addLog('error', `恢复失败: ${error.message}`)
  } finally {
    appStore.hideProgress()
  }
}

const deleteSelectedBackup = async () => {
  if (!driverStore.selectedBackup) {
    ElMessage.warning('请先选择要删除的备份')
    return
  }

  const result = await ElMessageBox.confirm(
    `确定要删除备份 "${driverStore.selectedBackup.name}" 吗？\n此操作无法撤销。`,
    '确认删除',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'error',
    }
  ).catch(() => false)

  if (!result) return

  try {
    await driverStore.deleteBackupRecord(driverStore.selectedBackup)
    ElMessage.success('备份已删除')
    appStore.addLog('info', `删除备份: ${driverStore.selectedBackup.name}`)
  } catch (error) {
    ElMessage.error('删除失败: ' + error.message)
    appStore.addLog('error', `删除失败: ${error.message}`)
  }
}

const selectBackup = (record) => {
  driverStore.selectBackup(record)
}

const handleBackupAction = async ({ action, record }) => {
  if (action === 'restore') {
    driverStore.selectBackup(record)
    await startRestore()
  } else if (action === 'delete') {
    driverStore.selectBackup(record)
    await deleteSelectedBackup()
  }
}

const refreshBackupHistory = async () => {
  try {
    await driverStore.loadBackupHistory()
    ElMessage.success('备份历史已刷新')
  } catch (error) {
    ElMessage.error('刷新失败: ' + error.message)
  }
}

const loadSystemDrivers = async () => {
  try {
    appStore.addLog('info', '开始扫描系统驱动')
    const drivers = await driverStore.loadSystemDrivers()
    systemDrivers.value = drivers
    ElMessage.success(`扫描完成，检测到 ${drivers.length} 个驱动`)
    appStore.addLog('info', `扫描完成，检测到 ${drivers.length} 个驱动`)
  } catch (error) {
    ElMessage.error('扫描驱动失败: ' + error.message)
    appStore.addLog('error', '扫描驱动失败: ' + error.message)
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadBackupHistory()
  loadSystemDrivers()
})
</script>

<style scoped>
.driver-manager {
  padding: 20px;
}

.stat-card {
  text-align: center;
  padding: 20px;
}

.stat-number {
  font-size: 2em;
  font-weight: bold;
  color: #409eff;
}

.stat-label {
  margin-top: 10px;
  color: #666;
}

.driver-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #eee;
}

.driver-info h4 {
  margin: 0 0 5px 0;
}

.driver-info p {
  margin: 0;
  color: #666;
  font-size: 0.9em;
}
</style>