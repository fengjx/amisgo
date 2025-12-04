package amisgo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sync"

	"github.com/fengjx/amisgo/internal/dbutil"
	"github.com/fengjx/amisgo/internal/kit"
	"github.com/fengjx/amisgo/internal/ql"
)

var (
	// 全局默认配置
	global *globalConfig

	// InputTypeMap 数据库类型映射到 amis 组件类型
	InputTypeMap = map[string]Type{
		dbutil.Date:          TypeInputDate,
		dbutil.DateTime:      TypeInputDatetime,
		dbutil.SmallDateTime: TypeInputDatetime,
		dbutil.Time:          TypeInputDatetime,
		dbutil.TimeStamp:     TypeInputDatetime,
		dbutil.TimeStampz:    TypeInputDatetime,
	}
)

func init() {
	global = &globalConfig{
		metaMap:         make(map[string]*dbutil.DBTable),
		omitEditColumns: make(map[string]struct{}),
	}
}

const (
	OpAnd Op = "and"
	OpOr  Op = "or"

	ConditionTypeEq      ConditionType = "eq"       // 等于
	ConditionTypeNotEq   ConditionType = "not_eq"   // 不等于
	ConditionTypeLike    ConditionType = "like"     // 模糊匹配
	ConditionTypeNotLike ConditionType = "not_like" // 不包含
	ConditionTypeIn      ConditionType = "in"       // in
	ConditionTypeNotIn   ConditionType = "not_in"   // not in
	ConditionTypeGt      ConditionType = "gt"       // 大于
	ConditionTypeLt      ConditionType = "lt"       // 小于
	ConditionTypeGte     ConditionType = "gte"      // 大于等于
	ConditionTypeLte     ConditionType = "lte"      // 小于等于

	OrderTypeAsc  OrderType = "asc"  // 升序
	OrderTypeDesc OrderType = "desc" // 降序
)

// ConditionType 查询条件类型
type ConditionType string

// getInputType sqlType 映射 input type
func getInputType(sqlType string) Type {
	if input, ok := InputTypeMap[sqlType]; ok {
		return input
	}
	return TypeInputText
}

// OrderField 排序字段
type OrderField struct {
	Field     string    `json:"field"`
	OrderType OrderType `json:"order_type"`
}

type OrderType string

// Condition 条件语句
type Condition struct {
	Disable       bool          `json:"disable"`        // true 禁用该条件
	Op            Op            `json:"op"`             // and or 连接符
	Field         string        `json:"field"`          // 查询条件字段
	Vals          []any         `json:"vals"`           // 查询字段值
	ConditionType ConditionType `json:"condition_type"` // 查找类型
}

// Op and or连接符
type Op string

// ToSQLArgs 返回 sql 语句和参数
func (q QueryReq) ToSQLArgs(offset, limit int64) (sql string, args []any, err error) {
	selector := q.buildSelector()
	selector.Offset(offset).Limit(limit)
	return selector.SQLArgs()
}

// ToCountSQLArgs 返回 count 查询 sql 语句和参数
func (q QueryReq) ToCountSQLArgs() (sql string, args []any, err error) {
	selector := q.buildSelector()
	return selector.CountSQLArgs()
}

func (q QueryReq) buildSelector() *ql.Selector {
	selector := ql.NewSelector(q.TableName)
	selector.Columns(q.Fields...)
	selector.Where(buildCondition(q.Conditions))
	if q.Page != nil {
		selector.Offset(q.Page.Offset).Limit(q.Page.Limit)
	}
	if len(q.OrderFields) > 0 {
		var orderBy []ql.OrderBy
		for _, orderField := range q.OrderFields {
			switch orderField.OrderType {
			case OrderTypeAsc:
				orderBy = append(orderBy, ql.Asc(orderField.Field))
			case OrderTypeDesc:
				orderBy = append(orderBy, ql.Desc(orderField.Field))
			}
		}
		selector.OrderBy(orderBy...)
	}
	return selector
}

func buildCondition(conditions []Condition) ql.ConditionBuilder {
	where := ql.C()
	for _, c := range conditions {
		if c.Disable {
			continue
		}
		switch {
		case c.ConditionType == ConditionTypeEq && c.Op == OpAnd:
			where.And(ql.Col(c.Field).EQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeEq && c.Op == OpOr:
			where.Or(ql.Col(c.Field).EQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeNotEq && c.Op == OpAnd:
			where.And(ql.Col(c.Field).NotEQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeNotEq && c.Op == OpOr:
			where.Or(ql.Col(c.Field).NotEQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeLike && c.Op == OpAnd:
			where.And(ql.Col(c.Field).Like(c.Vals[0]))
		case c.ConditionType == ConditionTypeLike && c.Op == OpOr:
			where.Or(ql.Col(c.Field).Like(c.Vals[0]))
		case c.ConditionType == ConditionTypeNotLike && c.Op == OpAnd:
			where.And(ql.Col(c.Field).NotLike(c.Vals[0]))
		case c.ConditionType == ConditionTypeNotLike && c.Op == OpOr:
			where.Or(ql.Col(c.Field).NotLike(c.Vals[0]))
		case c.ConditionType == ConditionTypeIn && c.Op == OpAnd:
			where.And(ql.Col(c.Field).In(c.Vals...))
		case c.ConditionType == ConditionTypeIn && c.Op == OpOr:
			where.Or(ql.Col(c.Field).In(c.Vals...))
		case c.ConditionType == ConditionTypeNotIn && c.Op == OpAnd:
			where.And(ql.Col(c.Field).NotIn(c.Vals...))
		case c.ConditionType == ConditionTypeNotIn && c.Op == OpOr:
			where.Or(ql.Col(c.Field).NotIn(c.Vals...))
		case c.ConditionType == ConditionTypeGt && c.Op == OpAnd:
			where.And(ql.Col(c.Field).GT(c.Vals[0]))
		case c.ConditionType == ConditionTypeGt && c.Op == OpOr:
			where.Or(ql.Col(c.Field).GT(c.Vals[0]))
		case c.ConditionType == ConditionTypeLt && c.Op == OpAnd:
			where.And(ql.Col(c.Field).LT(c.Vals[0]))
		case c.ConditionType == ConditionTypeLt && c.Op == OpOr:
			where.Or(ql.Col(c.Field).LT(c.Vals[0]))
		case c.ConditionType == ConditionTypeGte && c.Op == OpAnd:
			where.And(ql.Col(c.Field).GTEQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeGte && c.Op == OpOr:
			where.Or(ql.Col(c.Field).GTEQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeLte && c.Op == OpAnd:
			where.And(ql.Col(c.Field).LTEQ(c.Vals[0]))
		case c.ConditionType == ConditionTypeLte && c.Op == OpOr:
			where.Or(ql.Col(c.Field).LTEQ(c.Vals[0]))
		}
	}
	return where
}

// global 全局配置
type globalConfig struct {
	mux sync.Mutex
	// defaultExecutor 全局默认数据库执行器
	defaultExecutor dbutil.Executor
	// 所有表元信息
	metaMap map[string]*dbutil.DBTable
	// 保存和更新时默认忽略的字段，全局生效
	// 一般用户统一的开发规范
	omitEditColumns map[string]struct{}
}

func (g *globalConfig) registerMeta(meta *dbutil.DBTable) {
	g.mux.Lock()
	defer g.mux.Unlock()
	g.metaMap[meta.Name] = meta
}

// UseOmitsEditColumns 设置保存时全局默认忽略的字段
func UseOmitsEditColumns(omits ...string) {
	for _, omit := range omits {
		global.omitEditColumns[omit] = struct{}{}
	}
}

// UseDefaultExecutor 设置全局默认数据库执行器
func UseDefaultExecutor(executor dbutil.Executor) {
	global.defaultExecutor = executor
}

// Field 字段设置
type Field struct {
	IsPrimaryKey        bool          // 是否主键
	Name                string        // 字段名称
	Type                Type          // 字段类型，必须是 input 类型
	Label               string        // 字段标签
	ShowTable           bool          // 表格字段是否显示
	ShowView            bool          // 详情是否显示
	Required            bool          // 是否必填
	CreateAble          bool          // 创建记录是否显示字段
	UpdateAble          bool          // 修改记录是否显示字段
	SearchAble          bool          // 是否搜索字段
	SearchConditionType ConditionType // 搜索字段匹配类型 ConditionType
	Sortable            bool          // 是否可排序
	QuickEdit           bool          // 是否快速编辑
	Options             *Options      // 动态加载选项, select tree-select 类型才有意义
	Tpl                 string        // 显示模板字符串
	Component           Component     // 组件，优先级最高，如果不为空则忽略其他配置
}

func (f *Field) withRequired(required bool) *Field {
	nf := *f
	nf.Required = required
	return &nf
}

// toInput 转换为 input 组件
func (f *Field) toInput() Component {
	if f.Component != nil {
		return f.Component
	}
	switch f.Type {
	case TypeTextarea:
		in := NewTextarea()
		in.Label = f.Label
		in.Name = f.Name
		in.Required = f.Required
		in.Tpl = f.Tpl
		return in
	case TypeSelect:
		in := NewSelect()
		in.Label = f.Label
		in.Name = f.Name
		in.Required = f.Required
		in.Options = f.Options
		in.Tpl = f.Tpl
		return in
	case TypeTreeSelect:
		in := NewTreeSelect()
		in.Name = f.Name
		in.Label = f.Label
		in.Required = f.Required
		in.Options = f.Options
		in.Tpl = f.Tpl
		return in
	case TypeInputDatetime:
		in := NewInputDatetime()
		in.Label = f.Label
		in.Name = f.Name
		in.Required = f.Required
		in.Tpl = f.Tpl
		return in
	}
	in := NewInputText()
	in.Label = f.Label
	in.Name = f.Name
	in.Required = f.Required
	in.Tpl = f.Tpl
	return in
}

// toQuickEditInput 转换为快速编辑 input 组件
func (f *Field) toQuickEditInput() any {
	if f.Component != nil {
		return true
	}
	switch f.Type {
	case TypeSelect:
		cfg := QuickEditConfig{
			Mode: "inline",
			Type: f.Type,
			BaseOptionType: &BaseOptionType{
				Options: f.Options,
			},
		}
		return cfg
	default:
	}
	return true
}

// CreateHandler 创建数据
// 返回值会返回给前端
type CreateHandler func(ctx context.Context, data map[string]any) (*Resp, error)

// UpdateHandler 更新数据
// 返回值会返回给前端
type UpdateHandler func(ctx context.Context, id any, data map[string]any) (*Resp, error)

// BatchUpdateHandler 批量更新数据
// 返回值会返回给前端
type BatchUpdateHandler func(ctx context.Context, param *BatchUpdate) (*Resp, error)

// DeleteHandler 删除数据
// 返回值会返回给前端
type DeleteHandler func(ctx context.Context, ids ...any) (*Resp, error)

// QueryHandler 搜索数据
// 返回值会返回给前端
type QueryHandler func(ctx context.Context, param *QueryReq) (*Resp, error)

// CRUDOptions CRUD 页面配置项
type CRUDOptions struct {
	menuID             string    // 菜单ID
	apiPrefix          string    // http 接口前缀
	apiPath            string    // http 接口路径
	dbName             string    // 数据库名，使用MySQL时需要指定
	tableName          string    // 数据库表名
	tableLabel         string    // 数据库表名显示标签
	idField            string    // 主键字段
	fieldSort          []string  // 字段显示顺序
	operations         []*Action // 自定义操作按钮
	fields             []*Field  // 字段设置
	createHandler      CreateHandler
	updateHandler      UpdateHandler
	batchUpdateHandler BatchUpdateHandler
	deleteHandler      DeleteHandler
	queryHandler       QueryHandler
	createPath         string
	updatePath         string
	batchUpdatePath    string
	deletePath         string
	queryPath          string
	pagePath           string
	executor           dbutil.Executor
}

type CRUDOption func(*CRUDOptions)

func WithMenuID(menuID string) CRUDOption {
	return func(o *CRUDOptions) {
		o.menuID = menuID
	}
}

// WithAPIPrefix 设置 http 接口前缀
func WithAPIPrefix(apiPrefix string) CRUDOption {
	return func(o *CRUDOptions) {
		o.apiPrefix = apiPrefix
	}
}

// WithDBName 设置数据库名
func WithDBName(dbName string) CRUDOption {
	return func(o *CRUDOptions) {
		o.dbName = dbName
	}
}

// WithExecutor 设置数据库执行器
func WithExecutor(executor dbutil.Executor) CRUDOption {
	return func(o *CRUDOptions) {
		o.executor = executor
	}
}

// WithTableLabel 设置表名外显名称
func WithTableLabel(tableLabel string) CRUDOption {
	return func(o *CRUDOptions) {
		o.tableLabel = tableLabel
	}
}

// WithIDField 设置主键字段
func WithIDField(idField string) CRUDOption {
	return func(o *CRUDOptions) {
		o.idField = idField
	}
}

// WithFields 设置字段设置
func WithFields(fields []*Field) CRUDOption {
	return func(o *CRUDOptions) {
		o.fields = fields
	}
}

// WithFieldSort 设置字段显示顺序
func WithFieldSort(fieldSort []string) CRUDOption {
	return func(o *CRUDOptions) {
		o.fieldSort = fieldSort
	}
}

// WithCreateHandler 设置创建数据处理器
func WithCreateHandler(createHandler CreateHandler) CRUDOption {
	return func(o *CRUDOptions) {
		o.createHandler = createHandler
	}
}

// WithUpdateHandler 设置更新数据处理器
func WithUpdateHandler(updateHandler UpdateHandler) CRUDOption {
	return func(o *CRUDOptions) {
		o.updateHandler = updateHandler
	}
}

// WithBatchUpdateHandler 设置批量更新数据处理器
func WithBatchUpdateHandler(batchUpdateHandler BatchUpdateHandler) CRUDOption {
	return func(o *CRUDOptions) {
		o.batchUpdateHandler = batchUpdateHandler
	}
}

// WithDeleteHandler 设置删除数据处理器
func WithDeleteHandler(deleteHandler DeleteHandler) CRUDOption {
	return func(o *CRUDOptions) {
		o.deleteHandler = deleteHandler
	}
}

// WithQueryHandler 设置搜索数据处理器
func WithQueryHandler(queryHandler QueryHandler) CRUDOption {
	return func(o *CRUDOptions) {
		o.queryHandler = queryHandler
	}
}

// WithOperations 设置自定义操作按钮
func WithOperations(ops []*Action) CRUDOption {
	return func(o *CRUDOptions) {
		o.operations = ops
	}
}

// AdminCRUD CRUD 管理配置
type AdminCRUD struct {
	init            bool              // 是否初始化
	component       Component         // 渲染 amis 组件
	opt             *CRUDOptions      // 配置选项
	table           *dbutil.DBTable   // 表元信息
	primaryKeyField *Field            // 主键
	fieldMap        map[string]*Field // 字段映射: 数据库字段名和外显名称映射
	fields          []*Field          // 字段定义
	tableColumns    []*Column         // 表格显示字段
	viewColumns     []*Field          // 详情显示字段
	createFields    []*Field          // 新增表单字段
	updateFields    []*Field          // 修改表单字段
	searchFields    []*Field          // 搜索字段
}

// NewAdminCRUD 创建一个新的 CrudAdmin 实例
func NewAdminCRUD(tableName string, apiPath string, opts ...CRUDOption) (*AdminCRUD, error) {
	options := &CRUDOptions{
		tableName: tableName,
		apiPath:   apiPath,
	}
	for _, opt := range opts {
		opt(options)
	}
	if options.menuID == "" {
		options.menuID = options.tableName
	}
	if options.executor == nil {
		options.executor = global.defaultExecutor
	}
	options.createPath = joinPath(options.apiPath, "add")
	options.updatePath = joinPath(options.apiPath, "update")
	options.batchUpdatePath = joinPath(options.apiPath, "batch-update")
	options.deletePath = joinPath(options.apiPath, "del")
	options.queryPath = joinPath(options.apiPath, "query")
	options.pagePath = joinPath(options.apiPath, "page.json")
	a := &AdminCRUD{
		opt: options,
	}
	if err := a.Init(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *AdminCRUD) Init() error {
	if a.init {
		return nil
	}
	if err := a.initTable(); err != nil {
		return err
	}
	a.initFields()
	a.initComponent()
	return nil
}

// initTable 初始化表
func (a *AdminCRUD) initTable() error {
	if a.opt.executor == nil {
		a.fields = a.opt.fields
		return nil
	}
	table, err := dbutil.GetTableMeta(context.Background(), a.opt.executor, a.opt.dbName, a.opt.tableName)
	if err != nil {
		return err
	}
	a.table = table
	a.opt.idField = table.PrimaryKey.Name
	fields := make([]*Field, 0, len(table.Columns))
	fieldMap := make(map[string]*Field)
	for _, f := range a.opt.fields {
		fieldMap[f.Name] = f
	}
	// 补全其他字段映射
	for _, col := range table.Columns {
		f := fieldMap[col.Name]
		if f == nil {
			createAble := true
			updateAble := true
			if _, ok := global.omitEditColumns[col.Name]; ok {
				createAble = false
				updateAble = false
			}
			f = &Field{
				IsPrimaryKey: col.IsPrimaryKey,
				Name:         col.Name,
				Label:        col.GetComment(),
				Type:         getInputType(col.SQLType),
				CreateAble:   createAble,
				UpdateAble:   updateAble,
				SearchAble:   false,
				ShowTable:    true,
				ShowView:     true,
			}
		}
		f.IsPrimaryKey = col.IsPrimaryKey
		f.Name = col.Name
		// 补充字段信息
		if f.Label == "" {
			f.Label = col.GetComment()
		}
		if f.Type == "" {
			f.Type = getInputType(col.SQLType)
		}
		if f.SearchConditionType == "" {
			f.SearchConditionType = ConditionTypeEq
		}
		if f.IsPrimaryKey {
			a.primaryKeyField = f
		}
		fieldMap[col.Name] = f
		fields = append(fields, f)
	}
	if len(a.opt.fieldSort) > 0 {
		// 1. 先按 fieldSort 顺序添加
		sortedFields := make([]*Field, 0, len(fields))
		used := make(map[string]bool)
		for _, name := range a.opt.fieldSort {
			if f, ok := fieldMap[name]; ok && f != nil {
				sortedFields = append(sortedFields, f)
				used[name] = true
			}
		}
		// 2. 再把剩下的字段按原顺序追加
		for _, f := range fields {
			if !used[f.Name] {
				sortedFields = append(sortedFields, f)
			}
		}
		fields = sortedFields
	}
	a.fieldMap = fieldMap
	a.fields = fields
	global.registerMeta(table)
	return nil
}

// initFields 初始化字段
func (a *AdminCRUD) initFields() {
	for _, f := range a.fields {
		if f.ShowTable {
			c := &Column{
				Name:     f.Name,
				Label:    f.Label,
				Sortable: f.Sortable,
			}
			if f.QuickEdit {
				c.QuickEdit = f.toQuickEditInput()
			}
			a.tableColumns = append(a.tableColumns, c)
		}
		if f.ShowView {
			a.viewColumns = append(a.viewColumns, f)
		}
		if f.CreateAble {
			a.createFields = append(a.createFields, f)
		}
		if f.UpdateAble {
			a.updateFields = append(a.updateFields, f)
		}
		if f.SearchAble {
			a.searchFields = append(a.searchFields, f)
		}
	}
}

func (a *AdminCRUD) initComponent() {
	// 以 sys_user 为例，实际可根据 a.tableName 动态生成
	kbName := kit.KebabCase(a.opt.tableName)

	// 快速保存 API
	quickSaveApi := NewAPI()
	quickSaveApi.Method = http.MethodPatch
	quickSaveApi.URL = path.Join("${API_BASEURL}", a.opt.apiPrefix, a.opt.batchUpdatePath)
	quickSaveApi.Data = map[string]any{"rows": "${rowsDiff}"}

	// 列定义
	var columns []*Column
	columns = append(columns, a.tableColumns...)

	buttons := []*Action{
		a.buildViewBtn(),
		a.buildUpdateBtn(),
		a.buildDeleteBtn(),
	}
	if a.opt.operations != nil {
		buttons = append(buttons, a.opt.operations...)
	}

	// 操作按钮列
	columns = append(columns, &Column{
		Label: "操作", Type: "operation", Width: 350,
		Toggled: true,
		Buttons: buttons,
	})

	filter, api := a.buildFilter()

	bulkApi := NewAPI()
	bulkApi.Method = "DELETE"
	bulkApi.URL = path.Join("${API_BASEURL}", a.opt.apiPrefix, a.opt.deletePath)
	bulkApi.Data = map[string]any{"ids": "${ids|split}"}
	// 批量操作
	bulkActions := []*Action{
		NewAction().
			WithLabel("批量删除").
			WithLevel(ButtonLevelDanger).
			WithConfirmText("确定要批量删除?").
			WithAPI(bulkApi),
	}

	// 底部工具栏
	footerToolbar := []string{
		"statistics",
		"switch-per-page",
		"pagination",
	}

	crud := NewCRUD()
	crud.Name = fmt.Sprintf("%s-crud", kbName)
	crud.AffixHeader = true
	crud.ColumnsTogglable = "auto"
	crud.Placeholder = "暂无数据"
	crud.CombineNum = 0
	crud.BodyClassName = "panel-default"
	crud.Columns = columns
	crud.API = api
	crud.Filter = filter
	crud.QuickSaveApi = quickSaveApi
	crud.BulkActions = bulkActions
	crud.FooterToolbar = footerToolbar

	c := NewPage()
	c.Name = kbName
	c.Body = append(c.Body, crud)

	a.component = c
}

func (a *AdminCRUD) buildFilter() (*Form, *API) {
	var conds []map[string]any
	var inputs []Component
	for _, f := range a.searchFields {
		val := fmt.Sprintf("${%s}", f.Name)
		if f.SearchConditionType == ConditionTypeLike || f.SearchConditionType == ConditionTypeNotLike {
			val = fmt.Sprintf("%%%s%%", val)
		}
		conds = append(conds, map[string]any{
			"disable":        fmt.Sprintf("${!%s}", f.Name),
			"op":             "and",
			"field":          f.Name,
			"condition_type": f.SearchConditionType,
			"vals":           []string{val},
		})
		inputs = append(inputs, f.withRequired(false).toInput())
	}
	api := NewAPI()
	api.Method = "POST"
	api.URL = path.Join("${API_BASEURL}", a.opt.apiPrefix, a.opt.queryPath)
	api.Data = map[string]any{
		"order_fields": []map[string]any{
			{"field": "${orderBy}", "order_type": "${orderDir}"},
		},
		"conditions": conds,
	}

	actions := []Component{
		a.buildCreateBtn(),
		NewButton(ButtonActionTypeReset).WithLabel("重置"),
		NewButton(ButtonActionTypeSubmit).WithLabel("查询").WithLevel(ButtonLevelPrimary),
	}

	form := NewForm()
	form.Title = "搜索"
	form.Body = []Component{
		NewGroup().WithBody(inputs),
	}
	form.Actions = actions
	return form, api

}

func (a *AdminCRUD) buildCreateBtn() *Action {
	fields := a.createFields
	data := make(map[string]any, len(fields)+1)
	for _, field := range fields {
		if field.Name == a.opt.idField {
			continue
		}
		data[field.Name] = fmt.Sprintf("${%s}", field.Name)
	}
	api := NewAPI()
	api.Method = "POST"
	api.URL = path.Join("${API_BASEURL}", a.opt.apiPrefix, a.opt.createPath)
	api.Data = data

	body := make([]Component, 0, len(fields)+1)
	for _, field := range fields {
		if field.Name == a.opt.idField {
			continue
		}
		body = append(body, field.toInput())
	}

	form := NewForm()
	form.Name = "add-form"
	form.API = api
	form.Body = body

	createBtn := NewAction().WithIcon("fa fa-plus").WithLabel("新增").WithDrawer(
		NewDrawer().
			WithTitle(a.opt.tableLabel).
			WithSize("xl").
			WithBody(form),
	)
	return createBtn
}

func (a *AdminCRUD) buildViewBtn() *Action {
	fields := a.viewColumns
	// 动态生成新增表单内容
	body := make([]Component, 0, len(fields)+1)
	for _, field := range fields {
		if field.Name == a.opt.idField {
			continue
		}
		body = append(body, field.toInput())
	}
	form := NewForm()
	form.Name = "view-form"
	form.Static = true
	form.Body = body

	viewBtn := NewAction().WithIcon("fa fa-eye").WithLabel("详情").WithDrawer(
		NewDrawer().
			WithTitle(a.opt.tableLabel).
			WithSize("xl").
			WithBody(form),
	)
	return viewBtn
}

func (a *AdminCRUD) buildUpdateBtn() *Action {

	fields := a.updateFields
	data := make(map[string]any, len(fields)+1)
	for _, field := range fields {
		data[field.Name] = fmt.Sprintf("${%s}", field.Name)
	}
	// 主键字段
	data[a.opt.idField] = fmt.Sprintf("${%s}", a.opt.idField)
	api := NewAPI()
	api.Method = http.MethodPut
	api.URL = path.Join("${API_BASEURL}", a.opt.apiPrefix, a.opt.updatePath)
	api.Data = data

	body := make([]Component, 0, len(fields)+1)
	for _, field := range fields {
		body = append(body, field.toInput())
	}

	form := NewForm()
	form.Name = "update-form"
	form.API = api
	form.Body = body

	updateBtn := NewAction().WithIcon("fa fa-pencil").WithLabel("修改").WithDrawer(
		NewDrawer().
			WithTitle(a.opt.tableLabel).
			WithSize("xl").
			WithBody(form),
	)
	return updateBtn
}

func (a *AdminCRUD) buildDeleteBtn() *Action {
	api := NewAPI()
	api.Method = http.MethodDelete
	api.URL = path.Join("${API_BASEURL}", a.opt.apiPrefix, a.opt.deletePath)
	api.Data = map[string]any{
		"ids": fmt.Sprintf("${[${%s}]|toJson}", a.opt.idField),
	}
	delBtn := NewAction().
		WithIcon("fa fa-trash").
		WithLabel("删除").
		WithAPI(api).
		WithConfirmText("确定要删除?")

	return delBtn
}

func (a *AdminCRUD) create(ctx context.Context, data map[string]any) (*Resp, error) {
	if a.opt.createHandler != nil {
		return a.opt.createHandler(ctx, data)
	}
	r, err := dbutil.Insert(ctx, a.opt.executor, a.table, data)
	if err != nil {
		return nil, err
	}
	id, _ := r.LastInsertId()
	return &Resp{
		Status: StatusOK,
		Msg:    "创建成功",
		Data:   map[string]any{a.opt.idField: id},
	}, nil
}

func (a *AdminCRUD) update(ctx context.Context, id any, data map[string]any) (*Resp, error) {
	if len(data) == 0 {
		return nil, nil
	}
	if a.opt.updateHandler != nil {
		return a.opt.updateHandler(ctx, id, data)
	}
	r, err := dbutil.UpdateByID(ctx, a.opt.executor, a.table, id, data)
	if err != nil {
		return nil, err
	}
	i, err := r.RowsAffected()
	if err != nil {
		return nil, err
	}
	resp := &Resp{
		Status: StatusOK,
		Msg:    "更新成功",
		Data: map[string]any{
			"success": i > 0,
		},
	}
	return resp, nil
}

func (a *AdminCRUD) batchUpdate(ctx context.Context, param *BatchUpdate) (*Resp, error) {
	if param == nil || len(param.Rows) == 0 {
		return nil, nil
	}
	if a.opt.batchUpdateHandler != nil {
		return a.opt.batchUpdateHandler(ctx, param)
	}
	rows := make([]dbutil.RowData, 0, len(param.Rows))
	for _, row := range param.Rows {
		rows = append(rows, row)
	}
	err := dbutil.BatchUpdate(ctx, a.opt.executor, a.table, rows)
	if err != nil {
		return nil, err
	}
	return &Resp{
		Status: StatusOK,
		Msg:    "更新成功",
		Data:   map[string]any{"success": true},
	}, nil
}

// Delete 删除数据
func (a *AdminCRUD) delete(ctx context.Context, ids ...any) (*Resp, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if a.opt.deleteHandler != nil {
		return a.opt.deleteHandler(ctx, ids...)
	}
	r, err := dbutil.DeleteByIDs(ctx, a.opt.executor, a.table, ids...)
	if err != nil {
		return nil, err
	}
	ra, _ := r.RowsAffected()
	return &Resp{
		Status: StatusOK,
		Msg:    "删除成功",
		Data:   map[string]any{"success": ra > 0},
	}, nil
}

// query 分页查询封装
func (a *AdminCRUD) query(ctx context.Context, query *QueryReq) (*Resp, error) {
	if a.opt.queryHandler != nil {
		return a.opt.queryHandler(ctx, query)
	}
	query.TableName = a.opt.tableName
	sql, args, err := query.ToCountSQLArgs()
	if err != nil {
		return nil, err
	}
	if query.Page == nil {
		query.Page = &PageReq{
			Offset:     0,
			Limit:      10,
			QueryCount: true,
		}
	}
	offset := query.Page.Offset
	limit := query.Page.Limit
	total, err := dbutil.Count(ctx, a.opt.executor, sql, args...)
	if err != nil {
		return nil, err
	}

	page := &PageData[dbutil.RowData]{
		Total: total,
		Items: []dbutil.RowData{},
	}
	if total == 0 {
		return &Resp{
			Status: StatusOK,
			Data:   page,
		}, nil
	}
	sql, args, err = query.ToSQLArgs(offset, limit)
	if err != nil {
		return nil, err
	}
	list, err := dbutil.List(ctx, a.opt.executor, sql, args...)
	if err != nil {
		return nil, err
	}
	page.Items = list
	return &Resp{
		Status: StatusOK,
		Data:   page,
	}, nil
}

// Component 返回 amis 组件
func (a *AdminCRUD) Component() Component {
	return a.component
}

// ToJSON 返回 amis 组件 json
func (a *AdminCRUD) ToJSON() string {
	b, _ := json.Marshal(a.component)
	return string(b)
}
