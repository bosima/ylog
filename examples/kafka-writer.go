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

var (
	address        string = "localhost:9092"
	topic          string = "ylog"
	writeBatchSize int    = 8
)

func main() {
	go receiveData()
	go publishLog()
	select {}
}

func publishLog() {
	logger := ylog.NewYesLogger(
		ylog.CacheSize(3),
		ylog.Level(ylog.LevelInfo),
		ylog.Writer(ylog.NewKafkaWriter(address, topic, writeBatchSize, ylog.NewTextFormatter())),
	)

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 30)
			logger.Info("loop info log.")
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 60)
			logger.Warn("loop warn log.")
		}
	}()

	go func() {
		defer wg.Done()
		for {
			time.Sleep(time.Millisecond * 90)
			logger.Error("loop error log.")
		}
	}()

	wg.Wait()
	logger.Close()
	log.Println("Publish Done")
}

func receiveData() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		Topic:    topic,
		GroupID:  "ylog",
		MinBytes: 500, // 500B
		MaxBytes: 1e6, // 1MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d:%s\n", m.Offset, string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
