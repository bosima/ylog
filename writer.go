package ylog

import (
	"time"
)

type LoggerWriter interface {
	Ensure(time.Time) error
	Write([]byte) (int, error)
	Sync() error
	Close() error
}
