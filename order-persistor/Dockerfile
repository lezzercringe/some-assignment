FROM golang:latest AS builder
WORKDIR /order-persistor
COPY go.mod go.sum .
RUN go mod download
COPY . .
RUN go build -o /bin/order-persistor ./cmd/main.go 

FROM debian:latest
WORKDIR app
COPY --from=builder /bin/order-persistor .
COPY config.yaml .
EXPOSE 80
ENTRYPOINT ["./order-persistor", "-config", "config.yaml"]
