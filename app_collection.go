package main

import (
	"fmt"
	"time"
)

// Collection 集合结构体
type Collection struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SearchCount int        `json:"searchCount"`
	Os          []string   `json:"os"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

// CreateCollection 创建集合
func (a *App) CreateCollection(col *Collection) error {
	// 简单的ID生成

	return nil
}

// GetCollection 获取单个集合
func (a *App) GetCollection(id uint64) (*Collection, error) {

	return nil, nil
}

// GetCollections 获取所有集合
func (a *App) GetCommandsByCollectionID(option Option) []*Command {
	fmt.Printf("GetCommandsByCollectionID: %v\n", option)
	return a.commands
}

// UpdateCollection 更新集合
func (a *App) UpdateCollection(col *Collection) error {
	// 检查集合是否存在

	return nil
}

// DeleteCollection 删除集合
func (a *App) DeleteCollection(id uint64) error {
	// 检查集合是否存在

	return nil
}
