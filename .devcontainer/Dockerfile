# Base image with Go 1.24
FROM mcr.microsoft.com/devcontainers/go:1.24

# Install zsh, oh-my-zsh, and dependencies
RUN apt-get update && \
    apt-get install -y zsh curl git vim libpq-dev locales && \
    locale-gen en_US.UTF-8 && \
    usermod --shell /usr/bin/zsh vscode && \
    sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" "" --unattended && \
    git clone --depth=1 https://github.com/romkatv/powerlevel10k.git /home/vscode/.oh-my-zsh/custom/themes/powerlevel10k && \
    sed -i 's|ZSH_THEME=".*"|ZSH_THEME="powerlevel10k/powerlevel10k"|' /home/vscode/.zshrc && \
    chown -R vscode:vscode /home/vscode/.oh-my-zsh /home/vscode/.zshrc /home/vscode/.p10k.zsh || true

# Set working directory
WORKDIR /workspace/ape-go-services


# Install Go tools globally if needed (e.g., golang-migrate, mockgen, etc.)
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest
# RUN go install github.com/vektra/mockery/v2@latest
