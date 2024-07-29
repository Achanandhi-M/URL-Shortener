# Stage 1: Build the Go application
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Build the Go application

RUN go build -o url-shortener

#  smaller image with the compiled binary

FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/url-shortener /app/

EXPOSE 8080


CMD ["./url-shortener"]
