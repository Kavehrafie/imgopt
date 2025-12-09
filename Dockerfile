# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install build dependencies
# bimg requires libvips-dev and CGO enabled
RUN apk add --no-cache \
    build-base \
    pkgconfig \
    vips-dev \
    gcc \
    musl-dev

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
# CGO_ENABLED=1 is required for bimg
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o imgopt ./cmd/server/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    vips \
    ca-certificates \
    bash

# Copy binary from builder
COPY --from=builder /app/imgopt .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./imgopt"]
