package ylog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

var levelMap map[LogLevel][]byte = map[LogLevel][]byte{
	LevelTrace: []byte("Trace"),
	LevelDebug: []byte("Debug"),
	LevelInfo:  []byte("Info"),
	LevelWarn:  []byte("Warn"),
	LevelError: []byte("Error"),
	LevelFatal: []byte("Fatal"),
}

type Logger interface {
	Trace(format string, v ...any)
	Debug(format string, v ...any)
	Info(format string, v ...any)
	Warn(format string, v ...any)
	Error(format string, v ...any)
	Fatal(format string, v ...any)
}

type FileLogger struct {
	Level     LogLevel
	Path      string
	lastHour  int64
	file      *os.File
	mu        sync.Mutex
	pipe      chan *logEntry
	cacheSize uint16
	Layout    string
	exit      chan struct{}
}

type logEntry struct {
	Ts    *time.Time
	Msg   *string
	File  *string
	Line  *int
	Level *LogLevel
}

func NewFileLogger(opts ...Option) (logger *FileLogger) {
	logger = &FileLogger{}
	logger.Path = "logs"
	logger.pipe = make(chan *logEntry)

	for _, opt := range opts {
		opt(logger)
	}

	// todo: Layout

	logger.exit = make(chan struct{})
	go logger.write()

	return
}

func (l *FileLogger) Close() {
	close(l.exit)
}

func (l *FileLogger) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *FileLogger) GetLevel() LogLevel {
	return l.Level
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
		l.send(2, LevelTrace, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Debug(v ...any) {
	if l.CanDebug() {
		l.send(2, LevelDebug, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Info(v ...any) {
	if l.CanInfo() {
		l.send(2, LevelInfo, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Warn(v ...any) {
	if l.CanWarn() {
		l.send(2, LevelWarn, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Error(v ...any) {
	if l.CanError() {
		l.send(2, LevelError, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) Fatal(v ...any) {
	if l.CanFatal() {
		l.send(2, LevelFatal, fmt.Sprintln(v...))
	}
}

func (l *FileLogger) send(calldepth int, level LogLevel, s string) {
	now := time.Now()
	var file string
	var line int

	// todo: layout
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	entry := &logEntry{
		Level: &level,
		Msg:   &s,
		File:  &file,
		Line:  &line,
		Ts:    &now,
	}

	l.pipe <- entry
}

func (l *FileLogger) write() {
	var buf []byte
	for {
		select {
		case entry := <-l.pipe:
			l.ensureFile(entry.Ts)

			buf = buf[:0]
			formatTime(&buf, *entry.Ts)
			buf = append(buf, ' ')

			formatFile(&buf, *entry.File, *entry.Line)
			buf = append(buf, ' ')

			buf = append(buf, getLevelName(*entry.Level)...)
			buf = append(buf, ' ')

			buf = append(buf, *entry.Msg...)

			_, err := l.file.Write(buf)
			if err != nil {
				// todo: write to ylog.txt
				log.Println(err)
			}
		case <-l.exit:
			return
		}
	}
}

func formatFile(buf *[]byte, file string, line int) {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short

	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	itoa(buf, line, -1)
}

func (l *FileLogger) ensureFile(curTime *time.Time) (err error) {
	if l.file == nil {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.file == nil {
			l.file, err = createFile(&l.Path, curTime)
			l.lastHour = getTimeHour(curTime)
		}
		return
	}

	currentHour := getTimeHour(curTime)
	if l.lastHour != currentHour {
		l.mu.Lock()
		defer l.mu.Unlock()
		if l.lastHour != currentHour {
			_ = l.file.Close()
			l.file, err = createFile(&l.Path, curTime)
			l.lastHour = getTimeHour(curTime)
		}
	}

	return
}

func formatTime(buf *[]byte, t time.Time) {
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	hour, min, sec := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	*buf = append(*buf, ':')
	itoa(buf, sec, 2)
	*buf = append(*buf, '.')
	itoa(buf, t.Nanosecond()/1e6, 3)
}

// from log/log.go in standard library
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
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
		err = os.MkdirAll(dir, 0766) // for mac 766, for windows 666
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

func getLevelName(level LogLevel) []byte {
	if levelName, ok := levelMap[level]; ok {
		return levelName
	}

	return []byte{}
}
