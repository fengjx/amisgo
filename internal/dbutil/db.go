package dbutil

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fengjx/amisgo/internal/kit"
	"github.com/fengjx/amisgo/internal/log"
	"github.com/fengjx/amisgo/internal/ql"
)

// Exec sql 执行器
type Exec interface {
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, execSQL string, args ...any) (sql.Result, error)
}

// Executor sql 执行器，支持事务
type Executor interface {
	Exec
	// Begin 开启事务
	Begin() (*sql.Tx, error)
}

// DoInTx 事务执行函数
type DoInTx func(ctx context.Context, tx *sql.Tx) error

// ExecTx 执行事务
func ExecTx(ctx context.Context, executor Executor, fn DoInTx) error {
	tx, err := executor.Begin()
	if err != nil {
		return err
	}
	err = fn(ctx, tx)
	if err != nil {
		err1 := tx.Rollback()
		if err1 != nil {
			log.Errorf("rollback error: %v", err1)
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Errorf("commit error: %v", err)
	}
	return err
}

// RowData 行数据
type RowData map[string]any

// Insert 插入数据
func Insert(ctx context.Context, e Exec, table *DBTable, data RowData) (sql.Result, error) {
	query, args, err := ql.NewInserter(table.Name).SetMap(data).SQLArgs()
	if err != nil {
		return nil, err
	}
	return e.ExecContext(ctx, query, args...)
}

// UpdateByID 更新数据
func UpdateByID(ctx context.Context, e Exec, table *DBTable, id any, data RowData) (sql.Result, error) {
	query, args, err := ql.NewUpdater(table.Name).SetMap(data).Where(
		ql.C(ql.Col(table.PrimaryKey.Name).EQ(id)),
	).SQLArgs()
	if err != nil {
		return nil, err
	}
	return e.ExecContext(ctx, query, args...)
}

// BatchUpdate 批量更新数据
func BatchUpdate(ctx context.Context, executor Executor, table *DBTable, rows []RowData) error {
	return ExecTx(ctx, executor, func(ctx context.Context, tx *sql.Tx) error {
		for _, r := range rows {
			id := r[table.PrimaryKey.Name]
			if kit.IsEmptyOrDefault(id) {
				return errors.New("row miss primary key")
			}
			_, err := UpdateByID(ctx, tx, table, id, r)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// DeleteByIDs 删除数据
func DeleteByIDs(ctx context.Context, e Exec, table *DBTable, ids ...any) (sql.Result, error) {
	execSQL, args, err := ql.NewDeleter(table.Name).Where(
		ql.C(ql.Col(table.PrimaryKey.Name).In(ids...)),
	).SQLArgs()
	if err != nil {
		return nil, err
	}
	return e.ExecContext(ctx, execSQL, args...)
}

// GetByID 查询单条数据
func GetByID(ctx context.Context, e Exec, table *DBTable, cols []string, id int64) (RowData, error) {
	if len(cols) == 0 {
		for _, col := range table.Columns {
			cols = append(cols, col.Name)
		}
	}
	querySQL, args, err := ql.NewSelector(table.Name).
		Columns(cols...).
		Where(ql.C(ql.Col(table.PrimaryKey.Name).EQ(id))).
		SQLArgs()
	if err != nil {
		return nil, err
	}
	return Get(ctx, e, querySQL, args...)
}

// QueryList 查询数据
func QueryList(ctx context.Context, e Exec, table *DBTable, cols []string, where ql.ConditionBuilder) ([]RowData, error) {
	if len(cols) == 0 {
		for _, col := range table.Columns {
			cols = append(cols, col.Name)
		}
	}
	querySQL, args, err := ql.NewSelector(table.Name).
		Columns(cols...).
		Where(where).
		SQLArgs()
	if err != nil {
		return nil, err
	}
	return List(ctx, e, querySQL, args...)
}

// QueryOne 查询单条数据
func QueryOne(ctx context.Context, e Exec, table *DBTable, cols []string, where ql.ConditionBuilder) (map[string]any, error) {
	if len(cols) == 0 {
		for _, col := range table.Columns {
			cols = append(cols, col.Name)
		}
	}
	querySQL, args, err := ql.NewSelector(table.Name).
		Columns(cols...).
		Where(where).
		SQLArgs()
	if err != nil {
		return nil, err
	}
	return Get(ctx, e, querySQL, args...)
}

// QueryCount 查询数据总数
func QueryCount(ctx context.Context, e Exec, table *DBTable, where ql.ConditionBuilder) (int64, error) {
	querySQL, args, err := ql.NewSelector(table.Name).
		Columns("COUNT(*)").
		Where(where).
		SQLArgs()
	if err != nil {
		return 0, err
	}
	return Count(ctx, e, querySQL, args...)
}

func List(ctx context.Context, e Exec, querySQL string, args ...any) ([]RowData, error) {
	rows, err := e.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToList(rows)
}

// Get 查询单条数据
func Get(ctx context.Context, e Exec, querySQL string, args ...any) (RowData, error) {
	rows, err := e.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return rowsToData(rows)
}

// Count 查询数据总数
func Count(ctx context.Context, e Exec, querySQL string, args ...any) (int64, error) {
	rows, err := e.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	var cnt int64 = 0
	if rows.Next() {
		_ = rows.Scan(&cnt)
	}
	return cnt, nil
}

func rowsToData(r *sql.Rows) (RowData, error) {
	// 获取列名
	columns, err := r.Columns()
	if err != nil {
		return nil, err
	}

	// 获取列类型
	colTypes, err := r.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// 为每一列创建一个 any 切片来存储数据
	values := make([]any, len(colTypes))
	for i := range colTypes {
		values[i] = new(any)
	}

	// 创建结果 map
	result := make(map[string]any)

	// 只获取第一行数据
	if r.Next() {
		// 扫描当前行到 values 中
		if err := r.Scan(values...); err != nil {
			return nil, err
		}

		// 将值存储到结果 map 中
		for i, col := range columns {
			if v := *(values[i].(*any)); v != nil {
				result[col] = v
			} else {
				result[col] = nil
			}
		}
	}

	// 检查遍历过程中是否有错误
	if err = r.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func rowsToList(r *sql.Rows) ([]RowData, error) {
	// 获取列名
	columns, err := r.Columns()
	if err != nil {
		return nil, err
	}

	// 获取列类型
	colTypes, err := r.ColumnTypes()
	if err != nil {
		return nil, err
	}

	// 创建结果切片
	var result []RowData

	// 为每一列创建一个 any 切片来存储数据
	values := make([]any, len(colTypes))
	for i := range colTypes {
		values[i] = new(any)
	}

	// 遍历每一行数据
	for r.Next() {
		// 扫描当前行到 values 中
		if err := r.Scan(values...); err != nil {
			return nil, err
		}

		// 创建当前行的 map
		row := make(RowData)
		for i, col := range columns {
			if val := *(values[i].(*any)); val != nil {
				switch v := val.(type) {
				case []byte:
					row[col] = string(v)
				default:
					row[col] = v
				}
			} else {
				row[col] = nil
			}
		}
		result = append(result, row)
	}

	// 检查遍历过程中是否有错误
	if err = r.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
