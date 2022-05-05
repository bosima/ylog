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
    BenchmarkInfo-12             	  365077	      3073 ns/op	     368 B/op	       9 allocs/op
    BenchmarkInfo_Parallel-12    	  368134	      3170 ns/op	     368 B/op	       9 allocs/op
