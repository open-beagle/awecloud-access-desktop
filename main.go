package main

import (
	"embed"
	"fmt"
	"io/fs"
	"path"
	"runtime"
	"time"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/application"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/open-beagle/awecloud-access-desktop/pkg/util"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed all:build
var build embed.FS

var mainApp *application.Application
var app *App

const (
	menuSettingsText = "设置"
	menuQuitText     = "退出"
)

func main() {
	// 启动 systray
	go func() {
		runtime.LockOSThread()         // 锁定当前 goroutine 到操作系统线程
		defer runtime.UnlockOSThread() // 确保在退出时解锁
		systray.Run(onReady, onExit)
	}()

	// 启动 Wails App
	app = NewApp()
	mainApp = application.NewWithOptions(&options.App{
		Title:  "Beagle Access Desktop - v" + util.VERSION,
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		HideWindowOnClose: true,
	})
	err := mainApp.Run()
	if err != nil {
		println("Error:", err.Error())
	}
}

// onReady initializes the systray and its menu items
func onReady() {
	systray.SetIcon(getIcon())
	systray.SetTitle("比格访问")
	systray.SetTooltip("比格访问")

	mSettings := systray.AddMenuItem(menuSettingsText, "打开设置")
	mQuit := systray.AddMenuItem(menuQuitText, "退出应用")

	go handleMenuActions(mSettings, mQuit)
}

// handleMenuActions processes the menu item click events
func handleMenuActions(mSettings, mQuit *systray.MenuItem) {
	for {
		select {
		case <-mSettings.ClickedCh:
			handleSettingsClick()
		case <-mQuit.ClickedCh:
			handleQuit()
			return
		}
	}
}

// handleSettingsClick checks if the app is hidden and shows it if necessary
func handleSettingsClick() {
	app.Show()
}

// handleQuit handles the quit action
func handleQuit() {
	systray.Quit()
}

func onExit() {
	// 清理操作
	now := time.Now()
	fmt.Println("onExit() at", now)
	mainApp.Quit()
}

func getIcon() []byte {
	b, err := fs.ReadFile(build, path.Join("build", "windows", "icon.ico"))
	if err != nil {
		fmt.Print(err)
	}
	return b
}
