TARG=todo2
GOFILES=todo.go view.go consoleview.go legacyio.go jsonio.go main.go importer.go
PREFIX=/usr/local

all: $(TARG) $(TARG).1

$(TARG): $(GOFILES)
	go build -o $@ $^

$(TARG).1: $(TARG)
	./$(TARG) --create-manpage > $@

install: $(TARG) $(TARG).1
	install -m755 $(TARG) $(PREFIX)/bin
	install -m644 $(TARG).1 $(PREFIX)/share/man/man1

clean:
	rm -f $(TARG)
