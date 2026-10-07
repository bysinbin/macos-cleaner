VERSION ?= v2.0.0
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +"%Y-%m-%d")
LDFLAGS = -s -w -X disk-cleaner/internal/version.Version=$(VERSION) -X disk-cleaner/internal/version.GitCommit=$(COMMIT) -X disk-cleaner/internal/version.BuildDate=$(BUILD_DATE)

.PHONY: all build run test test-all app release clean help

all: build

build:
	@echo "-> Building disk-cleaner binary for host architecture..."
	go build -ldflags "$(LDFLAGS)" -o disk-cleaner .

run: build
	@echo "-> Running disk-cleaner..."
	./disk-cleaner

test:
	@echo "-> Running fast unit tests..."
	go test -short ./...

test-all:
	@echo "-> Running all tests (including live system benchmarks)..."
	go test -v ./...

app:
	@echo "-> Building DiskCleaner.app and Universal binary..."
	chmod +x scripts/build-macos-app.sh
	./scripts/build-macos-app.sh $(VERSION)

release: app
	@echo "-> Release package created under build/ directory."

clean:
	@echo "-> Cleaning build artifacts..."
	rm -rf build disk-cleaner

help:
	@echo "Available commands:"
	@echo "  make build     - Compile local binary for current architecture"
	@echo "  make run       - Build and launch disk-cleaner in browser"
	@echo "  make test      - Run fast unit tests"
	@echo "  make test-all  - Run all integration tests"
	@echo "  make app       - Create DiskCleaner.app macOS bundle and universal binary"
	@echo "  make release   - Generate release zip/tar.gz packages with checksums"
	@echo "  make clean     - Remove binaries and build directories"
