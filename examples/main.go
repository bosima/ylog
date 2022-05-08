package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/bosima/ylog"
	"github.com/segmentio/kafka-go"
)

func main() {
	logger := ylog.NewYesLogger(
		ylog.CacheSize(3),
		ylog.Level(ylog.LevelInfo),
		ylog.Writer(ylog.NewKafkaWriter("localhost:9092", "ylog", 8)),
	)

	logger.Trace("This is a trace log.")
	logger.Debug("This is a debug log.")
	logger.Info("This is a info log.")
	logger.Warn("This is a warn log.")
	logger.Error("This is a error log.")
	logger.Fatal("This is a fatal log.")

	go func() {
		conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "ylog", 0)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}

		conn.SetReadDeadline(time.Now().Add(10 * time.Second))

		b := make([]byte, 10e3) // 10KB max per message
		for {
			n, err := conn.Read(b)
			if err != nil {
				break
			}
			fmt.Print(string(b[:n]))
		}
	}()

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
