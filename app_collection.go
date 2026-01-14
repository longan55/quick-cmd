package main

import (
	"fmt"
	"log"
	"time"
)

// Collection 集合结构体
type Collection struct {
	ID          uint64   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	SearchCount int      `json:"searchCount,omitempty"`
	Os          []string `json:"os,omitempty"`
	CommandIDs  []uint64 `json:"commandIds,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
	DeletedAt   string   `json:"deletedAt,omitempty"`
}

// CreateCollection 创建集合
func (a *App) CreateCollection(col *Collection) error {
	log.Printf("CreateCollection: %+v\n", col)
	// 简单的ID生成
	col.ID = uint64(time.Now().UnixNano())
	err := CreateCollectionSQLite(col)
	if err != nil {
		log.Printf("创建集合失败: %v", err)
		return fmt.Errorf("创建集合失败: %v", err)
	}
	log.Printf("创建集合成功: %+v", col)
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
