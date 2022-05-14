package ylog

import (
	"fmt"
	"log"
	"runtime"
	"time"
)

type LogLevel byte

func (l LogLevel) MarshalJSON() ([]byte, error) {
	var lName string = levelNames[l]
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

var levelNames = []string{
	LevelTrace: "Trace",
	LevelDebug: "Debug",
	LevelInfo:  "Info",
	LevelWarn:  "Warn",
	LevelError: "Error",
	LevelFatal: "Fatal",
}

type YesLogger struct {
	level     LogLevel
	writer    LoggerWriter
	formatter LoggerFormatter
	pipe      chan *logEntry
	cacheSize int
	sync      chan struct{}
	exit      chan struct{}
}

type logEntry struct {
	Ts    time.Time `json:"ts"`
	File  string    `json:"file"`
	Line  int       `json:"line"`
	Level LogLevel  `json:"level"`
	Msg   string    `json:"msg"`
}

func NewYesLogger(opts ...Option) (logger *YesLogger) {
	logger = &YesLogger{}
	logger.pipe = make(chan *logEntry, runtime.NumCPU())
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
	l.level = level
}

func (l *YesLogger) GetLevel() LogLevel {
	return l.level
}

func (l *YesLogger) CanTrace() bool {
	return l.level <= LevelTrace
}

func (l *YesLogger) CanDebug() bool {
	return l.level <= LevelDebug
}

func (l *YesLogger) CanInfo() bool {
	return l.level <= LevelInfo
}

func (l *YesLogger) CanWarn() bool {
	return l.level <= LevelWarn
}

func (l *YesLogger) CanError() bool {
	return l.level <= LevelError
}

func (l *YesLogger) CanFatal() bool {
	return l.level <= LevelFatal
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

func (l *YesLogger) send(calldepth int, lev LogLevel, s string) {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}

	entry := &logEntry{
		Ts:    time.Now(),
		File:  file,
		Line:  line,
		Level: lev,
		Msg:   s,
	}

	l.pipe <- entry
}

func (l *YesLogger) write() {
	var buf []byte
	for {
		select {
		case entry := <-l.pipe:
			// reuse the slice memory
			buf = buf[:0]
			l.formatter.Format(entry, &buf)
			l.writer.Ensure(entry.Ts)
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
