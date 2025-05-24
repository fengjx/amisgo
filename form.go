package amisgo

// Form amis 表单组件
// https://baidu.github.io/amis/zh-CN/components/form/index
type Form struct {
	Base
	Title   string      `json:"title,omitempty"`   // 表单标题
	Static  bool        `json:"static,omitempty"`  // 是否静态
	API     *API        `json:"api,omitempty"`     // 表单提交 API
	Body    []Component `json:"body,omitempty"`    // 表单项配置
	Actions []Component `json:"actions,omitempty"` // 表单操作按钮配置
}

// NewForm 创建一个新的 Form 实例
func NewForm() *Form {
	return &Form{
		Base: Base{
			Type: TypeForm,
		},
	}
}

// WithTitle 设置表单标题
func (f *Form) WithTitle(title string) *Form {
	f.Title = title
	return f
}

// WithAPI 设置表单提交 API
func (f *Form) WithAPI(api *API) *Form {
	f.API = api
	return f
}

// WithBody 设置表单项配置
func (f *Form) WithBody(body ...Component) *Form {
	f.Body = body
	return f
}

// WithActions 设置表单操作按钮配置
func (f *Form) WithActions(actions ...Component) *Form {
	f.Actions = actions
	return f
}

// WithStatic 设置表单是否静态
func (f *Form) WithStatic(static bool) *Form {
	f.Static = static
	return f
}
