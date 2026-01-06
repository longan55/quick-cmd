package main

import (
	"context"
	"fmt"
	"log"
)

const (
	Windows string = "windows"
	Mac     string = "mac"
	Linux   string = "linux"
	NoneOs  string = ""
	AllOs   string = "all"
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
	Os   []string   `json:"os"`
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

func (a *App) GetOptions(option Option) (response Response) {
	fmt.Printf("GetOptions: %+v\n", option)
	switch option.Type {
	case "commands", "all":
		log.Printf("GetCommandsOptions")
		response.Data = getCommandsOptions(option)
	case "tags":
		log.Printf("GetTagsOptions")
		response.Data = getTagsOptions(option)
	case "collections":
		log.Printf("GetCollectionsOptions")
		response.Data = getCollectionsOptions(option)
	}
	fmt.Printf("response: %+v\n", response)
	fmt.Printf("response.data.tags: %+v\n", response.Data.(AllCommands).Tags)
	fmt.Printf("response.data.collections: %+v\n", response.Data.(AllCommands).Collections)
	fmt.Printf("response.data.commands: %+v\n", response.Data.(AllCommands).Commands)
	return response
}

func getTagsOptions(option Option) AllCommands {
	// 这里应该根据option参数查询数据库获取标签列表
	// 目前返回空列表，后续需要实现具体逻辑

	tags, err := GetTagsSQLite(option)
	if err != nil {
		log.Printf("GetTagsSQLite failed: %v", err)
		return AllCommands{
			Tags:        []*Tag{},
			Collections: []*Collection{},
			Commands:    []*Command{},
		}
	}
	// commands := make([]*Command, 0, 64)
	tagIDs := make([]uint64, 0, len(tags))
	for _, tag := range tags {
		tagIDs = append(tagIDs, tag.ID)
	}
	commands, err := GetCommandsByTagIDs(tagIDs)
	if err != nil {
		log.Printf("GetCommandsByTagIDs failed: %v", err)
		return AllCommands{
			Tags:        tags,
			Collections: []*Collection{},
			Commands:    commands,
		}
	}
	return AllCommands{
		Tags:        tags,
		Collections: []*Collection{},
		Commands:    commands,
	}
}

func getCollectionsOptions(option Option) AllCommands {
	// 这里应该根据option参数查询数据库获取集合列表
	// 目前返回空列表，后续需要实现具体逻辑
	collections, err := GetCollectionsSQLite(option)
	if err != nil {
		log.Printf("GetCollectionsSQLite failed: %v", err)
		return AllCommands{
			Tags:        []*Tag{},
			Collections: []*Collection{},
			Commands:    []*Command{},
		}
	}
	collectionIDs := make([]uint64, 0, len(collections))
	for _, collection := range collections {
		collectionIDs = append(collectionIDs, collection.ID)
	}
	commands, err := GetCommandByCollectionIds(collectionIDs)
	if err != nil {
		log.Printf("GetCommandByCollectionIds failed: %v", err)
		return AllCommands{
			Tags:        []*Tag{},
			Collections: collections,
			Commands:    []*Command{},
		}
	}
	return AllCommands{
		Tags:        []*Tag{},
		Collections: collections,
		Commands:    commands,
	}
}

func getCommandsOptions(option Option) AllCommands {
	// 这里应该根据option参数查询数据库获取命令列表
	// 目前返回空列表，后续需要实现具体逻辑
	commands, err := GetCommandsSQLite(option)
	if err != nil {
		log.Printf("GetCommandsSQLite failed: %v", err)
		return AllCommands{
			Tags:        []*Tag{},
			Collections: []*Collection{},
			Commands:    []*Command{},
		}
	}
	log.Printf("GetCommandsSQLite success: %v", commands)
	return AllCommands{
		Tags:        []*Tag{},
		Collections: []*Collection{},
		Commands:    commands,
	}
}

type AllCommands struct {
	Tags        []*Tag        `json:"tags"`
	Collections []*Collection `json:"collections"`
	Commands    []*Command    `json:"options"`
}

func (a *App) GetAllTagsIDAndName() ([]Tag, error) {
	return GetTagIDAndNameSQLite()
}

func (a *App) GetAllCollectionsIDAndName() ([]Collection, error) {
	return GetCollectionIDAndNameSQLite()
}
