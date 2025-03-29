FROM golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o neko-bot ./cmd

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/neko-bot .
CMD ["./neko-bot"]