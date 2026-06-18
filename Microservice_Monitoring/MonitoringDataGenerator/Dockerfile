FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /generator ./

FROM alpine:latest
WORKDIR /
COPY --from=builder /generator /generator
CMD ["/generator"]
