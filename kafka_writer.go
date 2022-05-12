package ylog

import (
	"time"

	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type kafkaWriter struct {
	Topic     string
	Address   string
	writer    *kafka.Writer
	batchSize int
}

func NewKafkaWriter(address string, topic string, batchSize int) LoggerWriter {
	return &kafkaWriter{
		Address:   address,
		Topic:     topic,
		batchSize: batchSize,
	}
}

func (w *kafkaWriter) Ensure(curTime time.Time) (err error) {
	if w.writer == nil {
		w.writer = &kafka.Writer{
			Addr:      kafka.TCP(w.Address),
			Topic:     w.Topic,
			BatchSize: w.batchSize,
			Async:     true,
		}
	}

	return
}

func (w *kafkaWriter) Write(buf []byte) (err error) {
	// buf will be reused by ylog when this method return,
	// with aysnc write, we need copy data to a new slice
	kbuf := append([]byte(nil), buf...)
	err = w.writer.WriteMessages(context.Background(),
		kafka.Message{Value: kbuf},
	)
	return
}

func (w *kafkaWriter) Sync() error {
	return nil
}

func (w *kafkaWriter) Close() error {
	return w.writer.Close()
}
