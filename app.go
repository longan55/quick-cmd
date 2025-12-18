package main

import (
	"context"
	"fmt"
	"time"
)

// Command 指令结构体
type Command struct {
	ID            string     `json:"id"`
	Name          string     `json:"name"`
	Content       string     `json:"content"`
	Description   string     `json:"description"`
	CopyCount     int        `json:"copyCount"`
	SearchCount   int        `json:"searchCount"`
	TagIDs        []string   `json:"tagIDs"`        // 标签ID列表
	CollectionIDs []string   `json:"collectionIDs"` // 集合ID列表
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

// Tag 标签结构体
type Tag struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SearchCount int        `json:"searchCount"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

// Collection 集合结构体
type Collection struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SearchCount int        `json:"searchCount"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

// App struct
type App struct {
	ctx         context.Context
	commands    map[string]*Command    // 模拟指令数据库
	tags        map[string]*Tag        // 模拟标签数据库
	collections map[string]*Collection // 模拟集合数据库
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		commands:    make(map[string]*Command),
		tags:        make(map[string]*Tag),
		collections: make(map[string]*Collection),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化标签数据
	tags := []struct {
		ID   string
		Name string
	}{
		{"tag1", "工作"},
		{"tag2", "学习"},
		{"tag3", "生活"},
		{"tag4", "娱乐"},
	}

	for _, t := range tags {
		tag := &Tag{
			ID:          t.ID,
			Name:        t.Name,
			Description: t.Name + "相关标签",
		}
		a.CreateTag(tag)
	}

	// 初始化集合数据
	collections := []struct {
		ID   string
		Name string
	}{
		{"collection1", "常用工具"},
		{"collection2", "开发资源"},
		{"collection3", "文档资料"},
		{"collection4", "项目管理"},
	}

	for _, c := range collections {
		col := &Collection{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Name + "相关集合",
		}
		a.CreateCollection(col)
	}

	// 初始化特殊命令（首页和设置）
	homeCmd := &Command{
		ID:            "home",
		Name:          "首页",
		Content:       "",
		Description:   "首页菜单",
		TagIDs:        []string{},
		CollectionIDs: []string{},
	}
	a.CreateCommand(homeCmd)

	settingsCmd := &Command{
		ID:            "settings",
		Name:          "设置",
		Content:       "",
		Description:   "系统设置",
		TagIDs:        []string{},
		CollectionIDs: []string{},
	}
	a.CreateCommand(settingsCmd)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// GetMenuItems returns the menu items for the application
func (a *App) GetMenuItems() map[string]interface{} {
	// 从模拟数据库获取标签
	tagItems := make([]map[string]string, 0)

	// 添加首页到标签列表
	tagItems = append(tagItems, map[string]string{
		"id":   "home",
		"name": "首页",
		"icon": "🏷️",
	})

	// 添加所有标签
	for _, tag := range a.GetTags() {
		tagItems = append(tagItems, map[string]string{
			"id":   tag.ID,
			"name": tag.Name,
			"icon": "🏷️",
		})
	}

	// 从模拟数据库获取集合
	collectionItems := make([]map[string]string, 0)

	// 添加首页到集合列表
	collectionItems = append(collectionItems, map[string]string{
		"id":   "home",
		"name": "首页",
		"icon": "{ }",
	})

	// 添加所有集合
	for _, col := range a.GetCollections() {
		collectionItems = append(collectionItems, map[string]string{
			"id":   col.ID,
			"name": col.Name,
			"icon": "{ }",
		})
	}

	// 构建"all"列表
	allItems := make([]map[string]string, 0)

	// 添加首页
	allItems = append(allItems, map[string]string{
		"id":   "home",
		"name": "首页",
		"icon": "🎇",
	})

	// 添加所有标签
	for _, tag := range a.GetTags() {
		allItems = append(allItems, map[string]string{
			"id":   tag.ID,
			"name": tag.Name,
			"icon": "🎇",
		})
	}

	// 添加所有集合
	for _, col := range a.GetCollections() {
		allItems = append(allItems, map[string]string{
			"id":   col.ID,
			"name": col.Name,
			"icon": "🎇",
		})
	}

	// 添加设置
	allItems = append(allItems, map[string]string{
		"id":   "settings",
		"name": "设置",
		"icon": "🎇",
	})

	return map[string]interface{}{
		"tags":        tagItems,
		"collections": collectionItems,
		"all":         allItems,
	}
}

// CreateCommand 创建指令
func (a *App) CreateCommand(cmd *Command) error {
	// 简单的ID生成（实际应用中应该使用更可靠的ID生成方式）
	if cmd.ID == "" {
		cmd.ID = fmt.Sprintf("cmd_%d", time.Now().UnixNano())
	}

	now := time.Now()
	cmd.CreatedAt = now
	cmd.UpdatedAt = now
	cmd.CopyCount = 0
	cmd.SearchCount = 0

	// 确保TagIDs和CollectionIDs不为nil
	if cmd.TagIDs == nil {
		cmd.TagIDs = []string{}
	}
	if cmd.CollectionIDs == nil {
		cmd.CollectionIDs = []string{}
	}

	// 保存到模拟数据库
	a.commands[cmd.ID] = cmd
	return nil
}

// GetCommand 获取单个指令
func (a *App) GetCommand(id string) (*Command, error) {
	cmd, exists := a.commands[id]
	if !exists || cmd.DeletedAt != nil {
		return nil, fmt.Errorf("command not found: %s", id)
	}
	return cmd, nil
}

// GetCommands 获取所有指令
func (a *App) GetCommands() []*Command {
	commands := make([]*Command, 0, len(a.commands))
	for _, cmd := range a.commands {
		if cmd.DeletedAt == nil {
			commands = append(commands, cmd)
		}
	}
	return commands
}

// UpdateCommand 更新指令
func (a *App) UpdateCommand(cmd *Command) error {
	// 检查指令是否存在
	_, exists := a.commands[cmd.ID]
	if !exists {
		return fmt.Errorf("command not found: %s", cmd.ID)
	}

	// 更新时间
	cmd.UpdatedAt = time.Now()

	// 保存到模拟数据库
	a.commands[cmd.ID] = cmd
	return nil
}

// DeleteCommand 删除指令
func (a *App) DeleteCommand(id string) error {
	// 检查指令是否存在
	cmd, exists := a.commands[id]
	if !exists {
		return fmt.Errorf("command not found: %s", id)
	}

	// 软删除
	now := time.Now()
	cmd.DeletedAt = &now
	cmd.UpdatedAt = now

	// 保存到模拟数据库
	a.commands[id] = cmd
	return nil
}

// CreateTag 创建标签
func (a *App) CreateTag(tag *Tag) error {
	// 简单的ID生成
	if tag.ID == "" {
		tag.ID = fmt.Sprintf("tag_%d", time.Now().UnixNano())
	}

	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	tag.SearchCount = 0

	// 保存到模拟数据库
	a.tags[tag.ID] = tag
	return nil
}

// GetTag 获取单个标签
func (a *App) GetTag(id string) (*Tag, error) {
	tag, exists := a.tags[id]
	if !exists || tag.DeletedAt != nil {
		return nil, fmt.Errorf("tag not found: %s", id)
	}
	return tag, nil
}

// GetTags 获取所有标签
func (a *App) GetTags() []*Tag {
	tags := make([]*Tag, 0, len(a.tags))
	for _, tag := range a.tags {
		if tag.DeletedAt == nil {
			tags = append(tags, tag)
		}
	}
	return tags
}

// UpdateTag 更新标签
func (a *App) UpdateTag(tag *Tag) error {
	// 检查标签是否存在
	_, exists := a.tags[tag.ID]
	if !exists {
		return fmt.Errorf("tag not found: %s", tag.ID)
	}

	// 更新时间
	tag.UpdatedAt = time.Now()

	// 保存到模拟数据库
	a.tags[tag.ID] = tag
	return nil
}

// DeleteTag 删除标签
func (a *App) DeleteTag(id string) error {
	// 检查标签是否存在
	tag, exists := a.tags[id]
	if !exists {
		return fmt.Errorf("tag not found: %s", id)
	}

	// 软删除
	now := time.Now()
	tag.DeletedAt = &now
	tag.UpdatedAt = now

	// 保存到模拟数据库
	a.tags[id] = tag
	return nil
}

// CreateCollection 创建集合
func (a *App) CreateCollection(col *Collection) error {
	// 简单的ID生成
	if col.ID == "" {
		col.ID = fmt.Sprintf("col_%d", time.Now().UnixNano())
	}

	now := time.Now()
	col.CreatedAt = now
	col.UpdatedAt = now
	col.SearchCount = 0

	// 保存到模拟数据库
	a.collections[col.ID] = col
	return nil
}

// GetCollection 获取单个集合
func (a *App) GetCollection(id string) (*Collection, error) {
	col, exists := a.collections[id]
	if !exists || col.DeletedAt != nil {
		return nil, fmt.Errorf("collection not found: %s", id)
	}
	return col, nil
}

// GetCollections 获取所有集合
func (a *App) GetCollections() []*Collection {
	cols := make([]*Collection, 0, len(a.collections))
	for _, col := range a.collections {
		if col.DeletedAt == nil {
			cols = append(cols, col)
		}
	}
	return cols
}

// UpdateCollection 更新集合
func (a *App) UpdateCollection(col *Collection) error {
	// 检查集合是否存在
	_, exists := a.collections[col.ID]
	if !exists {
		return fmt.Errorf("collection not found: %s", col.ID)
	}

	// 更新时间
	col.UpdatedAt = time.Now()

	// 保存到模拟数据库
	a.collections[col.ID] = col
	return nil
}

// DeleteCollection 删除集合
func (a *App) DeleteCollection(id string) error {
	// 检查集合是否存在
	col, exists := a.collections[id]
	if !exists {
		return fmt.Errorf("collection not found: %s", id)
	}

	// 软删除
	now := time.Now()
	col.DeletedAt = &now
	col.UpdatedAt = now

	// 保存到模拟数据库
	a.collections[id] = col
	return nil
}
