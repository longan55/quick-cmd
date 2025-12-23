package main

import (
	"fmt"
	"time"
)

// Tag 标签结构体
type Tag struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	SearchCount int        `json:"searchCount"`
	Os          OS         `json:"os"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

// CreateTag 创建标签
func (a *App) CreateTag(tag *Tag) error {

	return nil
}

// GetTag 获取单个标签
func (a *App) GetTag(id uint64) (*Tag, error) {
	return nil, fmt.Errorf("tag not found: %d", id)
}

func (a *App) GetCommandsByTagId(option Option) []*Command {
	fmt.Printf("GetCommandsByTagId: %v\n", option)
	return a.commands
}

// UpdateTag 更新标签
func (a *App) UpdateTag(tag *Tag) error {
	// 检查标签是否存在
	for _, t := range a.tags {
		if t.ID == tag.ID && t.DeletedAt == nil {
			t.Name = tag.Name
			t.Description = tag.Description
			t.UpdatedAt = time.Now()
			return nil
		}
	}
	// 更新时间
	tag.UpdatedAt = time.Now()

	// 保存到模拟数据库
	a.tags = append(a.tags, tag)
	return nil
}

// DeleteTag 删除标签
func (a *App) DeleteTag(id uint64) error {
	// 检查标签是否存在

	return nil
}
