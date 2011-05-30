include $(GOROOT)/src/Make.inc

everything: todo2 todo2.1

TARG=todo2
GOFILES=tiocgwinsz.go todo.go consoleview.go legacyio.go jsonio.go main.go

include $(GOROOT)/src/Make.cmd

tiocgwinsz: tiocgwinsz.c
	cc $< -o $@

tiocgwinsz.go: tiocgwinsz
	./tiocgwinsz

todo2.1: todo2
	./todo2 --create-manpage > $@
