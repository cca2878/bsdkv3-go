package bsdkv3

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger 由接入方实现，SDK 内部不直接向标准输出打印日志。
type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Warn(format string, args ...any)
	Error(format string, args ...any)
}

// nopLogger 丢弃所有日志，为未配置 Logger 时的默认实现（非导出）。
type nopLogger struct{}

func (nopLogger) Debug(string, ...any) {}
func (nopLogger) Info(string, ...any)  {}
func (nopLogger) Warn(string, ...any)  {}
func (nopLogger) Error(string, ...any) {}

// stdLogger 将日志写入指定 Writer（通常为 os.Stderr），供需要默认输出的场景使用。
type stdLogger struct {
	w     io.Writer
	level LogLevel
}

// LogLevel 定义了日志输出级别。
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// NewStdLogger 创建写入 w 的标准格式日志器，level 为最低输出级别。
func NewStdLogger(w io.Writer, level LogLevel) Logger {
	if w == nil {
		w = os.Stderr
	}
	if level < LogLevelDebug {
		level = LogLevelDebug
	}
	if level > LogLevelError {
		level = LogLevelError
	}
	return &stdLogger{w: w, level: level}
}

var levelNames = map[LogLevel]string{
	LogLevelDebug: "DEBUG",
	LogLevelInfo:  "INFO",
	LogLevelWarn:  "WARN",
	LogLevelError: "ERROR",
}

func (l *stdLogger) log(lvl LogLevel, format string, args ...any) {
	if lvl < l.level {
		return
	}
	name := levelNames[lvl]
	if name == "" {
		name = "UNKNOWN"
	}
	prefix := fmt.Sprintf("[%s][%s] ", time.Now().Format("2006-01-02 15:04:05"), name)
	_, _ = fmt.Fprintln(l.w, prefix+fmt.Sprintf(format, args...))
}

func (l *stdLogger) Debug(format string, args ...any) { l.log(LogLevelDebug, format, args...) }
func (l *stdLogger) Info(format string, args ...any)  { l.log(LogLevelInfo, format, args...) }
func (l *stdLogger) Warn(format string, args ...any)  { l.log(LogLevelWarn, format, args...) }
func (l *stdLogger) Error(format string, args ...any) { l.log(LogLevelError, format, args...) }

// discardLogger 供未注入 Logger 时使用。
var discardLogger Logger = nopLogger{}

var _ Logger = (*stdLogger)(nil)
