package main

import (
	"fmt"
	"time"
)

// Tag 标签结构体
type Tag struct {
	ID          string     `json:"id"`
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
	// 简单的ID生成
	if tag.ID == "" {
		tag.ID = fmt.Sprintf("tag_%d", time.Now().UnixNano())
	}

	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	tag.SearchCount = 0

	// 保存到模拟数据库
	a.tags = append(a.tags, tag)
	return nil
}

// GetTag 获取单个标签
func (a *App) GetTag(id string) (*Tag, error) {
	for _, tag := range a.tags {
		if tag.ID == id && tag.DeletedAt == nil {
			return tag, nil
		}
	}
	return nil, fmt.Errorf("tag not found: %s", id)
}

// GetTags 获取所有标签
func (a *App) GetTags() []*Tag {
	return []*Tag{
		{
			ID:          "0",
			Name:        "全部",
			Description: "全部",
			SearchCount: 1,
			Os:          AllOs,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "1",
			Name:        "存储",
			Description: "存储空间检查",
			SearchCount: 1,
			Os:          Linux,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "2",
			Name:        "进程",
			Description: "进程管理",
			SearchCount: 2,
			Os:          Linux,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          "3",
			Name:        "tcp工具",
			Description: "tcp工具",
			SearchCount: 3,
			Os:          Linux,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
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
func (a *App) DeleteTag(id string) error {
	// 检查标签是否存在

	return nil
}
