package ylog

import "runtime"

var std = NewYesLogger(
	Level(LevelInfo),
	CacheSize(runtime.NumCPU()),
)

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
