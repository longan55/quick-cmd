package main

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"time"
)

// CreateTagSQLite 创建标签
func CreateTagSQLite(tag *Tag) error {
	log.Printf("CreateTagSQLite: %+v", tag)
	now := time.Now().Format("2006-01-02 15:04:05")
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

	// 3. 保存指令关联关系
	for _, commandID := range tag.CommandIDs {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO command_tags (command_id, tag_id) VALUES (?, ?)",
			commandID, tag.ID,
		)
		if err != nil {
			log.Printf("添加标签指令关系失败: %v", err)
			tx.Rollback()
			return fmt.Errorf("添加标签指令关系失败: %v", err)
		}
	}
	log.Printf("添加标签指令关系成功: %+v", tag)

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

	// 使用参数化查询，防止SQL注入
	err := DB.QueryRow(
		"SELECT id, name, description, search_count, created_at, updated_at, deleted_at FROM tags WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(
		&tag.ID, &tag.Name, &tag.Description, &tag.SearchCount, &tag.CreatedAt, &tag.UpdatedAt, &deletedAt,
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
		tag.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
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
		args1, err := SQL_Slice_To_In_Args(option.Os)
		if err != nil {
			return nil, fmt.Errorf("转换OS切片为SQLite IN 子句参数失败: %v", err)
		}
		query += fmt.Sprintf(" AND tos.os IN %s", args1)
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
	var i int
	for rows.Next() {
		log.Printf("A%d", i)
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
		log.Printf("B%d", i)

		// 处理deletedAt字段
		if deletedAt.Valid {
			tag.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
		}
		log.Printf("C%d", i)
		// 从关联表获取指令关联关系
		if tag.ComandIdNames, err = GetCommandIDsByTagIDSQLite(tag.ID); err != nil {
			return nil, fmt.Errorf("获取标签指令关联关系失败: %v", err)
		} else {
			log.Printf("标签[%d]关联指令: %+v", tag.ID, tag.ComandIdNames)
		}
		log.Printf("D%d", i)

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

	tag.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	//0. 查询标签的OS关联关系
	rows, err := tx.Query("SELECT os FROM tag_os WHERE tag_id = ?", tag.ID)
	if err != nil {
		return fmt.Errorf("查询标签OS关联关系失败: %v", err)
	}
	defer rows.Close()
	var osList []string
	for rows.Next() {
		var os string
		if err = rows.Scan(&os); err != nil {
			return fmt.Errorf("扫描标签OS关联关系失败: %v", err)
		}
		osList = append(osList, os)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("遍历标签OS关联关系结果集失败: %v", err)
	}
	// 查询标签的指令关联关系
	rows, err = tx.Query("SELECT command_id FROM command_tags WHERE tag_id = ?", tag.ID)
	if err != nil {
		return fmt.Errorf("查询标签指令关联关系失败: %v", err)
	}
	defer rows.Close()
	var commandIDList []uint64
	for rows.Next() {
		var commandID uint64
		if err = rows.Scan(&commandID); err != nil {
			return fmt.Errorf("扫描标签指令关联关系失败: %v", err)
		}
		commandIDList = append(commandIDList, commandID)
	}
	if err = rows.Err(); err != nil {
		return fmt.Errorf("遍历标签指令关联关系结果集失败: %v", err)
	}

	// 1. 更新SQLite数据库中的标签
	_, err = tx.Exec(
		"UPDATE tags SET name = ?, description = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		tag.Name, tag.Description, tag.UpdatedAt, tag.ID,
	)
	if err != nil {
		return fmt.Errorf("更新标签失败: %v", err)
	}
	// 2. 更新标签的OS关联关系
	for _, os := range tag.Os {
		_, err = tx.Exec("INSERT OR IGNORE INTO tag_os (tag_id, os) VALUES (?, ?)", tag.ID, os)
		if err != nil {
			return fmt.Errorf("添加标签OS关联关系失败: %v", err)
		}
	}
	// 删除标签的OS关联关系中不在tag.Os中的OS
	for _, os := range osList {
		if !slices.Contains(tag.Os, os) {
			_, err = tx.Exec("DELETE FROM tag_os WHERE tag_id = ? AND os = ?", tag.ID, os)
			if err != nil {
				return fmt.Errorf("删除标签OS关联关系失败: %v", err)
			}
		}
	}
	//3. 更新标签的指令关联关系
	for _, commandID := range tag.CommandIDs {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO command_tags (command_id, tag_id) VALUES (?, ?)", commandID, tag.ID)
		if err != nil {
			return fmt.Errorf("添加标签指令关系失败: %v", err)
		}
	}
	// 删除标签的指令关联关系中不在tag.CommandIDs中的指令
	for _, commandID := range commandIDList {
		if !slices.Contains(tag.CommandIDs, commandID) {
			_, err = tx.Exec("DELETE FROM command_tags WHERE tag_id = ? AND command_id = ?", tag.ID, commandID)
			if err != nil {
				return fmt.Errorf("删除标签指令关联关系失败: %v", err)
			}
		}
	}
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	return nil
}

// DeleteTagSQLite 删除标签（软删除）
func DeleteTagSQLite(id uint64) error {
	// 输入验证
	if id == 0 {
		return fmt.Errorf("标签ID不能为空")
	}

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

	now := time.Now()
	// 1.更新标签的deleted_at字段
	_, err = tx.Exec(
		"UPDATE tags SET deleted_at = ?, updated_at = ? WHERE id = ? AND deleted_at IS NULL",
		now, now, id,
	)
	if err != nil {
		return fmt.Errorf("删除标签失败: %v", err)
	}
	//2.删除标签的OS关联关系
	_, err = tx.Exec(
		"DELETE FROM tag_os WHERE tag_id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("删除标签OS关联关系失败: %v", err)
	}
	//3.删除标签的指令关联关系
	_, err = tx.Exec(
		"DELETE FROM command_tags WHERE tag_id = ?",
		id,
	)
	if err != nil {
		return fmt.Errorf("删除标签指令关联关系失败: %v", err)
	}
	//4.提交事务
	return tx.Commit()
}

func GetTagIDAndNameSQLite() ([]Tag, error) {
	var tags []Tag
	query := "SELECT DISTINCT id, name FROM tags WHERE id IS NOT NULL AND deleted_at IS NULL"
	log.Printf("GetTagIDAndNameSQLite SQL: %s", query)

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("查询标签ID和名称失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var name string
		if err = rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("扫描标签ID和名称失败: %v", err)
		}
		tags = append(tags, Tag{ID: id, Name: name})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历标签ID和名称结果集失败: %v", err)
	}

	return tags, nil
}

func GetCommandIDsByTagIDSQLite(tagID uint64) ([]CommandIDName, error) {
	log.Printf("GetCommandIDsByTagIDSQLite tagID: %d", tagID)
	var commandIDs []CommandIDName
	query := "SELECT command_id, name FROM command_tags ct LEFT JOIN commands cmd ON ct.command_id = cmd.id WHERE tag_id = ?"
	log.Printf("GetCommandIDsByTagIDSQLite SQL: %s, tagID: %d", query, tagID)

	rows, err := DB.Query(query, tagID)
	if err != nil {
		log.Printf("查询标签指令关联关系失败: %v", err)
		return nil, fmt.Errorf("查询标签指令关联关系失败: %v", err)
	}
	defer rows.Close()
	log.Printf("查询标签指令关联关系成功: %+v", rows)
	for rows.Next() {
		var commandID uint64
		var name string
		if err = rows.Scan(&commandID, &name); err != nil {
			return nil, fmt.Errorf("扫描标签指令关联关系失败: %v", err)
		}
		commandIDs = append(commandIDs, CommandIDName{ID: commandID, Name: name})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历标签指令关联关系结果集失败: %v", err)
	}
	log.Printf("标签[%d]关联指令: %+v", tagID, commandIDs)
	return commandIDs, nil
}
