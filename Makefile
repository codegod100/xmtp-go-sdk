.PHONY: all build test clean example

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
EXAMPLE_BIN=example

# Rust parameters for building the FFI library
CARGO=cargo
RUST_TARGETS := libxmtp_ffi.so libxmtp_ffi.dylib xmtp_ffi.dll

all: build

# Build the Go SDK
build:
	$(GOBUILD) ./...

# Build the FFI library (requires Rust)
ffi:
	cd ffi && $(CARGO) build --release
	@echo "FFI library built at: ffi/target/release/"

# Copy the appropriate FFI library for your platform
install-ffi: ffi
	@if [ -f ffi/target/release/libxmtp_ffi.so ]; then \
		cp ffi/target/release/libxmtp_ffi.so .; \
		echo "Copied libxmtp_ffi.so"; \
	elif [ -f ffi/target/release/libxmtp_ffi.dylib ]; then \
		cp ffi/target/release/libxmtp_ffi.dylib .; \
		echo "Copied libxmtp_ffi.dylib"; \
	elif [ -f ffi/target/release/xmtp_ffi.dll ]; then \
		cp ffi/target/release/xmtp_ffi.dll .; \
		echo "Copied xmtp_ffi.dll"; \
	else \
		echo "FFI library not found. Run 'make ffi' first."; \
		exit 1; \
	fi

# Run tests
test:
	$(GOTEST) -v ./...

# Run tests with coverage
test-cov:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Build the example
example:
	cd examples/basic && $(GOBUILD) -o ../../$(EXAMPLE_BIN)

# Run the example
run-example: example
	XMTP_FFI_PATH=./libxmtp_ffi.so ./$(EXAMPLE_BIN)

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(EXAMPLE_BIN)
	rm -f coverage.out coverage.html

# Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Generate C header from Rust FFI
generate-header:
	cd ffi && $(CARGO) build
	@echo "Header generated at: ffi/include/xmtp_ffi.h"

# Format code
fmt:
	$(GOCMD) fmt ./...
	cd ffi && cargo fmt

# Lint code
lint:
	$(GOCMD) vet ./...

# Help
help:
	@echo "XMTP Go SDK Makefile"
	@echo ""
	@echo "Targets:"
	@echo "  all          - Build the SDK (default)"
	@echo "  build        - Build the Go SDK"
	@echo "  ffi          - Build the Rust FFI library"
	@echo "  install-ffi  - Build and copy the FFI library"
	@echo "  test         - Run tests"
	@echo "  test-cov     - Run tests with coverage"
	@echo "  example      - Build the example"
	@echo "  run-example  - Build and run the example"
	@echo "  clean        - Remove build artifacts"
	@echo "  deps         - Download dependencies"
	@echo "  generate-header - Generate C header from Rust"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
