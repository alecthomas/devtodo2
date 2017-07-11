TARG=todo2
GOFILES=todo.go view.go consoleview.go legacyio.go jsonio.go main.go importer.go
PREFIX=/usr/local
prefix=$(PREFIX)
bindir=$(prefix)/bin
mandir=$(prefix)/share/man
man1dir=$(mandir)/man1

all: $(TARG) $(TARG).1

$(TARG): $(GOFILES)
	go build -o $@ $^

$(TARG).1: $(TARG)
	./$(TARG) --help-man > $@

install:
	install -m755 $(TARG) $(bindir)/$(TARG)
	install -d $(man1dir)
	install -m644 $(TARG).1 $(man1dir)/$(TARG).1

uninstall:
	rm $(bindir)/$(TARG)
	rm $(man1dir)/$(TARG).1
