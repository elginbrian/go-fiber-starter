FROM golang:1.23-alpine as builder

# Install git for dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the vendor directory (after running go mod vendor)
COPY vendor/ ./vendor/

# Copy the rest of the application code
COPY . .

# Build the Go binary using the vendor directory
RUN go build -mod=vendor -o /app/fiber-starter ./cmd/main.go

# Start with a clean Alpine image for the runtime
FROM alpine:latest

# Set the working directory for runtime
WORKDIR /app

# Copy the built binary from the builder image
COPY --from=builder /app/fiber-starter /app/fiber-starter

# Make the binary executable
RUN chmod +x /app/fiber-starter

# Expose the application port
EXPOSE 3000

# Set the entry point for the container to run the application
CMD ["/app/fiber-starter"]