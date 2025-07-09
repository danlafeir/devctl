APP_NAME=devctl
BUILD_DIR=bin
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

# Default values (can be overridden)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

.PHONY: all build clean

# Usage:
#   make build           # builds for your current system, output: bin/devctl
#   GOOS=linux GOARCH=amd64 make build  # cross-compiles for linux/amd64

all: build

build:
	@echo "Building $(APP_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/$(APP_NAME) ./main.go

clean:
	rm -rf $(BUILD_DIR) 