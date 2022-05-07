# Introduction
Ylog is a log library of Go language.

She writes data to disk file, which are divided by hours.

# Usage

## Settings

    ylog.SetLevel(ylog.LevelTrace)

## Write Log

    ylog.Info("Hello Golang!")


# Benchamark Test 

    goos: darwin
    goarch: amd64
    pkg: github.com/bosima/ylog
    cpu: Intel(R) Core(TM) i5-8259U CPU @ 2.30GHz
    BenchmarkInfo-8            	  678504	      1762 ns/op	     328 B/op	       4 allocs/op
    BenchmarkInfo_Parallel-8   	 1322624	       879.5 ns/op	     328 B/op	       4 allocs/op
