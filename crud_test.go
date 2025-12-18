package main

import (
	"testing"
)

// TestCommandCRUD 测试指令的CURD功能
func TestCommandCRUD(t *testing.T) {
	app := NewApp()
	
	// 测试创建指令
	cmd := &Command{
		Name:        "测试指令",
		Content:     "echo test",
		Description: "测试用指令",
	}
	err := app.CreateCommand(cmd)
	if err != nil {
		t.Fatalf("创建指令失败: %v", err)
	}
	
	// 测试获取指令
	cmd2, err := app.GetCommand(cmd.ID)
	if err != nil {
		t.Fatalf("获取指令失败: %v", err)
	}
	if cmd2.Name != "测试指令" {
		t.Errorf("指令名称不匹配，期望: 测试指令, 实际: %s", cmd2.Name)
	}
	
	// 测试更新指令
	cmd2.Name = "更新后的指令"
	err = app.UpdateCommand(cmd2)
	if err != nil {
		t.Fatalf("更新指令失败: %v", err)
	}
	
	cmd3, _ := app.GetCommand(cmd.ID)
	if cmd3.Name != "更新后的指令" {
		t.Errorf("指令更新后名称不匹配，期望: 更新后的指令, 实际: %s", cmd3.Name)
	}
	
	// 测试删除指令
	err = app.DeleteCommand(cmd.ID)
	if err != nil {
		t.Fatalf("删除指令失败: %v", err)
	}
	
	_, err = app.GetCommand(cmd.ID)
	if err == nil {
		t.Error("指令删除后仍能获取到，删除失败")
	}
	
	t.Log("指令CRUD测试通过")
}

// TestTagCRUD 测试标签的CURD功能
func TestTagCRUD(t *testing.T) {
	app := NewApp()
	
	// 测试创建标签
	tag := &Tag{
		Name:        "测试标签",
		Description: "测试用标签",
	}
	err := app.CreateTag(tag)
	if err != nil {
		t.Fatalf("创建标签失败: %v", err)
	}
	
	// 测试获取标签
	tag2, err := app.GetTag(tag.ID)
	if err != nil {
		t.Fatalf("获取标签失败: %v", err)
	}
	if tag2.Name != "测试标签" {
		t.Errorf("标签名称不匹配，期望: 测试标签, 实际: %s", tag2.Name)
	}
	
	// 测试更新标签
	tag2.Name = "更新后的标签"
	err = app.UpdateTag(tag2)
	if err != nil {
		t.Fatalf("更新标签失败: %v", err)
	}
	
	tag3, _ := app.GetTag(tag.ID)
	if tag3.Name != "更新后的标签" {
		t.Errorf("标签更新后名称不匹配，期望: 更新后的标签, 实际: %s", tag3.Name)
	}
	
	// 测试删除标签
	err = app.DeleteTag(tag.ID)
	if err != nil {
		t.Fatalf("删除标签失败: %v", err)
	}
	
	_, err = app.GetTag(tag.ID)
	if err == nil {
		t.Error("标签删除后仍能获取到，删除失败")
	}
	
	t.Log("标签CRUD测试通过")
}

// TestCollectionCRUD 测试集合的CURD功能
func TestCollectionCRUD(t *testing.T) {
	app := NewApp()
	
	// 测试创建集合
	col := &Collection{
		Name:        "测试集合",
		Description: "测试用集合",
	}
	err := app.CreateCollection(col)
	if err != nil {
		t.Fatalf("创建集合失败: %v", err)
	}
	
	// 测试获取集合
	col2, err := app.GetCollection(col.ID)
	if err != nil {
		t.Fatalf("获取集合失败: %v", err)
	}
	if col2.Name != "测试集合" {
		t.Errorf("集合名称不匹配，期望: 测试集合, 实际: %s", col2.Name)
	}
	
	// 测试更新集合
	col2.Name = "更新后的集合"
	err = app.UpdateCollection(col2)
	if err != nil {
		t.Fatalf("更新集合失败: %v", err)
	}
	
	col3, _ := app.GetCollection(col.ID)
	if col3.Name != "更新后的集合" {
		t.Errorf("集合更新后名称不匹配，期望: 更新后的集合, 实际: %s", col3.Name)
	}
	
	// 测试删除集合
	err = app.DeleteCollection(col.ID)
	if err != nil {
		t.Fatalf("删除集合失败: %v", err)
	}
	
	_, err = app.GetCollection(col.ID)
	if err == nil {
		t.Error("集合删除后仍能获取到，删除失败")
	}
	
	t.Log("集合CRUD测试通过")
}

// TestCommandTagRelation 测试指令和标签的关联关系
func TestCommandTagRelation(t *testing.T) {
	app := NewApp()
	
	// 创建标签
	tag := &Tag{
		Name: "工作",
	}
	app.CreateTag(tag)
	
	// 创建带标签的指令
	cmd := &Command{
		Name:    "Git提交",
		Content: "git commit -m \"test\"",
		TagIDs:  []string{tag.ID},
	}
	app.CreateCommand(cmd)
	
	// 验证关联关系
	cmd2, _ := app.GetCommand(cmd.ID)
	if len(cmd2.TagIDs) != 1 || cmd2.TagIDs[0] != tag.ID {
		t.Error("指令和标签的关联关系不正确")
	}
	
	t.Log("指令和标签关联关系测试通过")
}

// TestCommandCollectionRelation 测试指令和集合的关联关系
func TestCommandCollectionRelation(t *testing.T) {
	app := NewApp()
	
	// 创建集合
	col := &Collection{
		Name: "开发工具",
	}
	app.CreateCollection(col)
	
	// 创建带集合的指令
	cmd := &Command{
		Name:          "Git拉取",
		Content:       "git pull",
		CollectionIDs: []string{col.ID},
	}
	app.CreateCommand(cmd)
	
	// 验证关联关系
	cmd2, _ := app.GetCommand(cmd.ID)
	if len(cmd2.CollectionIDs) != 1 || cmd2.CollectionIDs[0] != col.ID {
		t.Error("指令和集合的关联关系不正确")
	}
	
	t.Log("指令和集合关联关系测试通过")
}