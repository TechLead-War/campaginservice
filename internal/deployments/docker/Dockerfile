# syntax=docker/dockerfile:1
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o campaignservice ./cmd/api/main.go

# Start a new stage for a minimal image
FROM alpine:latest
WORKDIR /app

# Copy the built binary
COPY --from=builder /app/campaignservice .
# Copy .env if present (optional)
COPY .env .

# Expose the default app port
EXPOSE 8080

# Command to run
CMD ["./campaignservice"] 