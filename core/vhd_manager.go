package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// VHDManager VHD管理器
type VHDManager struct {
	vhdDir        string
	downloadDir   string
	workingDir    string
	progress      map[string]*DownloadProgress
	progressMutex sync.RWMutex
}

// VHDInfo VHD信息
type VHDInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Size        string `json:"size"`
	LocalPath   string `json:"local_path"`
	Downloaded  bool   `json:"downloaded"`
	Version     string `json:"version"`
	OS          string `json:"os"`
	Arch        string `json:"arch"`
}

// DownloadProgress 下载进度
type DownloadProgress struct {
	Percentage int    `json:"percentage"`
	Speed      string `json:"speed"`
	ETA        string `json:"eta"`
	Status     string `json:"status"` // downloading, completed, error, paused
	Message    string `json:"message"`
}

// NewVHDManager 创建VHD管理器
func NewVHDManager() *VHDManager {
	// 获取当前工作目录
	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = "." // 如果获取失败，使用当前目录
	}

	return &VHDManager{
		workingDir:  workingDir,
		vhdDir:      filepath.Join(workingDir, "vhd"),
		downloadDir: filepath.Join(workingDir, "downloads"),
		progress:    make(map[string]*DownloadProgress),
	}
}

// Initialize 初始化VHD管理器
func (vm *VHDManager) Initialize() error {
	// 创建目录
	os.MkdirAll(vm.vhdDir, 0755)
	os.MkdirAll(vm.downloadDir, 0755)
	return nil
}

// GetAvailableVHDs 获取可用的VHD列表
func (vm *VHDManager) GetAvailableVHDs() []VHDInfo {
	return []VHDInfo{
		{
			Name:        "Windows 10 LTSC 2021",
			Description: "Windows 10 LTSC 2021 精简版",
			URL:         "https://example.com/win10-ltsc-2021.vhd",
			Size:        "4.2GB",
			Version:     "21H2",
			OS:          "Windows",
			Arch:        "x64",
		},
		{
			Name:        "Windows 11 Pro",
			Description: "Windows 11 Pro 完整版",
			URL:         "https://example.com/win11-pro.vhd",
			Size:        "5.8GB",
			Version:     "23H2",
			OS:          "Windows",
			Arch:        "x64",
		},
		{
			Name:        "Ubuntu 22.04 Desktop",
			Description: "Ubuntu 22.04 桌面版",
			URL:         "https://example.com/ubuntu-22.04-desktop.vhd",
			Size:        "3.1GB",
			Version:     "22.04",
			OS:          "Linux",
			Arch:        "x64",
		},
		{
			Name:        "CentOS 9 Stream",
			Description: "CentOS 9 Stream 服务器版",
			URL:         "https://example.com/centos-9-stream.vhd",
			Size:        "2.4GB",
			Version:     "9",
			OS:          "Linux",
			Arch:        "x64",
		},
	}
}

// GetLocalVHDs 获取本地VHD列表
func (vm *VHDManager) GetLocalVHDs() ([]VHDInfo, error) {
	var vhds []VHDInfo

	// 确保目录存在
	if _, err := os.Stat(vm.vhdDir); os.IsNotExist(err) {
		os.MkdirAll(vm.vhdDir, 0755)
		return vhds, nil
	}

	files, err := os.ReadDir(vm.vhdDir)
	if err != nil {
		return vhds, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".vhd") {
			info, err := file.Info()
			if err != nil {
				continue
			}

			vhd := VHDInfo{
				Name:       strings.TrimSuffix(file.Name(), ".vhd"),
				LocalPath:  filepath.Join(vm.vhdDir, file.Name()),
				Downloaded: true,
				Size:       fmt.Sprintf("%.2f GB", float64(info.Size())/(1024*1024*1024)),
			}
			vhds = append(vhds, vhd)
		}
	}

	return vhds, nil
}

// DownloadVHD 下载VHD文件
func (vm *VHDManager) DownloadVHD(vhd VHDInfo) error {
	filename := fmt.Sprintf("%s.vhd", vhd.Name)
	localPath := filepath.Join(vm.vhdDir, filename)

	// 检查文件是否已存在
	if _, err := os.Stat(localPath); err == nil {
		return fmt.Errorf("文件已存在: %s", filename)
	}

	// 初始化进度
	vm.progressMutex.Lock()
	vm.progress[vhd.Name] = &DownloadProgress{
		Percentage: 0,
		Status:     "downloading",
		Message:    "开始下载...",
	}
	vm.progressMutex.Unlock()

	// 创建HTTP请求
	resp, err := http.Get(vhd.URL)
	if err != nil {
		vm.updateProgress(vhd.Name, 0, "error", fmt.Sprintf("下载失败: %v", err))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("HTTP错误: %d", resp.StatusCode)
		vm.updateProgress(vhd.Name, 0, "error", err.Error())
		return err
	}

	// 创建本地文件
	file, err := os.Create(localPath)
	if err != nil {
		vm.updateProgress(vhd.Name, 0, "error", fmt.Sprintf("创建文件失败: %v", err))
		return err
	}
	defer file.Close()

	// 下载文件
	contentLength := resp.ContentLength
	var downloaded int64

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			file.Write(buffer[:n])
			downloaded += int64(n)

			// 更新进度
			if contentLength > 0 {
				percentage := int((downloaded * 100) / contentLength)
				vm.updateProgress(vhd.Name, percentage, "downloading", "下载中...")
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			vm.updateProgress(vhd.Name, 0, "error", fmt.Sprintf("下载失败: %v", err))
			return err
		}
	}

	vm.updateProgress(vhd.Name, 100, "completed", "下载完成")
	return nil
}

// InstallVHD 安装VHD系统
func (vm *VHDManager) InstallVHD(vhdPath string, options InstallOptions) error {
	// 检查VHD文件是否存在
	if _, err := os.Stat(vhdPath); os.IsNotExist(err) {
		return fmt.Errorf("VHD文件不存在: %s", vhdPath)
	}

	// 使用DD方式安装VHD
	ddOptions := InstallOptions{
		OSType:   "dd",
		ImageURL: "file://" + vhdPath,
		Password: options.Password,
		SSHKey:   options.SSHKey,
		SSHPort:  options.SSHPort,
	}

	// 创建系统安装器
	installer := NewSystemInstaller()
	if err := installer.Initialize(); err != nil {
		return err
	}

	return installer.InstallSystem(ddOptions)
}

// DeleteVHD 删除VHD文件
func (vm *VHDManager) DeleteVHD(vhdName string) error {
	filename := fmt.Sprintf("%s.vhd", vhdName)
	localPath := filepath.Join(vm.vhdDir, filename)

	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		return fmt.Errorf("VHD文件不存在: %s", filename)
	}

	return os.Remove(localPath)
}

// GetDownloadProgress 获取下载进度
func (vm *VHDManager) GetDownloadProgress(vhdName string) *DownloadProgress {
	vm.progressMutex.RLock()
	defer vm.progressMutex.RUnlock()

	if progress, exists := vm.progress[vhdName]; exists {
		return progress
	}
	return nil
}

// updateProgress 更新下载进度
func (vm *VHDManager) updateProgress(vhdName string, percentage int, status, message string) {
	vm.progressMutex.Lock()
	defer vm.progressMutex.Unlock()

	if progress, exists := vm.progress[vhdName]; exists {
		progress.Percentage = percentage
		progress.Status = status
		progress.Message = message
	}
}

// MountVHD 挂载VHD文件 (Windows)
func (vm *VHDManager) MountVHD(vhdPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("VHD挂载仅支持Windows系统")
	}

	// 使用diskpart挂载VHD
	cmd := exec.Command("diskpart")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("select vdisk file=%s\nattach vdisk\n", vhdPath))
	return cmd.Run()
}

// UnmountVHD 卸载VHD文件 (Windows)
func (vm *VHDManager) UnmountVHD(vhdPath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("VHD卸载仅支持Windows系统")
	}

	// 使用diskpart卸载VHD
	cmd := exec.Command("diskpart")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("select vdisk file=%s\ndetach vdisk\n", vhdPath))
	return cmd.Run()
}
