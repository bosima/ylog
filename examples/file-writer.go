package main

import (
	"log"
	"sync"

	"github.com/bosima/ylog"
)

func main() {
	logger := ylog.NewYesLogger(
		ylog.CacheSize(3),
		ylog.Level(ylog.LevelInfo),
	)

	logger.Trace("This is a trace log.")
	logger.Debug("This is a debug log.")
	logger.Info("This is a info log.")
	logger.Warn("This is a warn log.")
	logger.Error("This is a error log.")
	logger.Fatal("This is a fatal log.")

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		for i := 0; i < 10000; i++ {
			logger.Info("loop info log.")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			logger.Warn("loop warn log.")
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < 10000; i++ {
			logger.Error("loop error log.")
		}
		wg.Done()
	}()

	wg.Wait()
	log.Println("Publish Done")

	select {}
}
