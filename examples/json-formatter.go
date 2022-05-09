package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/bosima/ylog"
)

func main() {
	logger := ylog.NewYesLogger(
		ylog.CacheSize(runtime.NumCPU()),
		ylog.Level(ylog.LevelInfo),
		ylog.Formatter(ylog.NewJsonFormatter()),
	)

	logger.Trace("This is a trace log.")
	logger.Debug("This is a debug log.")
	logger.Info("This is a info log.")
	logger.Warn("This is a warn log.")
	logger.Error("This is a error log.")
	logger.Fatal("This is a fatal log.")

	log.Println("Logging...")

	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 100)
			logger.Info("loop info log.")
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 300)
			logger.Warn("loop warn log.")
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 600)
			logger.Error("loop error log.")
		}
	}()
	wg.Wait()
	log.Println("Publish Done")

	select {}
}
