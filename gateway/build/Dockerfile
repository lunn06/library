FROM golang:1.24-alpine AS base

WORKDIR /src

COPY .. .
RUN go mod download

FROM base AS builder

RUN go build ./cmd/gateway

FROM alpine:3

WORKDIR /app

EXPOSE 8080

COPY --from=builder /src/gateway /usr/bin/gateway
CMD ["gateway"]