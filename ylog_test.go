package ylog

import (
	"testing"
)

func BenchmarkInfo(b *testing.B) {
	var log = NewYesLogger(
		Level(LevelInfo),
		Writer(&discardWriter{}),
	)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Info("This is a Benchmark info.")
		}
	})
}

// all Write calls succeed without doing anything.
type discardWriter struct {
}

func (w *discardWriter) Write(entry *logEntry) (err error) {
	return
}

func (w *discardWriter) Sync() error {
	return nil
}

func (w *discardWriter) Close() error {
	return nil
}
