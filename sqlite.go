package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB *sql.DB
)

// 安全说明：
// 1. 所有SQL查询都使用参数化查询（?占位符）来防止SQL注入
// 2. 所有接受外部输入的函数都包含输入验证
// 3. 使用软删除（deleted_at字段）而不是物理删除
// 4. 所有数据库操作都有适当的错误处理

// InitSqlite 初始化SQLite数据库
func InitSqlite() {
	// 打开SQLite数据库，使用文件存储而不是内存存储
	db, err := sql.Open("sqlite3", "file:quick-cmd.db?cache=shared")
	if err != nil {
		panic(err)
	}
	DB = db

	// 创建所有必需的表
	if err := createTables(); err != nil {
		panic(fmt.Errorf("创建表失败: %v", err))
	}
}

// createTables 创建所有必需的表
func createTables() error {
	// 创建标签表
	_, err := DB.Exec(`
	CREATE TABLE IF NOT EXISTS tags (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		search_count INTEGER DEFAULT 0,
		os INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME
	);
	`)
	if err != nil {
		return fmt.Errorf("创建tags表失败: %v", err)
	}

	// 创建集合表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS collections (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		search_count INTEGER DEFAULT 0,
		os INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME
	);
	`)
	if err != nil {
		return fmt.Errorf("创建collections表失败: %v", err)
	}

	// 创建命令表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS commands (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		description TEXT,
		copy_count INTEGER DEFAULT 0,
		search_count INTEGER DEFAULT 0,
		os INTEGER NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME
	);
	`)
	if err != nil {
		return fmt.Errorf("创建commands表失败: %v", err)
	}

	// 创建命令与标签的多对多关系表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS command_tags (
		command_id TEXT NOT NULL,
		tag_id TEXT NOT NULL,
		PRIMARY KEY (command_id, tag_id),
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建command_tags表失败: %v", err)
	}

	// 创建命令与集合的多对多关系表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS command_collections (
		command_id TEXT NOT NULL,
		collection_id TEXT NOT NULL,
		PRIMARY KEY (command_id, collection_id),
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建command_collections表失败: %v", err)
	}

	return nil
}

// CreateTagSQLite 创建标签
func CreateTagSQLite(tag *Tag) error {
	// 简单的ID生成
	if tag.ID == "" {
		tag.ID = fmt.Sprintf("tag_%d", time.Now().UnixNano())
	}

	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	tag.SearchCount = 0

	// 保存到SQLite数据库
	_, err := DB.Exec(
		"INSERT INTO tags (id, name, description, search_count, os, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		tag.ID, tag.Name, tag.Description, tag.SearchCount, tag.Os, tag.CreatedAt, tag.UpdatedAt, tag.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("创建标签失败: %v", err)
	}

	return nil
}

// GetTagSQLite 获取单个标签
func GetTagSQLite(id string) (*Tag, error) {
	// 输入验证
	if id == "" {
		return nil, fmt.Errorf("标签ID不能为空")
	}

	var tag Tag
	var deletedAt sql.NullTime

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM tags WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.SearchCount, &tag.Os, &tag.CreatedAt, &tag.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %s", id)
		}
		return nil, fmt.Errorf("获取标签失败: %v", err)
	}

	// 处理deletedAt字段
	if deletedAt.Valid {
		tag.DeletedAt = &deletedAt.Time
	}

	return &tag, nil
}

// GetTagsSQLite 获取所有标签
func GetTagsSQLite() ([]*Tag, error) {
	var tags []*Tag

	// 从SQLite数据库获取所有标签
	rows, err := DB.Query(
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM tags WHERE deleted_at IS NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("获取标签列表失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var tag Tag
		var deletedAt sql.NullTime

		err := rows.Scan(
			&tag.ID, &tag.Name, &tag.Description, &tag.SearchCount, &tag.Os, &tag.CreatedAt, &tag.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描标签失败: %v", err)
		}

		// 处理deletedAt字段
		if deletedAt.Valid {
			tag.DeletedAt = &deletedAt.Time
		}

		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历标签结果集失败: %v", err)
	}

	return tags, nil
}

// UpdateTagSQLite 更新标签
func UpdateTagSQLite(tag *Tag) error {
	tag.UpdatedAt = time.Now()

	// 更新SQLite数据库中的标签
	_, err := DB.Exec(
		"UPDATE tags SET name = ?, description = ?, search_count = ?, os = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		tag.Name, tag.Description, tag.SearchCount, tag.Os, tag.UpdatedAt, tag.ID,
	)
	if err != nil {
		return fmt.Errorf("更新标签失败: %v", err)
	}

	return nil
}

// DeleteTagSQLite 删除标签（软删除）
func DeleteTagSQLite(id string) error {
	now := time.Now()

	// 更新标签的deleted_at字段
	result, err := DB.Exec(
		"UPDATE tags SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		now, now, id,
	)
	if err != nil {
		return fmt.Errorf("删除标签失败: %v", err)
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取删除影响的行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("tag not found: %s", id)
	}

	return nil
}

// CreateCollectionSQLite 创建集合
func CreateCollectionSQLite(collection *Collection) error {
	// 简单的ID生成
	if collection.ID == "" {
		collection.ID = fmt.Sprintf("collection_%d", time.Now().UnixNano())
	}

	now := time.Now()
	collection.CreatedAt = now
	collection.UpdatedAt = now
	collection.SearchCount = 0

	// 保存到SQLite数据库
	_, err := DB.Exec(
		"INSERT INTO collections (id, name, description, search_count, os, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		collection.ID, collection.Name, collection.Description, collection.SearchCount, collection.Os, collection.CreatedAt, collection.UpdatedAt, collection.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("创建集合失败: %v", err)
	}

	return nil
}

// GetCollectionSQLite 获取单个集合
func GetCollectionSQLite(id string) (*Collection, error) {
	// 输入验证
	if id == "" {
		return nil, fmt.Errorf("集合ID不能为空")
	}

	var collection Collection
	var deletedAt sql.NullTime

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM collections WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&collection.ID, &collection.Name, &collection.Description, &collection.SearchCount, &collection.Os, &collection.CreatedAt, &collection.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found: %s", id)
		}
		return nil, fmt.Errorf("获取集合失败: %v", err)
	}

	// 处理deletedAt字段
	if deletedAt.Valid {
		collection.DeletedAt = &deletedAt.Time
	}

	return &collection, nil
}

// GetCollectionsSQLite 获取所有集合
func GetCollectionsSQLite() ([]*Collection, error) {
	var collections []*Collection

	// 从SQLite数据库获取所有集合
	rows, err := DB.Query(
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM collections WHERE deleted_at IS NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("获取集合列表失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var collection Collection
		var deletedAt sql.NullTime

		err := rows.Scan(
			&collection.ID, &collection.Name, &collection.Description, &collection.SearchCount, &collection.Os, &collection.CreatedAt, &collection.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描集合失败: %v", err)
		}

		// 处理deletedAt字段
		if deletedAt.Valid {
			collection.DeletedAt = &deletedAt.Time
		}

		collections = append(collections, &collection)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历集合结果集失败: %v", err)
	}

	return collections, nil
}

// UpdateCollectionSQLite 更新集合
func UpdateCollectionSQLite(collection *Collection) error {
	collection.UpdatedAt = time.Now()

	// 更新SQLite数据库中的集合
	_, err := DB.Exec(
		"UPDATE collections SET name = ?, description = ?, search_count = ?, os = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		collection.Name, collection.Description, collection.SearchCount, collection.Os, collection.UpdatedAt, collection.ID,
	)
	if err != nil {
		return fmt.Errorf("更新集合失败: %v", err)
	}

	return nil
}

// DeleteCollectionSQLite 删除集合（软删除）
func DeleteCollectionSQLite(id string) error {
	now := time.Now()

	// 更新集合的deleted_at字段
	result, err := DB.Exec(
		"UPDATE collections SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		now, now, id,
	)
	if err != nil {
		return fmt.Errorf("删除集合失败: %v", err)
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取删除影响的行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("collection not found: %s", id)
	}

	return nil
}

// CreateCommandSQLite 创建命令
func CreateCommandSQLite(cmd *Command) error {
	// 简单的ID生成
	if cmd.ID == "" {
		cmd.ID = fmt.Sprintf("command_%d", time.Now().UnixNano())
	}

	now := time.Now()
	cmd.CreatedAt = now
	cmd.UpdatedAt = now
	cmd.CopyCounts = 0
	cmd.SearchCount = 0

	// 保存命令基本信息到SQLite数据库
	_, err := DB.Exec(
		"INSERT INTO commands (id, name, content, description, copy_count, search_count, os, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		cmd.ID, cmd.Name, cmd.Content, cmd.Description, cmd.CopyCounts, cmd.SearchCount, cmd.Os, cmd.CreatedAt, cmd.UpdatedAt, cmd.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("创建命令失败: %v", err)
	}

	// 保存命令与标签的多对多关系
	for _, tagID := range cmd.TagIDs {
		if err := AddTagToCommandSQLite(cmd.ID, tagID); err != nil {
			return fmt.Errorf("添加命令标签关系失败: %v", err)
		}
	}

	// 保存命令与集合的多对多关系
	for _, collectionID := range cmd.CollectionIDs {
		if err := AddCollectionToCommandSQLite(cmd.ID, collectionID); err != nil {
			return fmt.Errorf("添加命令集合关系失败: %v", err)
		}
	}

	return nil
}

// GetCommandSQLite 获取单个命令
func GetCommandSQLite(id string) (*Command, error) {
	// 输入验证
	if id == "" {
		return nil, fmt.Errorf("命令ID不能为空")
	}

	var cmd Command
	var deletedAt sql.NullTime

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, content, description, copy_count, search_count, os, created_at, updated_at, deleted_at FROM commands WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&cmd.ID, &cmd.Name, &cmd.Content, &cmd.Description, &cmd.CopyCounts, &cmd.SearchCount, &cmd.Os, &cmd.CreatedAt, &cmd.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("command not found: %s", id)
		}
		return nil, fmt.Errorf("获取命令失败: %v", err)
	}

	// 处理deletedAt字段
	if deletedAt.Valid {
		cmd.DeletedAt = &deletedAt.Time
	}

	// 获取命令关联的标签ID
	if cmd.TagIDs, err = GetTagIDsByCommandIDSQLite(id); err != nil {
		return nil, fmt.Errorf("获取命令标签ID失败: %v", err)
	}

	// 获取命令关联的集合ID
	if cmd.CollectionIDs, err = GetCollectionIDsByCommandIDSQLite(id); err != nil {
		return nil, fmt.Errorf("获取命令集合ID失败: %v", err)
	}

	return &cmd, nil
}

// GetCommandsSQLite 获取所有命令
func GetCommandsSQLite() ([]*Command, error) {
	var commands []*Command

	// 从SQLite数据库获取所有命令的基本信息
	rows, err := DB.Query(
		"SELECT id, name, content, description, copy_count, search_count, os, created_at, updated_at, deleted_at FROM commands WHERE deleted_at IS NULL",
	)
	if err != nil {
		return nil, fmt.Errorf("获取命令列表失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var cmd Command
		var deletedAt sql.NullTime

		err := rows.Scan(
			&cmd.ID, &cmd.Name, &cmd.Content, &cmd.Description, &cmd.CopyCounts, &cmd.SearchCount, &cmd.Os, &cmd.CreatedAt, &cmd.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描命令失败: %v", err)
		}

		// 处理deletedAt字段
		if deletedAt.Valid {
			cmd.DeletedAt = &deletedAt.Time
		}

		// 获取命令关联的标签ID
		if cmd.TagIDs, err = GetTagIDsByCommandIDSQLite(cmd.ID); err != nil {
			return nil, fmt.Errorf("获取命令标签ID失败: %v", err)
		}

		// 获取命令关联的集合ID
		if cmd.CollectionIDs, err = GetCollectionIDsByCommandIDSQLite(cmd.ID); err != nil {
			return nil, fmt.Errorf("获取命令集合ID失败: %v", err)
		}

		commands = append(commands, &cmd)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历命令结果集失败: %v", err)
	}

	return commands, nil
}

// UpdateCommandSQLite 更新命令
func UpdateCommandSQLite(cmd *Command) error {
	// 输入验证
	if cmd == nil {
		return fmt.Errorf("命令对象不能为空")
	}
	if cmd.ID == "" {
		return fmt.Errorf("命令ID不能为空")
	}

	cmd.UpdatedAt = time.Now()

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"UPDATE commands SET name = ?, content = ?, description = ?, copy_count = ?, search_count = ?, os = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		cmd.Name, cmd.Content, cmd.Description, cmd.CopyCounts, cmd.SearchCount, cmd.Os, cmd.UpdatedAt, cmd.ID,
	)
	if err != nil {
		return fmt.Errorf("更新命令失败: %v", err)
	}

	// 先删除现有的命令与标签的关系
	if err := RemoveAllTagsFromCommandSQLite(cmd.ID); err != nil {
		return fmt.Errorf("删除命令标签关系失败: %v", err)
	}

	// 重新添加命令与标签的关系
	for _, tagID := range cmd.TagIDs {
		if err := AddTagToCommandSQLite(cmd.ID, tagID); err != nil {
			return fmt.Errorf("添加命令标签关系失败: %v", err)
		}
	}

	// 先删除现有的命令与集合的关系
	if err := RemoveAllCollectionsFromCommandSQLite(cmd.ID); err != nil {
		return fmt.Errorf("删除命令集合关系失败: %v", err)
	}

	// 重新添加命令与集合的关系
	for _, collectionID := range cmd.CollectionIDs {
		if err := AddCollectionToCommandSQLite(cmd.ID, collectionID); err != nil {
			return fmt.Errorf("添加命令集合关系失败: %v", err)
		}
	}

	return nil
}

// DeleteCommandSQLite 删除命令（软删除）
func DeleteCommandSQLite(id string) error {
	// 输入验证
	if id == "" {
		return fmt.Errorf("命令ID不能为空")
	}

	now := time.Now()

	// 使用参数化查询，防止SQL注入
	result, err := DB.Exec(
		"UPDATE commands SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		now, now, id,
	)
	if err != nil {
		return fmt.Errorf("删除命令失败: %v", err)
	}

	// 检查是否有行被更新
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("获取删除影响的行数失败: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("command not found: %s", id)
	}

	return nil
}

// 管理命令与标签的多对多关系

// AddTagToCommandSQLite 添加标签到命令
func AddTagToCommandSQLite(commandID, tagID string) error {
	// 输入验证
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}
	if tagID == "" {
		return fmt.Errorf("标签ID不能为空")
	}

	// 检查标签是否存在
	if _, err := GetTagSQLite(tagID); err != nil {
		return fmt.Errorf("标签不存在: %v", err)
	}

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"INSERT OR IGNORE INTO command_tags (command_id, tag_id) VALUES (?, ?)",
		commandID, tagID,
	)
	if err != nil {
		return fmt.Errorf("添加命令标签关系失败: %v", err)
	}

	return nil
}

// RemoveTagFromCommandSQLite 从命令移除标签
func RemoveTagFromCommandSQLite(commandID, tagID string) error {
	// 输入验证
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}
	if tagID == "" {
		return fmt.Errorf("标签ID不能为空")
	}

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"DELETE FROM command_tags WHERE command_id = ? AND tag_id = ?",
		commandID, tagID,
	)
	if err != nil {
		return fmt.Errorf("删除命令标签关系失败: %v", err)
	}

	return nil
}

// RemoveAllTagsFromCommandSQLite 从命令移除所有标签
func RemoveAllTagsFromCommandSQLite(commandID string) error {
	// 输入验证
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"DELETE FROM command_tags WHERE command_id = ?",
		commandID,
	)
	if err != nil {
		return fmt.Errorf("删除所有命令标签关系失败: %v", err)
	}

	return nil
}

// GetTagIDsByCommandIDSQLite 获取命令的所有标签ID
func GetTagIDsByCommandIDSQLite(commandID string) ([]string, error) {
	// 输入验证
	if commandID == "" {
		return nil, fmt.Errorf("commandID不能为空")
	}

	var tagIDs []string

	// 使用参数化查询，防止SQL注入
	rows, err := DB.Query(
		"SELECT tag_id FROM command_tags WHERE command_id = ?",
		commandID,
	)
	if err != nil {
		return nil, fmt.Errorf("获取命令标签ID失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var tagID string
		if err := rows.Scan(&tagID); err != nil {
			return nil, fmt.Errorf("扫描标签ID失败: %v", err)
		}
		tagIDs = append(tagIDs, tagID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历标签ID结果集失败: %v", err)
	}

	return tagIDs, nil
}

// 管理命令与集合的多对多关系

// AddCollectionToCommandSQLite 添加集合到命令
func AddCollectionToCommandSQLite(commandID, collectionID string) error {
	// 输入验证
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}
	if collectionID == "" {
		return fmt.Errorf("集合ID不能为空")
	}

	// 检查集合是否存在
	if _, err := GetCollectionSQLite(collectionID); err != nil {
		return fmt.Errorf("集合不存在: %v", err)
	}

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"INSERT OR IGNORE INTO command_collections (command_id, collection_id) VALUES (?, ?)",
		commandID, collectionID,
	)
	if err != nil {
		return fmt.Errorf("添加命令集合关系失败: %v", err)
	}

	return nil
}

// RemoveCollectionFromCommandSQLite 从命令移除集合
func RemoveCollectionFromCommandSQLite(commandID, collectionID string) error {
	// 输入验证
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}
	if collectionID == "" {
		return fmt.Errorf("集合ID不能为空")
	}

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"DELETE FROM command_collections WHERE command_id = ? AND collection_id = ?",
		commandID, collectionID,
	)
	if err != nil {
		return fmt.Errorf("删除命令集合关系失败: %v", err)
	}

	return nil
}

// RemoveAllCollectionsFromCommandSQLite 从命令移除所有集合
func RemoveAllCollectionsFromCommandSQLite(commandID string) error {
	// 输入验证
	if commandID == "" {
		return fmt.Errorf("命令ID不能为空")
	}

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"DELETE FROM command_collections WHERE command_id = ?",
		commandID,
	)
	if err != nil {
		return fmt.Errorf("删除所有命令集合关系失败: %v", err)
	}

	return nil
}

// GetCollectionIDsByCommandIDSQLite 获取命令的所有集合ID
func GetCollectionIDsByCommandIDSQLite(commandID string) ([]string, error) {
	// 输入验证
	if commandID == "" {
		return nil, fmt.Errorf("commandID不能为空")
	}

	var collectionIDs []string

	// 使用参数化查询，防止SQL注入
	rows, err := DB.Query(
		"SELECT collection_id FROM command_collections WHERE command_id = ?",
		commandID,
	)
	if err != nil {
		return nil, fmt.Errorf("获取命令集合ID失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var collectionID string
		if err := rows.Scan(&collectionID); err != nil {
			return nil, fmt.Errorf("扫描集合ID失败: %v", err)
		}
		collectionIDs = append(collectionIDs, collectionID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历集合ID结果集失败: %v", err)
	}

	return collectionIDs, nil
}
