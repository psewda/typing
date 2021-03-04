BASE_DIR=$(shell pwd)
OUTPUT_DIR=$(BASE_DIR)/bin
SERVER=$(BASE_DIR)/cmd/server/main.go
PKG=github.com/psewda/typing
BUILD_NUMBER=$(shell echo $${TRAVIS_BUILD_NUMBER:-1})
LDFLAGS="-s -w -X $(PKG).BuildNumber=$(BUILD_NUMBER) -X main.build=RELEASE"
LINT_INSTALL_SCRIPT=https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh
APP=typing

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64
	go build -ldflags $(LDFLAGS) -o $(OUTPUT_DIR)/linux-amd64/$(APP) $(SERVER)

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64
	go build -ldflags $(LDFLAGS) -o $(OUTPUT_DIR)/windows-amd64/$(APP).exe $(SERVER)

build: build-linux build-windows

test:
	ginkgo ./...

lint:
	golangci-lint run ./...

all: lint build test

install-ginkgo:
	go get -v github.com/onsi/ginkgo/ginkgo@v1.14.2

install-golangci-lint:
	curl -sSfL $(LINT_INSTALL_SCRIPT) | sh -s -- -b $(shell go env GOPATH)/bin v1.35.2

run:
	go run $(SERVER)
