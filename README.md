# ape-go-services

```bash
make run-service name=<service-name>
make run-all
build-service name=<service-name> #This will build a specific service.
test-service name=<service-name> #This runs the tests for a particular service.
test #This command will run the tests for all services.
fmt #Use this to format the Go code across all your services.
tidy #This command tidies the go.mod files for all services.

# Clean all build artifacts, test caches, and binaries
make clean 

## Running tests

# Test from the root directory
make test-service name=website-cms-service

# Or navigate to the service directory and run directly
cd services/website-cms-service
make test-mongo
```