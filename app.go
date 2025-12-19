package main

import (
	"context"
)

type OS string

const (
	Windows OS = "windows"
	Mac     OS = "mac"
	Linux   OS = "linux"
)

// App struct
type App struct {
	ctx         context.Context
	commands    []*Command    // 模拟指令数据库
	tags        []*Tag        // 模拟标签数据库
	collections []*Collection // 模拟集合数据库
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		commands:    make([]*Command, 0, 64),
		tags:        make([]*Tag, 0, 16),
		collections: make([]*Collection, 0, 16),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// GetMenuItems returns the menu items for the application
func (a *App) GetMenuItems() map[string]interface{} {
	return nil
}
