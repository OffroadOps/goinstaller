package main

import (
	"context"
	"fmt"
	"runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	fmt.Println("App started successfully")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetSystemInfo 获取系统信息
func (a *App) GetSystemInfo() map[string]interface{} {
	return map[string]interface{}{
		"os":           runtime.GOOS,
		"arch":         runtime.GOARCH,
		"version":      "开发版本",
		"hostname":     "DESKTOP-DEV",
		"memory":       "16 GB",
		"cpu":          "Intel Core i7",
		"computerName": "DESKTOP-DEV",
		"timestamp":    "2024-01-01 12:00:00",
		"osVersion":    fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetAvailableServers 获取可用服务器列表
func (a *App) GetAvailableServers() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"id":       "1",
			"name":     "官方服务器",
			"location": "北京",
			"type":     "主服务器",
			"status":   "online",
		},
		map[string]interface{}{
			"id":       "2",
			"name":     "镜像服务器",
			"location": "上海",
			"type":     "镜像服务器",
			"status":   "online",
		},
	}
}

// GetVHDListFromServer 从服务器获取VHD列表
func (a *App) GetVHDListFromServer(serverId interface{}) []interface{} {
	// 模拟VHD列表数据
	return []interface{}{
		map[string]interface{}{
			"filename":    "windows10_pro_x64_22h2.vhd",
			"displayName": "Windows 10 Pro x64 22H2",
			"system":      "Windows 10",
			"version":     "22H2",
			"language":    "中文(zh-cn)",
			"bootMode":    "UEFI",
			"size":        "4.2 GB",
		},
		map[string]interface{}{
			"filename":    "windows11_pro_x64_23h2.vhd",
			"displayName": "Windows 11 Pro x64 23H2",
			"system":      "Windows 11",
			"version":     "23H2",
			"language":    "中文(zh-cn)",
			"bootMode":    "UEFI",
			"size":        "4.8 GB",
		},
		map[string]interface{}{
			"filename":    "windows_server_2022_datacenter.vhd",
			"displayName": "Windows Server 2022 Datacenter",
			"system":      "Windows Server",
			"version":     "2022",
			"language":    "英文(en-us)",
			"bootMode":    "UEFI",
			"size":        "3.9 GB",
		},
	}
}

// DownloadVHD 下载VHD文件
func (a *App) DownloadVHD(vhdId interface{}, savePath string) map[string]interface{} {
	fmt.Printf("开始下载VHD: %v 到 %s\n", vhdId, savePath)

	return map[string]interface{}{
		"success": true,
		"message": "下载开始",
		"id":      vhdId,
		"path":    savePath,
	}
}

// GetSystemDrivers 获取系统驱动列表
func (a *App) GetSystemDrivers() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":        "NVIDIA Graphics Driver",
			"version":     "528.49",
			"driverClass": "Display",
			"date":        "2023-10-15",
		},
		map[string]interface{}{
			"name":        "Realtek Audio Driver",
			"version":     "6.0.9228.1",
			"driverClass": "Media",
			"date":        "2023-09-20",
		},
		map[string]interface{}{
			"name":        "Intel Network Adapter",
			"version":     "12.19.2.4",
			"driverClass": "Net",
			"date":        "2023-08-30",
		},
	}
}

// BackupDrivers 备份驱动
func (a *App) BackupDrivers(targetPath string) map[string]interface{} {
	fmt.Printf("备份驱动到: %s\n", targetPath)

	return map[string]interface{}{
		"success":     true,
		"message":     "备份完成",
		"driverCount": 15,
		"path":        targetPath,
	}
}

// RestoreDrivers 恢复驱动
func (a *App) RestoreDrivers(backupPath string) map[string]interface{} {
	fmt.Printf("从备份恢复驱动: %s\n", backupPath)

	return map[string]interface{}{
		"success": true,
		"message": "恢复完成",
		"path":    backupPath,
	}
}

// LoadBackupHistory 加载备份历史
func (a *App) LoadBackupHistory() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"name":         "备份_20240101_120000",
			"path":         "C:\\DriverBackup\\backup_20240101_120000.7z",
			"time":         1704096000000, // 时间戳
			"size":         "156 MB",
			"isCompressed": true,
			"driverCount":  15,
		},
		map[string]interface{}{
			"name":         "备份_20231220_150000",
			"path":         "C:\\DriverBackup\\backup_20231220_150000",
			"time":         1703059200000,
			"size":         "234 MB",
			"isCompressed": false,
			"driverCount":  12,
		},
	}
}

// DeleteBackup 删除备份
func (a *App) DeleteBackup(record interface{}) map[string]interface{} {
	fmt.Printf("删除备份: %v\n", record)

	return map[string]interface{}{
		"success": true,
		"message": "删除成功",
	}
}

// SelectFile 选择文件
func (a *App) SelectFile(filters interface{}) string {
	fmt.Printf("选择文件，过滤器: %v\n", filters)
	return ""
}

// SelectDirectory 选择目录
func (a *App) SelectDirectory() string {
	fmt.Println("选择目录")
	return ""
}
