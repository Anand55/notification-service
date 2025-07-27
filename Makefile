.PHONY: build run test clean examples

# Build the notification service
build:
	go build -o notification-service main.go

# Run the notification service
run:
	go run main.go

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -f notification-service
	rm -f test_examples
	rm -f *.db

# Build and run examples
examples:
	go build -o test_examples test_examples.go
	./test_examples

# Install dependencies
deps:
	go mod tidy
	go mod download

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	air

# Build for different platforms
build-linux:
	GOOS=linux GOARCH=amd64 go build -o notification-service-linux main.go

build-windows:
	GOOS=windows GOARCH=amd64 go build -o notification-service.exe main.go

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o notification-service-darwin main.go

# Help
help:
	@echo "Available commands:"
	@echo "  build        - Build the notification service"
	@echo "  run          - Run the notification service"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  examples     - Build and run examples"
	@echo "  deps         - Install dependencies"
	@echo "  dev          - Run with hot reload (requires air)"
	@echo "  build-linux  - Build for Linux"
	@echo "  build-windows- Build for Windows"
	@echo "  build-darwin - Build for macOS"
	@echo "  help         - Show this help" 