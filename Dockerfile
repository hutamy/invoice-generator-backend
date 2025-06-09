# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o invoice-generator ./cmd/main.go

# Run stage
FROM alpine:3.19

WORKDIR /app

RUN apk --no-cache upgrade

COPY --from=builder /app/invoice-generator .

EXPOSE 8000

CMD ["./invoice-generator"]