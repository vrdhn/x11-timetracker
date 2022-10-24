package main

//   #cgo CFLAGS : -Wall -Wpedantic
//   #cgo LDFLAGS: -lXss -lX11 -lXmu
//   #include "x-info.h"
import "C"

import (
	"strings"
	"time"
)

type info struct {
	Now   int64
	Idle  int64
	App   string
	Class string
	Title string
}

func init() {
	C.XInfoInitialize()
}

// Also force make a copy of the string, as C is returning a static array
func sanitize(in string) string {
	in = in[:]
	return strings.TrimSpace(in)
}
func getInfo() info {
	C.XInfoCalculate()
	return info{
		Now:   time.Now().Unix(),
		Idle:  int64(C.XInfoIdleTime()),
		App:   sanitize(C.GoString(C.XInfoFocussedWindowApp())),
		Class: sanitize(C.GoString(C.XInfoFocussedWindowClass())),
		Title: sanitize(C.GoString(C.XInfoFocussedWindowTitle())),
	}

}

func isInfoSameApp(a info, b info) bool {
	return a.App == b.App && a.Class == b.Class && a.Title == b.Title
}
