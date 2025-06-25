// Wails前端API封装
// 在开发环境中提供模拟数据，在生产环境中调用真实的Go函数

// 检查是否在Wails环境中运行
const isWailsEnv = typeof window !== 'undefined' && window.go;

// 系统信息相关
export const GetSystemInfo = async () => {
  if (isWailsEnv && window.go.main.App.GetSystemInfo) {
    return await window.go.main.App.GetSystemInfo();
  }
  // 开发环境模拟数据
  return {
    os: 'Windows 10',
    arch: 'amd64',
    version: '10.0.19045',
    hostname: 'DESKTOP-TEST',
    memory: '16 GB',
    cpu: 'Intel Core i7-8700K'
  };
};

// 驱动管理相关
export const GetSystemDrivers = async () => {
  if (isWailsEnv && window.go.main.App.GetSystemDrivers) {
    return await window.go.main.App.GetSystemDrivers();
  }
  // 开发环境模拟数据
  return [];
};

export const BackupDrivers = async (drivers) => {
  if (isWailsEnv && window.go.main.App.BackupDrivers) {
    return await window.go.main.App.BackupDrivers(drivers);
  }
  // 开发环境模拟
  console.log('模拟备份驱动:', drivers);
  return { success: true, message: '备份完成（模拟）' };
};

export const RestoreDrivers = async (backupPath) => {
  if (isWailsEnv && window.go.main.App.RestoreDrivers) {
    return await window.go.main.App.RestoreDrivers(backupPath);
  }
  // 开发环境模拟
  console.log('模拟恢复驱动:', backupPath);
  return { success: true, message: '恢复完成（模拟）' };
};

export const LoadBackupHistory = async () => {
  if (isWailsEnv && window.go.main.App.LoadBackupHistory) {
    return await window.go.main.App.LoadBackupHistory();
  }
  // 开发环境模拟数据
  return [];
};

export const DeleteBackup = async (backupId) => {
  if (isWailsEnv && window.go.main.App.DeleteBackup) {
    return await window.go.main.App.DeleteBackup(backupId);
  }
  // 开发环境模拟
  console.log('模拟删除备份:', backupId);
  return { success: true, message: '删除完成（模拟）' };
};

// 下载相关
export const GetAvailableServers = async () => {
  if (isWailsEnv && window.go.main.App.GetAvailableServers) {
    return await window.go.main.App.GetAvailableServers();
  }
  // 开发环境模拟数据
  return [
    { id: 1, name: '官方服务器', url: 'https://example.com', status: 'online' },
    { id: 2, name: '镜像服务器', url: 'https://mirror.example.com', status: 'online' }
  ];
};

export const GetVHDListFromServer = async (serverId) => {
  if (isWailsEnv && window.go.main.App.GetVHDListFromServer) {
    return await window.go.main.App.GetVHDListFromServer(serverId);
  }
  // 开发环境模拟数据
  return [
    { id: 1, name: 'Windows 10 Pro', size: '4.2 GB', version: '22H2' },
    { id: 2, name: 'Windows 11 Pro', size: '4.8 GB', version: '23H2' }
  ];
};

export const DownloadVHD = async (vhdId, savePath) => {
  if (isWailsEnv && window.go.main.App.DownloadVHD) {
    return await window.go.main.App.DownloadVHD(vhdId, savePath);
  }
  // 开发环境模拟
  console.log('模拟下载VHD:', vhdId, savePath);
  return { success: true, message: '下载开始（模拟）' };
};

// 文件操作相关
export const SelectFile = async (filters) => {
  if (isWailsEnv && window.go.main.App.SelectFile) {
    return await window.go.main.App.SelectFile(filters);
  }
  // 开发环境模拟
  console.log('模拟文件选择:', filters);
  return null;
};

export const SelectDirectory = async () => {
  if (isWailsEnv && window.go.main.App.SelectDirectory) {
    return await window.go.main.App.SelectDirectory();
  }
  // 开发环境模拟
  console.log('模拟目录选择');
  return null;
};