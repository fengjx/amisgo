package amisgo

// Page 表示 amis 页面组件
type Page struct {
	Base
	Title    string      `json:"title,omitempty"`    // 页面标题
	SubTitle string      `json:"subTitle,omitempty"` // 页面副标题
	Body     []Component `json:"body,omitempty"`     // 内容区域
}

// NewPage 创建一个新的 Page 实例
func NewPage() *Page {
	return &Page{
		Base: Base{
			Type: TypePage,
		},
	}
}
