.PHONY: build

MAINFILE := cmd/tui/main.go
BINNAME := nodeui

build:
	go build -o $(BINNAME) $(MAINFILE)

fmt:
	go fmt ./...

lint:
	golint -set_exit_status ./...
	go vet ./...
