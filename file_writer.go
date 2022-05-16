package ylog

import (
	"os"
	"path/filepath"
	"time"
)

type fileWriter struct {
	file     *os.File
	lastHour int64
	Path     string
}

func NewFileWriter(path string) LoggerWriter {
	return &fileWriter{Path: path}
}

func (w *fileWriter) Ensure(entry *logEntry) (err error) {
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

func (w *fileWriter) Write(buf []byte) (err error) {
	buf = append(buf, '\n')
	_, err = w.file.Write(buf)
	return
}

func (w *fileWriter) Sync() error {
	return w.file.Sync()
}

func (w *fileWriter) Close() error {
	return w.file.Close()
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
