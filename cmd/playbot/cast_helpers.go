package main

import (
	"unsafe"
)

func s2b(str string) (bs []byte) {
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func b2s(bs []byte) (str string) {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}
