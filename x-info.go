package main

//   #cgo CFLAGS : -Wall -Wpedantic
//   #cgo LDFLAGS: -lXss -lX11 -lXmu
//   #include "x-info.h"
import "C"

import (
	"strings"
	"time"
)

func init() {
	C.XInfoInitialize()
}

func sanitize(in string) string {
	return strings.TrimSpace(in)
}
func getInfo() Captured {
	C.XInfoCalculate()
	return Captured{
		Now:   time.Now().Unix(),
		Idle:  int64(C.XInfoIdleTime()),
		App:   sanitize(C.GoString(C.XInfoFocussedWindowApp())),
		Class: sanitize(C.GoString(C.XInfoFocussedWindowClass())),
		Title: sanitize(C.GoString(C.XInfoFocussedWindowTitle())),
	}

}
