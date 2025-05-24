package amisgo

// Component 基础接口定义
type Component interface {
	GetType() Type
}

// Base amis 组件
type Base struct {
	Type          Type           `json:"type,omitempty"`          // 组件类型
	Name          string         `json:"name,omitempty"`          // 组件自定义名称
	CssVars       Css            `json:"cssVars,omitempty"`       // 自定义 CSS 变量
	Css           map[string]Css `json:"css,omitempty"`           // 自定义 CSS
	ClassName     string         `json:"className,omitempty"`     // 自定义 CSS 类名
	BodyClassName string         `json:"bodyClassName,omitempty"` // 自定义 CSS 类名
	Tpl           string         `json:"tpl,omitempty"`           // 模板字符串
}

// Css 自定义 CSS
type Css map[string]string

func (c *Base) GetType() Type {
	return c.Type
}

// NativeComponent 原生 amis 组件，用来定义一些 amisgo 还未支持的组件
type NativeComponent struct {
	Type   Type
	Schema string // amis 组件 schema，必须符合 amis 协议规范
}

func (c *NativeComponent) GetType() Type {
	return c.Type
}

func (s *NativeComponent) MarshalJSON() ([]byte, error) {
	return []byte(s.Schema), nil
}

// NewNativeComponent 创建原生 amis 组件
func NewNativeComponent(t string, schema string) *NativeComponent {
	return &NativeComponent{
		Type:   Type(t),
		Schema: schema,
	}
}
