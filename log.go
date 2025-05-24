package amisgo

import "github.com/fengjx/amisgo/internal/log"

// UseLogger 使用日志输出
func UseLogger(l log.Logger) {
	log.UseLogger(l)
}
