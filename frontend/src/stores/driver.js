import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { GetSystemDrivers, BackupDrivers, RestoreDrivers, LoadBackupHistory, DeleteBackup } from '../utils/wails'

export const useDriverStore = defineStore('driver', () => {
  // 驱动状态
  const drivers = ref([])
  const backupRecords = ref([])
  const selectedBackup = ref(null)
  const backupPath = ref('driver_backup')
  
  // 操作状态
  const isBackingUp = ref(false)
  const isRestoring = ref(false)
  const isLoading = ref(false)
  
  // 备份配置
  const backupConfig = reactive({
    enableCompression: true,
    backupLocation: 'driver_backup',
    autoCleanup: false,
    maxBackups: 10
  })
  
  // Actions
  const loadSystemDrivers = async () => {
    isLoading.value = true
    try {
      const result = await GetSystemDrivers()
      drivers.value = result || []
      return result
    } catch (error) {
      console.error('获取系统驱动失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }
  
  const loadBackupHistory = async (path) => {
    isLoading.value = true
    try {
      const result = await LoadBackupHistory(path || backupPath.value)
      backupRecords.value = result || []
      return result
    } catch (error) {
      console.error('加载备份历史失败:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }
  
  const backupSystemDrivers = async (targetPath) => {
    isBackingUp.value = true
    try {
      const result = await BackupDrivers(targetPath || backupPath.value)
      
      // 刷新备份历史
      await loadBackupHistory()
      
      return result
    } catch (error) {
      console.error('备份驱动失败:', error)
      throw error
    } finally {
      isBackingUp.value = false
    }
  }
  
  const restoreDriversFromBackup = async (backupPath) => {
    isRestoring.value = true
    try {
      const result = await RestoreDrivers(backupPath)
      return result
    } catch (error) {
      console.error('恢复驱动失败:', error)
      throw error
    } finally {
      isRestoring.value = false
    }
  }
  
  const deleteBackupRecord = async (record) => {
    try {
      await DeleteBackup(record)
      
      // 刷新备份历史
      await loadBackupHistory()
      
      // 如果删除的是当前选中的备份，清除选择
      if (selectedBackup.value && selectedBackup.value.path === record.path) {
        selectedBackup.value = null
      }
    } catch (error) {
      console.error('删除备份失败:', error)
      throw error
    }
  }
  
  const selectBackup = (record) => {
    selectedBackup.value = record
  }
  
  const clearSelection = () => {
    selectedBackup.value = null
  }
  
  const setBackupPath = (path) => {
    backupPath.value = path
  }
  
  const updateBackupConfig = (config) => {
    Object.assign(backupConfig, config)
    
    // 保存到本地存储
    localStorage.setItem('driver-backup-config', JSON.stringify(backupConfig))
  }
  
  // 格式化文件大小
  const formatFileSize = (bytes) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
  }
  
  // 格式化时间
  const formatTime = (timestamp) => {
    return new Date(timestamp).toLocaleString('zh-CN')
  }
  
  // 获取备份统计信息
  const getBackupStats = () => {
    const total = backupRecords.value.length
    const compressed = backupRecords.value.filter(r => r.isCompressed).length
    const uncompressed = total - compressed
    
    return {
      total,
      compressed,
      uncompressed
    }
  }
  
  // 初始化
  const init = () => {
    // 从本地存储恢复配置
    const savedConfig = localStorage.getItem('driver-backup-config')
    if (savedConfig) {
      try {
        const config = JSON.parse(savedConfig)
        Object.assign(backupConfig, config)
      } catch (error) {
        console.error('恢复备份配置失败:', error)
      }
    }
    
    // 设置备份路径
    backupPath.value = backupConfig.backupLocation
  }
  
  return {
    // State
    drivers,
    backupRecords,
    selectedBackup,
    backupPath,
    isBackingUp,
    isRestoring,
    isLoading,
    backupConfig,
    
    // Actions
    loadSystemDrivers,
    loadBackupHistory,
    backupSystemDrivers,
    restoreDriversFromBackup,
    deleteBackupRecord,
    selectBackup,
    clearSelection,
    setBackupPath,
    updateBackupConfig,
    
    // Getters
    formatFileSize,
    formatTime,
    getBackupStats,
    
    // Init
    init
  }
})