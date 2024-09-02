package pool

import (
	"bytes"
	"sync"
)

var buffers sync.Pool

func GetBuffer() *bytes.Buffer {
	buf, ok := buffers.Get().(*bytes.Buffer)
	if !ok {
		buf = new(bytes.Buffer)
	}
	return buf
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	buffers.Put(buf)
}
