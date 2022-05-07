package ylog

import (
	"io/ioutil"
	"runtime"
	"testing"
	"time"
)

func BenchmarkInfo(b *testing.B) {
	var stdPath = "logs"
	var logger = NewFileLogger(
		Level(LevelInfo),
		Path(stdPath),
		CacheSize(uint16(runtime.NumCPU())),
		Writer(&discardWriter{}),
	)
	logger.Info("Ready")
	<-time.After(time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("This is a Benchmark info.")
	}
}

func BenchmarkInfo_Parallel(b *testing.B) {
	var stdPath = "logs"
	var logger = NewFileLogger(
		Level(LevelInfo),
		Path(stdPath),
		CacheSize(uint16(runtime.NumCPU())),
		Writer(&discardWriter{}),
	)
	logger.Info("Ready")
	<-time.After(time.Second)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("This is a Benchmark info.")
		}
	})
}

// all Write calls succeed without doing anything.
type discardWriter struct {
}

func (w *discardWriter) Ensure(curTime time.Time) (err error) {
	return
}

func (w *discardWriter) Write(buf []byte) (int, error) {
	return ioutil.Discard.Write(buf)
}

func (w *discardWriter) Sync() error {
	return nil
}
