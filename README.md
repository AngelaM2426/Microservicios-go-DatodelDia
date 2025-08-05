## ape-go-services

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

# Run boilerplate service
make run-service name=boilerplate-service

# Test boilerplate health path
curl http://localhost:8088/health

## Running tests

# Test from the root directory
make test-service name=boilerplate-service

# Or navigate to the service directory and run directly
cd services/boilerplate-service
make test
make test-integration
```

### Rebuild the Devcontainer in VS Code

- Open the Command Palette in VS Code (Ctrl+Shift+P or Cmd+Shift+P on Mac).
- Type Rebuild Container and select the command: "Dev Containers: Rebuild Container".
- VS Code will rebuild the container from your Dockerfile and then run the updated postCreateCommand.bash script.

```bash
# Install the Dev Containers CLI
npm install -g @devcontainers/cli

# Rebuild devcontainer from command line
devcontainer up --workspace-folder . --remove-existing-container
```

### You could have to install nvm and npm

```bash
# Install nvm
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.7/install.sh | bash

# Reload your shell to use nvm
source ~/.bashrc  # Or source ~/.zshrc if you use zsh

# Install a recent version of Node.js
nvm install --lts

# Now run the installation command globally inside your WSL terminal
npm install -g @devcontainers/cli

# Check installation
devcontainer --version
```
