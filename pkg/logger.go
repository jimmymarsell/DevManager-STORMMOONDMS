package pkg

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fatih/color"
)

var (
	Success = color.New(color.FgGreen, color.Bold)
	Error   = color.New(color.FgRed, color.Bold)
	Warning = color.New(color.FgYellow, color.Bold)
	Info    = color.New(color.FgCyan, color.Bold)
)

var (
	logFile *os.File
	logMut  sync.Mutex
)

func InitLogger(installDir string) error {
	logDir := filepath.Join(installDir, "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	logPath := filepath.Join(logDir, "stormmoondms.log")
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}

	logFile = f
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.Ltime)
	return nil
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

func writeLog(level, msg string) {
	if logFile != nil {
		logMut.Lock()
		log.Printf("[%s] %s", level, msg)
		logMut.Unlock()
	}
}

func PrintSuccess(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Success.Println(msg)
	writeLog("SUCCESS", msg)
}

func PrintError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Error.Println(msg)
	writeLog("ERROR", msg)
}

func PrintWarning(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Warning.Println(msg)
	writeLog("WARNING", msg)
}

func PrintInfo(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	Info.Println(msg)
	writeLog("INFO", msg)
}

func GetOS() string {
	return runtime.GOOS
}