# Root Makefile for ape-go-services

# Create a new service
create-service:
	@bash scripts/create_service.sh $(name)

# Clean up unnecessary .gitkeep files
clean-gitkeep:
	@bash scripts/clean-gitkeep.sh

# Run a specific service
run-service:
	@$(MAKE) -C services/$(name) run

run-all:
	@bash scripts/run_all_services.sh

# Build a specific service
build-service:
	@$(MAKE) -C services/$(name) build

# Test a specific service
test-service:
	@$(MAKE) -C services/$(name) test

# Format a specific service
fmt-service:
	@$(MAKE) -C services/$(name) fmt

# Tidy a specific service
tidy-service:
	@$(MAKE) -C services/$(name) tidy

# Format all Go code in the monorepo
fmt:
	@for dir in $(shell find services -name go.mod -exec dirname {} \;); do \
		echo "Formatting $$dir"; \
		go fmt $$dir/...; \
	done

# Tidy all modules
tidy:
	@for dir in $(shell find services -name go.mod -exec dirname {} \;); do \
		echo "🧹 Running 'go mod tidy' in $$dir"; \
		$(MAKE) -C $$dir tidy; \
	done

# Run all tests
test:
	@for dir in $(shell find services -name go.mod -exec dirname {} \;); do \
		echo "Running tests in $$dir"; \
		go test $$dir/... || exit 1; \
	done

clean:
	@echo "🧹 Cleaning all services..."
	@for dir in $(shell find services -name go.mod -exec dirname {} \;); do \
		echo "Cleaning $$dir"; \
		(cd $$dir && go clean -i -x -testcache -cache); \
	done
	@echo "Clean complete."

# Show help
help:
	@echo "Usage:"
	@echo "  make create-service name=<service-name>     # Create a new service"
	@echo "  make run-service name=<service-name>        # Run a service from root"
	@echo "  make build-service name=<service-name>      # Build a service from root"
	@echo "  make test-service name=<service-name>       # Test a service from root"
	@echo "  make fmt-service name=<service-name>        # Format a service from root"
	@echo "  make tidy-service name=<service-name>       # Tidy a service from root"
	@echo "  make fmt                                     # Format all services"
	@echo "  make tidy                                    # Tidy all modules"
	@echo "  make test                                    # Run all tests"
	@echo "  make clean-gitkeep                           # Remove unused .gitkeep files"
