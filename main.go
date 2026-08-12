package main

import (
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
	// Ensure proper scaling on Wayland compositors (Hyprland, Sway, etc.)
	// WebKitGTK needs these hints to respect the compositor's scale factor.
	ensureWaylandScaling()
}

// ensureWaylandScaling sets environment variables so that GTK4/WebKitGTK
// picks up the display scale from the Wayland compositor.
func ensureWaylandScaling() {
	// Force GDK to use Wayland backend if available
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		setIfEmpty("GDK_BACKEND", "wayland")
	}

	// Enable native Wayland support in WebKit
	setIfEmpty("WEBKIT_DISABLE_DMABUF_RENDERER", "0")

	// If the user has set a specific scale in their compositor (e.g. Hyprland monitor config),
	// GTK4 should pick it up automatically via gdk_monitor_get_scale().
	// However, if GDK_SCALE is explicitly set, honor it for the webview zoom.
	// For fractional scaling on older GTK, we need GDK_DPI_SCALE.
	if os.Getenv("GDK_SCALE") == "" {
		// Try to detect the scale from Hyprland via hyprctl
		if scale := detectHyprlandScale(); scale > 0 && scale != 1.0 {
			// For fractional scales, GTK4's gdk_monitor_get_scale handles it.
			// But we also set GDK_DPI_SCALE as a fallback for font rendering.
			os.Setenv("GDK_DPI_SCALE", fmt.Sprintf("%.2f", scale))
		}
	}
}

// detectHyprlandScale tries to get the monitor scale from Hyprland.
func detectHyprlandScale() float64 {
	// Check if we're running under Hyprland
	if os.Getenv("HYPRLAND_INSTANCE_SIGNATURE") == "" {
		return 0
	}

	out, err := exec.Command("hyprctl", "monitors", "-j").Output()
	if err != nil {
		return 0
	}

	// Simple parse - find "scale": X.XX in the JSON output
	// We look for the focused monitor's scale
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

func setIfEmpty(key, value string) {
	if os.Getenv(key) == "" {
		os.Setenv(key, value)
	}
}

func main() {
	// Detect system scale for window zoom level
	systemScale := detectSystemScale()

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

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Gapiro - API Client",
		Width:            1440,
		Height:           900,
		MinWidth:         900,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(17, 17, 21),
		URL:              "/",
		// Set zoom to match system scale so content isn't tiny on HiDPI
		Zoom:               systemScale,
		ZoomControlEnabled: true,
		Linux: application.LinuxWindow{
			// Always use GPU for smooth rendering on Wayland
			WebviewGpuPolicy: application.WebviewGpuPolicyAlways,
		},
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// detectSystemScale returns the display scale factor.
// On Wayland/Hyprland it reads from hyprctl or GDK_SCALE env var.
func detectSystemScale() float64 {
	// 1. Check explicit GDK_SCALE
	if s := os.Getenv("GDK_SCALE"); s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil && val > 0 {
			return val
		}
	}

	// 2. Check QT_SCALE_FACTOR (some users set this globally)
	if s := os.Getenv("QT_SCALE_FACTOR"); s != "" {
		if val, err := strconv.ParseFloat(s, 64); err == nil && val > 0 {
			return val
		}
	}

	// 3. Try Hyprland
	if scale := detectHyprlandScale(); scale > 0 {
		return scale
	}

	// 4. Try wlr-randr or swaymsg for other Wayland compositors
	if scale := detectWlrScale(); scale > 0 {
		return scale
	}

	// Default: no zoom adjustment (GTK should handle it natively)
	return 1.0
}

// detectWlrScale tries wlr-randr for wlroots-based compositors (Sway, etc.)
func detectWlrScale() float64 {
	out, err := exec.Command("wlr-randr").Output()
	if err != nil {
		return 0
	}

	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "Scale:") {
			valStr := strings.TrimSpace(strings.TrimPrefix(line, "Scale:"))
			if val, err := strconv.ParseFloat(valStr, 64); err == nil && val > 0 {
				return val
			}
		}
	}
	return 0
}
