package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// LogWriter writes log messages with timestamp, file info, level color, and the original message. Just for debugging / test purposes.
type LogWriter struct{}

func (lw LogWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimSpace(string(p))

	// Get file and line of the log call
	_, file, line, ok := runtime.Caller(3)
	fileInfo := ""
	if ok {
		fileInfo = fmt.Sprintf("[%s:%d]", filepath.Base(file), line)
	}

	// Detect level and color
	level := detectLogLevel(msg)
	coloredLevel := colorForLevel(level)

	timestamp := time.Now().Format(timeFormatStr)
	fmt.Printf("%s %s %s %s\n", timestamp, fileInfo, coloredLevel, msg)

	return len(p), nil
}

const timeFormatStr = "2006/01/02 15:04:05"

// Terminal color codes (ANSI)
type Color uint8

const (
	Black Color = iota + 30
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

func (c Color) apply(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", c, s)
}

// Log levels
type LogLevel string

const (
	LevelError LogLevel = "ERROR"
	LevelWarn  LogLevel = "WARN"
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
)

func colorForLevel(level LogLevel) string {
	switch level {
	case LevelError:
		return Red.apply("ERROR")
	case LevelWarn:
		return Yellow.apply("WARN")
	case LevelDebug:
		return Blue.apply("DEBUG")
	case LevelInfo:
		fallthrough
	default:
		return Cyan.apply("INFO")
	}
}

// detectLogLevel tries to determine the log level from message content.
func detectLogLevel(msg string) LogLevel {
	upper := strings.ToUpper(msg)
	switch {
	case strings.Contains(upper, "ERROR"),
		strings.Contains(upper, "FAILED"),
		strings.Contains(upper, "FAIL"),
		strings.Contains(upper, "FATAL"),
		strings.Contains(upper, "PANIC"),
		strings.Contains(upper, "CRITICAL"):
		return LevelError
	case strings.Contains(upper, "WARN"),
		strings.Contains(upper, "WARNING"):
		return LevelWarn
	case strings.Contains(upper, "DEBUG"),
		strings.Contains(upper, "DEBUGGING"),
		strings.Contains(upper, "TRACE"):
		return LevelDebug
	default:
		return LevelInfo
	}
}

func RecoverPanic() func() {
	return func() {
		if r := recover(); r != nil {
			msg := fmt.Sprintf("panic: %v", r)
			log.Print(msg)
			os.Exit(1)
		}
	}
}
