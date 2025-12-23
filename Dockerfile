# Build Stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main cmd/main.go

# Run Stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080
CMD ["./main"]
