# Build the service binary
build:
	go build -o bin/boilerplate-service ./cmd/server

# Run the service locally
run:
	go run ./cmd/server

# Run all tests (unit tests by default) in verbose mode
test:
	go test -v ./...

# Run only the integration tests (requires database and the 'integration' tag) in verbose mode
test-integration:
	go test -v -tags=integration ./...

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy