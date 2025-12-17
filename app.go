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

// GetMenuItems returns the menu items for the application
func (a *App) GetMenuItems() map[string]interface{} {
	// 返回固定的菜单项数据
	return map[string]interface{}{
		"tags": []map[string]string{
			{"id": "home", "name": "首页", "icon": "🏷️"},
			{"id": "tag1", "name": "工作", "icon": "🏷️"},
			{"id": "tag2", "name": "学习", "icon": "🏷️"},
			{"id": "tag3", "name": "生活", "icon": "🏷️"},
			{"id": "tag4", "name": "娱乐", "icon": "🏷️"},
		},
		"collections": []map[string]string{
			{"id": "home", "name": "首页", "icon": "{ }"},
			{"id": "collection1", "name": "常用工具", "icon": "{ }"},
			{"id": "collection2", "name": "开发资源", "icon": "{ }"},
			{"id": "collection3", "name": "文档资料", "icon": "{ }"},
			{"id": "collection4", "name": "项目管理", "icon": "{ }"},
		},
		"all": []map[string]string{
			{"id": "home", "name": "首页", "icon": "🎇"},
			{"id": "tag1", "name": "工作", "icon": "🎇"},
			{"id": "tag2", "name": "学习", "icon": "🎇"},
			{"id": "tag3", "name": "生活", "icon": "🎇"},
			{"id": "tag4", "name": "娱乐", "icon": "🎇"},
			{"id": "collection1", "name": "常用工具", "icon": "🎇"},
			{"id": "collection2", "name": "开发资源", "icon": "🎇"},
			{"id": "collection3", "name": "文档资料", "icon": "🎇"},
			{"id": "collection4", "name": "项目管理", "icon": "🎇"},
			{"id": "settings", "name": "设置", "icon": "🎇"},
		},
	}
}
