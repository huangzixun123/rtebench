GIT_COMMIT=$(shell git rev-parse --short HEAD)
OS=$(shell uname -s)
GOARCH=$(shell go env GOARCH)
ifeq ($(GOARCH),)
	GOARCH="amd64"
endif
GOOS=$(shell go env GOOS)
ifeq ($(GOOS),)
	GOOS="linux"
endif

build:
	mkdir -p bin
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o bin/rtectl main.go