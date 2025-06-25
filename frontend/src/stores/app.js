import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'

export const useAppStore = defineStore('app', () => {
  // 应用状态
  const loading = ref(false)
  const theme = ref('light')
  const currentTab = ref('vhd')
  
  // 系统信息
  const systemInfo = reactive({
    computerName: '',
    timestamp: '',
    osVersion: '',
    arch: ''
  })
  
  // 服务器列表
  const servers = ref([])
  const selectedServer = ref(null)
  
  // 进度状态
  const progressInfo = reactive({
    visible: false,
    percentage: 0,
    status: '',
    text: ''
  })
  
  // 日志
  const logs = ref([])
  
  // Actions
  const setLoading = (state) => {
    loading.value = state
  }
  
  const setTheme = (newTheme) => {
    theme.value = newTheme
    localStorage.setItem('app-theme', newTheme)
  }
  
  const setCurrentTab = (tab) => {
    currentTab.value = tab
  }
  
  const updateSystemInfo = (info) => {
    Object.assign(systemInfo, info)
  }
  
  const setServers = (serverList) => {
    servers.value = serverList
  }
  
  const selectServer = (server) => {
    selectedServer.value = server
  }
  
  const showProgress = (text, percentage = 0, status = '') => {
    progressInfo.visible = true
    progressInfo.text = text
    progressInfo.percentage = percentage
    progressInfo.status = status
  }
  
  const updateProgress = (percentage, text) => {
    progressInfo.percentage = percentage
    if (text) progressInfo.text = text
  }
  
  const hideProgress = () => {
    progressInfo.visible = false
    progressInfo.percentage = 0
    progressInfo.text = ''
    progressInfo.status = ''
  }
  
  const addLog = (level, message) => {
    const log = {
      id: Date.now() + Math.random(),
      time: new Date().toLocaleTimeString(),
      level: level.toUpperCase(),
      message,
      timestamp: Date.now()
    }
    logs.value.push(log)
    
    // 限制日志数量，避免内存泄漏
    if (logs.value.length > 1000) {
      logs.value.splice(0, 100)
    }
    
    console.log(`[${log.level}] ${message}`)
  }
  
  const clearLogs = () => {
    logs.value = []
    console.clear()
  }
  
  const exportLogs = () => {
    const logText = logs.value.map(log => 
      `[${log.time}] [${log.level}] ${log.message}`
    ).join('\n')
    
    const blob = new Blob([logText], { type: 'text/plain;charset=utf-8' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `system-log-${new Date().toISOString().slice(0, 10)}.txt`
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    URL.revokeObjectURL(url)
  }
  
  // 初始化
  const init = () => {
    console.log('App store initializing...')
    
    // 从本地存储恢复主题
    const savedTheme = localStorage.getItem('app-theme')
    if (savedTheme) {
      theme.value = savedTheme
    }
    
    // 模拟获取系统信息
    updateSystemInfo({
      computerName: 'DESKTOP-DEV',
      timestamp: new Date().toLocaleString(),
      osVersion: 'Windows 10',
      arch: 'x64'
    })
    
    addLog('info', '应用初始化完成')
    console.log('App store initialized')
  }
  
  return {
    // State
    loading,
    theme,
    currentTab,
    systemInfo,
    servers,
    selectedServer,
    progressInfo,
    logs,
    
    // Actions
    setLoading,
    setTheme,
    setCurrentTab,
    updateSystemInfo,
    setServers,
    selectServer,
    showProgress,
    updateProgress,
    hideProgress,
    addLog,
    clearLogs,
    exportLogs,
    init
  }
})