all: build install

build:
	go build .

install:
	cp syncrepos ${HOME}/.local/bin/sr

.PHONY: build install
