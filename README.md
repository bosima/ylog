# Introduction
Ylog is a log library of Go language.

She writes data to disk file, which are divided by hours.

# Usage

## Settings

    ylog.SetLevel(ylog.LevelTrace)

## Write Log

    ylog.Info("Hello Golang!")


# Benchamark Test 

    goos: windows
    goarch: amd64
    pkg: github.com/bosima/ylog
    cpu: Intel(R) Core(TM) i5-10400F CPU @ 2.90GHz
    BenchmarkInfo-12             	  831427	      1318 ns/op	     328 B/op	       4 allocs/op
    BenchmarkInfo_Parallel-12    	 1669800	       734.5 ns/op	     328 B/op	       4 allocs/op
