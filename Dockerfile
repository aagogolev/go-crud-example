FROM golang:1.22-alpine AS builder
WORKDIR /app

COPY go-crud-example/go.mod go-crud-example/go.sum ./
RUN go mod download
COPY . .
RUN cd go-crud-example && CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

FROM alpine:3.18
WORKDIR /app

COPY --from=builder /app/go-crud-example/main .
COPY --from=builder /app/.env .
CMD ["./main"]
