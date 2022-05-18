package main

import (
	"github.com/bosima/ylog"
)

func main() {
	// Log with the default instance of the 'FileLogger'
	ylog.SetLevel(ylog.LevelInfo)
	ylog.Info("I am a info log.")

	// Log with a 'FileLogger' instance that created by NewXXX
	log := ylog.NewFileLogger(ylog.LevelInfo, "logs2")
	log.Info("I am a info log too.")
}
