package main

import (
	"time"

	"github.com/bosima/ylog"
)

func main() {
	// Log with the default instance of the 'FileLogger'
	ylog.SetLevel(ylog.LevelInfo)

	go func() {
		for i := 0; i < 1000; i++ {
			ylog.Info("I am a info log 1.")
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			ylog.Info("I am a info log 2.")
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			ylog.Info("I am a info log 3.")
		}
	}()

	// Log with a 'FileLogger' instance that created by NewXXX
	log := ylog.NewFileLogger(ylog.LevelInfo, "logs2")
	log.Info("I am a info log too.")

	time.Sleep(time.Second)
}
