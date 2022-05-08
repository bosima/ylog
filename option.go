package ylog

type Option func(logger *YesLogger)

func Level(level LogLevel) Option {
	return func(logger *YesLogger) {
		logger.Level = level
	}
}

func CacheSize(size uint16) Option {
	return func(logger *YesLogger) {
		logger.cacheSize = size
		if size > 0 {
			logger.pipe = make(chan *logEntry, logger.cacheSize)
		}
	}
}

func Writer(writer LoggerWriter) Option {
	return func(logger *YesLogger) {
		logger.writer = writer
	}
}

func Formatter(formatter LoggerFormatter) Option {
	return func(logger *YesLogger) {
		logger.formatter = formatter
	}
}
