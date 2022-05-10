package main

import (
	"github.com/bosima/ylog"
)

func main() {
	ylog.SetLevel(ylog.ToInfo)
	ylog.Info("I am a info log.")

	var stdPath = "logs2"
	var logger = ylog.NewFileLogger(ylog.ToInfo, stdPath)
	logger.Info("I am a info log too.")
}
