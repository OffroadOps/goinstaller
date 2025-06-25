package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "系统重装助手 v2.0",
		Width:  1200,
		Height: 800,
		MinWidth: 1000,
		MinHeight: 700,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 250, G: 250, B: 250, A: 1}, // 改为浅灰色
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		// 添加这些选项让它更像原生应用
		DisableResize: false,
		Fullscreen:   false,
		Frameless:    false, // 保持系统标题栏
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
