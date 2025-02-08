FROM golang:1.23-alpine as builder 

RUN apk add --no-cache git bash

WORKDIR /app

COPY go.mod go.sum ./ 
RUN go mod download

COPY vendor/ ./vendor/

COPY . .

COPY database/migrations /app/database/migrations

COPY public/ /app/public/

RUN go build -mod=vendor -o /app/fiber-starter ./cmd/app/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache bash

COPY --from=builder /app/fiber-starter /app/fiber-starter

EXPOSE 8080

CMD ["/app/fiber-starter"]