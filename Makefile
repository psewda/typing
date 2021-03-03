BASE_DIR=$(shell pwd)
OUTPUT_DIR=$(BASE_DIR)/bin
SERVER=$(BASE_DIR)/cmd/server/main.go
PKG=github.com/psewda/typing
BUILD_NUMBER=$(shell echo $${TRAVIS_BUILD_NUMBER:-1})
APP=typing

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -o $(OUTPUT_DIR)/linux-amd64/$(APP) $(SERVER)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64
	go build -o $(OUTPUT_DIR)/windows-amd64/$(APP).exe $(SERVER)

build: build-linux build-windows

all: build

run:
	go run $(SERVER)
