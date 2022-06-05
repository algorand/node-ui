.PHONY: build

MAINFILE := cmd/tui/main.go
BINNAME := nodeui

build:
	go build -o $(BINNAME) $(MAINFILE)
