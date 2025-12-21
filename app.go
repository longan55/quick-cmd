package main

import (
	"context"
	"slices"
	"strings"
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

func (a *App) GetTagsOptions(option Option) []*Tag {
	if len(option.Os) == 0 {
		return nil
	}
	osTags := make([]*Tag, 0, len(allTags))
	if len(option.Os) == 3 {
		osTags = allTags
	} else {
		for _, os := range option.Os {
			for _, tag := range allTags {
				if tag.Os == os {
					osTags = append(osTags, tag)
				}
			}
		}
	}
	//过滤名称
	if option.Name == "" {
		return osTags
	}
	nameTags := make([]*Tag, 0, len(osTags))
	for _, tag := range osTags {
		if strings.Contains(tag.Name, option.Name) {
			nameTags = append(nameTags, tag)
		}
	}
	//排序
	if option.Sort.Name != nil {
		slices.SortFunc(nameTags, func(a, b *Tag) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	//排序
	if option.Sort.CreateTime != nil {
		slices.SortFunc(nameTags, func(a, b *Tag) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	}
	//排序
	if option.Sort.CopyCounts != nil {
		slices.SortFunc(nameTags, func(a, b *Tag) int {
			return a.SearchCount - b.SearchCount
		})
	}
	return nameTags
}

func (a *App) GetCollectionsOptions(option Option) []*Collection {
	if len(option.Os) == 0 {
		return nil
	}
	osCollections := make([]*Collection, 0, len(allCollections))
	if len(option.Os) == 3 {
		osCollections = allCollections
	} else {
		for _, os := range option.Os {
			for _, col := range allCollections {
				if col.Os == os {
					osCollections = append(osCollections, col)
				}
			}
		}
	}
	//过滤名称
	if option.Name == "" {
		return osCollections
	}
	nameCollections := make([]*Collection, 0, len(osCollections))
	for _, col := range osCollections {
		if strings.Contains(col.Name, option.Name) {
			nameCollections = append(nameCollections, col)
		}
	}
	//排序
	if option.Sort.Name != nil {
		slices.SortFunc(nameCollections, func(a, b *Collection) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	//排序
	if option.Sort.CreateTime != nil {
		slices.SortFunc(nameCollections, func(a, b *Collection) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	}
	//排序
	if option.Sort.CopyCounts != nil {
		slices.SortFunc(nameCollections, func(a, b *Collection) int {
			return a.SearchCount - b.SearchCount
		})
	}
	return nameCollections
}

func (a *App) GetCommandsOptions(option Option) []*Command {
	if len(option.Os) == 0 {
		return nil
	}
	//过滤系统
	ostempCommands := make([]*Command, 0, len(allCommands))
	if len(option.Os) == 3 {
		ostempCommands = allCommands
	} else {
		for _, os := range option.Os {
			for _, cmd := range allCommands {
				if cmd.Os == os {
					ostempCommands = append(ostempCommands, cmd)
				}
			}
		}
	}
	//过滤名称
	if option.Name == "" {
		return ostempCommands
	}
	nameCommands := make([]*Command, 0, len(ostempCommands))
	for _, cmd := range ostempCommands {
		if strings.Contains(cmd.Name, option.Name) {
			nameCommands = append(nameCommands, cmd)
		}
	}

	//排序
	if option.Sort.Name != nil {
		slices.SortFunc(nameCommands, func(a, b *Command) int {
			return strings.Compare(a.Name, b.Name)
		})
	}
	//排序
	if option.Sort.CreateTime != nil {
		slices.SortFunc(nameCommands, func(a, b *Command) int {
			return a.CreatedAt.Compare(b.CreatedAt)
		})
	}
	//排序
	if option.Sort.CopyCounts != nil {
		slices.SortFunc(nameCommands, func(a, b *Command) int {
			return a.CopyCounts - b.CopyCounts
		})
	}
	// if option.Sort.SortValue != nil {
	// 	slices.SortFunc(tempCommands, func(a, b *Command) int {
	// 		return strings.Compare(a.SortValue, b.SortValue)
	// 	})
	// }
	return nameCommands
}

type Option struct {
	Name string     `json:"name"`
	Os   []OS       `json:"os"`
	Sort SortOption `json:"sort"`
}

type SortOption struct {
	Name       *string `json:"name"`
	CreateTime *string `json:"create_time"`
	CopyCounts *string `json:"copy_counts"`
	// SortValue  *string `json:"sort_value"`
}
