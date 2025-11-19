# Step 1: Build the Go binary
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install git (needed for Go modules) and bash if needed
RUN apk add --no-cache git bash

# Copy go.mod and go.sum first (for caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
RUN go build -o user-service .

# Step 2: Minimal image for running
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/user-service .

# Copy migrations if needed
COPY ./migrations ./migrations

EXPOSE 8080

CMD ["./user-service"]
