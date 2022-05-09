package ylog

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type fileWriter struct {
	file     *os.File
	mu       sync.Mutex
	lastHour int64
	Path     string
}

func NewFileWriter(path string) *fileWriter {
	return &fileWriter{Path: path}
}

func (w *fileWriter) Ensure(curTime time.Time) (err error) {
	if w.file == nil {
		w.mu.Lock()
		defer w.mu.Unlock()
		if w.file == nil {
			w.file, err = w.createFile(w.Path, curTime)
			w.lastHour = w.getTimeHour(curTime)
		}
		return
	}

	currentHour := w.getTimeHour(curTime)
	if w.lastHour != currentHour {
		w.mu.Lock()
		defer w.mu.Unlock()
		if w.lastHour != currentHour {
			_ = w.file.Close()
			w.file, err = w.createFile(w.Path, curTime)
			w.lastHour = currentHour
		}
	}

	return
}

func (w *fileWriter) Write(buf []byte) (err error) {
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
	if err != nil {
		return
	}

	return
}

func (w *fileWriter) getFileName(t time.Time) string {
	return t.Format("2006-01-02_15")
}

func (w *fileWriter) getTimeHour(t time.Time) int64 {
	return t.Unix() / 3600
}
