# Introduction
Ylog is a log library of Go language.

It can write data to disk file or Kafka, with plain text or Json.

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
    BenchmarkInfo-8   	 1332333	       871.6 ns/op	     328 B/op	       4 allocs/op
