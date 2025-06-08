package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Level тип для уровней логирования
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
	FATAL
)

var levelStrings = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// Color codes для терминала
var levelColors = []string{
	"\033[36m", // DEBUG - cyan
	"\033[32m", // INFO - green
	"\033[33m", // WARN - yellow
	"\033[31m", // ERROR - red
	"\033[35m", // FATAL - magenta
}

const colorReset = "\033[0m"

// Logger структура логгера
type Logger struct {
	mu         sync.Mutex
	out        *log.Logger
	level      Level
	jsonOutput bool
	showCaller bool
	color      bool
}

// Config структура для настройки логгера
type Config struct {
	Level      Level
	JsonOutput bool
	ShowCaller bool
	Color      bool
	OutputFile string // если пустая строка — вывод в stdout
	MaxSizeMB  int    // макс размер файла для ротации (MB)
	MaxBackups int    // кол-во резервных файлов
	MaxAgeDays int    // максимальный возраст файла в днях
	Compress   bool   // сжимать старые файлы
}

// New создаёт новый логгер по конфигу
func New(cfg Config) *Logger {
	var writer io.Writer

	if cfg.OutputFile != "" {
		writer = &lumberjack.Logger{
			Filename:   cfg.OutputFile,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
	} else {
		writer = os.Stdout
	}

	return &Logger{
		out:        log.New(writer, "", 0), // форматирование
		level:      cfg.Level,
		jsonOutput: cfg.JsonOutput,
		showCaller: cfg.ShowCaller,
		color:      cfg.Color,
	}
}

func (l *Logger) log(level Level, msg string) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	var callerInfo string
	if l.showCaller {
		_, file, line, ok := runtime.Caller(3)
		if ok {
			shortFile := file[strings.LastIndex(file, "/")+1:]
			callerInfo = fmt.Sprintf("%s:%d", shortFile, line)
		}
	}

	now := time.Now().Format(time.RFC3339)

	if l.jsonOutput {
		entry := map[string]interface{}{
			"time":    now,
			"level":   levelStrings[level],
			"message": msg,
		}
		if l.showCaller && callerInfo != "" {
			entry["caller"] = callerInfo
		}
		data, _ := json.Marshal(entry)
		l.out.Println(string(data))
		return
	}

	// Текстовый формат
	prefix := fmt.Sprintf("[%s] %s", levelStrings[level], now)
	if l.showCaller && callerInfo != "" {
		prefix += " " + callerInfo
	}

	if l.color {
		color := levelColors[level]
		l.out.Println(color + prefix + " " + msg + colorReset)
	} else {
		l.out.Println(prefix + " " + msg)
	}

}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, fmt.Sprintf(format, args...))
}
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, fmt.Sprintf(format, args...))
}
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, fmt.Sprintf(format, args...))
}
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, fmt.Sprintf(format, args...))
}
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, fmt.Sprintf(format, args...))
	os.Exit(1)
}
