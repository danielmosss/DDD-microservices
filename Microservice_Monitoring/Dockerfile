# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /monitoring-server .

# Run stage
FROM alpine:latest
WORKDIR /

COPY --from=builder /monitoring-server /monitoring-server

EXPOSE 8080

CMD ["/monitoring-server"]