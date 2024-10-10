# Define variables
BINARY_NAME=tracer  # Change to your desired binary name
LIBBPF_PATH=/usr/local/lib # Path to libbpf, update if installed elsewhere
LIBBPF_INCLUDE_PATH=/usr/local/include # Path to libbpf include files
GOFLAGS=-tags=cgo       # Specify cgo tags if needed

# Build target
.PHONY: build
build:
	@echo "Building ${BINARY_NAME}..."
	@CGO_CFLAGS="-I$(LIBBPF_INCLUDE_PATH)" CGO_LDFLAGS="-L$(LIBBPF_PATH) -lbpf" go build $(GOFLAGS) -o $(BINARY_NAME) .

# Run target
.PHONY: run
run: build
	@echo "Running ${BINARY_NAME}..."
	@./$(BINARY_NAME)

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -f $(BINARY_NAME)

# Install target (if needed)
.PHONY: install
install: build
	@echo "Installing ${BINARY_NAME}..."
	@cp $(BINARY_NAME) /usr/local/bin/

# Help target
.PHONY: help
help:
	@echo "Makefile commands:"
	@echo "  make build   - Build the Go binary"
	@echo "  make run     - Build and run the application"
	@echo "  make clean   - Remove build artifacts"
	@echo "  make install  - Install the binary to /usr/local/bin"
	@echo "  make help    - Display this help message"

