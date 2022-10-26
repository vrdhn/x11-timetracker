package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
)

// manages storing and archiving Captured

var (
	// Captured is immutable, so we can send pointer.
	chanCaptured = make(chan *Captured)
	// just send true to archive whatever is there.
	chanArchive = make(chan bool)

	fd *os.File = nil
)

func init() {
	go file_select_loop()
}
func exists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	if info.IsDir() {
		log.Fatalf("'%s' should  not be a directory", filename)
	}
	return true
}

func archive() {
	home, err := os.UserHomeDir()
	e(err)
	archiveFile := filepath.Join(home, ".x11-timetracker-archive.log")
	currentFile := filepath.Join(home, ".x11-timetracker.log")

	// 1. Close existing one.
	if fd != nil {
		err := fd.Close()
		e(err)
	}

	// 2. Copy content to archive file, don't worry about memory for now.
	if exists(currentFile) {
		cur, err := os.Open(currentFile)
		e(err)
		arc, err := os.OpenFile(archiveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		e(err)
		_, err = io.Copy(arc, cur)
		e(err)
		cur.Close()
		arc.Close()
	}

	//3. Start with a new clean
	fd, err = os.Create(currentFile)
	e(err)
}

func file_select_loop() {

	for {
		select {

		case rec := <-chanCaptured:
			if fd == nil {
				archive()
			}
			s, err := json.Marshal(rec)
			e(err)
			s = append(s, byte('\n'))
			log.Printf("%s", s)
			_, err = fd.Write(s)
			e(err)

		case <-chanArchive:
			archive()
		}
	}
}

func saveCaptured(i Captured) {
	chanCaptured <- &i
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
