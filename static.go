package amisgo

// Static 静态展示
// https://aisuda.bce.baidu.com/amis/zh-CN/components/form/static
type Static struct {
	Base
	Name  string `json:"name,omitempty"`  // 名称
	Label string `json:"label,omitempty"` // 标签
}

func NewStatic() *Static {
	return &Static{
		Base: Base{
			Type: TypeStatic,
		},
	}
}

// WithName 设置静态展示的名称
func (s *Static) WithName(name string) *Static {
	s.Name = name
	return s
}

// WithLabel 设置静态展示的标签
func (s *Static) WithLabel(label string) *Static {
	s.Label = label
	return s
}
