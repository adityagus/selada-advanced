# Stage 1: Build the Go binary
FROM golang:1.23-alpine AS builder

# Install essential dependencies
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy dependencies manifest
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked release binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go-api

# Stage 2: Production container
FROM alpine:latest

# Install tzdata for timezone configurations and ca-certificates for secure HTTPS calls
RUN apk add --no-cache ca-certificates tzdata

# Create a non-root system group and user for container security hardening (SonarQube)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy built executable from build stage
COPY --from=builder /go-api /app/go-api

# Configure file ownership permissions
RUN chown -R appuser:appgroup /app

# Enforce running the container under non-root context
USER appuser

EXPOSE 3000

ENV GIN_MODE=release

ENTRYPOINT ["/app/go-api"]
