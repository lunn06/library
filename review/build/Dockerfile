FROM golang:1.24-alpine AS base

WORKDIR /src

COPY .. .
RUN go mod download

FROM base AS builder

RUN go build ./cmd/review

FROM alpine:3

WORKDIR /app

COPY --from=builder /src/review /usr/bin/review
CMD ["review"]