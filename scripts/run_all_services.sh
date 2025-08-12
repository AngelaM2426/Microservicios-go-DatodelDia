#!/bin/bash

set -e

SERVICES=$(ls services)
PIDS=()

echo "🚀 Starting all services..."

for service in $SERVICES; do
  (
    cd services/$service
    echo "▶️  Running $service"
    go run ./cmd/server
  ) &
  PIDS+=($!)
done

# Trap Ctrl+C and terminate all
trap "echo 'Shutting down...'; for pid in \${PIDS[@]}; do kill \$pid 2>/dev/null; done; exit" SIGINT

# Wait for all background jobs
wait
