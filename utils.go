package ylog

import (
	"time"
)

func toShort(file string) string {
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			return file[i+1:]
		}
	}
	return file
}

// from log/log.go in standard library
func formatTime(buf *[]byte, t time.Time) {
	year, month, day := t.Date()
	itoa4(buf, year)
	*buf = append(*buf, '/')
	itoa2(buf, int(month))
	*buf = append(*buf, '/')
	itoa2(buf, day)
	*buf = append(*buf, ' ')
	hour, min, sec := t.Clock()
	itoa2(buf, hour)
	*buf = append(*buf, ':')
	itoa2(buf, min)
	*buf = append(*buf, ':')
	itoa2(buf, sec)
	*buf = append(*buf, '.')
	itoa3(buf, t.Nanosecond()/1e6)
}

// from log/log.go in standard library
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte // the max number of characters for Int64
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func itoa2(buf *[]byte, i int) {
	q := i / 10
	s := byte('0' + i - q*10)
	f := byte('0' + q)
	*buf = append(*buf, f, s)
}

func itoa3(buf *[]byte, i int) {
	tq := i / 10
	sq := i / 100
	t := byte('0' + i - tq*10)
	s := byte('0' + tq - sq*10)
	f := byte('0' + sq)
	*buf = append(*buf, f, s, t)
}

func itoa4(buf *[]byte, i int) {
	foq := i / 10
	tq := i / 100
	sq := i / 1000
	fo := byte('0' + i - foq*10)
	t := byte('0' + foq - tq*10)
	s := byte('0' + tq - sq*10)
	fi := byte('0' + sq)
	*buf = append(*buf, fi, s, t, fo)
}
