# Build stage
FROM golang:1.24 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy && go build -o main ./cmd/server

# Final stage
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/.env .

CMD ["./main"]
