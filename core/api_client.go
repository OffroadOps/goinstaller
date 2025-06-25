package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// APIClient API客户端
type APIClient struct {
	baseURL    string
	apiKey     string
	clientID   string
	httpClient *http.Client
}

// ServerInfo 服务器信息
type ServerInfo struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Location     string                 `json:"location"`
	DownloadURLs map[string]interface{} `json:"download_urls"`
}

// ServerListResponse 服务器列表响应
type ServerListResponse struct {
	Success bool         `json:"success"`
	Data    ServerData   `json:"data"`
	Error   string       `json:"error,omitempty"`
}

// ServerData 服务器数据
type ServerData struct {
	Servers []ServerInfo `json:"servers"`
}

// VHDListResponse VHD列表响应
type VHDListResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   string                 `json:"error,omitempty"`
}

// NewAPIClient 创建API客户端
func NewAPIClient() *APIClient {
	return &APIClient{
		baseURL:  "https://autoinstaller.qkdny.com/api/autoinstaller",
		clientID: "go-client",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SetAPIKey 设置API密钥
func (ac *APIClient) SetAPIKey(apiKey string) {
	ac.apiKey = apiKey
}

// GetServerList 获取服务器列表
func (ac *APIClient) GetServerList() map[string]interface{} {
	url := fmt.Sprintf("%s/iso-download/", ac.baseURL)
	
	resp, err := ac.httpClient.Get(url)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	
	var response ServerListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	
	if !response.Success {
		return map[string]interface{}{
			"success": false,
			"error":   response.Error,
		}
	}
	
	// 转换为标准格式
	formattedServers := make([]map[string]interface{}, 0)
	for _, server := range response.Data.Servers {
		formattedServers = append(formattedServers, map[string]interface{}{
			"id":            server.ID,
			"name":          server.Name,
			"location":      server.Location,
			"download_urls": server.DownloadURLs,
		})
	}
	
	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"servers": formattedServers,
		},
	}
}

// GetVHDList 获取VHD列表
func (ac *APIClient) GetVHDList(serverID string) map[string]interface{} {
	// 首先获取服务器列表
	serverResult := ac.GetServerList()
	if !serverResult["success"].(bool) {
		return serverResult
	}
	
	data := serverResult["data"].(map[string]interface{})
	servers := data["servers"].([]map[string]interface{})
	
	if len(servers) == 0 {
		return map[string]interface{}{
			"success": false,
			"error":   "没有可用的服务器",
		}
	}
	
	// 查找指定服务器
	var selectedServer map[string]interface{}
	if serverID != "" {
		for _, server := range servers {
			if server["id"] == serverID {
				selectedServer = server
				break
			}
		}
		
		if selectedServer == nil {
			return map[string]interface{}{
				"success": false,
				"error":   fmt.Sprintf("未找到服务器: %s", serverID),
			}
		}
	} else {
		// 如果没有指定服务器ID，使用第一个服务器
		selectedServer = servers[0]
	}
	
	// 返回VHD列表
	downloadURLs := selectedServer["download_urls"].(map[string]interface{})
	vhdList := make([]map[string]interface{}, 0)
	
	for name, url := range downloadURLs {
		vhdList = append(vhdList, map[string]interface{}{
			"name": name,
			"url":  url,
			"size": "未知", // 实际应用中可以获取文件大小
		})
	}
	
	return map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"server": selectedServer,
			"vhds":   vhdList,
		},
	}
}

// DownloadVHD 下载VHD文件
func (ac *APIClient) DownloadVHD(serverID, vhdName, savePath string) bool {
	// 获取VHD列表
	vhdResult := ac.GetVHDList(serverID)
	if !vhdResult["success"].(bool) {
		return false
	}
	
	data := vhdResult["data"].(map[string]interface{})
	vhds := data["vhds"].([]map[string]interface{})
	
	// 查找指定VHD
	var downloadURL string
	for _, vhd := range vhds {
		if vhd["name"] == vhdName {
			downloadURL = vhd["url"].(string)
			break
		}
	}
	
	if downloadURL == "" {
		return false
	}
	
	// 下载文件（这里是简化版本，实际应用中需要进度回调等）
	resp, err := ac.httpClient.Get(downloadURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	// 这里应该实现实际的文件下载逻辑
	// 包括进度回调、断点续传等功能
	
	return true
}

// UploadLog 上传日志
func (ac *APIClient) UploadLog(logData map[string]interface{}) map[string]interface{} {
	url := fmt.Sprintf("%s/upload-log/", ac.baseURL)
	
	jsonData, err := json.Marshal(logData)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	
	resp, err := ac.httpClient.Post(url, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}
	
	return result
}