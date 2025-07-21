# WebAssembly build configuration
GOROOT ?= $(shell go env GOROOT)
WASM_EXEC = $(GOROOT)/misc/wasm/wasm_exec.js
GOOS = js
GOARCH = wasm
BUILD_DIR = dist

.PHONY: all build clean serve

all: build

# Build WebAssembly binary
build:
	@echo "Building WebAssembly application..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BUILD_DIR)/main.wasm main.go
	@cp $(WASM_EXEC) $(BUILD_DIR)/wasm_exec.js
	@cp index.html $(BUILD_DIR)/
	@cp styles.css $(BUILD_DIR)/
	@echo "Build complete!"

# Clean build artifacts
clean:
	@echo "Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "Clean complete!"

# Serve the application for local development
serve: build
	@echo "Serving WebAssembly app at http://localhost:8000"
	@cd $(BUILD_DIR) && python3 -m http.server
