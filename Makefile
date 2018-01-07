
_SOURCES = $(wildcard *.go)

all: Physick

run: Physick
	./Physick

Physick: $(_SOURCES)
	go build .
