# Makefile for Wayland Color Picker GTK3 Application

# Variables
BINARY_NAME = wayland-color-picker-gtk4
MAIN_FILE = main.go
BUILD_DIR = .
INSTALL_PREFIX = /usr
DESKTOP_DIR = $(INSTALL_PREFIX)/share/applications
ICON_DIR = $(INSTALL_PREFIX)/share/pixmaps
HICOLOR_ICON_DIR = $(INSTALL_PREFIX)/share/icons/hicolor
BIN_DIR = $(INSTALL_PREFIX)/bin

# Go build flags
GO_BUILD_FLAGS = -ldflags="-s -w"

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	go build $(GO_BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_FILE)
	@echo "Build complete: $(BINARY_NAME)"

# Run the application
.PHONY: run
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BINARY_NAME)

# Build and install in one step (recommended)
.PHONY: build-install
build-install: build
	@echo "Building and installing $(BINARY_NAME)..."
	sudo $(MAKE) install

# Install the application system-wide
.PHONY: install
install:
	@if [ ! -f $(BINARY_NAME) ]; then \
		echo "Error: $(BINARY_NAME) not found. Please run 'make build' first."; \
		exit 1; \
	fi
	@echo "Installing $(BINARY_NAME) to $(INSTALL_PREFIX)..."

	# Create directories if they don't exist
	install -d $(BIN_DIR)
	install -d $(DESKTOP_DIR)
	install -d $(ICON_DIR)
	install -d $(HICOLOR_ICON_DIR)/48x48/apps
	install -d $(HICOLOR_ICON_DIR)/64x64/apps
	install -d $(HICOLOR_ICON_DIR)/128x128/apps
	install -d $(HICOLOR_ICON_DIR)/256x256/apps
	install -d $(HICOLOR_ICON_DIR)/scalable/apps

	# Install binary
	install -m 755 $(BINARY_NAME) $(BIN_DIR)/$(BINARY_NAME)

	# Install desktop file
	install -m 644 $(BINARY_NAME).desktop $(DESKTOP_DIR)/$(BINARY_NAME).desktop

	# Install icons
	install -m 644 icon.png $(ICON_DIR)/$(BINARY_NAME).png
	install -m 644 icon.png $(HICOLOR_ICON_DIR)/48x48/apps/$(BINARY_NAME).png
	install -m 644 icon.png $(HICOLOR_ICON_DIR)/64x64/apps/$(BINARY_NAME).png
	install -m 644 icon.png $(HICOLOR_ICON_DIR)/128x128/apps/$(BINARY_NAME).png
	install -m 644 icon.png $(HICOLOR_ICON_DIR)/256x256/apps/$(BINARY_NAME).png
	install -m 644 icon.png $(HICOLOR_ICON_DIR)/scalable/apps/$(BINARY_NAME).png

	# Update icon cache
	-gtk-update-icon-cache -q -t -f $(HICOLOR_ICON_DIR) 2>/dev/null || true
	-update-desktop-database $(DESKTOP_DIR) 2>/dev/null || true

	@echo "Installation complete!"
	@echo "You can now run '$(BINARY_NAME)' from anywhere or find it in your applications menu."
	@echo ""
	@echo "Note: If you experience permission issues, ensure you built as a regular user first."

# Uninstall the application
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(BIN_DIR)/$(BINARY_NAME)
	rm -f $(DESKTOP_DIR)/$(BINARY_NAME).desktop
	rm -f $(ICON_DIR)/$(BINARY_NAME).png
	rm -f $(HICOLOR_ICON_DIR)/48x48/apps/$(BINARY_NAME).png
	rm -f $(HICOLOR_ICON_DIR)/64x64/apps/$(BINARY_NAME).png
	rm -f $(HICOLOR_ICON_DIR)/128x128/apps/$(BINARY_NAME).png
	rm -f $(HICOLOR_ICON_DIR)/256x256/apps/$(BINARY_NAME).png
	rm -f $(HICOLOR_ICON_DIR)/scalable/apps/$(BINARY_NAME).png

	# Update icon cache after removal
	-gtk-update-icon-cache -q -t -f $(HICOLOR_ICON_DIR) 2>/dev/null || true
	-update-desktop-database $(DESKTOP_DIR) 2>/dev/null || true
	@echo "Uninstall complete!"

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f $(BINARY_NAME)
	rm -f main
	go clean
	@echo "Clean complete!"

# Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod tidy
	@echo "Dependencies updated!"

# Check code
.PHONY: check
check:
	@echo "Running go vet..."
	go vet ./...
	@echo "Code check complete!"

# Development build (with debug info)
.PHONY: debug
debug:
	@echo "Building debug version..."
	go build -gcflags="all=-N -l" -o $(BINARY_NAME)-debug $(MAIN_FILE)
	@echo "Debug build complete: $(BINARY_NAME)-debug"

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build         - Build the application as regular user"
	@echo "  run           - Build and run the application"
	@echo "  build-install - Build then install (recommended method)"
	@echo "  install       - Install the application system-wide (run 'make build' first)"
	@echo "  uninstall     - Remove the installed application"
	@echo "  clean         - Remove build artifacts"
	@echo "  deps          - Download and tidy dependencies"
	@echo "  check         - Run code checks (go vet)"
	@echo "  debug         - Build with debug information"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Recommended installation process:"
	@echo "  1. make clean"
	@echo "  2. make build-install"
	@echo ""
