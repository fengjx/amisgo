package amisgo

import "encoding/json"

type BaseInput struct {
	Base
	Label        string `json:"label,omitempty"`        // 标签
	Clearable    bool   `json:"clearable,omitempty"`    // 是否可清空
	Placeholder  string `json:"placeholder,omitempty"`  // 占位符
	ReadOnly     bool   `json:"readOnly,omitempty"`     // 是否只读
	Required     bool   `json:"required,omitempty"`     // 是否必填
	TrimContents bool   `json:"trimContents,omitempty"` // 是否自动去除前后空格
}

// InputText 输入框
type InputText struct {
	*BaseInput
}

// NewInputText 创建一个新的 InputText 实例
func NewInputText() *InputText {
	return &InputText{
		BaseInput: &BaseInput{
			Base: Base{
				Type: TypeInputText,
			},
		},
	}
}

type Textarea struct {
	*BaseInput
	MinRows int32 `json:"minRows,omitempty"` // 最小行数
	MaxRows int32 `json:"maxRows,omitempty"` // 最大行数
}

// NewTextarea 创建一个新的 Textarea 实例
func NewTextarea() *Textarea {
	return &Textarea{
		BaseInput: &BaseInput{
			Base: Base{
				Type: TypeTextarea,
			},
		},
	}
}

// InputDatetime 日期时间选择器
// https://baidu.github.io/amis/zh-CN/components/form/input-datetime
type InputDatetime struct {
	*BaseInput
	ValueFormat   string `json:"valueFormat,omitempty"`   // 值格式化
	DisplayFormat string `json:"displayFormat,omitempty"` // 显示格式化
	MinDate       string `json:"minDate,omitempty"`       // 最小日期
	MaxDate       string `json:"maxDate,omitempty"`       // 最大日期
	IsEndDate     bool   `json:"isEndDate,omitempty"`     // 如果配置为 true，会自动默认为 23:59:59 秒
	DisabledDate  string `json:"disabledDate,omitempty"`  // 用字符函数来控制哪些天不可以被点选
}

// NewInputDatetime 创建一个新的 InputDatetime 实例
func NewInputDatetime() *InputDatetime {
	return &InputDatetime{
		BaseInput: &BaseInput{
			Base: Base{
				Type: TypeInputDatetime,
			},
		},
	}
}

// Option 选项
type Option struct {
	Value    string    `json:"value,omitempty"`    // 值
	Label    string    `json:"label,omitempty"`    // 标签
	Children []*Option `json:"children,omitempty"` // 子选项
}

// GetOptions 动态加载选项
type GetOptions func() ([]*Option, error)

// Options 选项列表
type Options struct {
	Options    []*Option  `json:"-"` // 选项
	GetOptions GetOptions `json:"-"` // 动态加载选项
}

func (s *Options) MarshalJSON() ([]byte, error) {
	if s.GetOptions != nil {
		options, err := s.GetOptions()
		if err != nil {
			return nil, err
		}
		s.Options = options
	}
	return json.Marshal(s.Options)
}

// BaseOptionType 基础选项类型
type BaseOptionType struct {
	Options *Options `json:"options,omitempty"` // 选项
	Source  any      `json:"source,omitempty"`  // 数据源 API 或者 string (数据映射)
}

// Select 下拉框
type Select struct {
	BaseInput
	BaseOptionType
	Label    string `json:"label,omitempty"`    // 标签
	ReadOnly bool   `json:"readOnly,omitempty"` // 是否只读
	Required bool   `json:"required,omitempty"` // 是否必填
}

// NewSelect 创建一个新的 Select 实例
func NewSelect() *Select {
	return &Select{
		BaseInput: BaseInput{
			Base: Base{
				Type: TypeSelect,
			},
		},
		BaseOptionType: BaseOptionType{},
	}
}

// WithLabel 设置标签
func (s *Select) WithLabel(label string) *Select {
	s.Label = label
	return s
}

// WithRequired 设置是否必填
func (s *Select) WithRequired(required bool) *Select {
	s.Required = required
	return s
}

// WithSource 设置数据源
func (s *Select) WithSource(source any) *Select {
	s.Source = source
	return s
}

// InputTree 树形选择器
// https://aisuda.bce.baidu.com/amis/zh-CN/components/form/input-tree
type InputTree struct {
	*BaseInput
	*BaseOptionType
}

// NewInputTree 创建一个新的 InputTree 实例
func NewInputTree() *InputTree {
	return &InputTree{
		BaseInput: &BaseInput{
			Base: Base{
				Type: TypeInputTree,
			},
		},
		BaseOptionType: &BaseOptionType{},
	}
}

// WithSource 设置数据源
func (s *InputTree) WithSource(source any) *InputTree {
	s.Source = source
	return s
}

// TreeSelect 树形选择器
// https://aisuda.bce.baidu.com/amis/zh-CN/components/form/treeselect
type TreeSelect struct {
	*BaseInput
	*BaseOptionType
	HideNodePathLabel bool `json:"hideNodePathLabel,omitempty"` // 是否隐藏选择框中已选择节点的路径 label 信息
	OnlyLeaf          bool `json:"onlyLeaf,omitempty"`          // 只允许选择叶子节点
	Searchable        bool `json:"searchable,omitempty"`        // 是否可检索，仅在 type 为 tree-select 的时候生效
}

// NewTreeSelect 创建一个新的 TreeSelect 实例
func NewTreeSelect() *TreeSelect {
	return &TreeSelect{
		BaseInput: &BaseInput{
			Base: Base{
				Type: TypeTreeSelect,
			},
		},
		BaseOptionType: &BaseOptionType{},
	}
}

// WithSource 设置数据源
func (s *TreeSelect) WithSource(source any) *TreeSelect {
	s.Source = source
	return s
}
