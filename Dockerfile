FROM golang:1.23-alpine as builder

RUN apk add --no-cache git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY vendor/ ./vendor/
COPY . .
RUN go build -mod=vendor -o /app/fiber-starter ./cmd/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/fiber-starter /app/fiber-starter
RUN chmod +x /app/fiber-starter
EXPOSE 3000
CMD ["/app/fiber-starter"]