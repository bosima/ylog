package ylog

type Option func(logger *yesLogger)

func Level(level LogLevel) Option {
	return func(logger *yesLogger) {
		logger.level = level
	}
}

func CacheSize(size int) Option {
	return func(logger *yesLogger) {
		logger.cacheSize = size
		if size > 0 {
			logger.pipe = make(chan *logEntry, logger.cacheSize)
		}
	}
}

func Writer(writer LoggerWriter) Option {
	return func(logger *yesLogger) {
		logger.writer = writer
	}
}

func Formatter(formatter LoggerFormatter) Option {
	return func(logger *yesLogger) {
		logger.formatter = formatter
	}
}
