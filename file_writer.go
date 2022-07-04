package ylog

import (
	"os"
	"path/filepath"
	"time"
)

type fileWriter struct {
	lastHour  int64
	file      *os.File
	formatter LoggerFormatter
	Path      string
	buf       []byte
}

func NewFileWriter(path string, formatter LoggerFormatter) LoggerWriter {
	return &fileWriter{Path: path, formatter: formatter}
}

func (w *fileWriter) Write(entry *logEntry) (err error) {
	w.ensure(entry)

	// reuse the slice memory
	w.buf = w.buf[:0]
	w.formatter.Format(entry, &w.buf)
	w.buf = append(w.buf, '\n')

	_, err = w.file.Write(w.buf)
	return
}

func (w *fileWriter) Sync() error {
	return w.file.Sync()
}

func (w *fileWriter) Close() error {
	return w.file.Close()
}

func (w *fileWriter) ensure(entry *logEntry) (err error) {
	if w.file == nil {
		f, err := w.createFile(w.Path, entry.Ts)
		if err != nil {
			return err
		}
		w.lastHour = w.getTimeHour(entry.Ts)
		w.file = f
		return nil
	}

	currentHour := w.getTimeHour(entry.Ts)
	if w.lastHour != currentHour {
		_ = w.file.Close()
		f, err := w.createFile(w.Path, entry.Ts)
		if err != nil {
			return err
		}
		w.lastHour = currentHour
		w.file = f
	}

	return
}

func (w *fileWriter) createFile(path string, t time.Time) (file *os.File, err error) {
	dir := filepath.Join(path, t.Format("200601"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766) // for mac 766, for windows 666
		if err != nil {
			return nil, err
		}
	}

	filePath := filepath.Join(dir, w.getFileName(t)+".txt")
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	return
}

func (w *fileWriter) getFileName(t time.Time) string {
	return t.Format("2006-01-02_15")
}

func (w *fileWriter) getTimeHour(t time.Time) int64 {
	return t.Unix() / 3600
}
