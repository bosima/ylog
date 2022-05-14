package ylog

type LoggerWriter interface {
	Ensure(*logEntry) error
	Write([]byte) error
	Sync() error
	Close() error
}
