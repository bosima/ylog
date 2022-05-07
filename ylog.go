package ylog

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type LogLevel = byte

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// Reference: https://github.com/golang/glog/blob/9ef845f417d839250ceabbc25c1b26101e772dd7/glog.go#L110
var levelName = []string{
	LevelTrace: "Trace",
	LevelDebug: "Debug",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelFatal: "Fatal",
}

type Logger interface {
	Trace(v ...any)
	Debug(v ...any)
	Info(v ...any)
	Warn(v ...any)
	Error(v ...any)
	Fatal(v ...any)
}

type FileLogger struct {
	Level     LogLevel
	pipe      chan *logEntry
	cacheSize uint16
	Layout    string
	Path      string
	writer    LoggerWriter
	sync      chan struct{}
	exit      chan struct{}
}

type logEntry struct {
	Ts    time.Time
	Msg   string
	File  string
	Line  int
	Level LogLevel
}

func NewFileLogger(opts ...Option) (logger *FileLogger) {
	logger = &FileLogger{}
	logger.Path = "logs"
	logger.pipe = make(chan *logEntry)
	logger.writer = &fileWriter{Path: logger.Path}

	for _, opt := range opts {
		opt(logger)
	}

	logger.sync = make(chan struct{})
	logger.exit = make(chan struct{})
	go logger.write()

	return
}

func (l *FileLogger) Close() {
	close(l.exit)
}

func (l *FileLogger) Sync() {
	close(l.sync)
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
		Level: level,
		Msg:   s,
		File:  file,
		Line:  line,
		Ts:    now,
	}

	l.pipe <- entry
}

func (l *FileLogger) write() {
	var buf []byte
	for {
		select {
		case entry := <-l.pipe:
			l.writer.Ensure(entry.Ts)

			// resue the slice memory
			buf = buf[:0]
			formatTime(&buf, entry.Ts)
			buf = append(buf, ' ')

			formatFile(&buf, entry.File, entry.Line)
			buf = append(buf, ' ')

			buf = append(buf, levelName[entry.Level]...)
			buf = append(buf, ' ')

			buf = append(buf, entry.Msg...)

			_, err := l.writer.Write(buf)
			if err != nil {
				// todo: write to ylog.txt
				log.Println(err)
			}
		case _, ok := <-l.sync:
			if ok {
				l.writer.Sync()
				l.sync = make(chan struct{})
			}
		case <-l.exit:
			return
		}
	}
}

// from log/log.go in standard library
func formatFile(buf *[]byte, file string, line int) {
	// todo reuse filename
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

// from log/log.go in standard library
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

func getTimeHour(t time.Time) int64 {
	return t.Unix() / 3600
}

func getFileName(t time.Time) string {
	return t.Format("2006-01-02_15")
}
