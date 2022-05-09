package ylog

import "encoding/json"

type jsonFormatter struct {
}

func NewJsonFormatter() *jsonFormatter {
	return &jsonFormatter{}
}

func (f *jsonFormatter) Format(entry *logEntry, buf *[]byte) (err error) {
	entry.File = toShortFile(entry.File)
	*buf, err = json.Marshal(entry)
	return
}
