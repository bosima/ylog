package ylog

import (
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
	)
	logger.Info("Ready")
	<-time.After(time.Second)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("This is a Benchmark info.")
	}
	b.StopTimer()
	<-time.After(time.Second)
	logger.Sync()
}

func BenchmarkInfo_Parallel(b *testing.B) {
	var stdPath = "logs"
	var logger = NewFileLogger(
		Level(LevelInfo),
		Path(stdPath),
		CacheSize(uint16(runtime.NumCPU())),
	)
	logger.Info("Ready")
	<-time.After(time.Second)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("This is a Benchmark info.")
		}
	})
	b.StopTimer()
	<-time.After(time.Second)
	logger.Sync()
}
