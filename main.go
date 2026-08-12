package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "Gapiro",
		Description: "High-performance API testing client",
		Services: []application.Service{
			application.NewService(&HttpService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Main window - optimized for Wayland/Hyprland
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Gapiro - API Client",
		Width:            1440,
		Height:           900,
		MinWidth:         900,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(17, 17, 21),
		URL:              "/",
		Linux: application.LinuxWindow{
			// Enable GPU acceleration for WebKitGTK on Wayland
			WebviewGpuPolicy: application.WebviewGpuPolicyAlways,
		},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
