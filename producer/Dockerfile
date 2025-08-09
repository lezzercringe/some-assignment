FROM golang:latest AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/producer ./cmd/main.go

FROM debian:latest
WORKDIR /app
COPY --from=builder /bin/producer /bin/producer
