package ylog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

type FileLogger struct {
	lastHour int64
	file     *os.File
	Level    LogLevel
	mu       sync.Mutex
	iLogger  *log.Logger
	Path     string
}

func getTimeHour(t *time.Time) int64 {
	return t.Unix() / 3600
}

func getFileName(t *time.Time) string {
	return t.Format("2006-01-02_15")
}

func createFile(path *string, t *time.Time) (file *os.File, err error) {
	dir := filepath.Join(*path, t.Format("200601"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766)
		if err != nil {
			return nil, err
		}
	}

	filePath := filepath.Join(dir, getFileName(t)+".txt")
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	return
}

var stdPath = "logs"
var std = NewFileLogger(LevelInfo, stdPath)

func SetLevel(level LogLevel) {
	std.SetLevel(level)
}

func Trace(v ...any) {
	if std.CanTrace() {
		std.ensureFile()
		v = append([]any{"Trace"}, v...)
		std.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Debug(v ...any) {
	if std.CanDebug() {
		std.ensureFile()
		v = append([]any{"Debug"}, v...)
		std.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Info(v ...any) {
	if std.CanInfo() {
		std.ensureFile()
		v = append([]any{"Info"}, v...)
		std.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Warn(v ...any) {
	if std.CanWarn() {
		std.ensureFile()
		v = append([]any{"Warn"}, v...)
		std.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Error(v ...any) {
	if std.CanError() {
		std.ensureFile()
		v = append([]any{"Error"}, v...)
		std.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func Fatal(v ...any) {
	if std.CanFatal() {
		std.ensureFile()
		v = append([]any{"Fatal"}, v...)
		std.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func NewFileLogger(level LogLevel, path string) (logger *FileLogger) {
	logger = &FileLogger{}
	logger.iLogger = log.New(os.Stderr, "", log.LstdFlags)
	logger.Level = level
	logger.Path = path
	return
}

func (l *FileLogger) ensureFile() (err error) {
	currentTime := time.Now()
	if l.file == nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.file == nil {
			l.file, err = createFile(&l.Path, &currentTime)
			l.iLogger.SetOutput(l.file)
			l.iLogger.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)
			l.lastHour = getTimeHour(&currentTime)
		}
		return
	}

	currentHour := getTimeHour(&currentTime)
	if l.lastHour != currentHour {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.lastHour != currentHour {
			_ = l.file.Close()
			l.file, err = createFile(&l.Path, &currentTime)
			l.iLogger.SetOutput(l.file)
			l.iLogger.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
			l.lastHour = getTimeHour(&currentTime)
		}
	}

	return
}

func (l *FileLogger) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *FileLogger) CanTrace() bool {
	return l.Level <= LevelTrace
}

func (l *FileLogger) CanDebug() bool {
	return l.Level <= LevelDebug
}

func (l *FileLogger) CanInfo() bool {
	return l.Level <= LevelInfo
}

func (l *FileLogger) CanWarn() bool {
	return l.Level <= LevelWarn
}

func (l *FileLogger) CanError() bool {
	return l.Level <= LevelError
}

func (l *FileLogger) CanFatal() bool {
	return l.Level <= LevelFatal
}

func (l *FileLogger) Trace(v ...any) {
	if l.CanTrace() {
		l.ensureFile()
		v = append([]any{"Trace "}, v...)
		l.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Debug(v ...any) {
	if l.CanDebug() {
		l.ensureFile()
		v = append([]any{"Debug "}, v...)
		l.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Info(v ...any) {
	if l.CanInfo() {
		l.ensureFile()
		v = append([]any{"Info"}, v...)
		l.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Warn(v ...any) {
	if l.CanWarn() {
		l.ensureFile()
		v = append([]any{"Warn"}, v...)
		l.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Error(v ...any) {
	if l.CanError() {
		l.ensureFile()
		v = append([]any{"Error"}, v...)
		l.iLogger.Output(2, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Fatal(v ...any) {
	if l.CanFatal() {
		l.ensureFile()
		v = append([]any{"Fatal"}, v...)
		l.iLogger.Output(2, fmt.Sprintln(v...))
	}
}
