package ylog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

type LogLevel int

const (
	ToTrace LogLevel = iota
	ToDebug
	ToInfo
	ToWarn
	ToError
	ToFatal
)

type Logger interface {
	Trace(format string, v ...any)
	Debug(format string, v ...any)
	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
	Fatal(format string, v ...any)
}

type FileLogger struct {
	lastHour int64
	file     *os.File
	Level    LogLevel
	mu       sync.Mutex
	iLogger  *log.Logger
	Path     *string
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
		err = os.MkdirAll(dir, 0666)
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
var std = NewFileLogger(ToInfo, &stdPath)

func SetLevel(level LogLevel) {
	std.SetLevel(level)
}

func Trace(v ...any) {
	std.Trace(v...)
}

func Debug(v ...any) {
	std.Debug(v...)
}

func Info(v ...any) {
	std.Info(v...)
}

func Warn(v ...any) {
	std.Warn(v...)
}

func Error(v ...any) {
	std.Error(v...)
}

func Fatal(v ...any) {
	std.Fatal(v...)
}

func NewFileLogger(level LogLevel, path *string) (logger *FileLogger) {
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
			l.file, err = createFile(l.Path, &currentTime)
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
			l.file, err = createFile(l.Path, &currentTime)
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
	return l.Level <= ToTrace
}

func (l *FileLogger) CanDebug() bool {
	return l.Level <= ToDebug
}

func (l *FileLogger) CanInfo() bool {
	return l.Level <= ToInfo
}

func (l *FileLogger) CanWarn() bool {
	return l.Level <= ToWarn
}

func (l *FileLogger) CanError() bool {
	return l.Level <= ToError
}

func (l *FileLogger) CanFatal() bool {
	return l.Level <= ToFatal
}

func (l *FileLogger) Trace(v ...any) {
	if l.CanTrace() {
		l.ensureFile()
		v[0] = "[Trace] " + toString(v[0])
		l.iLogger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Debug(v ...any) {
	if l.CanDebug() {
		l.ensureFile()
		v[0] = "[Debug] " + toString(v[0])
		l.iLogger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Info(v ...any) {
	if l.CanInfo() {
		l.ensureFile()
		v[0] = "[Info] " + toString(v[0])
		l.iLogger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Warn(v ...any) {
	if l.CanWarn() {
		l.ensureFile()
		v[0] = "[Warn] " + toString(v[0])
		l.iLogger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Error(v ...any) {
	if l.CanError() {
		l.ensureFile()
		v[0] = "[Error] " + toString(v[0])
		l.iLogger.Output(3, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Fatal(v ...any) {
	if l.CanFatal() {
		l.ensureFile()
		v[0] = "[Fatal] " + toString(v[0])
		l.iLogger.Output(3, fmt.Sprintln(v...))
	}
}

func toString(v any) string {
	switch v := v.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return reflect.TypeOf(v).String()
	}
}
