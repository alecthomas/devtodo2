include $(GOROOT)/src/Make.inc

TARG=todo2
GOFILES=tiocgwinsz.go todo.go consoleview.go legacyio.go jsonio.go main.go

include $(GOROOT)/src/Make.cmd

tiocgwinsz: tiocgwinsz.c
	cc $< -o $@

tiocgwinsz.go: tiocgwinsz
	./tiocgwinsz
