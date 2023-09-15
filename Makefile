.POSIX:
PREFIX ?= /usr/local
MANPREFIX ?= $(PREFIX)/share/man
GO ?= go
GOFLAGS ?=

all: clean build

build:
	go build

clean:
	rm -f battery

install: build
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp -f battery $(DESTDIR)$(PREFIX)/bin
	chmod 755 $(DESTDIR)$(PREFIX)/bin/battery

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/battery

.PHONY: all build clean install uninstall
