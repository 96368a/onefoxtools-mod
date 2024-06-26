package main

import (
	"changeme/common"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/exp/slog"
	"time"
)

var wailsContext *context.Context

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	wailsContext = &ctx
	common.InitEnv()
	if common.Paths.Title != "" {
		runtime.WindowSetTitle(ctx, common.Paths.Title)
		slog.Info("已从配置中读取标题")
	} else {
		slog.Error("未设置标题")
	}
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
	cost := time.Since(startTime)
	slog.Info("初始化耗时:", cost.Seconds())
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

func (a *App) setTitle(name string) {
	runtime.WindowSetTitle(a.ctx, name)
}
