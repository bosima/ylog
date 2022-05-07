package ylog

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LoggerWriter interface {
	Ensure(time.Time) error
	Write([]byte) (int, error)
	Sync() error
}

type fileWriter struct {
	file     *os.File
	mu       sync.Mutex
	lastHour int64
	Path     string
}

func (w *fileWriter) Ensure(curTime time.Time) (err error) {
	if w.file == nil {
		w.mu.Lock()
		defer w.mu.Unlock()
		if w.file == nil {
			w.file, err = createFile(w.Path, curTime)
			w.lastHour = getTimeHour(curTime)
		}
		return
	}

	currentHour := getTimeHour(curTime)
	if w.lastHour != currentHour {
		w.mu.Lock()
		defer w.mu.Unlock()
		if w.lastHour != currentHour {
			_ = w.file.Close()
			w.file, err = createFile(w.Path, curTime)
			w.lastHour = getTimeHour(curTime)
		}
	}

	return
}

func (w *fileWriter) Write(buf []byte) (int, error) {
	return w.file.Write(buf)
}

func (w *fileWriter) Sync() error {
	return w.file.Sync()
}

func createFile(path string, t time.Time) (file *os.File, err error) {
	dir := filepath.Join(path, t.Format("200601"))
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0766) // for mac 766, for windows 666
		if err != nil {
			return nil, err
		}
	}

	filePath := filepath.Join(dir, getFileName(t)+".txt")
	file, err = os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	return
}
