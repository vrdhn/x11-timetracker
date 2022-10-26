package main

import (
	"log"
	"time"
)

const (
	// 5 min or more is a official ''break''
	BREAK_TIME = 5 * 60
	// don't hold print more than this.
	DUMP_INTERVAL = 1 * 10
)

/* The idea is to minimize logged lines, while taking care of:
 *   suspends -> time suddently jumps.
 *   walkaway -> idle time keeps growing monotonically
 *   switch  -> window id changes.
 *
 */
func dumpBoth(prev Captured, printprev bool, next Captured, msg string) (Captured, bool) {
	if printprev {
		saveCaptured(prev)
	}
	if msg != "" {
		log.Printf("# %s", msg)
	}
	saveCaptured(next)
	return next, false
}

var nextPrintAt = int64(0)

func smoothen(prev Captured, printprev bool, next Captured) (Captured, bool) {
	sameApp := prev.isSameApp(next)
	if !sameApp {
		return dumpBoth(prev, printprev, next, "switch")
	}
	// are we back from sleep ?
	if next.Now-prev.Now > BREAK_TIME {
		return dumpBoth(prev, printprev, next, "break")
	}
	// if idling, or even if intensely working .. print the curr only...
	// the Time is in seconds .. so this has to be creative.
	if nextPrintAt < next.Now {
		nextPrintAt = next.Now + DUMP_INTERVAL
		return dumpBoth(prev, printprev && !sameApp, next, "continue")
	}
	// else we can skip printing !!
	return next, true

}

func main() {

	prev := getInfo()
	printprev := true
	for {
		time.Sleep(100 * time.Millisecond)
		prev, printprev = smoothen(prev, printprev, getInfo())
	}
}
