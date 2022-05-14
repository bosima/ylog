package ylog

type textFormatter struct {
}

func NewTextFormatter() *textFormatter {
	return &textFormatter{}
}

func (f *textFormatter) Format(entry *logEntry, buf *[]byte) error {
	formatTime(buf, entry.Ts)
	*buf = append(*buf, ' ')

	file := toShort(entry.File)
	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	itoa(buf, entry.Line, -1)
	*buf = append(*buf, ' ')

	*buf = append(*buf, levelNames[entry.Level]...)
	*buf = append(*buf, ' ')

	*buf = append(*buf, entry.Msg...)

	return nil
}
