package ylog

type Option func(logger *FileLogger)

func Level(level LogLevel) Option {
	return func(logger *FileLogger) {
		logger.Level = level
	}
}

func CacheSize(size uint16) Option {
	return func(logger *FileLogger) {
		logger.cacheSize = size
		if size > 0 {
			logger.pipe = make(chan *logEntry, logger.cacheSize)
		}
	}
}

func Path(path string) Option {
	return func(logger *FileLogger) {
		logger.Path = path
	}
}

func Layout(layout string) Option {
	return func(logger *FileLogger) {
		logger.Layout = layout
	}
}

func Writer(writer LoggerWriter) Option {
	return func(logger *FileLogger) {
		logger.writer = writer
	}
}
