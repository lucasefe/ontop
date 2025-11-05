.PHONY: build test run install clean lint help

# Build the ontop binary
build:
	go build -o ontop cmd/ontop/main.go

# Run all tests
test:
	go test ./... -v

# Run ontop directly (TUI mode)
run:
	go run cmd/ontop/main.go

# Install ontop to $GOPATH/bin
install:
	go install cmd/ontop/main.go

# Remove binary and test database
clean:
	rm -f ontop
	rm -f ~/.ontop/tasks.db

# Format and lint the code
lint:
	@echo "Formatting code..."
	go fmt ./...
	@echo "Running golangci-lint..."
	golangci-lint run
	@echo "Linting complete!"

# Show available targets
help:
	@echo "Available targets:"
	@echo "  build   - Build the ontop binary"
	@echo "  test    - Run all tests"
	@echo "  run     - Run ontop directly (TUI mode)"
	@echo "  install - Install ontop to \$$GOPATH/bin"
	@echo "  clean   - Remove binary and test database"
	@echo "  lint    - Format and lint code (go fmt, go vet, golangci-lint)"
