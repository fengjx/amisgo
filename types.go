package amisgo

type Type string

const (
	TypeTemplate      Type = "template"
	TypePage          Type = "page"
	TypeCRUD          Type = "crud"
	TypeGroup         Type = "group"
	TypeForm          Type = "form"
	TypeInputText     Type = "input-text"
	TypeTextarea      Type = "textarea"
	TypeInputNumber   Type = "input-number"
	TypeInputDate     Type = "input-date"
	TypeInputDatetime Type = "input-datetime"
	TypeSelect        Type = "select"
	TypeInputTree     Type = "input-tree"
	TypeTreeSelect    Type = "tree-select"
	TypeButton        Type = "button"
	TypeReset         Type = "reset"
	TypeSubmit        Type = "submit"
	TypeClear         Type = "clear"
	TypeStatic        Type = "static"
)

// Messages 消息提示覆写，默认消息读取的是 API 返回的消息，但是在此可以覆写它。
type Messages struct {
	FetchSuccess string `json:"fetchSuccess,omitempty"` // 获取成功时提示
	FetchFailed  string `json:"fetchFailed,omitempty"`  // 获取失败时提示
	SaveSuccess  string `json:"saveSuccess,omitempty"`  // 保存成功时提示
	SaveFailed   string `json:"saveFailed,omitempty"`   // 保存失败时提示
}

// Horizontal 水平布局配置
type Horizontal struct {
	Left    int  `json:"left,omitempty"`    // 左侧宽度占比
	Right   int  `json:"right,omitempty"`   // 右侧宽度占比
	Justify bool `json:"justify,omitempty"` // 是否开启水平居中
}

// Size 抽屉大小
type SizeLabel string

const (
	SizeXs SizeLabel = "xs"
	SizeSm SizeLabel = "sm"
	SizeMd SizeLabel = "md"
	SizeLg SizeLabel = "lg"
	SizeXl SizeLabel = "xl"
)
