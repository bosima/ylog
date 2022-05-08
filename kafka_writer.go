package ylog

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaWriter struct {
	Topic         string
	Address       string
	conn          *kafka.Conn
	batchSize     int
	buf           []kafka.Message
	lastWriteTime time.Time
}

func NewKafkaWriter(address string, topic string, batchSize int) *KafkaWriter {
	return &KafkaWriter{
		Address:   address,
		Topic:     topic,
		batchSize: batchSize,
	}
}

func (w *KafkaWriter) Ensure(curTime time.Time) (err error) {
	if w.conn == nil {
		w.buf = make([]kafka.Message, 0, w.batchSize)
		partition := 0
		w.conn, err = kafka.DialLeader(context.Background(), "tcp", w.Address, w.Topic, partition)
		w.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	}

	return
}

func (w *KafkaWriter) Write(buf []byte) (bts int, err error) {
	w.buf = append(w.buf, kafka.Message{Value: buf})

	now := time.Now()
	if now.UnixMilli()-w.lastWriteTime.UnixMilli() > 1000 {
		bts, err = w.conn.WriteMessages(
			w.buf...,
		)
		w.buf = w.buf[:0]
		w.lastWriteTime = now
		return
	}

	if len(w.buf) <= w.batchSize-1 {
		return -1, nil
	}

	bts, err = w.conn.WriteMessages(
		w.buf...,
	)
	w.buf = w.buf[:0]
	w.lastWriteTime = now
	return
}

func (w *KafkaWriter) Sync() error {
	return nil
}

func (w *KafkaWriter) Close() error {
	return w.conn.Close()
}
