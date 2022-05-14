package ylog

import (
	"io/ioutil"
	"testing"
	"time"
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

func (w *discardWriter) Ensure(curTime time.Time) (err error) {
	return
}

func (w *discardWriter) Write(buf []byte) (err error) {
	_, err = ioutil.Discard.Write(buf)
	return err
}

func (w *discardWriter) Sync() error {
	return nil
}

func (w *discardWriter) Close() error {
	return nil
}
