package main

import (
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
	Os            OS         `json:"os"`
	TagIDs        []string   `json:"tagIDs"`        // 标签ID列表
	CollectionIDs []string   `json:"collectionIDs"` // 集合ID列表
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

// CreateCommand 创建指令
func (a *App) CreateCommand(cmd *Command) error {
	// 简单的ID生成（实际应用中应该使用更可靠的ID生成方式）

	return nil
}

// GetCommand 获取单个指令
func (a *App) GetCommand(id string) (*Command, error) {

	return nil, nil
}

var allCommands = []*Command{
	{
		ID:            "1",
		Name:          "当前存储空间",
		Content:       "df -h",
		Description:   "列出当前存储空间",
		CopyCount:     0,
		SearchCount:   0,
		Os:            Linux,
		TagIDs:        []string{"1"},
		CollectionIDs: []string{"1"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		ID:            "2",
		Name:          "进程号查询",
		Content:       "ps aux",
		Description:   "列出所有进程",
		CopyCount:     0,
		SearchCount:   0,
		Os:            Linux,
		TagIDs:        []string{"2"},
		CollectionIDs: []string{"1"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		ID:            "3",
		Name:          "tcpdump",
		Content:       "tcpdump -i eth0",
		Description:   "监听eth0网卡",
		CopyCount:     0,
		SearchCount:   0,
		Os:            Linux,
		TagIDs:        []string{"3"},
		CollectionIDs: []string{"1"},
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
}

// GetCommands 获取所有指令
func (a *App) GetCommands() []*Command {
	return allCommands
}

// GetCommandsByTagID 根据标签ID获取指令
func (a *App) GetCommandsByTagID(tagID string) []*Command {
	// 过滤出包含指定标签ID的命令
	filteredCommands := make([]*Command, 0)
	for _, cmd := range allCommands {
		for _, id := range cmd.TagIDs {
			if id == tagID {
				filteredCommands = append(filteredCommands, cmd)
				break
			}
		}
	}
	return filteredCommands
}
func (a *App) GetCommandsByCollectionID(collectionID string) []*Command {
	filteredCommands := make([]*Command, 0)
	for _, cmd := range allCommands {
		if cmd.DeletedAt == nil {
			for _, id := range cmd.CollectionIDs {
				if id == collectionID {
					filteredCommands = append(filteredCommands, cmd)
					break
				}
			}
		}
	}
	return filteredCommands
}

// UpdateCommand 更新指令
func (a *App) UpdateCommand(cmd *Command) error {
	// 检查指令是否存在

	return nil
}

// DeleteCommand 删除指令
func (a *App) DeleteCommand(id string) error {
	// 检查指令是否存在

	return nil
}
