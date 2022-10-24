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
- [ ] strategy for managing multiple files.
- [ ] a simple web server showing usages.
- [ ] a grading engine to classify time 
- [ ] a library of well-known-applications and window titles.
