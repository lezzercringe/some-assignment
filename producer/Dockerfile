FROM golang:latest AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /bin/produce ./cmd/main.go

FROM debian:latest
WORKDIR /app
COPY --from=builder /bin/produce /bin/produce
COPY config.yaml .
CMD ["tail", "-f", "/dev/null"]
