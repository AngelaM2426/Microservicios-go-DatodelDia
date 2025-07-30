# Root Makefile for ape-go-services

# Create a new service
create-service:
	#Example: make create-service name=admin-service
	@bash scripts/create_service.sh $(name)

# Clean up unnecessary .gitkeep files
clean-gitkeep:
	@bash scripts/clean-gitkeep.sh

# Format all Go code
fmt:
	@go fmt ./...

# Tidy all modules
tidy:
	@go work sync
	@find services -name go.mod -execdir go mod tidy \;

# Run tests recursively
test:
	@go test ./...

# Show help
help:
	@echo "Usage:"
	@echo "  make create-service name=<service-name>   # Create a new service"
	@echo "  make clean-gitkeep                        # Remove unused .gitkeep files"
	@echo "  make fmt                                  # Format Go code"
	@echo "  make tidy                                 # Tidy all go.mod files"
	@echo "  make test                                 # Run all tests"
