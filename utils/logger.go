package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// Logger 日志记录器
type Logger struct {
	logFile *os.File
	logger  *log.Logger
}

// NewLogger 创建新的日志记录器
func NewLogger() *Logger {
	// 创建日志目录
	logDir := "logs"
	os.MkdirAll(logDir, 0755)
	
	// 创建日志文件
	timestamp := time.Now().Format("20060102_150405")
	logFileName := fmt.Sprintf("system_reinstaller_%s.log", timestamp)
	logPath := filepath.Join(logDir, logFileName)
	
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("无法创建日志文件: %v", err)
		return &Logger{
			logger: log.New(os.Stdout, "", log.LstdFlags),
		}
	}
	
	return &Logger{
		logFile: logFile,
		logger:  log.New(logFile, "", log.LstdFlags),
	}
}

// Info 记录信息日志
func (l *Logger) Info(message string) {
	l.logger.Printf("[INFO] %s", message)
	fmt.Printf("[INFO] %s\n", message)
}

// Error 记录错误日志
func (l *Logger) Error(message string) {
	l.logger.Printf("[ERROR] %s", message)
	fmt.Printf("[ERROR] %s\n", message)
}

// Warning 记录警告日志
func (l *Logger) Warning(message string) {
	l.logger.Printf("[WARNING] %s", message)
	fmt.Printf("[WARNING] %s\n", message)
}

// Debug 记录调试日志
func (l *Logger) Debug(message string) {
	l.logger.Printf("[DEBUG] %s", message)
	fmt.Printf("[DEBUG] %s\n", message)
}

// Close 关闭日志文件
func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}