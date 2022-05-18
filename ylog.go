package ylog

import (
	"fmt"
	"os"
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

type yesLogger struct {
	level     LogLevel
	writer    LoggerWriter
	formatter LoggerFormatter
	pipe      chan *logEntry
	cacheSize int
	errFile   *os.File
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

func NewYesLogger(opts ...Option) (logger *yesLogger) {
	logger = &yesLogger{}
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

func (l *yesLogger) Close() {
	close(l.exit)
}

func (l *yesLogger) Sync() {
	close(l.sync)
}

func (l *yesLogger) SetLevel(level LogLevel) {
	l.level = level
}

func (l *yesLogger) GetLevel() LogLevel {
	return l.level
}

func (l *yesLogger) CanTrace() bool {
	return l.level <= LevelTrace
}

func (l *yesLogger) CanDebug() bool {
	return l.level <= LevelDebug
}

func (l *yesLogger) CanInfo() bool {
	return l.level <= LevelInfo
}

func (l *yesLogger) CanWarn() bool {
	return l.level <= LevelWarn
}

func (l *yesLogger) CanError() bool {
	return l.level <= LevelError
}

func (l *yesLogger) CanFatal() bool {
	return l.level <= LevelFatal
}

func (l *yesLogger) Trace(v ...any) {
	if l.CanTrace() {
		l.send(2, LevelTrace, fmt.Sprint(v...))
	}
}

func (l *yesLogger) Debug(v ...any) {
	if l.CanDebug() {
		l.send(2, LevelDebug, fmt.Sprint(v...))
	}
}

func (l *yesLogger) Info(v ...any) {
	if l.CanInfo() {
		l.send(2, LevelInfo, fmt.Sprint(v...))
	}
}

func (l *yesLogger) Warn(v ...any) {
	if l.CanWarn() {
		l.send(2, LevelWarn, fmt.Sprint(v...))
	}
}

func (l *yesLogger) Error(v ...any) {
	if l.CanError() {
		l.send(2, LevelError, fmt.Sprint(v...))
	}
}

func (l *yesLogger) Fatal(v ...any) {
	if l.CanFatal() {
		l.send(2, LevelFatal, fmt.Sprint(v...))
	}
}

func (l *yesLogger) send(calldepth int, lev LogLevel, s string) {
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

func (l *yesLogger) write() {
	var buf []byte
	for {
		select {
		case entry := <-l.pipe:
			// reuse the slice memory
			buf = buf[:0]
			l.formatter.Format(entry, &buf)
			l.writer.Ensure(entry)
			err := l.writer.Write(buf)
			if err != nil {
				l.writeError(err)
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

func (l *yesLogger) writeError(err error) {
	if l.errFile == nil {
		l.errFile, _ = os.OpenFile("ylog.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	}

	var buf []byte
	formatTime(&buf, time.Now())
	buf = append(buf, ' ')
	buf = append(buf, "Error"...)
	buf = append(buf, ' ')
	buf = append(buf, err.Error()...)
	buf = append(buf, '\n')

	l.errFile.Write(buf)
}
