package log

import (
	"sync"

	"github.com/fengjx/go-halo/logger"
)

var (
	clog     Logger
	clogOnce sync.Once
)

// Logger 日志输出接口
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Error(msg string)
	Panic(msg string)

	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Panicf(format string, args ...any)
}

// UseLogger 使用日志输出
func UseLogger(l Logger) {
	clog = l
}

// getLogger 获取日志输出
func getLogger() Logger {
	clogOnce.Do(func() {
		if clog == nil {
			clog = &consoleLogger{
				log: logger.NewConsole(),
			}
			clog.Info("use console logger")
		}
	})
	return clog
}

type consoleLogger struct {
	log logger.Logger
}

func (l *consoleLogger) Debug(msg string) {
	l.log.Debug(msg)
}

func (l *consoleLogger) Info(msg string) {
	l.log.Info(msg)
}

func (l *consoleLogger) Error(msg string) {
	l.log.Error(msg)
}

func (l *consoleLogger) Panic(msg string) {
	l.log.Panic(msg)
}

func (l *consoleLogger) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *consoleLogger) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *consoleLogger) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *consoleLogger) Panicf(format string, args ...any) {
	l.log.Panicf(format, args...)
}

func Debug(msg string) {
	getLogger().Debug(msg)
}

func Debugf(format string, args ...any) {
	getLogger().Debugf(format, args...)
}

func Info(msg string) {
	getLogger().Info(msg)
}

func Infof(format string, args ...any) {
	getLogger().Infof(format, args...)
}

func Error(msg string) {
	getLogger().Error(msg)
}

func Errorf(format string, args ...any) {
	getLogger().Errorf(format, args...)
}

func Panic(msg string) {
	getLogger().Panic(msg)
}

func Panicf(format string, args ...any) {
	getLogger().Panicf(format, args...)
}
