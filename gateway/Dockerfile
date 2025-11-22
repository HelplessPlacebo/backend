FROM golang:1.25.2-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/gateway
RUN go build -o /bin/gateway ./cmd

FROM alpine:latest
WORKDIR /app
COPY --from=builder /bin/gateway .
CMD ["./gateway"]
