# Introduction
Ylog is a log library of Go language.

It can write data to disk file or Kafka, with plain text or Json.

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
    BenchmarkInfo-12    	 1899787	       613.3 ns/op	     328 B/op	       4 allocs/op

# TODO

Write a log message to multiple writers.
