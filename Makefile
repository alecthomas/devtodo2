include $(GOROOT)/src/Make.inc

all: todo2 todo2.1

TARG=todo2
GOFILES=todo.go view.go consoleview.go legacyio.go jsonio.go main.go importer.go

include $(GOROOT)/src/Make.cmd

todo2.1: todo2
	./todo2 --create-manpage > $@
