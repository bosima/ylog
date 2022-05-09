package ylog

import (
	"time"

	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type KafkaWriter struct {
	Topic     string
	Address   string
	writer    *kafka.Writer
	batchSize int
	buf       []kafka.Message
}

func NewKafkaWriter(address string, topic string, batchSize int) *KafkaWriter {
	return &KafkaWriter{
		Address:   address,
		Topic:     topic,
		batchSize: batchSize,
	}
}

func (w *KafkaWriter) Ensure(curTime time.Time) (err error) {
	if w.writer == nil {
		w.buf = make([]kafka.Message, 0, w.batchSize)

		w.writer = &kafka.Writer{
			Addr:      kafka.TCP(w.Address),
			Topic:     w.Topic,
			BatchSize: w.batchSize,
		}
	}

	return
}

func (w *KafkaWriter) Write(buf []byte) (err error) {

	if len(w.buf) < w.batchSize-1 {
		// buf will be reused by ylog, so need copy to a new slice
		kbuf := append([]byte(nil), buf...)
		w.buf = append(w.buf, kafka.Message{Value: kbuf})
		return
	}
	w.buf = append(w.buf, kafka.Message{Value: buf})
	err = w.writer.WriteMessages(context.Background(),
		w.buf...,
	)
	w.buf = w.buf[:0]
	return
}

func (w *KafkaWriter) Sync() error {
	return nil
}

func (w *KafkaWriter) Close() error {
	return w.writer.Close()
}
