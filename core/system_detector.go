package core

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// SystemDetector 系统检测器
type SystemDetector struct {
	isAdmin    bool
	workingDir string
}

// SystemDetectionResult 系统检测结果
type SystemDetectionResult struct {
	OSInfo         map[string]interface{} `json:"os_info"`
	BootMode       string                 `json:"boot_mode"`
	PartitionTable map[string]interface{} `json:"partition_table"`
	EFIPartition   map[string]interface{} `json:"efi_partition"`
	HardwareInfo   map[string]interface{} `json:"hardware_info"`
	DiskInfo       map[string]interface{} `json:"disk_info"`
	Compatibility  map[string]interface{} `json:"compatibility"`
}

// NewSystemDetector 创建系统检测器
func NewSystemDetector() *SystemDetector {
	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "." // 如果获取失败，使用当前目录
	}

	return &SystemDetector{
		isAdmin:    checkAdminPrivileges(),
		workingDir: workingDir,
	}
}

// Initialize 初始化
func (sd *SystemDetector) Initialize() error {
	return nil
}

// checkAdminPrivileges 检查管理员权限
func checkAdminPrivileges() bool {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		return err == nil
	} else {
		return os.Geteuid() == 0
	}
}

// GetCompleteSystemInfo 获取完整系统信息
func (sd *SystemDetector) GetCompleteSystemInfo() (*SystemDetectionResult, error) {
	result := &SystemDetectionResult{}

	// 获取操作系统信息
	result.OSInfo = sd.getOSInfo()

	// 检测启动模式
	result.BootMode = sd.detectBootMode()

	// 检测分区表
	result.PartitionTable = sd.detectPartitionTable()

	// 检测EFI分区
	result.EFIPartition = sd.detectEFIPartition()

	// 获取硬件信息
	result.HardwareInfo = sd.getHardwareInfo()

	// 获取磁盘信息
	result.DiskInfo = sd.getDiskInfo()

	// 计算兼容性
	result.Compatibility = sd.calculateCompatibility(result)

	if !sd.isAdmin {
		// 如果没有管理员权限，提供基础信息但标记需要权限
		result.Compatibility["warning"] = "需要管理员权限获取完整信息"
	}

	return result, nil
}

// getOSInfo 获取操作系统信息
func (sd *SystemDetector) getOSInfo() map[string]interface{} {
	return map[string]interface{}{
		"system":  runtime.GOOS,
		"arch":    runtime.GOARCH,
		"version": getOSVersion(),
	}
}

// getOSVersion 获取操作系统版本
func getOSVersion() string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "ver")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	} else {
		cmd := exec.Command("uname", "-r")
		output, err := cmd.Output()
		if err == nil {
			return strings.TrimSpace(string(output))
		}
	}
	return "Unknown"
}

// detectBootMode 检测启动模式
func (sd *SystemDetector) detectBootMode() string {
	if !sd.isAdmin {
		return "Unknown (需要管理员权限)"
	}

	if runtime.GOOS == "windows" {
		// Windows检测
		cmd := exec.Command("bcdedit", "/enum", "{current}")
		output, err := cmd.Output()
		if err == nil {
			outputStr := strings.ToLower(string(output))
			if strings.Contains(outputStr, "winload.efi") {
				return "UEFI"
			} else if strings.Contains(outputStr, "winload.exe") {
				return "Legacy"
			}
		}

		// 备用检测方法
		cmd = exec.Command("powershell", "-Command", "Get-ComputerInfo | Select-Object BiosFirmwareType")
		output, err = cmd.Output()
		if err == nil {
			if strings.Contains(string(output), "Uefi") {
				return "UEFI"
			} else if strings.Contains(string(output), "Bios") {
				return "Legacy"
			}
		}
	} else {
		// Linux检测
		if _, err := os.Stat("/sys/firmware/efi"); err == nil {
			return "UEFI"
		} else {
			return "Legacy"
		}
	}

	return "Unknown"
}

// detectPartitionTable 检测分区表类型
func (sd *SystemDetector) detectPartitionTable() map[string]interface{} {
	if !sd.isAdmin {
		return map[string]interface{}{
			"type":    "Unknown",
			"details": "需要管理员权限",
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", "Get-Disk | Select-Object PartitionStyle")
		output, err := cmd.Output()
		if err == nil {
			if strings.Contains(string(output), "GPT") {
				return map[string]interface{}{
					"type":    "GPT",
					"details": "GPT partition style detected",
				}
			} else if strings.Contains(string(output), "MBR") {
				return map[string]interface{}{
					"type":    "MBR",
					"details": "MBR partition style detected",
				}
			}
		}
	} else {
		cmd := exec.Command("lsblk", "-o", "NAME,PTTYPE")
		output, err := cmd.Output()
		if err == nil {
			outputStr := strings.ToLower(string(output))
			if strings.Contains(outputStr, "gpt") {
				return map[string]interface{}{
					"type":    "GPT",
					"details": "GPT partition table detected",
				}
			} else if strings.Contains(outputStr, "dos") {
				return map[string]interface{}{
					"type":    "MBR",
					"details": "MBR partition table detected",
				}
			}
		}
	}

	return map[string]interface{}{
		"type":    "Unknown",
		"details": "Unable to detect partition table",
	}
}

// detectEFIPartition 检测EFI分区
func (sd *SystemDetector) detectEFIPartition() map[string]interface{} {
	if !sd.isAdmin {
		return map[string]interface{}{
			"exists":  false,
			"details": "需要管理员权限",
		}
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("powershell", "-Command", "Get-Partition | Where-Object {$_.Type -eq 'System'}")
		output, err := cmd.Output()
		if err == nil && len(strings.TrimSpace(string(output))) > 0 {
			return map[string]interface{}{
				"exists":  true,
				"details": "EFI System Partition found",
			}
		}
	} else {
		cmd := exec.Command("lsblk", "-o", "NAME,FSTYPE,MOUNTPOINT")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.Contains(line, "/boot/efi") || strings.Contains(line, "vfat") {
					return map[string]interface{}{
						"exists":  true,
						"details": "EFI partition found",
					}
				}
			}
		}
	}

	return map[string]interface{}{
		"exists":  false,
		"details": "No EFI partition found",
	}
}

// getHardwareInfo 获取硬件信息
func (sd *SystemDetector) getHardwareInfo() map[string]interface{} {
	hardwareInfo := map[string]interface{}{}

	if runtime.GOOS == "windows" {
		// 获取CPU信息
		cmd := exec.Command("wmic", "cpu", "get", "name", "/value")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "Name=") {
					hardwareInfo["cpu"] = strings.TrimSpace(strings.TrimPrefix(line, "Name="))
					break
				}
			}
		}

		// 获取内存信息
		cmd = exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/value")
		output, err = cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "TotalPhysicalMemory=") {
					memStr := strings.TrimSpace(strings.TrimPrefix(line, "TotalPhysicalMemory="))
					if mem, err := strconv.ParseInt(memStr, 10, 64); err == nil {
						hardwareInfo["memory_bytes"] = mem
						hardwareInfo["memory_gb"] = mem / (1024 * 1024 * 1024)
					}
					break
				}
			}
		}
	} else {
		// Linux硬件信息获取
		if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "model name") {
					parts := strings.Split(line, ":")
					if len(parts) > 1 {
						hardwareInfo["cpu"] = strings.TrimSpace(parts[1])
						break
					}
				}
			}
		}

		if data, err := os.ReadFile("/proc/meminfo"); err == nil {
			lines := strings.Split(string(data), "\n")
			for _, line := range lines {
				if strings.HasPrefix(line, "MemTotal:") {
					parts := strings.Fields(line)
					if len(parts) > 1 {
						if mem, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
							hardwareInfo["memory_kb"] = mem
							hardwareInfo["memory_gb"] = mem / (1024 * 1024)
						}
					}
					break
				}
			}
		}
	}

	return hardwareInfo
}

// getDiskInfo 获取磁盘信息
func (sd *SystemDetector) getDiskInfo() map[string]interface{} {
	diskInfo := map[string]interface{}{}

	if runtime.GOOS == "windows" {
		cmd := exec.Command("wmic", "logicaldisk", "get", "size,freespace,caption", "/value")
		output, err := cmd.Output()
		if err == nil {
			disks := []map[string]interface{}{}
			lines := strings.Split(string(output), "\n")
			var currentDisk map[string]interface{}

			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == "" {
					if currentDisk != nil && len(currentDisk) > 0 {
						disks = append(disks, currentDisk)
						currentDisk = nil
					}
					continue
				}

				if currentDisk == nil {
					currentDisk = make(map[string]interface{})
				}

				if strings.HasPrefix(line, "Caption=") {
					currentDisk["drive"] = strings.TrimPrefix(line, "Caption=")
				} else if strings.HasPrefix(line, "FreeSpace=") {
					freeStr := strings.TrimPrefix(line, "FreeSpace=")
					if freeStr != "" {
						if free, err := strconv.ParseInt(freeStr, 10, 64); err == nil {
							currentDisk["free_bytes"] = free
						}
					}
				} else if strings.HasPrefix(line, "Size=") {
					sizeStr := strings.TrimPrefix(line, "Size=")
					if sizeStr != "" {
						if size, err := strconv.ParseInt(sizeStr, 10, 64); err == nil {
							currentDisk["total_bytes"] = size
						}
					}
				}
			}

			if currentDisk != nil && len(currentDisk) > 0 {
				disks = append(disks, currentDisk)
			}

			diskInfo["disks"] = disks
		}
	} else {
		cmd := exec.Command("df", "-h")
		output, err := cmd.Output()
		if err == nil {
			diskInfo["df_output"] = string(output)
		}
	}

	return diskInfo
}

// calculateCompatibility 计算兼容性评分
func (sd *SystemDetector) calculateCompatibility(result *SystemDetectionResult) map[string]interface{} {
	score := 0
	level := "不兼容"
	issues := []string{}

	if !sd.isAdmin {
		return map[string]interface{}{
			"overall_score": 50,
			"level":         "需要管理员权限进行完整检测",
			"issues":        []string{"请以管理员身份运行程序以获取完整的兼容性信息"},
		}
	}

	// 检查启动模式
	if result.BootMode == "UEFI" {
		score += 30
	} else if result.BootMode == "Legacy" {
		score += 20
		issues = append(issues, "建议使用UEFI启动模式")
	} else {
		issues = append(issues, "无法检测启动模式")
	}

	// 检查分区表
	if ptType, ok := result.PartitionTable["type"].(string); ok {
		if ptType == "GPT" {
			score += 25
		} else if ptType == "MBR" {
			score += 15
			issues = append(issues, "建议使用GPT分区表")
		} else {
			issues = append(issues, "无法检测分区表类型")
		}
	}

	// 检查EFI分区
	if efiExists, ok := result.EFIPartition["exists"].(bool); ok && efiExists {
		score += 20
	} else {
		issues = append(issues, "未找到EFI分区")
	}

	// 检查硬件信息
	if result.HardwareInfo != nil {
		score += 15

		// 检查内存
		if memGB, ok := result.HardwareInfo["memory_gb"].(int64); ok {
			if memGB >= 4 {
				score += 10
			} else {
				issues = append(issues, "内存不足4GB，可能影响性能")
			}
		}
	}

	// 确定兼容性等级
	if score >= 80 {
		level = "完全兼容"
	} else if score >= 60 {
		level = "基本兼容"
	} else if score >= 40 {
		level = "部分兼容"
	} else if score >= 20 {
		level = "兼容性差"
	} else {
		level = "不兼容"
	}

	return map[string]interface{}{
		"overall_score": score,
		"level":         level,
		"issues":        issues,
	}
}
