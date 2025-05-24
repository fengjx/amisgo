package ql

import (
	"strings"
)

type intoType string

const (
	intoTypeDefault intoType = "default"
	intoTypeReplace intoType = "replace"
	intoTypeIgnore  intoType = "ignore"
)

// Inserter insert 语句构造器
type Inserter struct {
	sqlBuilder
	tableName                  string
	insertFields               []Field // insert 字段
	duplicateFields            []Field // OnDuplicateKeyUpdate 更新的字段
	onDuplicateKeyUpdateString string
	intoType                   intoType
}

// NewInserter 创建 insert 语句构造器
func NewInserter(tableName string) *Inserter {
	inserter := &Inserter{
		tableName: tableName,
		intoType:  intoTypeDefault,
	}
	return inserter
}

// Columns insert 字段
func (ins *Inserter) Columns(columns ...string) *Inserter {
	for _, col := range columns {
		ins.insertFields = append(ins.insertFields, F(col))
	}
	return ins
}

// Fields 设置 insert 字段和值
func (ins *Inserter) Fields(fields ...Field) *Inserter {
	for _, field := range fields {
		if !field.isUse {
			continue
		}
		ins.insertFields = append(ins.insertFields, field)
	}
	return ins
}

// SetMap 设置 insert 字段和值
func (ins *Inserter) SetMap(data map[string]any) *Inserter {
	for k, v := range data {
		ins.insertFields = append(ins.insertFields, F(k).Val(v))
	}
	return ins
}

// IsReplaceInto 是否使用  replace into
func (ins *Inserter) IsReplaceInto(replaceInto bool) *Inserter {
	if replaceInto {
		ins.intoType = intoTypeReplace
	}
	return ins
}

// IsIgnoreInto 是否使用  ignore into
func (ins *Inserter) IsIgnoreInto(ignoreInto bool) *Inserter {
	if ignoreInto {
		ins.intoType = intoTypeIgnore
	}
	return ins
}

// OnDuplicateKeyUpdateString 设置 on duplicate key update 字段
func (ins *Inserter) OnDuplicateKeyUpdateString(updateString string) *Inserter {
	ins.onDuplicateKeyUpdateString = updateString
	return ins
}

// OnDuplicateKeyUpdate 设置 on duplicate key update 字段
func (ins *Inserter) OnDuplicateKeyUpdate(fields ...Field) *Inserter {
	for _, field := range fields {
		if !field.isUse {
			continue
		}
		ins.duplicateFields = append(ins.duplicateFields, field)
	}
	return ins
}

// NameSQL 返回名称风格的 sql
func (ins *Inserter) NameSQL() (string, error) {
	if len(ins.insertFields) == 0 {
		return "", ErrColumnsRequire
	}
	ins.reset()
	switch ins.intoType {
	case intoTypeDefault:
		ins.writeString("INSERT INTO ")
	case intoTypeReplace:
		ins.writeString("REPLACE INTO ")
	case intoTypeIgnore:
		ins.writeString("INSERT IGNORE INTO ")
	}
	ins.quote(strings.TrimSpace(ins.tableName))
	ins.writeByte('(')
	for i, f := range ins.insertFields {
		ins.quote(f.col)
		if i != len(ins.insertFields)-1 {
			ins.writeString(", ")
		}
	}
	ins.writeByte(')')
	ins.writeString(" VALUES (")
	for i, f := range ins.insertFields {
		ins.writeByte(':')
		ins.writeString(f.col)
		if i != len(ins.insertFields)-1 {
			ins.writeString(", ")
		}
	}
	ins.writeString(")")
	if ins.onDuplicateKeyUpdateString != "" {
		ins.writeString(" ON DUPLICATE KEY UPDATE ")
		ins.writeString(ins.onDuplicateKeyUpdateString)
	} else if len(ins.duplicateFields) > 0 {
		ins.writeString(" ON DUPLICATE KEY UPDATE ")
		ins.setNameFields(ins.duplicateFields)
	}
	ins.end()
	return ins.sb.String(), nil
}

// SQL 返回数组参数风格的 sql
func (ins *Inserter) SQL() (string, error) {
	if len(ins.insertFields) == 0 {
		return "", ErrColumnsRequire
	}
	ins.reset()
	switch ins.intoType {
	case intoTypeDefault:
		ins.writeString("INSERT INTO ")
	case intoTypeReplace:
		ins.writeString("REPLACE INTO ")
	case intoTypeIgnore:
		ins.writeString("INSERT IGNORE INTO ")
	}
	ins.quote(ins.tableName)
	ins.writeByte('(')
	for i, f := range ins.insertFields {
		ins.quote(f.col)
		if i != len(ins.insertFields)-1 {
			ins.writeString(", ")
		}
	}
	ins.writeString(") VALUES (")
	for i := range ins.insertFields {
		ins.writeByte('?')
		if i != len(ins.insertFields)-1 {
			ins.writeString(", ")
		}
	}
	ins.writeString(")")
	if ins.onDuplicateKeyUpdateString != "" {
		ins.writeString(" ON DUPLICATE KEY UPDATE ")
		ins.writeString(ins.onDuplicateKeyUpdateString)
	} else if len(ins.duplicateFields) > 0 {
		ins.writeString(" ON DUPLICATE KEY UPDATE ")
		ins.setFields(ins.duplicateFields)
	}
	ins.end()
	return ins.sb.String(), nil
}

// SQLArgs 构造 sql 并返回数组类型参数
// 需要通过 Fields 方法赋值，否则使用 NameSQL
func (ins *Inserter) SQLArgs() (string, []any, error) {
	execSQL, err := ins.SQL()
	if err != nil {
		return "", nil, err
	}
	var args []any
	for _, f := range ins.insertFields {
		if f.val != nil {
			args = append(args, f.val)
		}
	}
	for _, f := range ins.duplicateFields {
		if f.val != nil {
			args = append(args, f.val)
		}
	}
	return execSQL, args, nil
}
