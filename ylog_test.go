package ylog

import (
	"testing"
	"time"
)

func BenchmarkInfo(b *testing.B) {
	var stdPath = "logs"
	var logger = NewFileLogger(
		Level(LevelInfo),
		Path(stdPath),
		CacheSize(12),
	)
	logger.Info("Ready")
	<-time.After(time.Second * 2)

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
		CacheSize(12),
	)
	logger.Info("Ready")
	<-time.After(time.Second * 2)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("This is a Benchmark info.")
		}
	})
}
