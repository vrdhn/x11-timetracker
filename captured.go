package main

// the fields that we actually capture
// so tags are not part of this.
// THIS IS IMMUTABLE.
type Captured struct {
	Now   int64
	Idle  int64
	App   string
	Class string
	Title string
}

func (a Captured) isSameApp(b Captured) bool {
	return a.App == b.App && a.Class == b.Class && a.Title == b.Title
}
