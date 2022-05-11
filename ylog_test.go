package ylog

import (
	"testing"
)

func BenchmarkInfo_Parallel(b *testing.B) {
	var logger = NewFileLogger(LevelInfo, "logs")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("This is a Benchmark info.")
		}
	})
}
