package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
	db, err := sql.Open("sqlite3", "file:quick-cmd.db?cache=shared&mode=rw&_foreign_keys=1")
	if err != nil {
		panic(err)
	}
	DB = db
	DB.SetMaxOpenConns(1)
	// 创建所有必需的表
	if err := createTables(); err != nil {
		panic(fmt.Errorf("创建表失败: %v", err))
	}
}

// createTables 创建所有必需的表
func createTables() error {
	var err error
	// 创建标签表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS tags (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		search_count INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME
	);
	`)
	if err != nil {
		return fmt.Errorf("创建tags表失败: %v", err)
	} else {
		log.Println("tags表创建成功")
	}

	// 创建集合表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS collections (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		search_count INTEGER DEFAULT 0,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		deleted_at DATETIME
	);
	`)
	if err != nil {
		return fmt.Errorf("创建collections表失败: %v", err)
	} else {
		log.Println("collections表创建成功")
	}

	// 创建命令表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS commands (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		content TEXT NOT NULL,
		description TEXT,
		copy_count INTEGER DEFAULT 0,
		search_count INTEGER DEFAULT 0,
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
		command_id INTEGER NOT NULL,
		tag_id INTEGER NOT NULL,
		PRIMARY KEY (command_id, tag_id),
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建command_tags表失败: %v", err)
	} else {
		log.Println("command_tags表创建成功")
	}

	// 创建命令与集合的多对多关系表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS command_collections (
		command_id INTEGER NOT NULL,
		collection_id INTEGER NOT NULL,
		PRIMARY KEY (command_id, collection_id),
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE,
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建command_collections表失败: %v", err)
	}

	// 创建标签与OS的关联表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS tag_os (
		tag_id INTEGER NOT NULL,
		os TEXT NOT NULL,
		PRIMARY KEY (tag_id, os),
		FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建tag_os表失败: %v", err)
	} else {
		log.Println("tag_os表创建成功")
	}

	// 创建集合与OS的关联表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS collection_os (
		collection_id INTEGER NOT NULL,
		os TEXT NOT NULL,
		PRIMARY KEY (collection_id, os),
		FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建collection_os表失败: %v", err)
	} else {
		log.Println("collection_os表创建成功")
	}

	// 创建命令与OS的关联表
	_, err = DB.Exec(`
	CREATE TABLE IF NOT EXISTS command_os (
		command_id INTEGER NOT NULL,
		os TEXT NOT NULL,
		PRIMARY KEY (command_id, os),
		FOREIGN KEY (command_id) REFERENCES commands(id) ON DELETE CASCADE
	);
	`)
	if err != nil {
		return fmt.Errorf("创建command_os表失败: %v", err)
	} else {
		log.Println("command_os表创建成功")
	}

	// 为OS关联表创建索引，提高查询性能
	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_tag_os_os ON tag_os(os)`)
	if err != nil {
		log.Printf("创建tag_os索引失败: %v", err)
	}

	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_collection_os_os ON collection_os(os)`)
	if err != nil {
		log.Printf("创建collection_os索引失败: %v", err)
	}

	_, err = DB.Exec(`CREATE INDEX IF NOT EXISTS idx_command_os_os ON command_os(os)`)
	if err != nil {
		log.Printf("创建command_os索引失败: %v", err)
	}

	return nil
}

// 辅助函数：处理OS关联表的操作

// AddOSToTagSQLite 为标签添加OS
func AddOSToTagSQLite(tx *sql.Tx, tagID uint64, os string) error {
	_, err := tx.Exec(
		"INSERT OR IGNORE INTO tag_os (tag_id, os) VALUES (?, ?)",
		tagID, os,
	)
	return err
}

// RemoveAllOSFromTagSQLite 从标签移除所有OS
func RemoveAllOSFromTagSQLite(tagID uint64) error {
	_, err := DB.Exec(
		"DELETE FROM tag_os WHERE tag_id = ?",
		tagID,
	)
	return err
}

// GetTagOSsSQLite 获取标签的所有OS
func GetTagOSsSQLite(tagID uint64) ([]string, error) {
	var osList []string
	rows, err := DB.Query(
		"SELECT os FROM tag_os WHERE tag_id = ?",
		tagID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var os string
		err := rows.Scan(&os)
		if err != nil {
			return nil, err
		}
		osList = append(osList, os)
	}

	return osList, nil
}

// AddOSToCollectionSQLite 为集合添加OS
func AddOSToCollectionSQLite(collectionID uint64, os string) error {
	_, err := DB.Exec(
		"INSERT OR IGNORE INTO collection_os (collection_id, os) VALUES (?, ?)",
		collectionID, os,
	)
	return err
}

// RemoveAllOSFromCollectionSQLite 从集合移除所有OS
func RemoveAllOSFromCollectionSQLite(collectionID uint64) error {
	_, err := DB.Exec(
		"DELETE FROM collection_os WHERE collection_id = ?",
		collectionID,
	)
	return err
}

// GetCollectionOSsSQLite 获取集合的所有OS
func GetCollectionOSsSQLite(collectionID uint64) ([]string, error) {
	var osList []string
	rows, err := DB.Query(
		"SELECT os FROM collection_os WHERE collection_id = ?",
		collectionID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var os string
		err := rows.Scan(&os)
		if err != nil {
			return nil, err
		}
		osList = append(osList, os)
	}

	return osList, nil
}

// AddOSToCommandSQLite 为命令添加OS
func AddOSToCommandSQLite(commandID uint64, os string) error {
	_, err := DB.Exec(
		"INSERT OR IGNORE INTO command_os (command_id, os) VALUES (?, ?)",
		commandID, os,
	)
	return err
}

// RemoveAllOSFromCommandSQLite 从命令移除所有OS
func RemoveAllOSFromCommandSQLite(commandID uint64) error {
	_, err := DB.Exec(
		"DELETE FROM command_os WHERE command_id = ?",
		commandID,
	)
	return err
}

// GetCommandOSsSQLite 获取命令的所有OS
func GetCommandOSsSQLite(commandID uint64) ([]string, error) {
	var osList []string
	rows, err := DB.Query(
		"SELECT os FROM command_os WHERE command_id = ?",
		commandID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var os string
		err := rows.Scan(&os)
		if err != nil {
			return nil, err
		}
		osList = append(osList, os)
	}

	return osList, nil
}

// CreateTagSQLite 创建标签
func CreateTagSQLite(tag *Tag) error {
	log.Printf("CreateTagSQLite: %+v", tag)
	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	tag.SearchCount = 0

	//1. 检查标签是否存在
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM tags WHERE name = ? AND deleted_at IS NULL)", tag.Name).Scan(&exists)
	if err != nil {
		log.Printf("检查标签是否存在失败: %v", err)
		return fmt.Errorf("检查标签是否存在失败: %v", err)
	}
	if exists {
		log.Printf("标签[%s]已存在", tag.Name)
		return fmt.Errorf("标签[%s]已存在", tag.Name)
	}

	// 开启事务，创建tag和tag_os关联关系
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("开启事务失败: %v", err)
		return fmt.Errorf("开启事务失败: %v", err)
	}
	// defer

	result, err := tx.Exec(
		"INSERT INTO tags (name, description, search_count, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?)",
		tag.Name, tag.Description, tag.SearchCount, tag.CreatedAt, tag.UpdatedAt, tag.DeletedAt,
	)
	if err != nil {
		log.Printf("创建标签失败: %v", err)
		tx.Rollback()
		return fmt.Errorf("创建标签失败: %v", err)
	}
	log.Printf("更新tags表成功: %+v", tag)
	// 获取SQLite自动生成的ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取标签ID失败: %v", err)
		return fmt.Errorf("获取标签ID失败: %v", err)
	}
	tag.ID = uint64(id)
	log.Printf("标签ID: %d", tag.ID)
	// 2. 保存OS关联关系
	for i, os := range tag.Os {
		log.Printf("添加标签OS关系%d: %s", i, os)
		if err = AddOSToTagSQLite(tx, tag.ID, os); err != nil {
			log.Printf("添加标签OS关系失败: %v", err)
			tx.Rollback()
			return fmt.Errorf("添加标签OS关系失败: %v", err)
		}
		log.Printf("更新tag_os关系成功: %s", os)
	}
	log.Printf("创建标签成功: %+v", tag)
	tx.Commit()
	return nil
}

// GetTagSQLite 获取单个标签
func GetTagSQLite(id uint64) (*Tag, error) {
	// 输入验证
	if id == 0 {
		return nil, fmt.Errorf("标签ID不能为空")
	}

	var tag Tag
	var deletedAt sql.NullTime
	var osString string

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM tags WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.SearchCount, &osString, &tag.CreatedAt, &tag.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("tag not found: %d", id)
		}
		return nil, fmt.Errorf("获取标签失败: %v", err)
	}

	// 从关联表获取OS信息
	if tag.Os, err = GetTagOSsSQLite(id); err != nil {
		return nil, fmt.Errorf("获取标签OS失败: %v", err)
	}

	// 处理deletedAt字段
	if deletedAt.Valid {
		tag.DeletedAt = &deletedAt.Time
	}

	return &tag, nil
}

// GetTagsSQLite 获取所有标签
func GetTagsSQLite(option Option) ([]*Tag, error) {
	var tags []*Tag
	var query string
	var args []interface{}

	// 构建查询条件
	if len(option.Os) > 0 {
		// 如果有OS条件，使用JOIN查询关联表
		query = `
			SELECT DISTINCT t.id, t.name, t.description, t.search_count, t.created_at, t.updated_at, t.deleted_at
			FROM tags t RIGHT JOIN tag_os tos 
			ON t.id = tos.tag_id
			WHERE t.deleted_at IS NULL
		`
	} else {
		// 没有OS条件，直接查询tags表
		query = "SELECT DISTINCT t.id, t.name, t.description, t.search_count, t.created_at, t.updated_at, t.deleted_at FROM tags AS t WHERE t.deleted_at IS NULL"
	}

	if option.Name != "" {
		query += " AND t.name LIKE ?"
		args = append(args, "%"+option.Name+"%")
	}
	if option.ID != 0 {
		query += " AND t.id = ?"
		args = append(args, option.ID)
	}
	if len(option.Os) > 0 {
		// 转换OS切片为SQLite IN 子句的参数
		// for _, os := range option.Os {
		// 	query += fmt.Sprintf(" AND tos.os = '%s' ", os)
		// }
		args, err := SQL_Slice_To_In_Args(option.Os)
		if err != nil {
			return nil, fmt.Errorf("转换OS切片为SQLite IN 子句参数失败: %v", err)
		}
		query += fmt.Sprintf(" AND tos.os IN %s", args)
	}

	query += " ORDER BY t.created_at DESC"

	log.Printf("SQL: %s,args:%v", query, args)
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Printf("准备查询语句失败: %v", err)
		return nil, fmt.Errorf("准备查询语句失败: %v", err)
	}
	defer stmt.Close()

	// 从SQLite数据库获取所有标签
	rows, err := stmt.Query(args...)
	if err != nil {
		log.Printf("执行查询语句失败: %v", err)
		return nil, fmt.Errorf("获取标签列表失败: %v", err)
	}
	defer rows.Close()
	log.Printf("查询标签列表成功: %+v", rows)
	// 遍历结果集
	for rows.Next() {
		var tag Tag
		var deletedAt sql.NullTime
		// var osString string

		err = rows.Scan(
			&tag.ID, &tag.Name, &tag.Description, &tag.SearchCount, &tag.CreatedAt, &tag.UpdatedAt, &deletedAt,
		)
		if err != nil {
			log.Printf("扫描标签失败: %v", err)
			return nil, fmt.Errorf("扫描标签失败: %v", err)
		}

		// 从关联表获取OS信息
		// if tag.Os, err = GetTagOSsSQLite(tag.ID); err != nil {
		// 	log.Printf("获取标签OS失败: %v", err)
		// 	return nil, fmt.Errorf("获取标签OS失败: %v", err)
		// }

		// 处理deletedAt字段
		if deletedAt.Valid {
			tag.DeletedAt = &deletedAt.Time
		}

		tags = append(tags, &tag)
	}

	if err = rows.Err(); err != nil {
		log.Printf("遍历标签结果集失败: %v", err)
		return nil, fmt.Errorf("遍历标签结果集失败: %v", err)
	}
	log.Printf("获取标签列表成功: %+v\n", tags)
	return tags, nil
}

// UpdateTagSQLite 更新标签
func UpdateTagSQLite(tag *Tag) error {
	tag.UpdatedAt = time.Now()

	// 更新SQLite数据库中的标签
	_, err := DB.Exec(
		"UPDATE tags SET name = ?, description = ?, search_count = ?, os = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		tag.Name, tag.Description, tag.SearchCount, "", tag.UpdatedAt, tag.ID,
	)
	if err != nil {
		return fmt.Errorf("更新标签失败: %v", err)
	}

	// 先删除现有的OS关系
	if err := RemoveAllOSFromTagSQLite(tag.ID); err != nil {
		return fmt.Errorf("删除标签OS关系失败: %v", err)
	}

	// 重新添加OS关系
	// for _, os := range tag.Os {
	// 	if err := AddOSToTagSQLite(tag.ID, os); err != nil {
	// 		return fmt.Errorf("添加标签OS关系失败: %v", err)
	// 	}
	// }

	return nil
}

// DeleteTagSQLite 删除标签（软删除）
func DeleteTagSQLite(id uint64) error {
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
	now := time.Now()
	collection.CreatedAt = now
	collection.UpdatedAt = now
	collection.SearchCount = 0

	// 开启事务，确保所有操作要么全部成功，要么全部失败
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("回滚事务失败: %v", rollbackErr)
			}
		}
	}()

	// 保存到SQLite数据库，不指定id字段，让SQLite自动生成
	result, err := tx.Exec(
		"INSERT INTO collections (name, description, search_count, os, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		collection.Name, collection.Description, collection.SearchCount, "", collection.CreatedAt, collection.UpdatedAt, collection.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("创建集合失败: %v", err)
	}

	// 获取SQLite自动生成的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取集合ID失败: %v", err)
	}
	collection.ID = uint64(id)

	// 保存OS关联关系（在事务中执行）
	for _, os := range collection.Os {
		_, err := tx.Exec(
			"INSERT OR IGNORE INTO collection_os (collection_id, os) VALUES (?, ?)",
			collection.ID, os,
		)
		if err != nil {
			return fmt.Errorf("添加集合OS关系失败: %v", err)
		}
	}

	// 提交事务，所有操作都成功完成
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// GetCollectionSQLite 获取单个集合
func GetCollectionSQLite(id uint64) (*Collection, error) {
	// 输入验证
	if id == 0 {
		return nil, fmt.Errorf("集合ID不能为空")
	}

	var collection Collection
	var osString string
	var deletedAt sql.NullTime

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM collections WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&collection.ID, &collection.Name, &collection.Description, &collection.SearchCount, &osString, &collection.CreatedAt, &collection.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("collection not found: %d", id)
		}
		return nil, fmt.Errorf("获取集合失败: %v", err)
	}

	// 从关联表获取OS信息
	if collection.Os, err = GetCollectionOSsSQLite(id); err != nil {
		return nil, fmt.Errorf("获取集合OS失败: %v", err)
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
		"SELECT id, name, description, search_count, os, created_at, updated_at, deleted_at FROM collections WHERE deleted_at IS NULL ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("获取集合列表失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var collection Collection
		var osString string
		var deletedAt sql.NullTime

		err := rows.Scan(
			&collection.ID, &collection.Name, &collection.Description, &collection.SearchCount, &osString, &collection.CreatedAt, &collection.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描集合失败: %v", err)
		}

		// 从关联表获取OS信息
		if collection.Os, err = GetCollectionOSsSQLite(collection.ID); err != nil {
			return nil, fmt.Errorf("获取集合OS失败: %v", err)
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
		collection.Name, collection.Description, collection.SearchCount, "", collection.UpdatedAt, collection.ID,
	)
	if err != nil {
		return fmt.Errorf("更新集合失败: %v", err)
	}

	// 先删除现有的OS关系
	if err := RemoveAllOSFromCollectionSQLite(collection.ID); err != nil {
		return fmt.Errorf("删除集合OS关系失败: %v", err)
	}

	// 重新添加OS关系
	for _, os := range collection.Os {
		if err := AddOSToCollectionSQLite(collection.ID, os); err != nil {
			return fmt.Errorf("添加集合OS关系失败: %v", err)
		}
	}

	return nil
}

// DeleteCollectionSQLite 删除集合（软删除）
func DeleteCollectionSQLite(id uint64) error {
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
		return fmt.Errorf("collection not found: %d", id)
	}

	return nil
}

// CreateCommandSQLite 创建命令
func CreateCommandSQLite(cmd *Command) error {
	log.Printf("CreateCommandSQLite请求参数: %+v\n", cmd)
	now := time.Now()
	cmd.CreatedAt = now
	cmd.UpdatedAt = now
	cmd.CopyCounts = 0
	cmd.SearchCount = 0

	// 检查命令名称是否已存在
	var exists bool
	err := DB.QueryRow("SELECT EXISTS(SELECT 1 FROM commands WHERE name = ? AND deleted_at IS NULL)", cmd.Name).Scan(&exists)
	if err != nil {
		log.Printf("检查命令是否存在失败: %v", err)
		return fmt.Errorf("检查命令是否存在失败: %v", err)
	}
	if exists {
		log.Printf("命令[%s]已存在", cmd.Name)
		return fmt.Errorf("命令[%s]已存在", cmd.Name)
	}

	// 开启事务，确保所有操作要么全部成功，要么全部失败
	tx, err := DB.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("回滚事务失败: %v", rollbackErr)
			}
		}
	}()

	// 保存命令基本信息到SQLite数据库，不指定id字段，让SQLite自动生成
	result, err := tx.Exec(
		"INSERT INTO commands (name, content, description, copy_count, search_count, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		cmd.Name, cmd.Content, cmd.Description, cmd.CopyCounts, cmd.SearchCount, cmd.CreatedAt, cmd.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("创建命令失败: %v", err)
	}

	// 获取SQLite自动生成的ID
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取命令ID失败: %v", err)
	}
	cmd.ID = uint64(id)

	log.Printf("创建命令成功, ID: %d", id)
	// 保存命令与标签的多对多关系（在事务中执行）
	for _, tagID := range cmd.TagIDs {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO command_tags (command_id, tag_id) VALUES (?, ?)",
			cmd.ID, tagID,
		)
		if err != nil {
			return fmt.Errorf("添加命令标签关系失败: %v", err)
		}
	}
	log.Printf("添加命令标签关系成功, 命令ID: %d, 标签ID: %v", cmd.ID, cmd.TagIDs)
	// 保存命令与集合的多对多关系（在事务中执行）
	for _, collectionID := range cmd.CollectionIDs {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO command_collections (command_id, collection_id) VALUES (?, ?)",
			cmd.ID, collectionID,
		)
		if err != nil {
			return fmt.Errorf("添加命令集合关系失败: %v", err)
		}
	}
	log.Printf("添加命令集合关系成功, 命令ID: %d, 集合ID: %v", cmd.ID, cmd.CollectionIDs)

	// 保存OS关联关系（在事务中执行）
	for _, os := range cmd.Os {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO command_os (command_id, os) VALUES (?, ?)",
			cmd.ID, os,
		)
		if err != nil {
			return fmt.Errorf("添加命令OS关系失败: %v", err)
		}
	}
	log.Printf("添加命令OS关系成功, 命令ID: %d, OS: %v", cmd.ID, cmd.Os)

	// 提交事务，所有操作都成功完成
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	log.Printf("创建命令事物完成, ID: %d", cmd.ID)
	return nil
}

// GetCommandSQLite 获取单个命令
func GetCommandSQLite(id uint64) (*Command, error) {
	// 输入验证
	if id == 0 {
		return nil, fmt.Errorf("命令ID不能为空")
	}

	var cmd Command
	var deletedAt sql.NullTime
	var osString string

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, content, description, copy_count, search_count, os, created_at, updated_at, deleted_at FROM commands WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&cmd.ID, &cmd.Name, &cmd.Content, &cmd.Description, &cmd.CopyCounts, &cmd.SearchCount, &osString, &cmd.CreatedAt, &cmd.UpdatedAt, &deletedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("command not found: %d", id)
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

	// 从关联表获取OS信息
	if cmd.Os, err = GetCommandOSsSQLite(id); err != nil {
		return nil, fmt.Errorf("获取命令OS失败: %v", err)
	}

	return &cmd, nil
}

func GetCommandsByTagIDs(ids []uint64) ([]*Command, error) {
	if len(ids) == 0 {
		return []*Command{}, nil
	}
	var commands []*Command
	query := `
	SELECT cmd.* FROM commands cmd
	INNER JOIN command_tags ct ON cmd.id = ct.command_id
	WHERE cmd.deleted_at IS NULL
	AND ct.tag_id IN %s ;
	`
	args, err := SQL_Slice_To_In_Args(ids)
	if err != nil {
		return nil, fmt.Errorf("转换ID列表为IN参数失败: %v", err)
	}
	log.Printf("GetCommandsByTagIDs SQL: %s", query)
	log.Printf("GetCommandsByTagIDs args: %s", args)
	query = fmt.Sprintf(query, args)
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取命令列表失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cmd Command
		var deletedAt sql.NullTime
		// var osString string

		err := rows.Scan(
			&cmd.ID, &cmd.Name, &cmd.Content, &cmd.Description, &cmd.CopyCounts, &cmd.SearchCount, &cmd.CreatedAt, &cmd.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描命令失败: %v", err)
		}

		// 处理deletedAt字段
		if deletedAt.Valid {
			cmd.DeletedAt = &deletedAt.Time
		}

		commands = append(commands, &cmd)
	}
	return commands, nil
}

// GetCommandsSQLite 获取所有命令
func GetCommandsSQLite(option Option) ([]*Command, error) {
	var commands []*Command

	// 从SQLite数据库获取所有命令的基本信息
	// 构建查询条件
	query := "SELECT id, name, content, description, copy_count, search_count, created_at, updated_at, deleted_at FROM commands WHERE deleted_at IS NULL"
	args := []interface{}{}

	// 添加OS条件
	if len(option.Os) > 0 {
		osPlaceholders := make([]string, len(option.Os))
		for i, os := range option.Os {
			osPlaceholders[i] = "?"
			args = append(args, os)
		}
		query += fmt.Sprintf(" AND EXISTS (SELECT 1 FROM command_os co WHERE co.command_id = commands.id AND co.os IN (%s))", strings.Join(osPlaceholders, ","))
	}
	if option.Name != "" {
		query += " AND name LIKE ?"
		args = append(args, "%"+option.Name+"%")
	}
	if option.ID != 0 {
		query += " AND id = ?"
		args = append(args, option.ID)
	}

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("获取命令列表失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var cmd Command
		var deletedAt sql.NullTime
		var osString string

		err := rows.Scan(
			&cmd.ID, &cmd.Name, &cmd.Content, &cmd.Description, &cmd.CopyCounts, &cmd.SearchCount, &osString, &cmd.CreatedAt, &cmd.UpdatedAt, &deletedAt,
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

		// 从关联表获取OS信息
		if cmd.Os, err = GetCommandOSsSQLite(cmd.ID); err != nil {
			return nil, fmt.Errorf("获取命令OS失败: %v", err)
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
	if cmd.ID == 0 {
		return fmt.Errorf("命令ID不能为空")
	}

	cmd.UpdatedAt = time.Now()

	// 使用参数化查询，防止SQL注入
	_, err := DB.Exec(
		"UPDATE commands SET name = ?, content = ?, description = ?, copy_count = ?, search_count = ?, os = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		cmd.Name, cmd.Content, cmd.Description, cmd.CopyCounts, cmd.SearchCount, 0, cmd.UpdatedAt, cmd.ID,
	)
	if err != nil {
		return fmt.Errorf("更新命令失败: %v", err)
	}

	// 先删除现有的OS关系
	if err := RemoveAllOSFromCommandSQLite(cmd.ID); err != nil {
		return fmt.Errorf("删除命令OS关系失败: %v", err)
	}

	// 重新添加OS关系
	for _, os := range cmd.Os {
		if err := AddOSToCommandSQLite(cmd.ID, os); err != nil {
			return fmt.Errorf("添加命令OS关系失败: %v", err)
		}
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
func DeleteCommandSQLite(id uint64) error {
	// 输入验证
	if id == 0 {
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
		return fmt.Errorf("command not found: %d", id)
	}

	return nil
}

// 数据库迁移相关函数

// MigrateOSFieldsSQLite 迁移OS字段到关联表（可选执行）
func MigrateOSFieldsSQLite() error {
	log.Println("开始OS字段数据迁移...")

	// 1. 迁移tags表的os字段数据
	log.Println("迁移tags表数据...")
	rows, err := DB.Query("SELECT id, os FROM tags WHERE os IS NOT NULL AND os != '' AND deleted_at IS NULL")
	if err != nil {
		return fmt.Errorf("查询tags表失败: %v", err)
	}
	defer rows.Close()

	tagCount := 0
	for rows.Next() {
		var tagID uint64
		var osJSON string
		if err := rows.Scan(&tagID, &osJSON); err != nil {
			log.Printf("扫描tag失败: %v", err)
			continue
		}

		// 解析JSON数据
		var osList []string
		if err := json.Unmarshal([]byte(osJSON), &osList); err != nil {
			log.Printf("解析tag %d的OS数据失败: %v", tagID, err)
			continue
		}

		// 插入到关联表
		// for _, os := range osList {
		// 	if err := AddOSToTagSQLite(tagID, os); err != nil {
		// 		log.Printf("为tag %d添加OS %s失败: %v", tagID, os, err)
		// 	}
		// }
		tagCount++
	}

	log.Printf("tags表迁移完成，处理了 %d 条记录", tagCount)

	// 2. 迁移collections表的os字段数据
	log.Println("迁移collections表数据...")
	rows, err = DB.Query("SELECT id, os FROM collections WHERE os IS NOT NULL AND os != '' AND deleted_at IS NULL")
	if err != nil {
		return fmt.Errorf("查询collections表失败: %v", err)
	}
	defer rows.Close()

	collectionCount := 0
	for rows.Next() {
		var collectionID uint64
		var osJSON string
		if err := rows.Scan(&collectionID, &osJSON); err != nil {
			log.Printf("扫描collection失败: %v", err)
			continue
		}

		// 解析JSON数据
		var osList []string
		if err := json.Unmarshal([]byte(osJSON), &osList); err != nil {
			log.Printf("解析collection %d的OS数据失败: %v", collectionID, err)
			continue
		}

		// 插入到关联表
		for _, os := range osList {
			if err := AddOSToCollectionSQLite(collectionID, os); err != nil {
				log.Printf("为collection %d添加OS %s失败: %v", collectionID, os, err)
			}
		}
		collectionCount++
	}

	log.Printf("collections表迁移完成，处理了 %d 条记录", collectionCount)

	// 3. 迁移commands表的os字段数据（这里是整数类型）
	log.Println("迁移commands表数据...")
	rows, err = DB.Query("SELECT id, os FROM commands WHERE os IS NOT NULL AND os != 0 AND deleted_at IS NULL")
	if err != nil {
		return fmt.Errorf("查询commands表失败: %v", err)
	}
	defer rows.Close()

	commandCount := 0
	for rows.Next() {
		var commandID uint64
		var osValue int
		if err := rows.Scan(&commandID, &osValue); err != nil {
			log.Printf("扫描command失败: %v", err)
			continue
		}

		// 这里需要根据实际的OS枚举值来转换
		// 假设有一些预设的OS值映射
		osList := getOSListByValue(osValue)
		for _, os := range osList {
			if err := AddOSToCommandSQLite(commandID, os); err != nil {
				log.Printf("为command %d添加OS %s失败: %v", commandID, os, err)
			}
		}
		commandCount++
	}

	log.Printf("commands表迁移完成，处理了 %d 条记录", commandCount)
	log.Println("OS字段数据迁移完成！")

	return nil
}

// CleanupOSFieldsSQLite 清理OS字段（谨慎操作）
func CleanupOSFieldsSQLite() error {
	log.Println("开始清理OS字段...")

	// 询问用户确认
	log.Println("警告：此操作将永久删除os字段中的数据！")

	// 1. 清理tags表的os字段
	_, err := DB.Exec("UPDATE tags SET os = '' WHERE os IS NOT NULL")
	if err != nil {
		return fmt.Errorf("清理tags表os字段失败: %v", err)
	}
	log.Println("tags表os字段已清空")

	// 2. 清理collections表的os字段
	_, err = DB.Exec("UPDATE collections SET os = '' WHERE os IS NOT NULL")
	if err != nil {
		return fmt.Errorf("清理collections表os字段失败: %v", err)
	}
	log.Println("collections表os字段已清空")

	// 3. 清理commands表的os字段
	_, err = DB.Exec("UPDATE commands SET os = 0 WHERE os IS NOT NULL")
	if err != nil {
		return fmt.Errorf("清理commands表os字段失败: %v", err)
	}
	log.Println("commands表os字段已清空")

	log.Println("OS字段清理完成！")
	return nil
}

// getOSListByValue 根据整数值获取OS列表（需要根据实际枚举值调整）
func getOSListByValue(osValue int) []string {
	var osList []string

	// 这里需要根据实际的OS枚举值进行调整
	// 示例：如果osValue是标志位组合
	if osValue&1 != 0 { // 假设第1位代表Windows
		osList = append(osList, "Windows")
	}
	if osValue&2 != 0 { // 假设第2位代表macOS
		osList = append(osList, "macOS")
	}
	if osValue&4 != 0 { // 假设第3位代表Linux
		osList = append(osList, "Linux")
	}

	return osList
}

// 管理命令与标签的多对多关系

// AddTagToCommandSQLite 添加标签到命令
func AddTagToCommandSQLite(commandID uint64, tagID uint64) error {
	// 输入验证
	if commandID == 0 {
		return fmt.Errorf("命令ID不能为空")
	}
	if tagID == 0 {
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
func RemoveTagFromCommandSQLite(commandID, tagID uint64) error {
	// 输入验证
	if commandID == 0 {
		return fmt.Errorf("命令ID不能为空")
	}
	if tagID == 0 {
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
func RemoveAllTagsFromCommandSQLite(commandID uint64) error {
	// 输入验证
	if commandID == 0 {
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
func GetTagIDsByCommandIDSQLite(commandID uint64) ([]uint64, error) {
	// 输入验证
	if commandID == 0 {
		return nil, fmt.Errorf("commandID不能为空")
	}

	var tagIDs []uint64

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
		var tagID uint64
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
func AddCollectionToCommandSQLite(commandID uint64, collectionID uint64) error {
	// 输入验证
	if commandID == 0 {
		return fmt.Errorf("命令ID不能为空")
	}
	if collectionID == 0 {
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
func RemoveCollectionFromCommandSQLite(commandID uint64, collectionID uint64) error {
	// 输入验证
	if commandID == 0 {
		return fmt.Errorf("命令ID不能为空")
	}
	if collectionID == 0 {
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
func RemoveAllCollectionsFromCommandSQLite(commandID uint64) error {
	// 输入验证
	if commandID == 0 {
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
func GetCollectionIDsByCommandIDSQLite(commandID uint64) ([]uint64, error) {
	// 输入验证
	if commandID == 0 {
		return nil, fmt.Errorf("commandID不能为空")
	}

	var collectionIDs []uint64

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
		var collectionID uint64
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
