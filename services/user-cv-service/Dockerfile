# Build stage
FROM golang:1.24 AS builder

# Set working dir before copying
WORKDIR /app

# Copy the entire monorepo (from root)
COPY . .

# Move into the service directory
WORKDIR /app/services/user-cv-service

RUN go mod tidy && go build -o main ./cmd/server

# Final stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/services/user-cv-service/main .
COPY --from=builder /app/services/user-cv-service/.env .

CMD ["./main"]
