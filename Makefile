# Define variables
BINARY_NAME=nftscout
BUILD_DIR=build

# Default target
all: clean build-linux build-windows build-macos

# Build for Linux
build-linux:
	@echo "Building for Linux"
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/...

# Build for Windows
build-windows:
	@echo "Building for Windows"
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/...

# Build for macOS
build-macos:
	@echo "Building for macOS"
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts"
	rm -rf $(BUILD_DIR)

# Create build directory
$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

# Ensure build directory exists before building
build-linux build-windows build-macos: $(BUILD_DIR)

# Phony targets
.PHONY: all build-linux build-windows build-macos clean
