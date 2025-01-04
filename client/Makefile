# Makefile for the indicum project

# Project directories
CLIENT_DIR := client

# Binary output
CLIENT_BIN := client-indicum

# Go command
GO_CMD := go

# Build commands
GO_BUILD := $(GO_CMD) build
GO_CLEAN := $(GO_CMD) clean
GO_TEST := $(GO_CMD) test

# Default target
all: client

# Build client
client:
	$(GO_BUILD) -o $(CLIENT_BIN)

# Clean up binaries
clean:
	$(GO_CLEAN)
	rm -f $(CLIENT_BIN)

# Run tests
test:
	$(GO_TEST) ./...


.PHONY: all client clean test
