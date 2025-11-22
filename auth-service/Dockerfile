FROM golang:1.25.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/auth-service
RUN go build -o /bin/auth-service ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/auth-service .
CMD ["./auth-service"]
