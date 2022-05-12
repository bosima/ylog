package ylog

import "encoding/json"

type jsonFormatter struct {
}

func NewJsonFormatter() LoggerFormatter {
	return &jsonFormatter{}
}

func (f *jsonFormatter) Format(entry *logEntry, buf *[]byte) (err error) {
	entry.File = toShortFile(entry.File)
	jsonBuf, err := json.Marshal(entry)
	*buf = append(*buf, jsonBuf...)
	return
}
