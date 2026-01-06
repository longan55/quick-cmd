package main

import (
	"fmt"
	"log"
	"time"
)

// Tag 标签结构体
type Tag struct {
	ID          uint64     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	SearchCount int        `json:"searchCount,omitempty"`
	Os          []string   `json:"os,omitempty"`
	CommandIDs  []uint64   `json:"commandIds,omitempty"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

// CreateTag 创建标签
func (a *App) CreateTag(tag *Tag) error {
	log.Printf("CreateTag: %+v\n", tag)
	// 输入验证
	if tag.Name == "" {
		return fmt.Errorf("标签名称不能为空")
	}
	err := CreateTagSQLite(tag)
	if err != nil {
		return fmt.Errorf("创建标签失败: %v", err)
	}
	return nil
}

// GetTag 获取单个标签
func (a *App) GetTag(id uint64) (*Tag, error) {
	log.Printf("GetTag: %d\n", id)
	tag, err := GetTagSQLite(id)
	if err != nil {
		return nil, fmt.Errorf("获取标签失败: %v", err)
	}
	return tag, nil
}

func (a *App) GetCommandsByTagId(option Option) []*Command {
	fmt.Printf("GetCommandsByTagId: %v\n", option)
	return a.commands
}

// UpdateTag 更新标签
func (a *App) UpdateTag(tag *Tag) error {
	log.Printf("UpdateTag: %+v\n", tag)
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
