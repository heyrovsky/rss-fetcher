BINARY_NAME := rsscurator
SRC_FILES   := $(wildcard *.go)
LINUX_OUT   := bin/linux/$(BINARY_NAME)


all: deps build

deps:
	@go mod download

build: $(SRC_FILES)
	@go build -o $(BINARY_NAME) cmd/rss/main.go

linux: $(SRC_FILES)
	@GOOS=linux GOARCH=amd64 go build -o $(LINUX_OUT) cmd/rss/main.go

clean:
	@rm -rf $(BINARY_NAME) $(LINUX_OUT) $(DARWIN_OUT) $(WINDOWS_OUT)

.PHONY: all deps build linux darwin windows clean