.PHONY: build install clean test fmt vet

# Binary name
BINARY_NAME=docmind

# Install directory (defaults to ~/bin, can be overridden)
INSTALL_DIR?=$(HOME)/bin

# Build the binary
build:
	go build -o $(BINARY_NAME) .

# Install the binary to ~/bin
install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"

# Clean build artifacts
clean:
	go clean
	rm -f $(BINARY_NAME)

# Run tests
test:
	go test -v ./...

# Format code
fmt:
	go fmt ./...

# Run go vet
vet:
	go vet ./...

# Run all checks
check: fmt vet test

# Uninstall the binary
uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstalled $(BINARY_NAME) from $(INSTALL_DIR)"
