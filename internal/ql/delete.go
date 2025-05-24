package ql

import (
	"strconv"
)

// Deleter delete 语句构造器
type Deleter struct {
	sqlBuilder
	tableName string
	where     ConditionBuilder
	limit     *int
}

// NewDeleter
// tableName 数据库表名
func NewDeleter(tableName string) *Deleter {
	return &Deleter{
		tableName: tableName,
	}
}

// Where 条件
// condition 可以通过 ql.C() 方法创建
func (d *Deleter) Where(where ConditionBuilder) *Deleter {
	d.where = where
	return d
}

// Limit 限制删除数量
func (d *Deleter) Limit(limit int) *Deleter {
	d.limit = &limit
	return d
}

// SQL 输出sql语句
func (d *Deleter) SQL() (string, error) {
	if d.where == nil || len(d.where.getPredicates()) == 0 {
		return "", ErrDeleteMissWhere
	}
	d.reset()
	d.writeString("DELETE FROM ")
	d.quote(d.tableName)
	d.whereSQL(d.where)
	if d.limit != nil {
		d.writeString(" LIMIT ")
		d.writeString(strconv.Itoa(*d.limit))
	}
	d.end()
	return d.sb.String(), nil
}

// SQLArgs 构造 sql 并返回对应参数
func (d *Deleter) SQLArgs() (string, []any, error) {
	if d.where == nil || len(d.where.getPredicates()) == 0 {
		return "", nil, ErrDeleteMissWhere
	}
	execSQL, err := d.SQL()
	if err != nil {
		return "", nil, err
	}
	args, hasInSQL := d.whereArgs(d.where)
	if !hasInSQL {
		return execSQL, args, err
	}
	return parseIn(execSQL, args...)
}
