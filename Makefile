include $(GOROOT)/src/Make.inc

TARG=todo
GOFILES=tiocgwinsz.go todo.go filter.go consoleview.go legacyloader.go main.go

include $(GOROOT)/src/Make.cmd

tiocgwinsz: tiocgwinsz.c
	gcc $< -o $@

tiocgwinsz.go: tiocgwinsz
	./tiocgwinsz
