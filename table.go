package amisgo

// ColumnType 列类型
type ColumnType string

const (
	ColumnTypeOperation ColumnType = "operation"
)

// Column 表格列
// https://baidu.github.io/amis/zh-CN/components/table2#%E5%B1%9E%E6%80%A7%E8%A1%A8
type Column struct {
	Name      string    `json:"name,omitempty"`      // 列名
	Type      string    `json:"type,omitempty"`      // 列类型
	Label     string    `json:"label,omitempty"`     // 列标题
	Sortable  bool      `json:"sortable,omitempty"`  // 是否可排序
	Width     int32     `json:"width,omitempty"`     // 列宽
	Toggled   bool      `json:"toggled,omitempty"`   // 是否默认展开
	Buttons   []*Action `json:"buttons,omitempty"`   // 按钮
	Tpl       string    `json:"tpl,omitempty"`       // 显示模板字符串
	QuickEdit any       `json:"quickEdit,omitempty"` // boolean | QuickEditConfig，快速编辑，一般需要配合quickSaveApi接口使用
}

type QuickEditConfig struct {
	*BaseOptionType
	Mode            string      `json:"mode,omitempty"`            // 编辑模式，inline 为行内编辑，popOver 为浮层编辑
	Type            Type        `json:"type,omitempty"`            // 表单项组件类型
	Icon            string      `json:"icon,omitempty"`            // 自定义快速编辑按钮的图标
	SaveImmediately bool        `json:"saveImmediately,omitempty"` // 是否立即保存
	Body            []Component `json:"body,omitempty"`            // 编辑内容
}

const (
	CRUDModeTable = "table"
	CRUDModeCards = "cards"
	CRUDModeList  = "list"
)

// CRUD amis CRUD 组件
// https://baidu.github.io/amis/zh-CN/components/crud
type CRUD struct {
	Base
	Columns          []*Column `json:"columns,omitempty"`          // 表格列配置
	ColumnsTogglable string    `json:"columnsTogglable,omitempty"` // 展示列显示开关, 自动即：列数量大于或等于 5 个时自动开启，auto 或者 boolean
	API              *API      `json:"api,omitempty"`              // 用来获取列表数据的 api
	QuickSaveApi     *API      `json:"quickSaveApi,omitempty"`     // 快速保存 api
	Filter           *Form     `json:"filter,omitempty"`           // 过滤器配置
	BulkActions      []*Action `json:"bulkActions,omitempty"`      // 批量操作配置
	AffixHeader      bool      `json:"affixHeader,omitempty"`      // 是否固定表头
	CombineNum       int32     `json:"combineNum,omitempty"`       // 合并行数
	Placeholder      string    `json:"placeholder,omitempty"`      // 空数据时显示的文案
	FooterToolbar    []string  `json:"footerToolbar,omitempty"`    // 底部工具栏配置, ['statistics', 'pagination']
}

// NewCRUD 创建一个新的 CRUD 实例
func NewCRUD() *CRUD {
	return &CRUD{
		Base: Base{
			Type: TypeCRUD,
		},
	}
}
