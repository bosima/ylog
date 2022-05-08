package ylog

type LoggerFormatter interface {
	Format(*logEntry, *[]byte) error
}
