package amisgo

type ActionType string

// 【必填】这是 action 最核心的配置，来指定该 action 的作用类型，支持：ajax、link、url、drawer、dialog、confirm、cancel、prev、next、copy、close。
const (
	ActionTypeAjax     ActionType = "ajax"
	ActionTypeLink     ActionType = "link"
	ActionTypeURL      ActionType = "url"
	ActionTypeDrawer   ActionType = "drawer"
	ActionTypeDialog   ActionType = "dialog"
	ActionTypeConfirm  ActionType = "confirm"
	ActionTypeCancel   ActionType = "cancel"
	ActionTypePrev     ActionType = "prev"
	ActionTypeNext     ActionType = "next"
	ActionTypeCopy     ActionType = "copy"
	ActionTypeClose    ActionType = "close"
	ActionTypeDownload ActionType = "download"
	ActionTypeReload   ActionType = "reload"
)

// Action 行为按钮
// https://baidu.github.io/amis/zh-CN/components/action?page=1#%E6%B8%85%E7%A9%BA%E8%A1%A8%E5%8D%95
type Action struct {
	Base
	ActionType  ActionType  `json:"actionType"`            // 按钮类型
	Icon        string      `json:"icon,omitempty"`        // 按钮图标
	Label       string      `json:"label,omitempty"`       // 按钮文本
	Tooltip     string      `json:"tooltip,omitempty"`     // 按钮提示
	Level       ButtonLevel `json:"level,omitempty"`       // 按钮级别
	ConfirmText string      `json:"confirmText,omitempty"` // 确认操作文案提示
	API         *API        `json:"api,omitempty"`         // 按钮 API
	Drawer      *Drawer     `json:"drawer,omitempty"`      // 抽屉配置
	Dialog      *Dialog     `json:"dialog,omitempty"`      // 对话框配置
}

// NewAction 创建一个新的 Action 实例
func NewAction() *Action {
	return &Action{
		Base: Base{
			Type: TypeButton,
		},
	}
}

func (a *Action) WithAPI(api *API) *Action {
	a.ActionType = ActionTypeAjax
	a.API = api
	return a
}

func (a *Action) WithIcon(icon string) *Action {
	a.Icon = icon
	return a
}

func (a *Action) WithLabel(label string) *Action {
	a.Label = label
	return a
}

func (a *Action) WithLevel(level ButtonLevel) *Action {
	a.Level = level
	return a
}

func (a *Action) WithConfirmText(confirmText string) *Action {
	a.ConfirmText = confirmText
	return a
}

func (a *Action) WithTooltip(tooltip string) *Action {
	a.Tooltip = tooltip
	return a
}

func (a *Action) WithDrawer(drawer *Drawer) *Action {
	a.ActionType = ActionTypeDrawer
	a.Drawer = drawer
	return a
}

func (a *Action) WithDialog(dialog *Dialog) *Action {
	a.ActionType = ActionTypeDialog
	a.Dialog = dialog
	return a
}

type Drawer struct {
	Title   string    `json:"title,omitempty"`   // 标题
	Size    SizeLabel `json:"size,omitempty"`    // 指定 Drawer 大小，支持: xs、sm、md、lg、xl
	Body    Component `json:"body,omitempty"`    // 抽屉渲染内容
	Actions []*Action `json:"actions,omitempty"` // 对话框底部按钮
}

// NewDrawer 创建一个新的 Drawer 实例
func NewDrawer() *Drawer {
	return &Drawer{}
}

// WithTitle 设置 Drawer 标题
func (d *Drawer) WithTitle(title string) *Drawer {
	d.Title = title
	return d
}

// WithSize 设置 Drawer 大小
func (d *Drawer) WithSize(size SizeLabel) *Drawer {
	d.Size = size
	return d
}

// WithBody 设置 Drawer 内容
func (d *Drawer) WithBody(body Component) *Drawer {
	d.Body = body
	return d
}

// Dialog 对话框
// https://baidu.github.io/amis/zh-CN/components/dialog?page=1
type Dialog struct {
	Title   string    `json:"title,omitempty"`   // 标题
	Size    SizeLabel `json:"size,omitempty"`    // 指定 Drawer 大小，支持: xs、sm、md、lg、xl
	Body    Component `json:"body,omitempty"`    // 对话框渲染内容
	Actions []*Action `json:"actions,omitempty"` // 对话框底部按钮
}

// NewDialog 创建一个新的 Dialog 实例
func NewDialog() *Dialog {
	return &Dialog{}
}

// WithTitle 设置 Dialog 标题
func (d *Dialog) WithTitle(title string) *Dialog {
	d.Title = title
	return d
}

// WithSize 设置 Dialog 大小
func (d *Dialog) WithSize(size SizeLabel) *Dialog {
	d.Size = size
	return d
}

// WithBody 设置 Dialog 内容
func (d *Dialog) WithBody(body Component) *Dialog {
	d.Body = body
	return d
}
