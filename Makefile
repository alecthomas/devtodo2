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
	./$(TARG) --create-manpage > $@

install: $(TARG) $(TARG).1
	mkdir -p $(bindir)
	mkdir -p $(mandir)
	install -m755 $(TARG) $(bindir)
	install -m644 $(TARG).1 $(man1dir)

clean:
	rm -f $(TARG) $(TARG).1
