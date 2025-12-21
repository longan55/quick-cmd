package main

import (
	"time"
)

// Collection 集合结构体
type Collection struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SearchCount int        `json:"searchCount"`
	Os          OS         `json:"os"`
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
func (a *App) GetCollection(id string) (*Collection, error) {

	return nil, nil
}

var allCollections = []*Collection{
	{
		ID:          "0",
		Name:        "全部",
		Description: "全部",
		SearchCount: 0,
		Os:          AllOs,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
	{
		ID:          "1",
		Name:        "常用指令",
		Description: "常用指令集合",
		SearchCount: 0,
		Os:          Linux,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

// GetCollections 获取所有集合
func (a *App) GetCollections() []*Collection {
	return allCollections
}

// UpdateCollection 更新集合
func (a *App) UpdateCollection(col *Collection) error {
	// 检查集合是否存在

	return nil
}

// DeleteCollection 删除集合
func (a *App) DeleteCollection(id string) error {
	// 检查集合是否存在

	return nil
}
