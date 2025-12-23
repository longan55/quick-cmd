package main

import (
	"context"
	"fmt"
)

type OS string

const (
	Windows OS = "windows"
	Mac     OS = "mac"
	Linux   OS = "linux"
	NoneOs  OS = ""
	AllOs   OS = "all"
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

type Option struct {
	Name string     `json:"name"`
	Os   []OS       `json:"os"`
	Type string     `json:"type"` //command,tag,collection
	ID   uint64     `json:"id"`
	Sort SortOption `json:"sort"`
}

type SortOption struct {
	Name       *string `json:"name"`
	CreateTime *string `json:"create_time"`
	CopyCounts *string `json:"copy_counts"`
	// SortValue  *string `json:"sort_value"`
}

func (a *App) GetOptions(option Option) Response {
	fmt.Printf("GetOptions: %+v\n", option)
	switch option.Type {
	case "command":
		return getCommandsOptions(option)
	case "tag":
		return getTagsOptions(option)
	case "collection":
		return getCollectionsOptions(option)
	}
	return Response{}
}

func getTagsOptions(option Option) Response {
	// 这里应该根据option参数查询数据库获取标签列表
	// 目前返回空列表，后续需要实现具体逻辑
	return Response{
		Tags:        []*Tag{},
		Collections: []*Collection{},
		Commands:    []*Command{},
	}
}

func getCollectionsOptions(option Option) Response {
	// 这里应该根据option参数查询数据库获取集合列表
	// 目前返回空列表，后续需要实现具体逻辑
	return Response{
		Tags:        []*Tag{},
		Collections: []*Collection{},
		Commands:    []*Command{},
	}
}

func getCommandsOptions(option Option) Response {
	// 这里应该根据option参数查询数据库获取命令列表
	// 目前返回空列表，后续需要实现具体逻辑
	return Response{
		Tags:        []*Tag{},
		Collections: []*Collection{},
		Commands:    []*Command{},
	}
}

type Response struct {
	Tags        []*Tag        `json:"tags"`
	Collections []*Collection `json:"collections"`
	Commands    []*Command    `json:"options"`
}
