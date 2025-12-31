package main

import (
	"fmt"
	"log"
	"time"
)

// Command 指令结构体
type Command struct {
	ID            uint64     `json:"id"`
	Name          string     `json:"name"`
	Content       string     `json:"content"`
	Description   string     `json:"description"`
	CopyCounts    int        `json:"copyCount"`
	SearchCount   int        `json:"searchCount"`
	Os            []string   `json:"os"`
	TagIDs        []uint64   `json:"tagIDs"`        // 标签ID列表
	CollectionIDs []uint64   `json:"collectionIDs"` // 集合ID列表
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

// CreateCommand 创建指令
func (a *App) CreateCommand(cmd *Command) error {
	// 简单的ID生成（实际应用中应该使用更可靠的ID生成方式）
	err := CreateCommandSQLite(cmd)
	if err != nil {
		return fmt.Errorf("创建指令失败: %v", err)
	}
	return nil
}

// GetCommand 获取单个指令
func (a *App) GetCommand(id uint64) (*Command, error) {
	log.Printf("GetCommand: %d\n", id)
	cmd, err := GetCommandSQLite(id)
	if err != nil {
		return nil, fmt.Errorf("获取指令失败: %v", err)
	}
	return cmd, nil
}

// UpdateCommand 更新指令
func (a *App) UpdateCommand(cmd *Command) error {
	// 检查指令是否存在

	return nil
}

// DeleteCommand 删除指令
func (a *App) DeleteCommand(id uint64) error {
	// 检查指令是否存在

	return nil
}
