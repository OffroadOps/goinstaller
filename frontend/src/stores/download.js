import { defineStore } from 'pinia'
import { ref, reactive, computed } from 'vue'
import { GetAvailableServers, GetVHDListFromServer, DownloadVHD } from '../utils/wails'

export const useDownloadStore = defineStore('download', () => {
  // 服务器和VHD数据
  const servers = ref([])
  const selectedServer = ref(null)
  const allVHDs = ref([])
  const filteredVHDs = ref([])
  const selectedVHD = ref(null)
  
  // 筛选条件
  const filters = reactive({
    systemType: '全部',
    language: '全部',
    bootMode: '自动检测'
  })
  
  // 下载状态
  const downloadQueue = ref([])
  const activeDownloads = ref([])
  const downloadHistory = ref([])
  
  // 加载状态
  const isLoadingServers = ref(false)
  const isLoadingVHDs = ref(false)
  
  // Getters
  const hasSelectedServer = computed(() => selectedServer.value !== null)
  const hasSelectedVHD = computed(() => selectedVHD.value !== null)
  const isDownloading = computed(() => activeDownloads.value.length > 0)
  
  // Actions
  const loadServers = async () => {
    isLoadingServers.value = true
    try {
      const result = await GetAvailableServers()
      servers.value = result || []
      
      // 自动选择第一个服务器
      if (servers.value.length > 0 && !selectedServer.value) {
        await selectServer(servers.value[0])
      }
      
      return result
    } catch (error) {
      console.error('获取服务器列表失败:', error)
      throw error
    } finally {
      isLoadingServers.value = false
    }
  }
  
  const selectServer = async (server) => {
    selectedServer.value = server
    
    // 清除之前的VHD数据
    allVHDs.value = []
    filteredVHDs.value = []
    selectedVHD.value = null
    
    // 加载VHD列表
    await loadVHDList()
  }
  
  const loadVHDList = async () => {
    if (!selectedServer.value) return
    
    isLoadingVHDs.value = true
    try {
      const result = await GetVHDListFromServer(selectedServer.value)
      allVHDs.value = result || []
      
      // 应用筛选
      applyFilters()
      
      return result
    } catch (error) {
      console.error('获取VHD列表失败:', error)
      throw error
    } finally {
      isLoadingVHDs.value = false
    }
  }
  
  const applyFilters = () => {
    let filtered = [...allVHDs.value]
    
    // 按系统类型筛选
    if (filters.systemType !== '全部') {
      filtered = filtered.filter(vhd => 
        vhd.system.includes(filters.systemType)
      )
    }
    
    // 按语言筛选
    if (filters.language !== '全部') {
      filtered = filtered.filter(vhd => 
        vhd.language === filters.language
      )
    }
    
    // 按启动模式筛选
    if (filters.bootMode !== '自动检测') {
      filtered = filtered.filter(vhd => 
        vhd.bootMode === filters.bootMode
      )
    }
    
    filteredVHDs.value = filtered
    
    // 如果当前选中的VHD不在筛选结果中，清除选择
    if (selectedVHD.value && !filtered.find(vhd => vhd.filename === selectedVHD.value.filename)) {
      selectedVHD.value = null
    }
  }
  
  const setFilter = (key, value) => {
    filters[key] = value
    applyFilters()
  }
  
  const selectVHD = (vhd) => {
    selectedVHD.value = vhd
  }
  
  const startDownload = async (vhd) => {
    if (!vhd) vhd = selectedVHD.value
    if (!vhd) return
    
    // 检查是否已在下载队列中
    if (downloadQueue.value.find(item => item.filename === vhd.filename)) {
      throw new Error('该文件已在下载队列中')
    }
    
    // 检查是否正在下载
    if (activeDownloads.value.find(item => item.filename === vhd.filename)) {
      throw new Error('该文件正在下载中')
    }
    
    const downloadItem = {
      ...vhd,
      id: Date.now() + Math.random(),
      status: 'downloading',
      progress: 0,
      speed: '',
      startTime: Date.now(),
      error: null
    }
    
    // 添加到活动下载列表
    activeDownloads.value.push(downloadItem)
    
    try {
      await DownloadVHD(vhd)
      
      // 下载完成
      downloadItem.status = 'completed'
      downloadItem.endTime = Date.now()
      
      // 移动到历史记录
      downloadHistory.value.unshift({ ...downloadItem })
      
    } catch (error) {
      downloadItem.status = 'failed'
      downloadItem.error = error.message
      console.error('下载失败:', error)
      throw error
    } finally {
      // 从活动下载列表移除
      const index = activeDownloads.value.findIndex(item => item.id === downloadItem.id)
      if (index !== -1) {
        activeDownloads.value.splice(index, 1)
      }
    }
  }
  
  const cancelDownload = (downloadId) => {
    const index = activeDownloads.value.findIndex(item => item.id === downloadId)
    if (index !== -1) {
      activeDownloads.value.splice(index, 1)
    }
  }
  
  const clearDownloadHistory = () => {
    downloadHistory.value = []
  }
  
  const retryDownload = async (historyItem) => {
    const vhd = {
      filename: historyItem.filename,
      downloadURL: historyItem.downloadURL,
      displayName: historyItem.displayName
    }
    
    await startDownload(vhd)
  }
  
  // 获取系统类型选项
  const getSystemTypeOptions = () => {
    const types = new Set(['全部'])
    allVHDs.value.forEach(vhd => {
      if (vhd.system) types.add(vhd.system)
    })
    return Array.from(types)
  }
  
  // 获取语言选项
  const getLanguageOptions = () => {
    const languages = new Set(['全部'])
    allVHDs.value.forEach(vhd => {
      if (vhd.language) languages.add(vhd.language)
    })
    return Array.from(languages)
  }
  
  // 获取启动模式选项
  const getBootModeOptions = () => {
    return ['自动检测', 'UEFI', 'Legacy']
  }
  
  // 初始化
  const init = async () => {
    try {
      await loadServers()
    } catch (error) {
      console.error('初始化下载管理器失败:', error)
    }
  }
  
  return {
    // State
    servers,
    selectedServer,
    allVHDs,
    filteredVHDs,
    selectedVHD,
    filters,
    downloadQueue,
    activeDownloads,
    downloadHistory,
    isLoadingServers,
    isLoadingVHDs,
    
    // Getters
    hasSelectedServer,
    hasSelectedVHD,
    isDownloading,
    
    // Actions
    loadServers,
    selectServer,
    loadVHDList,
    applyFilters,
    setFilter,
    selectVHD,
    startDownload,
    cancelDownload,
    clearDownloadHistory,
    retryDownload,
    
    // Options
    getSystemTypeOptions,
    getLanguageOptions,
    getBootModeOptions,
    
    // Init
    init
  }
})