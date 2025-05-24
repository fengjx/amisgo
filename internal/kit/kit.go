package kit

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unsafe"
)

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func toASCIIUpper(r rune) rune {
	if 'a' <= r && r <= 'z' {
		r -= 'a' - 'A'
	}
	return r
}

func LineString(str string) string {
	if str == "" {
		return "-"
	}
	return str
}

// FirstUpper 首字母大写
func FirstUpper(str string) string {
	if str == "" {
		return ""
	}
	newstr := make([]byte, 0, len(str)+1)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if i == 0 {
			c = strings.ToUpper(str[:1])[0]
		}
		newstr = append(newstr, c)
	}
	return b2s(newstr)
}

// FirstLower 首字母小写
func FirstLower(str string) string {
	if str == "" {
		return ""
	}
	newstr := make([]byte, 0, len(str)+1)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if i == 0 {
			c = strings.ToLower(str[:1])[0]
		}
		newstr = append(newstr, c)
	}
	return b2s(newstr)
}

// SnakeCase 下划线命名
func SnakeCase(str string) string {
	newstr := make([]byte, 0, len(str)+1)
	for i := 0; i < len(str); i++ {
		c := str[i]
		if isUpper := 'A' <= c && c <= 'Z'; isUpper {
			if i > 0 {
				newstr = append(newstr, '_')
			}
			c += 'a' - 'A'
		}
		newstr = append(newstr, c)
	}
	return b2s(newstr)
}

// TitleCase 下划线转驼峰
func TitleCase(str string) string {
	newstr := make([]byte, 0, len(str))
	upNextChar := true

	str = strings.ToLower(str)

	for i := 0; i < len(str); i++ {
		c := str[i]
		switch {
		case upNextChar:
			upNextChar = false
			if 'a' <= c && c <= 'z' {
				c -= 'a' - 'A'
			}
		case c == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, c)
	}

	return b2s(newstr)
}

// GonicCase golang 风格的驼峰命名
func GonicCase(str string) string {
	newstr := make([]rune, 0)
	str = strings.ToLower(str)
	parts := strings.Split(str, "_")
	for _, p := range parts {
		_, isInitialism := LintGonicMapper[strings.ToUpper(p)]
		for i, r := range p {
			if i == 0 || isInitialism {
				r = toASCIIUpper(r)
			}
			newstr = append(newstr, r)
		}
	}
	return string(newstr)
}

// KebabCase kebab case 转换函数
func KebabCase(s string) string {
	if s == "" {
		return ""
	}
	var result strings.Builder
	var prev rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 && ((prev != 0 && (unicode.IsLower(prev) || unicode.IsDigit(prev))) || (prev != 0 && unicode.IsUpper(prev) && i+1 < len(s) && unicode.IsLower(rune(s[i+1])))) {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
		prev = r
	}
	return result.String()
}

// ToLowerAndTrim 将字符串转换为小写，并移除特殊字符（如 - 和 _）
func ToLowerAndTrim(s string) string {
	var result strings.Builder
	for _, r := range strings.ToLower(s) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ToKebabCase 通用转横线分割（kebab-case）
// 支持驼峰、下划线等格式
func ToKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if r == '_' {
			result.WriteRune('-')
			continue
		}
		if i > 0 && unicode.IsUpper(r) && (unicode.IsLower(rune(s[i-1])) || unicode.IsDigit(rune(s[i-1]))) {
			result.WriteRune('-')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

var LintGonicMapper = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SLA":   true,
	"SMTP":  true,
	"SSH":   true,
	"TLS":   true,
	"TTL":   true,
	"UI":    true,
	"UID":   true,
	"UUID":  true,
	"URI":   true,
	"URL":   true,
	"UTF8":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

// IsEmptyOrDefault 判断是否为空或者默认值
func IsEmptyOrDefault(id any) bool {
	if id == nil {
		return true
	}
	idStr := ToString(id)
	return idStr == "" || idStr == "0"
}

func ToString(src any) string {
	if src == nil {
		return ""
	}

	switch v := src.(type) {
	case string:
		return src.(string)
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", src)
	case float32, float64:
		bs, _ := json.Marshal(v)
		return string(bs)
	case bool:
		if b, ok := src.(bool); ok && b {
			return "true"
		} else {
			return "false"
		}
	default:
		return fmt.Sprintf("%v", v)
	}
}

// Contains 判断元素是否在集合中
func Contains[T comparable](collection []T, element T) bool {
	for _, item := range collection {
		if item == element {
			return true
		}
	}
	return false
}

// ToAnySlice 将切片转换为 []any
func ToAnySlice[T any](collection []T) []any {
	result := make([]any, len(collection))
	for i, item := range collection {
		result[i] = item
	}
	return result
}

// SplitTrim 分割字符串并去掉空格
func SplitTrim(input, sep string) []string {
	slc := strings.Split(input, sep)
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return slc
}

// SplitToSlice 分割字符串，去掉空格并遍历
func SplitToSlice[T any](input, sep string, fn func(item string) T) []T {
	slc := strings.Split(input, sep)
	var res []T
	for i := range slc {
		res = append(res, fn(strings.TrimSpace(slc[i])))
	}
	return res
}

// Omit 忽略指定键
func Omit(data map[string]any, keys ...string) map[string]any {
	if len(keys) == 0 {
		return data
	}
	m := make(map[string]any)
	for k, v := range data {
		if Contains(keys, k) {
			continue
		}
		m[k] = v
	}
	return m
}

// QuicklyClose 关闭 io.Closer 资源，忽略错误
func QuicklyClose(c io.Closer) {
	if c != nil {
		_ = c.Close()
	}
}
