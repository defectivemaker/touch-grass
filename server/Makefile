# Project directories
CMD_SERVER_DIR := cmd/server
INTERNAL_DIR := internal

# Binary output
SERVER_BIN := server-indicum

# Go command
GO_CMD := go

# Build commands
GO_BUILD := $(GO_CMD) build
GO_CLEAN := $(GO_CMD) clean
GO_TEST := $(GO_CMD) test
GO_FMT := $(GO_CMD) fmt

DOCKER_FLAGS := CGO_ENABLED=0 GOOS=linux

# Default target
all: server

# Build server
server:
	cd $(CMD_SERVER_DIR) && $(GO_BUILD) -o ../../$(SERVER_BIN)

docker-server:
	cd $(CMD_SERVER_DIR) && $(DOCKER_FLAGS) $(GO_BUILD) -o ../../$(SERVER_BIN)

# Clean up binaries
clean:
	$(GO_CLEAN)
	rm -f $(SERVER_BIN)

# Run tests
test:
	$(GO_TEST) ./...

# Format code
fmt:
	$(GO_FMT) ./...

.PHONY: all server clean test fmt
