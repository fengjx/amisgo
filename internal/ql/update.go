package ql

// Updater update 语句构造器
type Updater struct {
	sqlBuilder
	tableName string
	fields    []Field
	where     ConditionBuilder
}

// NewUpdater 创建一个 update 语句构造器
func NewUpdater(tableName string) *Updater {
	return &Updater{
		tableName: tableName,
	}
}

// Fields 设置字段值
func (u *Updater) Fields(fields ...Field) *Updater {
	for _, field := range fields {
		if !field.isUse {
			continue
		}
		u.fields = append(u.fields, field)
	}
	return u
}

// Set 设置字段值
func (u *Updater) Set(column string, val any) *Updater {
	u.fields = append(u.fields, F(column).Val(val))
	return u
}

// SetMap 设置字段值
func (u *Updater) SetMap(data map[string]any) *Updater {
	for k, v := range data {
		u.Set(k, v)
	}
	return u
}

// Columns update 的数据库字段
func (u *Updater) Columns(columns ...string) *Updater {
	for _, col := range columns {
		u.fields = append(u.fields, F(col))
	}
	return u
}

// Incr 数值增加，eg: set a = a + 1
func (u *Updater) Incr(column string, n int64) *Updater {
	u.fields = append(u.fields, F(column).Incr(n))
	return u
}

// Where 条件
// condition 可以通过 sqlbuilder.C() 方法创建
func (u *Updater) Where(where ConditionBuilder) *Updater {
	u.where = where
	return u
}

// SQL 输出sql语句
func (u *Updater) SQL() (string, error) {
	if len(u.fields) == 0 {
		return "", ErrColumnsRequire
	}
	u.reset()
	u.writeString("UPDATE ")
	u.quote(u.tableName)
	u.writeString(" SET ")
	u.setFields(u.fields)
	u.whereSQL(u.where)
	u.end()
	return u.sb.String(), nil
}

// SQLArgs 构造 sql 并返回对应参数
func (u *Updater) SQLArgs() (string, []any, error) {
	if len(u.fields) == 0 {
		return "", nil, ErrColumnsRequire
	}
	if u.where != nil && len(u.where.getPredicates()) == 0 {
		return "", nil, ErrUpdateMissWhere
	}
	execSQL, err := u.SQL()
	var args []any
	for _, f := range u.fields {
		if f.val != nil {
			args = append(args, f.val)
		}
	}
	wargs, hasInSQL := u.whereArgs(u.where)
	if len(wargs) > 0 {
		args = append(args, wargs...)
	}
	if !hasInSQL {
		return execSQL, args, err
	}
	return parseIn(execSQL, args...)
}

func (u *Updater) NameSQL() (string, error) {
	if len(u.fields) == 0 {
		return "", ErrColumnsRequire
	}
	if u.where != nil && len(u.where.getPredicates()) == 0 {
		return "", ErrUpdateMissWhere
	}
	u.reset()
	u.writeString("UPDATE ")
	u.quote(u.tableName)
	u.writeString(" SET ")
	u.setNameFields(u.fields)
	u.whereSQL(u.where)
	u.end()
	return u.sb.String(), nil
}
