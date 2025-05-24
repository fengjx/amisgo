package kit

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

const (
	contentTypeJSON = "application/json"
	contentTypeForm = "application/x-www-form-urlencoded"
)

var (
	trueClientIP  = http.CanonicalHeaderKey("True-Client-IP")
	xForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	xRealIP       = http.CanonicalHeaderKey("X-Real-IP")
)

// ShouldBindJSON 绑定JSON请求
func ShouldBindJSON(r *http.Request, obj any) error {
	return json.NewDecoder(r.Body).Decode(obj)
}

// ParseParams 解析请求参数
func ParseParams(r *http.Request) (map[string]string, error) {
	m := make(map[string]string)
	for k, v := range r.URL.Query() {
		m[k] = v[0]
	}
	ct := r.Header.Get("Content-Type")
	if ct == contentTypeForm {
		err := r.ParseForm()
		if err != nil {
			return nil, err
		}
		for k, v := range r.Form {
			m[k] = v[0]
		}
	}
	return m, nil
}

// GetParam 获取URL查询参数，如果参数不存在返回默认值
func GetParam(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value != "" {
		return value
	}
	return GetForm(r, key, defaultValue)
}

// GetQuery 获取URL查询参数，如果参数不存在返回默认值
func GetQuery(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetForm 获取表单参数，如果参数不存在返回默认值
func GetForm(r *http.Request, key, defaultValue string) string {
	_ = r.ParseForm()
	value := r.FormValue(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetHeader 获取请求头，如果不存在返回空字符串
func GetHeader(r *http.Request, key string) string {
	return r.Header.Get(key)
}

// GetHeaderDefault 获取请求头，如果不存在返回默认值
func GetHeaderDefault(r *http.Request, key, defaultValue string) string {
	value := r.Header.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetCookie 获取Cookie值，如果不存在返回空字符串和错误
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetCookieDefault 获取Cookie值，如果不存在返回默认值
func GetCookieDefault(r *http.Request, name, defaultValue string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		return defaultValue
	}
	return cookie.Value
}

// Map is a map of string to any.
type Map map[string]any

// GetRealIP 从 request 获取真实客户端ip
func GetRealIP(r *http.Request) string {
	var ip string

	if tcip := r.Header.Get(trueClientIP); tcip != "" {
		ip = tcip
	} else if xrip := r.Header.Get(xRealIP); xrip != "" {
		ip = xrip
	} else if xff := r.Header.Get(xForwardedFor); xff != "" {
		i := strings.Index(xff, ",")
		if i == -1 {
			i = len(xff)
		}
		ip = xff[:i]
	}
	if ip == "" || net.ParseIP(ip) == nil {
		return ""
	}
	return ip
}

// Write 写入响应内容
func Write(w http.ResponseWriter, code int, contentType string, message string) error {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(code)
	_, err := w.Write([]byte(message))
	return err
}

// WriteJSON 写入JSON响应
func WriteJSON(w http.ResponseWriter, code int, data any) error {
	w.Header().Set("Content-Type", contentTypeJSON)
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(data)
}

// WriteNoContent 只返回响应码，不返回内容
func WriteNoContent(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return nil
}
