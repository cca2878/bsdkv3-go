package log

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 日志级别
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

var (
	// 默认日志级别
	level = LevelInfo

	// 标准库日志实例
	logger = log.New(os.Stdout, "", 0)

	// 级别名称映射
	levelNames = map[int]string{
		LevelDebug: "DEBUG",
		LevelInfo:  "INFO",
		LevelWarn:  "WARN",
		LevelError: "ERROR",
	}
)

// SetLevel 设置日志级别
func SetLevel(l int) {
	if l >= LevelDebug && l <= LevelError {
		level = l
	}
}

// SetOutput 设置日志输出位置
func SetOutput(w *os.File) {
	logger.SetOutput(w)
}

// 格式化日志前缀
func formatPrefix(lvl int) string {
	levelName, ok := levelNames[lvl]
	if !ok {
		levelName = "UNKNOWN"
	}
	return fmt.Sprintf("[%s][%s] ", time.Now().Format("2006-01-02 15:04:05"), levelName)
}

// Debug 记录调试日志
func Debug(format string, v ...interface{}) {
	if level <= LevelDebug {
		logger.Print(formatPrefix(LevelDebug) + fmt.Sprintf(format, v...))
	}
}

// Info 记录信息日志
func Info(format string, v ...interface{}) {
	if level <= LevelInfo {
		logger.Print(formatPrefix(LevelInfo) + fmt.Sprintf(format, v...))
	}
}

// Warn 记录警告日志
func Warn(format string, v ...interface{}) {
	if level <= LevelWarn {
		logger.Print(formatPrefix(LevelWarn) + fmt.Sprintf(format, v...))
	}
}

// Error 记录错误日志
func Error(format string, v ...interface{}) {
	if level <= LevelError {
		logger.Print(formatPrefix(LevelError) + fmt.Sprintf(format, v...))
	}
}
