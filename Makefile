.PHONY: clean

TARG = todo2
GOFILES = todo.go view.go consoleview.go legacyio.go jsonio.go main.go importer.go
PREFIX ?= /usr
prefix = $(PREFIX)
bindir = $(prefix)/bin
mandir = $(prefix)/share/man
man1dir = $(mandir)/man1

all: $(TARG) $(TARG).1

$(TARG): $(GOFILES)
	go build -o $@ $^

$(TARG).1: $(TARG)
	./$(TARG) --help-man > $@

install:
	install -Dm755 "$(TARG)" "$(bindir)/$(TARG)"
	install -d -m644 "$(man1dir)"
	install -m644 "$(TARG).1" "$(man1dir)/$(TARG).1"

uninstall:
	rm "$(bindir)/$(TARG)" "$(man1dir)/$(TARG).1"

clean:
	rm -f "$(TARG)"
