# Use Go base image to build the binary
FROM golang:1.21 AS builder

# Set working directory inside container
WORKDIR /app

# Copy Go files
COPY . .

# Download dependencies
RUN go mod init text-to-video-api && go mod tidy

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# Final lightweight image
FROM alpine:latest

# Working dir
WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/server .

# Expose API port
EXPOSE 8080

# Run the binary
CMD ["./server"]
