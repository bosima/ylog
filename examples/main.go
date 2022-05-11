package main

import (
	"github.com/bosima/ylog"
)

func main() {
	// Log with the default instance of the 'FileLogger'
	ylog.SetLevel(ylog.LevelInfo)
	ylog.Info("I am a info log.")

	// Log with a 'FileLogger' instance created by NewXXX
	var logPath = "logs2"
	var logger = ylog.NewFileLogger(ylog.LevelInfo, logPath)
	logger.Info("I am a info log too.")
}
