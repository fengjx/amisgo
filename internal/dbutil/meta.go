package dbutil

import (
	"context"
	"strings"

	"github.com/fengjx/amisgo/internal/kit"
)

var (
	// 判断 columnType 是否是时间字段
	timeTypes = []string{"DATE", "DATETIME", "TIMESTAMP", "TIME", "YEAR"}
)

// DBTable 数据库表信息
type DBTable struct {
	Name          string
	StructName    string
	Columns       []*DBCol
	PrimaryKey    *DBCol
	AutoIncrement bool
	Comment       string
	StoreEngine   string
}

// GetColumn 获取字段
func (t *DBTable) GetColumn(name string) *DBCol {
	if t == nil {
		return nil
	}
	for _, col := range t.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

// DBCol 数据库字段信息
type DBCol struct {
	Name         string
	SQLType      string
	Comment      string
	IsPrimaryKey bool
	IsTimeType   bool
	DefaultValue string
	Extra        string
}

// GetComment 获取字段注释
func (c *DBCol) GetComment() string {
	if c.Comment == "" {
		return ""
	}
	return c.Comment
}

// GetTableMeta 获取表的元数据
func GetTableMeta(ctx context.Context, executor Executor, dbName, tableName string) (*DBTable, error) {
	args := []any{dbName, tableName}
	querySQL := "SELECT `TABLE_NAME`, `ENGINE`, `AUTO_INCREMENT`, `TABLE_COMMENT` from" +
		" `INFORMATION_SCHEMA`.`TABLES` WHERE `TABLE_SCHEMA`=? AND TABLE_NAME = ?" +
		" AND (`ENGINE`='MyISAM' OR `ENGINE` = 'InnoDB' OR `ENGINE` = 'TokuDB')"

	rows, err := executor.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	table := new(DBTable)
	for rows.Next() {
		var name, engine string
		var comment *string
		var autoIncr *int
		err = rows.Scan(&name, &engine, &autoIncr, &comment)
		if err != nil {
			return nil, err
		}
		table.Name = name
		table.StructName = kit.GonicCase(name)
		if comment != nil {
			table.Comment = *comment
		}
		if table.Comment == "" {
			table.Comment = table.Name
		}
		table.StoreEngine = engine
		if autoIncr != nil {
			table.AutoIncrement = true
		}
	}
	if rows.Err() != nil {
		return nil, err
	}
	columns, primaryKey, err := loadColumnMeta(ctx, executor, dbName, tableName)
	if err != nil {
		return nil, err
	}
	table.Columns = columns
	table.PrimaryKey = primaryKey
	return table, nil
}

// loadColumnMeta
// []*DBCol table column meta
// *DBCol PrimaryKey column
func loadColumnMeta(ctx context.Context, executor Executor, dbName, tableName string) ([]*DBCol, *DBCol, error) {
	args := []any{dbName, tableName}
	querySQL := `SELECT
			column_name,
			column_type,
			column_comment, 
			column_key,
			ifnull(column_default, '') as col_default,
			ifnull(extra, '') as extra
		FROM information_schema.columns
		WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION`
	rows, err := executor.QueryContext(ctx, querySQL, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()
	var columns []*DBCol
	var primaryKey *DBCol
	for rows.Next() {
		var columnName string
		var columnType string
		var columnComment string
		var columnKey string
		var columnDefault string
		var extra string
		err = rows.Scan(
			&columnName,
			&columnType,
			&columnComment,
			&columnKey,
			&columnDefault,
			&extra,
		)
		if err != nil {
			return nil, nil, err
		}
		col := &DBCol{
			Name:         strings.Trim(columnName, "` "),
			Comment:      columnComment,
			DefaultValue: columnDefault,
			Extra:        extra,
		}
		if col.Comment == "" {
			col.Comment = col.Name
		}
		fields := strings.Fields(columnType)
		columnType = fields[0]
		cts := strings.Split(columnType, "(")
		colName := cts[0]
		// Remove the /* mariadb-5.3 */ suffix from coltypes
		colName = strings.TrimSuffix(colName, "/* mariadb-5.3 */")
		col.SQLType = strings.ToUpper(colName)

		if kit.Contains(timeTypes, col.SQLType) {
			col.IsTimeType = true
		}

		if columnKey == "PRI" {
			col.IsPrimaryKey = true
			primaryKey = col
		}
		columns = append(columns, col)
	}
	return columns, primaryKey, nil
}
