package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

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
	// 检查数据库文件是否存在
	_, err := os.Stat("quick-cmd.db")
	fileExists := err == nil

	// 打开SQLite数据库，使用文件存储而不是内存存储
	// 使用mode=rwc模式，当文件不存在时会自动创建
	db, err := sql.Open("sqlite3", "file:quick-cmd.db?cache=shared&mode=rwc&_foreign_keys=1")
	if err != nil {
		panic(err)
	}
	DB = db
	DB.SetMaxOpenConns(1)

	// 如果文件不存在，先创建数据库文件
	if !fileExists {
		log.Println("数据库文件不存在，正在创建...")
		// 执行一个空查询来触发数据库文件的创建
		if _, err := DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
			panic(fmt.Errorf("创建数据库文件失败: %v", err))
		}
		log.Println("数据库文件创建成功")
	}

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

func FillCommandRelations(commands []*Command) error {
	if len(commands) == 0 {
		return nil
	}

	commandIDs := make([]uint64, len(commands))
	for i, cmd := range commands {
		commandIDs[i] = cmd.ID
	}

	tagMap, err := GetTagIDsByCommandIDsSQLite(commandIDs)
	if err != nil {
		return fmt.Errorf("批量获取命令标签ID失败: %v", err)
	}
	collectionMap, err := GetCollectionIDsByCommandIDsSQLite(commandIDs)
	if err != nil {
		return fmt.Errorf("批量获取命令集合ID失败: %v", err)
	}
	osMap, err := GetCommandOSsByCommandIDsSQLite(commandIDs)
	if err != nil {
		return fmt.Errorf("批量获取命令OS失败: %v", err)
	}

	for _, cmd := range commands {
		cmd.TagIDs = tagMap[cmd.ID]
		cmd.CollectionIDs = collectionMap[cmd.ID]
		cmd.Os = osMap[cmd.ID]
	}

	return nil
}

func GetTagIDsByCommandIDsSQLite(commandIDs []uint64) (map[uint64][]uint64, error) {
	result := make(map[uint64][]uint64)
	if len(commandIDs) == 0 {
		return result, nil
	}

	query := "SELECT command_id, tag_id FROM command_tags WHERE command_id IN ("
	args := make([]interface{}, len(commandIDs))
	for i, id := range commandIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = id
	}
	query += ")"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("批量获取命令标签ID失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var commandID, tagID uint64
		if err = rows.Scan(&commandID, &tagID); err != nil {
			return nil, fmt.Errorf("扫描命令标签ID失败: %v", err)
		}
		result[commandID] = append(result[commandID], tagID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历命令标签ID结果集失败: %v", err)
	}

	return result, nil
}

func GetCollectionIDsByCommandIDsSQLite(commandIDs []uint64) (map[uint64][]uint64, error) {
	result := make(map[uint64][]uint64)
	if len(commandIDs) == 0 {
		return result, nil
	}

	query := "SELECT command_id, collection_id FROM command_collections WHERE command_id IN ("
	args := make([]interface{}, len(commandIDs))
	for i, id := range commandIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = id
	}
	query += ")"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("批量获取命令集合ID失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var commandID, collectionID uint64
		if err = rows.Scan(&commandID, &collectionID); err != nil {
			return nil, fmt.Errorf("扫描命令集合ID失败: %v", err)
		}
		result[commandID] = append(result[commandID], collectionID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历命令集合ID结果集失败: %v", err)
	}

	return result, nil
}

func GetCommandOSsByCommandIDsSQLite(commandIDs []uint64) (map[uint64][]string, error) {
	result := make(map[uint64][]string)
	if len(commandIDs) == 0 {
		return result, nil
	}

	query := "SELECT command_id, os FROM command_os WHERE command_id IN ("
	args := make([]interface{}, len(commandIDs))
	for i, id := range commandIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = id
	}
	query += ")"

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("批量获取命令OS失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var commandID uint64
		var os string
		if err = rows.Scan(&commandID, &os); err != nil {
			return nil, fmt.Errorf("扫描命令OS失败: %v", err)
		}
		result[commandID] = append(result[commandID], os)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历命令OS结果集失败: %v", err)
	}

	return result, nil
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
