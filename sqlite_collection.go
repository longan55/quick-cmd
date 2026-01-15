package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// CreateCollectionSQLite 创建集合
func CreateCollectionSQLite(collection *Collection) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	collection.CreatedAt = now
	collection.UpdatedAt = now
	collection.SearchCount = 0
	log.Printf("创建集合: %+v", collection)
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
		"INSERT INTO collections (name, description, search_count, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?)",
		collection.Name, collection.Description, collection.SearchCount, collection.CreatedAt, collection.UpdatedAt, collection.DeletedAt,
	)
	if err != nil {
		log.Printf("创建集合失败: %v", err)
		return fmt.Errorf("创建集合失败: %v", err)
	}
	log.Printf("创建集合成功: %+v", collection)
	// 获取SQLite自动生成的ID
	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("获取集合ID失败: %v", err)
		return fmt.Errorf("获取集合ID失败: %v", err)
	}
	collection.ID = uint64(id)

	// 保存OS关联关系（在事务中执行）
	for _, os := range collection.Os {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO collection_os (collection_id, os) VALUES (?, ?)",
			collection.ID, os,
		)
		if err != nil {
			log.Printf("添加集合OS关系失败: %v", err)
			return fmt.Errorf("添加集合OS关系失败: %v", err)
		}
	}
	log.Printf("添加集合OS关系成功: %+v", collection)

	// 保存指令关联关系
	for _, commandID := range collection.CommandIDs {
		_, err = tx.Exec(
			"INSERT OR IGNORE INTO command_collections (command_id, collection_id) VALUES (?, ?)",
			commandID, collection.ID,
		)
		if err != nil {
			log.Printf("添加集合指令关系失败: %v", err)
			return fmt.Errorf("添加集合指令关系失败: %v", err)
		}
	}
	log.Printf("添加集合指令关系成功: %+v", collection)

	// 提交事务，所有操作都成功完成
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	log.Printf("提交事务成功: %+v", collection)
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
		collection.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
	}

	return &collection, nil
}

// GetCollectionsSQLite 获取所有集合
func GetCollectionsSQLite(option Option) ([]*Collection, error) {
	var collections []*Collection
	query := `SELECT id, name, description, search_count, created_at, updated_at, deleted_at 
	FROM collections c RIGHT JOIN collection_os cs 
	ON c.id = cs.collection_id
	WHERE id IS NOT NULL AND deleted_at IS NULL`
	if option.Name != "" {
		query += fmt.Sprintf(" AND name LIKE '%%%s%%'", option.Name)
	}
	if option.ID != 0 {
		query += fmt.Sprintf(" AND id = %d", option.ID)
	}
	if len(option.Os) > 0 {
		args, err := SQL_Slice_To_In_Args(option.Os)
		if err != nil {
			return nil, fmt.Errorf("转换OS列表为IN参数失败: %v", err)
		}
		query += fmt.Sprintf(" AND os IN %s", args)
	}
	query += ";"
	log.Printf("GetCollectionsSQLite SQL: %s", query)

	// 从SQLite数据库获取所有集合
	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取集合列表失败: %v", err)
	}
	defer rows.Close()

	// 遍历结果集
	for rows.Next() {
		var collection Collection
		var deletedAt sql.NullTime

		err = rows.Scan(
			&collection.ID, &collection.Name, &collection.Description, &collection.SearchCount, &collection.CreatedAt, &collection.UpdatedAt, &deletedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("扫描集合失败: %v", err)
		}

		// 处理deletedAt字段
		if deletedAt.Valid {
			collection.DeletedAt = deletedAt.Time.Format("2006-01-02 15:04:05")
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
	collection.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

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

func GetCollectionIDAndNameSQLite() ([]Collection, error) {
	var collections []Collection
	query := "SELECT DISTINCT id, name FROM collections WHERE id IS NOT NULL AND deleted_at IS NULL"
	log.Printf("GetCollectionIDAndNameSQLite SQL: %s", query)

	rows, err := DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("获取集合ID和名称失败: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uint64
		var name string
		if err = rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("扫描集合ID和名称失败: %v", err)
		}
		collections = append(collections, Collection{ID: id, Name: name})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("遍历集合ID和名称结果集失败: %v", err)
	}

	return collections, nil
}
