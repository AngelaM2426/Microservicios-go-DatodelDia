#!/bin/bash

# Grant the vscode user permissions to use the Docker socket
sudo usermod -aG docker vscode
echo "postCreateCommand script finished successfully."

set -e

echo "Initializing .env files for all services..."

BASE_PORT=8080
POSTGRES_PORT=5432
DB_USER="ape_user"
DB_PASSWORD="123456"

SERVICES=(
  "admin-service"
  "analytics-service"
  "document-service"
  "job-offers-service"
  "job-posting-service"
  "public-data-service"
  "user-cv-service"
  "user-service"
  "website-cms-service"
)

for index in "${!SERVICES[@]}"; do
  SERVICE="${SERVICES[$index]}"
  PORT=$((BASE_PORT + index))
  DB_NAME=$(echo "$SERVICE" | sed 's/-/_/g')
  DB_HOST="postgres-$SERVICE"
  ENV_FILE="services/$SERVICE/.env"

  if [ -f "$ENV_FILE" ]; then
    echo "$ENV_FILE already exists — skipping"
  else
    echo "Generating $ENV_FILE"
    cat <<EOF > "$ENV_FILE"
PORT=$PORT
DB_HOST=$DB_HOST
DB_PORT=$POSTGRES_PORT
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_NAME=$DB_NAME
EOF
  fi
done

echo "✅ .env setup complete!"


# Copy VSCode settings
mkdir -p .vscode
cp .devcontainer/settings.template.json .vscode/settings.json
echo "🛠️  Copied settings.template.json to .vscode/settings.json"