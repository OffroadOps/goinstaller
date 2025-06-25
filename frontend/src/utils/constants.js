// 应用常量
export const APP_NAME = 'SystemReinstaller'
export const APP_VERSION = '2.0.0'
export const APP_DESCRIPTION = '系统重装助手 - 驱动备份与系统重装工具'

// 主题常量
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark'
}

// 日志级别
export const LOG_LEVELS = {
  DEBUG: 'debug',
  INFO: 'info',
  WARN: 'warn',
  ERROR: 'error'
}

// 下载状态
export const DOWNLOAD_STATUS = {
  PENDING: 'pending',
  DOWNLOADING: 'downloading',
  COMPLETED: 'completed',
  FAILED: 'failed',
  CANCELLED: 'cancelled'
}

// 系统类型
export const SYSTEM_TYPES = {
  WINDOWS: 'Windows',
  WINDOWS_SERVER: 'Windows Server',
  TINY_WINDOWS: 'Tiny Windows',
  LINUX: 'Linux'
}

// 启动模式
export const BOOT_MODES = {
  AUTO: '自动检测',
  UEFI: 'UEFI',
  LEGACY: 'Legacy'
}

// 语言选项
export const LANGUAGES = {
  ALL: '全部',
  ZH_CN: '中文(zh-cn)',
  EN_US: '英文(en-us)',
  JA_JP: '日文(ja-jp)'
}

// 文件类型
export const FILE_TYPES = {
  VHD: 'vhd',
  ISO: 'iso',
  IMG: 'img',
  '7Z': '7z',
  ZIP: 'zip'
}

// 备份类型
export const BACKUP_TYPES = {
  FOLDER: 'folder',
  COMPRESSED: 'compressed'
}

// 操作状态
export const OPERATION_STATUS = {
  IDLE: 'idle',
  RUNNING: 'running',
  SUCCESS: 'success',
  FAILED: 'failed'
}

// 默认配置
export const DEFAULT_CONFIG = {
  backup: {
    enableCompression: true,
    backupLocation: 'driver_backup',
    autoCleanup: false,
    maxBackups: 10
  },
  download: {
    maxConcurrent: 3,
    retryCount: 3,
    timeout: 300000 // 5分钟
  },
  ui: {
    theme: THEMES.LIGHT,
    language: 'zh-CN',
    autoRefresh: true
  }
}

// API端点
export const API_ENDPOINTS = {
  DEFAULT_BASE_URL: 'https://autoinstaller.qkdny.com',
  FALLBACK_VHD_LIST: 'https://www.ipa.gr/vhd/vhdlist.txt',
  FALLBACK_ISO_LIST: 'http://92.112.124.60:8080/filelist.txt'
}

// 错误消息
export const ERROR_MESSAGES = {
  NETWORK_ERROR: '网络连接失败，请检查网络设置',
  SERVER_ERROR: '服务器错误，请稍后重试',
  FILE_NOT_FOUND: '文件未找到',
  PERMISSION_DENIED: '权限不足，请以管理员身份运行',
  DISK_SPACE_INSUFFICIENT: '磁盘空间不足',
  OPERATION_CANCELLED: '操作已取消',
  INVALID_PARAMETERS: '参数无效',
  UNKNOWN_ERROR: '未知错误'
}

// 成功消息
export const SUCCESS_MESSAGES = {
  DOWNLOAD_COMPLETED: '下载完成',
  BACKUP_COMPLETED: '备份完成',
  RESTORE_COMPLETED: '恢复完成',
  OPERATION_SUCCESS: '操作成功'
}

// 文件大小单位
export const FILE_SIZE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB']

// 正则表达式
export const REGEX = {
  VHD_FILENAME: /^(.+)_(\d{8}_\d{6})\.(?:xz|vhd)$/,
  IP_ADDRESS: /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/,
  PORT: /^([1-9][0-9]{0,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])$/
}

// 时间格式
export const TIME_FORMATS = {
  DATE_TIME: 'YYYY-MM-DD HH:mm:ss',
  DATE: 'YYYY-MM-DD',
  TIME: 'HH:mm:ss',
  FILE_NAME: 'YYYYMMDD_HHmmss'
}

// 存储键名
export const STORAGE_KEYS = {
  THEME: 'app-theme',
  LANGUAGE: 'app-language',
  BACKUP_CONFIG: 'driver-backup-config',
  DOWNLOAD_CONFIG: 'download-config',
  LAST_SERVER: 'last-selected-server',
  WINDOW_SIZE: 'window-size',
  WINDOW_POSITION: 'window-position'
}