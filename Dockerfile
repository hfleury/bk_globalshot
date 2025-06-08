FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /build/app cmd/main.go

# Final image
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /build/app /app/

# Install necessary libraries for pgx
RUN apk add --no-cache tzdata postgresql-libs

EXPOSE 8080

CMD ["/app/app"]