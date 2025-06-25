package core

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// SystemInstaller 系统安装器
type SystemInstaller struct {
	config          map[string]interface{}
	licenseVerified bool
	userCredits     int
	userType        string
	progress        InstallProgress
	progressMutex   sync.RWMutex
	installRunning  bool
	stopChannel     chan bool
	reinstallPath   string
	workingDir      string
}

// InstallProgress 安装进度
type InstallProgress struct {
	Percentage int    `json:"percentage"`
	Message    string `json:"message"`
	Status     string `json:"status"` // running, success, error, stopped
}

// InstallOptions 安装选项
type InstallOptions struct {
	OSType       string            `json:"os_type"`       // linux, windows, dd
	System       string            `json:"system"`        // 系统名称
	Version      string            `json:"version"`       // 版本
	Username     string            `json:"username"`      // 用户名
	Password     string            `json:"password"`      // 密码
	SSHKey       string            `json:"ssh_key"`       // SSH密钥
	SSHPort      int               `json:"ssh_port"`      // SSH端口
	WebPort      int               `json:"web_port"`      // Web端口
	ImageURL     string            `json:"image_url"`     // 镜像URL (DD安装)
	ISOURL       string            `json:"iso_url"`       // ISO URL (Windows)
	ImageName    string            `json:"image_name"`    // Windows镜像名称
	Language     string            `json:"language"`      // 语言
	Minimal      bool              `json:"minimal"`       // 最小安装
	AllowPing    bool              `json:"allow_ping"`    // 允许ping
	RDPPort      int               `json:"rdp_port"`      // RDP端口
	Drivers      []string          `json:"drivers"`       // 驱动列表
	ExtraOptions map[string]string `json:"extra_options"` // 额外选项
}

// DDImageInfo DD镜像信息
type DDImageInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Size        string `json:"size"`
	Format      string `json:"format"`      // raw, vhd, qcow2
	Compression string `json:"compression"` // gz, xz, zst, tar
	Supported   bool   `json:"supported"`
}

// WindowsISOInfo Windows ISO信息
type WindowsISOInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Languages   []string `json:"languages"`
	ImageNames  []string `json:"image_names"`
	DownloadURL string   `json:"download_url"`
	Size        string   `json:"size"`
	Supported   bool     `json:"supported"`
}

// NewSystemInstaller 创建系统安装器
func NewSystemInstaller() *SystemInstaller {
	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "." // 如果获取失败，使用当前目录
	}

	return &SystemInstaller{
		config:      make(map[string]interface{}),
		stopChannel: make(chan bool, 1),
		progress: InstallProgress{
			Percentage: 0,
			Message:    "就绪",
			Status:     "ready",
		},
		workingDir:    workingDir,
		reinstallPath: filepath.Join(workingDir, "reinstall"),
	}
}

// Initialize 初始化安装器
func (si *SystemInstaller) Initialize() error {
	si.loadConfig()
	return si.setupReinstallScript()
}

// setupReinstallScript 设置reinstall脚本
func (si *SystemInstaller) setupReinstallScript() error {
	// 检查reinstall脚本是否存在
	scriptPath := filepath.Join(si.reinstallPath, "reinstall.sh")
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(si.reinstallPath, "reinstall.bat")
	}

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		// 如果不存在，尝试从offline目录复制
		offlineScript := filepath.Join(si.workingDir, "AutoInstaller", "reinstall-offline", "reinstall", "reinstall.sh")
		if runtime.GOOS == "windows" {
			offlineScript = filepath.Join(si.workingDir, "AutoInstaller", "reinstall-offline", "reinstall", "reinstall.bat")
		}

		if _, err := os.Stat(offlineScript); err == nil {
			// 创建目录
			os.MkdirAll(si.reinstallPath, 0755)
			// 复制脚本
			return si.copyFile(offlineScript, scriptPath)
		}

		// 如果找不到offline脚本，创建一个基本的脚本文件以避免错误
		os.MkdirAll(si.reinstallPath, 0755)
		if runtime.GOOS == "windows" {
			// 创建基本的Windows bat脚本
			content := "@echo off\necho SystemReinstaller Script\necho Args: %*\npause\n"
			return os.WriteFile(scriptPath, []byte(content), 0755)
		} else {
			// 创建基本的Linux shell脚本
			content := "#!/bin/bash\necho \"SystemReinstaller Script\"\necho \"Args: $@\"\nread -p \"Press enter to continue...\"\n"
			return os.WriteFile(scriptPath, []byte(content), 0755)
		}
	}

	return nil
}

// copyFile 复制文件
func (si *SystemInstaller) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// InstallSystem 安装系统主方法
func (si *SystemInstaller) InstallSystem(options InstallOptions) error {
	si.progressMutex.Lock()
	si.installRunning = true
	si.progress = InstallProgress{
		Percentage: 0,
		Message:    "开始安装...",
		Status:     "running",
	}
	si.progressMutex.Unlock()

	defer func() {
		si.progressMutex.Lock()
		si.installRunning = false
		si.progressMutex.Unlock()
	}()

	switch options.OSType {
	case "linux":
		return si.installLinuxSystem(options)
	case "windows":
		return si.installWindowsSystem(options)
	case "dd":
		return si.installDDImage(options)
	default:
		return fmt.Errorf("不支持的安装类型: %s", options.OSType)
	}
}

// installLinuxSystem 安装Linux系统
func (si *SystemInstaller) installLinuxSystem(options InstallOptions) error {
	si.updateProgress(10, "准备Linux安装...")

	// 构建reinstall命令
	args := []string{}

	// 添加系统和版本
	if options.Version != "" {
		args = append(args, options.System, options.Version)
	} else {
		args = append(args, options.System)
	}

	// 添加选项
	if options.Password != "" {
		args = append(args, "--password", options.Password)
	}
	if options.SSHKey != "" {
		args = append(args, "--ssh-key", options.SSHKey)
	}
	if options.SSHPort > 0 {
		args = append(args, "--ssh-port", fmt.Sprintf("%d", options.SSHPort))
	}
	if options.WebPort > 0 {
		args = append(args, "--web-port", fmt.Sprintf("%d", options.WebPort))
	}
	if options.Minimal {
		args = append(args, "--minimal")
	}

	return si.executeReinstallScript(args)
}

// installWindowsSystem 安装Windows系统
func (si *SystemInstaller) installWindowsSystem(options InstallOptions) error {
	si.updateProgress(10, "准备Windows安装...")

	// 构建reinstall命令
	args := []string{"windows"}

	// 添加镜像名称
	if options.ImageName != "" {
		args = append(args, "--image-name", options.ImageName)
	}

	// 添加ISO URL
	if options.ISOURL != "" {
		args = append(args, "--iso", options.ISOURL)
	}

	// 添加语言
	if options.Language != "" {
		args = append(args, "--lang", options.Language)
	}

	// 添加通用选项
	if options.Password != "" {
		args = append(args, "--password", options.Password)
	}
	if options.SSHPort > 0 {
		args = append(args, "--ssh-port", fmt.Sprintf("%d", options.SSHPort))
	}
	if options.RDPPort > 0 {
		args = append(args, "--rdp-port", fmt.Sprintf("%d", options.RDPPort))
	}
	if options.AllowPing {
		args = append(args, "--allow-ping")
	}

	// 添加驱动
	for _, driver := range options.Drivers {
		args = append(args, "--add-driver", driver)
	}

	return si.executeReinstallScript(args)
}

// installDDImage 安装DD镜像
func (si *SystemInstaller) installDDImage(options InstallOptions) error {
	si.updateProgress(10, "准备DD安装...")

	// 构建reinstall命令
	args := []string{"dd", "--img", options.ImageURL}

	// 添加通用选项
	if options.Password != "" {
		args = append(args, "--password", options.Password)
	}
	if options.SSHKey != "" {
		args = append(args, "--ssh-key", options.SSHKey)
	}
	if options.SSHPort > 0 {
		args = append(args, "--ssh-port", fmt.Sprintf("%d", options.SSHPort))
	}

	return si.executeReinstallScript(args)
}

// executeReinstallScript 执行reinstall脚本
func (si *SystemInstaller) executeReinstallScript(args []string) error {
	si.updateProgress(20, "执行安装脚本...")

	scriptPath := filepath.Join(si.reinstallPath, "reinstall.sh")
	if runtime.GOOS == "windows" {
		scriptPath = filepath.Join(si.reinstallPath, "reinstall.bat")
	}

	// 检查脚本是否存在
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("reinstall脚本不存在: %s", scriptPath)
	}

	// 创建命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command(scriptPath)
		cmd.Args = append(cmd.Args, args...)
	} else {
		cmd = exec.Command("bash", scriptPath)
		cmd.Args = append(cmd.Args, args...)
	}

	// 设置工作目录
	cmd.Dir = si.reinstallPath

	// 创建管道获取输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// 启动命令
	if err := cmd.Start(); err != nil {
		return err
	}

	// 监控输出
	go si.monitorOutput(stdout, stderr)

	// 等待完成
	err = cmd.Wait()
	if err != nil {
		si.updateProgress(0, fmt.Sprintf("安装失败: %v", err))
		si.progressMutex.Lock()
		si.progress.Status = "error"
		si.progressMutex.Unlock()
		return err
	}

	si.updateProgress(100, "安装完成")
	si.progressMutex.Lock()
	si.progress.Status = "success"
	si.progressMutex.Unlock()

	return nil
}

// monitorOutput 监控命令输出
func (si *SystemInstaller) monitorOutput(stdout, stderr io.ReadCloser) {
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if err != nil {
				break
			}
			output := string(buf[:n])
			si.parseProgress(output)
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if err != nil {
				break
			}
			output := string(buf[:n])
			si.parseProgress(output)
		}
	}()
}

// parseProgress 解析进度信息
func (si *SystemInstaller) parseProgress(output string) {
	// 根据输出解析进度
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// 解析不同阶段的进度
		if strings.Contains(line, "Downloading") {
			si.updateProgress(30, "下载中...")
		} else if strings.Contains(line, "Extracting") {
			si.updateProgress(50, "解压中...")
		} else if strings.Contains(line, "Installing") {
			si.updateProgress(70, "安装中...")
		} else if strings.Contains(line, "Configuring") {
			si.updateProgress(90, "配置中...")
		} else if strings.Contains(line, "Rebooting") {
			si.updateProgress(95, "重启中...")
		}
	}
}

// updateProgress 更新进度
func (si *SystemInstaller) updateProgress(percentage int, message string) {
	si.progressMutex.Lock()
	defer si.progressMutex.Unlock()

	si.progress.Percentage = percentage
	si.progress.Message = message
}

// GetProgress 获取安装进度
func (si *SystemInstaller) GetProgress() InstallProgress {
	si.progressMutex.RLock()
	defer si.progressMutex.RUnlock()
	return si.progress
}

// StopInstallation 停止安装
func (si *SystemInstaller) StopInstallation() error {
	if !si.installRunning {
		return fmt.Errorf("没有正在运行的安装任务")
	}

	select {
	case si.stopChannel <- true:
		si.progressMutex.Lock()
		si.progress.Status = "stopped"
		si.progress.Message = "安装已停止"
		si.installRunning = false
		si.progressMutex.Unlock()
		return nil
	default:
		return fmt.Errorf("无法停止安装")
	}
}

// GetSupportedLinuxSystems 获取支持的Linux系统列表
func (si *SystemInstaller) GetSupportedLinuxSystems() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"name":        "ubuntu",
			"distro":      "Ubuntu",
			"versions":    []string{"16.04", "18.04", "20.04", "22.04", "24.04", "25.04"},
			"description": "Ubuntu Linux",
			"supported":   true,
		},
		{
			"name":        "debian",
			"distro":      "Debian",
			"versions":    []string{"9", "10", "11", "12"},
			"description": "Debian Linux",
			"supported":   true,
		},
		{
			"name":        "centos",
			"distro":      "CentOS",
			"versions":    []string{"9", "10"},
			"description": "CentOS Linux",
			"supported":   true,
		},
		{
			"name":        "rocky",
			"distro":      "Rocky Linux",
			"versions":    []string{"8", "9", "10"},
			"description": "Rocky Linux",
			"supported":   true,
		},
		{
			"name":        "almalinux",
			"distro":      "AlmaLinux",
			"versions":    []string{"8", "9", "10"},
			"description": "AlmaLinux",
			"supported":   true,
		},
		{
			"name":        "fedora",
			"distro":      "Fedora",
			"versions":    []string{"41", "42"},
			"description": "Fedora Linux",
			"supported":   true,
		},
		{
			"name":        "alpine",
			"distro":      "Alpine",
			"versions":    []string{"3.19", "3.20", "3.21", "3.22"},
			"description": "Alpine Linux",
			"supported":   true,
		},
		{
			"name":        "arch",
			"distro":      "Arch",
			"versions":    []string{"latest"},
			"description": "Arch Linux",
			"supported":   true,
		},
		{
			"name":        "kali",
			"distro":      "Kali",
			"versions":    []string{"latest"},
			"description": "Kali Linux",
			"supported":   true,
		},
		{
			"name":        "opensuse",
			"distro":      "openSUSE",
			"versions":    []string{"15.6", "tumbleweed"},
			"description": "openSUSE Linux",
			"supported":   true,
		},
	}
}

// GetSupportedWindowsSystems 获取支持的Windows系统列表
func (si *SystemInstaller) GetSupportedWindowsSystems() []WindowsISOInfo {
	return []WindowsISOInfo{
		{
			Name:        "Windows 10",
			Version:     "10",
			Languages:   []string{"zh-cn", "en-us", "ja-jp"},
			ImageNames:  []string{"Windows 10 Pro", "Windows 10 Enterprise", "Windows 10 Home"},
			DownloadURL: "https://www.microsoft.com/software-download/windows10",
			Supported:   true,
		},
		{
			Name:        "Windows 11",
			Version:     "11",
			Languages:   []string{"zh-cn", "en-us", "ja-jp"},
			ImageNames:  []string{"Windows 11 Pro", "Windows 11 Enterprise", "Windows 11 Home"},
			DownloadURL: "https://www.microsoft.com/software-download/windows11",
			Supported:   true,
		},
		{
			Name:        "Windows Server 2019",
			Version:     "2019",
			Languages:   []string{"zh-cn", "en-us"},
			ImageNames:  []string{"Windows Server 2019 Standard", "Windows Server 2019 Datacenter"},
			DownloadURL: "https://www.microsoft.com/evalcenter/download-windows-server-2019",
			Supported:   true,
		},
		{
			Name:        "Windows Server 2022",
			Version:     "2022",
			Languages:   []string{"zh-cn", "en-us"},
			ImageNames:  []string{"Windows Server 2022 Standard", "Windows Server 2022 Datacenter"},
			DownloadURL: "https://www.microsoft.com/evalcenter/download-windows-server-2022",
			Supported:   true,
		},
	}
}

// GetSupportedDDImages 获取支持的DD镜像列表
func (si *SystemInstaller) GetSupportedDDImages() []DDImageInfo {
	return []DDImageInfo{
		{
			Name:        "Windows 10 LTSC",
			Description: "Windows 10 LTSC 精简版",
			URL:         "https://example.com/win10-ltsc.vhd.gz",
			Size:        "4.2GB",
			Format:      "vhd",
			Compression: "gz",
			Supported:   true,
		},
		{
			Name:        "Ubuntu 22.04 Server",
			Description: "Ubuntu 22.04 服务器版",
			URL:         "https://example.com/ubuntu-22.04-server.raw.xz",
			Size:        "2.1GB",
			Format:      "raw",
			Compression: "xz",
			Supported:   true,
		},
		{
			Name:        "Debian 12 Minimal",
			Description: "Debian 12 最小化安装",
			URL:         "https://example.com/debian-12-minimal.qcow2.zst",
			Size:        "1.8GB",
			Format:      "qcow2",
			Compression: "zst",
			Supported:   true,
		},
	}
}

// ValidateInstallOptions 验证安装选项
func (si *SystemInstaller) ValidateInstallOptions(options InstallOptions) error {
	if options.OSType == "" {
		return fmt.Errorf("必须指定操作系统类型")
	}

	switch options.OSType {
	case "linux":
		if options.System == "" {
			return fmt.Errorf("必须指定Linux系统名称")
		}
	case "windows":
		if options.ImageName == "" && options.ISOURL == "" {
			return fmt.Errorf("必须指定Windows镜像名称或ISO URL")
		}
	case "dd":
		if options.ImageURL == "" {
			return fmt.Errorf("必须指定DD镜像URL")
		}
	default:
		return fmt.Errorf("不支持的操作系统类型: %s", options.OSType)
	}

	return nil
}

// loadConfig 加载配置
func (si *SystemInstaller) loadConfig() {
	configFile := filepath.Join(si.workingDir, "config.json")
	if data, err := os.ReadFile(configFile); err == nil {
		json.Unmarshal(data, &si.config)
	}
}

// SaveConfig 保存配置
func (si *SystemInstaller) SaveConfig() error {
	configFile := filepath.Join(si.workingDir, "config.json")
	data, err := json.MarshalIndent(si.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0644)
}
