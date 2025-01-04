FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/fiber-starter ./cmd/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fiber-starter /app/fiber-starter
RUN chmod +x /app/fiber-starter
EXPOSE 3000
CMD ["/app/fiber-starter"]