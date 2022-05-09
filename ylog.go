package ylog

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type LogLevel byte

func (l LogLevel) MarshalJSON() ([]byte, error) {
	var lName string = levelName[l]
	return []byte(`"` + lName + `"`), nil
}

const (
	LevelTrace LogLevel = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

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

type YesLogger struct {
	Level     LogLevel
	writer    LoggerWriter
	formatter LoggerFormatter
	pipe      chan *logEntry
	cacheSize int
	sync      chan struct{}
	exit      chan struct{}
}

type logEntry struct {
	Ts    time.Time `json:"ts"`
	Msg   string    `json:"msg"`
	File  string    `json:"file"`
	Line  int       `json:"line"`
	Level LogLevel  `json:"level"`
}

func NewYesLogger(opts ...Option) (logger *YesLogger) {
	logger = &YesLogger{}
	logger.pipe = make(chan *logEntry)
	logger.writer = NewFileWriter("logs")
	logger.formatter = NewTextFormatter()

	for _, opt := range opts {
		opt(logger)
	}

	logger.sync = make(chan struct{})
	logger.exit = make(chan struct{})
	go logger.write()

	return
}

func (l *YesLogger) Close() {
	close(l.exit)
}

func (l *YesLogger) Sync() {
	close(l.sync)
}

func (l *YesLogger) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *YesLogger) GetLevel() LogLevel {
	return l.Level
}

func (l *YesLogger) CanTrace() bool {
	return l.Level <= LevelTrace
}

func (l *YesLogger) CanDebug() bool {
	return l.Level <= LevelDebug
}

func (l *YesLogger) CanInfo() bool {
	return l.Level <= LevelInfo
}

func (l *YesLogger) CanWarn() bool {
	return l.Level <= LevelWarn
}

func (l *YesLogger) CanError() bool {
	return l.Level <= LevelError
}

func (l *YesLogger) CanFatal() bool {
	return l.Level <= LevelFatal
}

func (l *YesLogger) Trace(v ...any) {
	if l.CanTrace() {
		l.send(2, LevelTrace, fmt.Sprint(v...))
	}
}

func (l *YesLogger) Debug(v ...any) {
	if l.CanDebug() {
		l.send(2, LevelDebug, fmt.Sprint(v...))
	}
}

func (l *YesLogger) Info(v ...any) {
	if l.CanInfo() {
		l.send(2, LevelInfo, fmt.Sprint(v...))
	}
}

func (l *YesLogger) Warn(v ...any) {
	if l.CanWarn() {
		l.send(2, LevelWarn, fmt.Sprint(v...))
	}
}

func (l *YesLogger) Error(v ...any) {
	if l.CanError() {
		l.send(2, LevelError, fmt.Sprint(v...))
	}
}

func (l *YesLogger) Fatal(v ...any) {
	if l.CanFatal() {
		l.send(2, LevelFatal, fmt.Sprint(v...))
	}
}

func (l *YesLogger) send(calldepth int, level LogLevel, s string) {
	now := time.Now()
	var file string
	var line int

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

func (l *YesLogger) write() {
	var buf []byte
	for {
		select {
		case entry := <-l.pipe:
			l.writer.Ensure(entry.Ts)

			// reuse the slice memory
			buf = buf[:0]
			l.formatter.Format(entry, &buf)
			err := l.writer.Write(buf)
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
			l.writer.Close()
			return
		}
	}
}
