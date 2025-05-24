package amisgo

// Group 用于实现表单项分组功能
// https://baidu.github.io/amis/zh-CN/components/form/group
type Group struct {
	Base
	Body []Component `json:"body,omitempty"` // 表单项集合
}

// NewGroup 创建一个新的 Group 实例
func NewGroup() *Group {
	return &Group{
		Base: Base{
			Type: TypeGroup,
		},
	}
}

// WithBody 设置表单项集合
func (g *Group) WithBody(body []Component) *Group {
	g.Body = body
	return g
}
