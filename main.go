package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	// 5 min or more is a official ''break''
	BREAK_TIME = 5 * 60
	// don't hold print more than this.
	DUMP_INTERVAL = 1 * 10
)

func dump(f *os.File, i info) {
	s, err := json.Marshal(i)
	e(err)
	s = append(s, byte('\n'))
	log.Printf("%s", s)
	_, err = f.Write(s)
	e(err)
}

/* The idea is to minimize logged lines, while taking care of:
 *   suspends -> time suddently jumps.
 *   walkaway -> idle time keeps growing monotonically
 *   switch  -> window id changes.
 *
 */
func dumpBoth(f *os.File, prev info, printprev bool, next info, msg string) (info, bool) {
	if printprev {
		dump(f, prev)
	}
	if msg != "" {
		log.Printf("# %s", msg)
	}
	dump(f, next)
	return next, false
}

var nextPrintAt = int64(0)

func smoothen(f *os.File, prev info, printprev bool, next info) (info, bool) {
	sameApp := isInfoSameApp(prev, next)
	if !sameApp {
		return dumpBoth(f, prev, printprev, next, "switch")
	}
	// are we back from sleep ?
	if next.Now-prev.Now > BREAK_TIME {
		return dumpBoth(f, prev, printprev, next, "break")
	}
	// if idling, or even if intensely working .. print the curr only...
	// the Time is in seconds .. so this has to be creative.
	if nextPrintAt < next.Now {
		nextPrintAt = next.Now + DUMP_INTERVAL
		return dumpBoth(f, prev, printprev && !sameApp, next, "continue")
	}
	// else we can skip printing !!
	return next, true

}

func main() {

	home, err := os.UserHomeDir()
	e(err)
	fd, err := os.OpenFile(filepath.Join(home, ".x11-timetracker.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	e(err)
	prev := getInfo()
	printprev := true
	for {
		time.Sleep(100 * time.Millisecond)
		prev, printprev = smoothen(fd, prev, printprev, getInfo())
	}
}

func e(err error) {
	if err != nil {
		panic(err)
	}
}
