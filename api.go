package amisgo

// API amis API 配置
// https://baidu.github.io/amis/zh-CN/docs/types/api
type API struct {
	Method      string            `json:"method,omitempty"`      // 请求方式 GET/POST/PUT/DELETE
	URL         string            `json:"url,omitempty"`         // 请求地址
	Data        any               `json:"data,omitempty"`        // 请求数据，支持数据映射
	Headers     map[string]string `json:"headers,omitempty"`     // 请求头
	SendOn      string            `json:"sendOn,omitempty"`      // 发送条件
	AutoRefresh bool              `json:"autoRefresh,omitempty"` // 是否自动刷新
	DataType    string            `json:"dataType,omitempty"`    // 数据体格式，默认为 json
	Cache       int               `json:"cache,omitempty"`       // 接口缓存时间，单位毫秒
}

// NewAPI 创建一个新的 API 实例
func NewAPI() *API {
	return &API{
		Method:   "GET",
		DataType: "json",
	}
}

// SetMethod 设置请求方法
func (a *API) SetMethod(method string) *API {
	a.Method = method
	return a
}

// SetURL 设置请求地址
func (a *API) SetURL(url string) *API {
	a.URL = url
	return a
}

// SetData 设置请求数据
func (a *API) SetData(data any) *API {
	a.Data = data
	return a
}

// SetHeaders 设置请求头
func (a *API) SetHeaders(headers map[string]string) *API {
	a.Headers = headers
	return a
}

// SetSendOn 设置发送条件
func (a *API) SetSendOn(condition string) *API {
	a.SendOn = condition
	return a
}

// SetAutoRefresh 设置是否自动刷新
func (a *API) SetAutoRefresh(autoRefresh bool) *API {
	a.AutoRefresh = autoRefresh
	return a
}

// SetDataType 设置数据类型
func (a *API) SetDataType(dataType string) *API {
	a.DataType = dataType
	return a
}

// SetCache 设置缓存时间
func (a *API) SetCache(milliseconds int) *API {
	a.Cache = milliseconds
	return a
}
