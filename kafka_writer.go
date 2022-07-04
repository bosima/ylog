package ylog

import (
	"github.com/segmentio/kafka-go"
	"golang.org/x/net/context"
)

type kafkaWriter struct {
	Topic     string
	Address   string
	writer    *kafka.Writer
	batchSize int
	formatter LoggerFormatter
}

func NewKafkaWriter(address string, topic string, batchSize int, formatter LoggerFormatter) LoggerWriter {
	return &kafkaWriter{
		Address:   address,
		Topic:     topic,
		batchSize: batchSize,
		formatter: formatter,
	}
}

func (w *kafkaWriter) Write(entry *logEntry) (err error) {
	w.ensure(entry)

	var buf []byte
	w.formatter.Format(entry, &buf)

	err = w.writer.WriteMessages(context.Background(),
		kafka.Message{Value: buf},
	)
	return
}

func (w *kafkaWriter) ensure(_ *logEntry) (err error) {
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

func (w *kafkaWriter) Sync() error {
	return nil
}

func (w *kafkaWriter) Close() error {
	return w.writer.Close()
}
