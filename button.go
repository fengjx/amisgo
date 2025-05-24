package amisgo

// ButtonLevel 按钮级别
// 'link' | 'primary' | 'enhance' | 'secondary' | 'info'|'success' | 'warning' | 'danger' | 'light'| 'dark' | 'default'
type ButtonLevel string

const (
	ButtonLevelDefault   ButtonLevel = "default"
	ButtonLevelPrimary   ButtonLevel = "primary"
	ButtonLevelDanger    ButtonLevel = "danger"
	ButtonLevelLink      ButtonLevel = "link"
	ButtonLevelEnhance   ButtonLevel = "enhance"
	ButtonLevelSecondary ButtonLevel = "secondary"
	ButtonLevelInfo      ButtonLevel = "info"
	ButtonLevelSuccess   ButtonLevel = "success"
	ButtonLevelWarning   ButtonLevel = "warning"
	ButtonLevelLight     ButtonLevel = "light"
	ButtonLevelDark      ButtonLevel = "dark"
)

type ButtonActionType string

const (
	ButtonActionTypeButton ButtonActionType = "button"
	ButtonActionTypeReset  ButtonActionType = "reset"
	ButtonActionTypeSubmit ButtonActionType = "submit"
	ButtonActionTypeClear  ButtonActionType = "clear"
	ButtonActionTypeURL    ButtonActionType = "url"
)

// Button 按钮
// https://baidu.github.io/amis/zh-CN/components/button
type Button struct {
	Base
	ActionType  ButtonActionType `json:"actionType,omitempty"`  // 按钮类型 'button' | 'reset' | 'submit'| 'clear'| 'url'
	Icon        string           `json:"icon,omitempty"`        // 按钮图标
	Label       string           `json:"label,omitempty"`       // 按钮文本
	Tooltip     string           `json:"tooltip,omitempty"`     // 按钮提示
	Level       ButtonLevel      `json:"level,omitempty"`       // 按钮级别
	ConfirmText string           `json:"confirmText,omitempty"` // 确认文本
}

// NewButton 创建一个新的 Button 实例
func NewButton(at ButtonActionType) *Button {
	return &Button{
		Base: Base{
			Type: TypeButton,
		},
		ActionType: at,
	}
}

func (b *Button) WithIcon(icon string) *Button {
	b.Icon = icon
	return b
}

func (b *Button) WithLabel(label string) *Button {
	b.Label = label
	return b
}

func (b *Button) WithLevel(level ButtonLevel) *Button {
	b.Level = level
	return b
}

func (b *Button) WithConfirmText(confirmText string) *Button {
	b.ConfirmText = confirmText
	return b
}

func (b *Button) WithTooltip(tooltip string) *Button {
	b.Tooltip = tooltip
	return b
}
