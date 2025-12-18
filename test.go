package main

import (
	"fmt"
)

// TestCRUD 测试CURD功能
func TestCRUD() {
	fmt.Println("=== 测试指令、标签、集合的CURD功能 ===")

	// 创建App实例
	app := NewApp()

	// 测试创建标签
	fmt.Println("\n1. 测试创建标签：")
	tag1 := &Tag{
		Name:        "工作",
		Description: "工作相关的指令",
	}
	err := app.CreateTag(tag1)
	if err != nil {
		fmt.Printf("创建标签失败: %v\n", err)
		return
	}
	fmt.Printf("创建标签成功: %+v\n", tag1)

	// 测试创建集合
	fmt.Println("\n2. 测试创建集合：")
	col1 := &Collection{
		Name:        "开发工具",
		Description: "开发常用的指令集合",
	}
	err = app.CreateCollection(col1)
	if err != nil {
		fmt.Printf("创建集合失败: %v\n", err)
		return
	}
	fmt.Printf("创建集合成功: %+v\n", col1)

	// 测试创建指令
	fmt.Println("\n3. 测试创建指令：")
	cmd1 := &Command{
		Name:          "Git状态",
		Content:       "git status",
		Description:   "查看Git仓库状态",
		TagIDs:        []string{tag1.ID},
		CollectionIDs: []string{col1.ID},
	}
	err = app.CreateCommand(cmd1)
	if err != nil {
		fmt.Printf("创建指令失败: %v\n", err)
		return
	}
	fmt.Printf("创建指令成功: %+v\n", cmd1)

	// 测试获取单个指令
	fmt.Println("\n4. 测试获取单个指令：")
	cmd, err := app.GetCommand(cmd1.ID)
	if err != nil {
		fmt.Printf("获取指令失败: %v\n", err)
		return
	}
	fmt.Printf("获取指令成功: %+v\n", cmd)

	// 测试更新指令
	fmt.Println("\n5. 测试更新指令：")
	cmd.Content = "git status -s"
	cmd.Description = "查看Git仓库精简状态"
	err = app.UpdateCommand(cmd)
	if err != nil {
		fmt.Printf("更新指令失败: %v\n", err)
		return
	}
	// 获取更新后的指令
	updatedCmd, _ := app.GetCommand(cmd.ID)
	fmt.Printf("更新指令成功: %+v\n", updatedCmd)

	// 测试获取所有指令
	fmt.Println("\n6. 测试获取所有指令：")
	cmds := app.GetCommands()
	fmt.Printf("共获取到 %d 条指令\n", len(cmds))
	for _, c := range cmds {
		fmt.Printf("  - %+v\n", c)
	}

	// 测试获取所有标签
	fmt.Println("\n7. 测试获取所有标签：")
	tags := app.GetTags()
	fmt.Printf("共获取到 %d 个标签\n", len(tags))
	for _, t := range tags {
		fmt.Printf("  - %+v\n", t)
	}

	// 测试获取所有集合
	fmt.Println("\n8. 测试获取所有集合：")
	cols := app.GetCollections()
	fmt.Printf("共获取到 %d 个集合\n", len(cols))
	for _, c := range cols {
		fmt.Printf("  - %+v\n", c)
	}

	// 测试删除指令
	fmt.Println("\n9. 测试删除指令：")
	err = app.DeleteCommand(cmd.ID)
	if err != nil {
		fmt.Printf("删除指令失败: %v\n", err)
		return
	}
	fmt.Printf("删除指令成功\n")

	// 验证指令是否被删除
	cmdsAfterDelete := app.GetCommands()
	fmt.Printf("删除后还剩 %d 条指令\n", len(cmdsAfterDelete))

	fmt.Println("\n=== 所有测试完成 ===")
}

// 可以在main.go的main函数中调用TestCRUD()来测试功能
// 或者创建一个单独的测试文件使用Go的测试框架
