# Build the service binary
build:
	go build -o bin/website-cms-service ./cmd/server

# Run the service locally
run:
	go run ./cmd/server

# Run MongoDB test
test-mongo:
	go run ./cmd/test_mongo

# Run tests
test:
	go test ./...

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy