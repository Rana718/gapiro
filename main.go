package main

import (
	"changeme/db"
	"changeme/services"
	"embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	ensureWaylandScaling()
}

func ensureWaylandScaling() {
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		setIfEmpty("GDK_BACKEND", "wayland")
	}
	setIfEmpty("WEBKIT_DISABLE_DMABUF_RENDERER", "0")

	if os.Getenv("GDK_SCALE") == "" {
		if scale := detectHyprlandScale(); scale > 0 && scale != 1.0 {
			os.Setenv("GDK_DPI_SCALE", fmt.Sprintf("%.2f", scale))
		}
	}
}

func detectHyprlandScale() float64 {
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		return 0
	}
	out, err := exec.Command("hyprctl", "monitors", "-j").Output()
	if err != nil {
		return 0
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "\"scale\":") {
			valStr := strings.TrimPrefix(line, "\"scale\":")
			valStr = strings.TrimSuffix(strings.TrimSpace(valStr), ",")
			if val, err := strconv.ParseFloat(valStr, 64); err == nil && val > 0 {
				return val
			}
		}
	}
	return 0
}

func detectSystemScale() float64 {
	if s := os.Getenv("GDK_SCALE"); s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil && val > 0 {
			return val
		}
	}
	if s := os.Getenv("QT_SCALE_FACTOR"); s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil && val > 0 {
			return val
		}
	}
	if scale := detectHyprlandScale(); scale > 0 {
		return scale
	}
	return 1.0
}

func setIfEmpty(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

func main() {
	// Initialize SQLite database
	if err := db.Initialize(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	systemScale := detectSystemScale()

	app := application.New(application.Options{
		Name:        "Gapiro",
		Description: "High-performance API testing client",
		Services: []application.Service{
			application.NewService(&services.HttpService{}),
			application.NewService(&services.CollectionService{}),
			application.NewService(&services.GrpcService{}),
			application.NewService(&services.GraphQLService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:              "Gapiro - API Client",
		Width:              1440,
		Height:             900,
		MinWidth:           900,
		MinHeight:          600,
		BackgroundColour:   application.NewRGB(17, 17, 21),
		URL:                "/",
		Zoom:               systemScale,
		ZoomControlEnabled: true,
		Linux: application.LinuxWindow{
			WebviewGpuPolicy: application.WebviewGpuPolicyAlways,
		},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
