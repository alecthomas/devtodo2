include $(GOROOT)/src/Make.inc

#everything: todo2 todo2.1

TARG=todo2
GOFILES=todo.go consoleview.go legacyio.go jsonio.go main.go

include $(GOROOT)/src/Make.cmd

todo2.1: $(GOFILES)
	./todo2 --create-manpage > $@
