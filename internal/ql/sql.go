package ql

import (
	"strconv"
	"strings"
)

type Builder struct {
	tableName string
}

// New 创建 sql builder
func New(tableName string) *Builder {
	builder := &Builder{
		tableName: tableName,
	}
	return builder
}

// Select 创建 select 语句构造器
func (b *Builder) Select(columns ...string) *Selector {
	return NewSelector(b.tableName).Columns(columns...)
}

// Insert 创建 insert 语句构造器
func (b *Builder) Insert(columns ...string) *Inserter {
	return NewInserter(b.tableName).Columns(columns...)
}

// Update 创建 update 语句构造器
func (b *Builder) Update(columns ...string) *Updater {
	return NewUpdater(b.tableName).Columns(columns...)
}

// Delete 创建 delete 语句构造器
func (b *Builder) Delete() *Deleter {
	return NewDeleter(b.tableName)
}

type sqlBuilder struct {
	sb strings.Builder
}

func (b *sqlBuilder) reset() {
	b.sb.Reset()
}

func (b *sqlBuilder) writeString(val string) {
	_, _ = b.sb.WriteString(val)
}

func (b *sqlBuilder) writeByte(c byte) {
	_ = b.sb.WriteByte(c)
}

func (b *sqlBuilder) quote(val string) {
	b.writeByte('`')
	b.writeString(strings.TrimSpace(val))
	b.writeByte('`')
}

func (b *sqlBuilder) col(col column) {
	if col.alias != "" {
		b.writeString(col.alias)
		b.writeByte('.')
	}
	b.writeByte('`')
	b.writeString(strings.TrimSpace(col.name))
	b.writeByte('`')
}

func (b *sqlBuilder) ifNullCol(col column, val string) {
	b.writeString("IFNULL(")
	b.col(col)
	b.writeString(", '")
	b.writeString(val)
	b.writeString("') as ")
	b.quote(col.name)
}

func (b *sqlBuilder) space() {
	b.writeByte(' ')
}

func (b *sqlBuilder) end() {
	b.writeByte(';')
}

func (b *sqlBuilder) comma() {
	b.writeByte(',')
}

// whereSQL 拼接 where 条件
func (b *sqlBuilder) whereSQL(where ConditionBuilder) {
	if where != nil && len(where.getPredicates()) > 0 {
		b.writeString(" WHERE ")
		for i, predicate := range where.getPredicates() {
			if i > 0 {
				b.writeString(predicate.Op.Text)
			}
			b.writeString(predicate.Express)
		}
	}
}

// whereArgs where 条件中的参数
func (b *sqlBuilder) whereArgs(where ConditionBuilder) (args []any, hasInSQL bool) {
	if where != nil && len(where.getPredicates()) > 0 {
		b.writeString(" WHERE ")
		for _, predicate := range where.getPredicates() {
			args = append(args, predicate.Args...)
			if predicate.HasInSQL {
				hasInSQL = true
			}
		}
	}
	return
}

// setFields 字段赋值语句
func (b *sqlBuilder) setFields(fields []Field) {
	n := len(fields)
	for i, f := range fields {
		b.quote(f.col)
		b.writeString(" = ")
		if f.incrVal != nil {
			b.quote(f.col)
			b.writeString(" + ")
			b.writeString(strconv.FormatInt(*f.incrVal, 10))
		} else {
			b.writeString("?")
		}
		if i != n-1 {
			b.writeString(", ")
		}
	}
}

// setFields 字段赋值语句
func (b *sqlBuilder) setNameFields(fields []Field) {
	n := len(fields)
	for i, f := range fields {
		b.quote(f.col)
		b.writeString(" = ")
		if f.incrVal != nil {
			b.quote(f.col)
			b.writeString(" + ")
			b.writeString(strconv.FormatInt(*f.incrVal, 10))
		} else {
			b.writeString(":")
			b.writeString(f.col)
		}
		if i != n-1 {
			b.writeString(", ")
		}
	}
}
