package ylog

import (
	"time"

	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type kafkaSyncWriter struct {
	Topic     string
	Address   string
	writer    *kafka.Writer
	batchSize int
	buf       []kafka.Message
	lastTime  time.Time
}

func NewKafkaSyncWriter(address string, topic string, batchSize int) LoggerWriter {
	return &kafkaSyncWriter{
		Address:   address,
		Topic:     topic,
		batchSize: batchSize,
	}
}

func (w *kafkaSyncWriter) Ensure(curTime time.Time) (err error) {
	if w.writer == nil {
		w.buf = make([]kafka.Message, 0, w.batchSize)

		w.writer = &kafka.Writer{
			Addr:      kafka.TCP(w.Address),
			Topic:     w.Topic,
			BatchSize: 1,
			Async:     true,
		}
	}

	return
}

func (w *kafkaSyncWriter) Write(buf []byte) (err error) {

	// buf will be reused by ylog, so need copy to a new slice
	kbuf := append([]byte(nil), buf...)
	w.buf = append(w.buf, kafka.Message{Value: kbuf})

	now := time.Now()
	if now.UnixMilli()-w.lastTime.UnixMilli() > 1000 {
		err = w.writer.WriteMessages(context.Background(),
			w.buf...,
		)
		w.buf = w.buf[:0]
		w.lastTime = now
		return
	}

	if len(w.buf) >= w.batchSize {
		err = w.writer.WriteMessages(context.Background(),
			w.buf...,
		)
		w.buf = w.buf[:0]
		w.lastTime = now
	}
	return
}

func (w *kafkaSyncWriter) Sync() error {
	return nil
}

func (w *kafkaSyncWriter) Close() error {
	return w.writer.Close()
}
