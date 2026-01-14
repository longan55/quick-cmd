package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateCommandSQLite 创建命令
func CreateCommandSQLite(cmd *Command) error {
	log.Printf("CreateCommandSQLite请求参数: %+v\n", cmd)
	now := time.Now().Format("2006-01-02 15:04:05")
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

	log.Printf("开始添加命令集合关系, 命令ID: %d, 集合ID: %v", cmd.ID, cmd.CollectionIDs)
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
		cmd.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
	}

	// // 获取命令关联的标签ID
	// if cmd.TagIDs, err = GetTagIDsByCommandIDSQLite(id); err != nil {
	// 	return nil, fmt.Errorf("获取命令标签ID失败: %v", err)
	// }

	// // 获取命令关联的集合ID
	// if cmd.CollectionIDs, err = GetCollectionIDsByCommandIDSQLite(id); err != nil {
	// 	return nil, fmt.Errorf("获取命令集合ID失败: %v", err)
	// }

	// // 从关联表获取OS信息
	// if cmd.Os, err = GetCommandOSsSQLite(id); err != nil {
	// 	return nil, fmt.Errorf("获取命令OS失败: %v", err)
	// }

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
			cmd.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
		}

		commands = append(commands, &cmd)
	}
	return commands, nil
}

func GetCommandByCollectionIds(ids []uint64) ([]*Command, error) {
	if len(ids) == 0 {
		return []*Command{}, nil
	}
	var commands []*Command
	query := `
	SELECT cmd.* FROM commands cmd
	INNER JOIN command_collections ct ON cmd.id = ct.command_id
	WHERE cmd.deleted_at IS NULL
	AND ct.collection_id IN %s ;
	`
	args, err := SQL_Slice_To_In_Args(ids)
	if err != nil {
		return nil, fmt.Errorf("转换ID列表为IN参数失败: %v", err)
	}
	log.Printf("GetCommandByCollectionIds SQL: %s", query)
	log.Printf("GetCommandByCollectionIds args: %s", args)
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
		cmd.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
	}

		commands = append(commands, &cmd)
	}
	return commands, nil
}

// GetCommandsSQLite 获取所有命令
func GetCommandsSQLite(option Option) ([]*Command, error) {
	var commands []*Command

	query := "SELECT DISTINCT id, name, content, description, copy_count, search_count, created_at, updated_at, deleted_at FROM commands c JOIN command_os cs WHERE deleted_at IS NULL"
	// args := []interface{}{}

	if option.Name != "" {
		query += fmt.Sprintf(" AND name LIKE '%%%s%%'", option.Name)
		// args = append(args, "%"+option.Name+"%")
	}
	if option.ID != 0 {
		query += fmt.Sprintf(" AND id = %d", option.ID)
		// args = append(args, option.ID)
	}
	if len(option.Os) > 0 {
		args, err := SQL_Slice_To_In_Args(option.Os)
		if err != nil {
			return nil, fmt.Errorf("转换OS列表为IN参数失败: %v", err)
		}
		query += fmt.Sprintf(" AND cs.os IN %s", args)
	}
	query += ";"

	log.Printf("GetCommandsSQLite SQL: %s", query)
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取命令列表失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cmd Command
		var deletedAt sql.NullTime

		err = rows.Scan(
			&cmd.ID, &cmd.Name, &cmd.Content, &cmd.Description, &cmd.CopyCounts, &cmd.SearchCount, &cmd.CreatedAt, &cmd.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描命令失败: %v", err)
		}

		if deletedAt.Valid {
		cmd.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
	}

		commands = append(commands, &cmd)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历命令结果集失败: %v", err)
	}

	if err = FillCommandRelations(commands); err != nil {
		return nil, fmt.Errorf("填充命令关联数据失败: %v", err)
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

	cmd.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

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

// 获取所有命令的id和name
func GetAllCommandsIDAndNameSQLite() ([]*Command, error) {
	var commands []*Command

	query := "SELECT id, name FROM commands WHERE deleted_at IS NULL"
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取命令列表失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cmd Command

		err = rows.Scan(
			&cmd.ID, &cmd.Name,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描命令失败: %v", err)
		}

		commands = append(commands, &cmd)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历命令结果集失败: %v", err)
	}

	return commands, nil
}
