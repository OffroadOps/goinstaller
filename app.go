package main

import (
	"context"
	"fmt"
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
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// 需要添加这些方法到App结构体
func (a *App) GetAvailableServers() []interface{} {
    // 实现获取服务器列表的逻辑
    return []interface{}{}
}

func (a *App) GetVHDListFromServer(serverId int) []interface{} {
    // 实现获取VHD列表的逻辑
    return []interface{}{}
}

func (a *App) DownloadVHD(vhdId int, savePath string) map[string]interface{} {
    // 实现VHD下载的逻辑
    return map[string]interface{}{
        "success": true,
        "message": "下载开始",
    }
}

// 其他需要的函数...
