x11-timetracker
===============


A simple X11 time tracker, polls for window having keyboard focus,
and logs it in a file.

install
=======

Needs cgo

```
go build
./x11-timetracker &
```


Task List
=========
- [X] basic poll and log to a file.
- [X] strategy for managing multiple files.
- [ ] find already running instance, and communicate with it.
- [ ] a grading engine to classify time 
- [ ] a library of well-known-applications and window titles.
- [ ] a simple web server showing usages.
- [ ] plugin system to capture more context.
