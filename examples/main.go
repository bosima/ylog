package main

import (
	"log"
	"time"

	"github.com/bosima/ylog"
)

func main() {
	logger := ylog.NewFileLogger(
		ylog.CacheSize(24),
		ylog.Level(ylog.LevelInfo),
	)

	logger.Trace("This is a trace log.")
	logger.Debug("This is a debug log.")
	logger.Info("This is a info log.")
	logger.Warn("This is a warn log.")
	logger.Error("This is a error log.")
	logger.Fatal("This is a fatal log.")

	go func() {
		for i := 0; i < 10000; i++ {
			logger.Info("loop info log.")
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			logger.Warn("loop warn log.")
		}
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			logger.Error("loop error log.")
		}
	}()

	<-time.After(time.Second * 5)
	log.Println("Done")
}
