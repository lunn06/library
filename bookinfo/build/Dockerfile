FROM golang:1.24-alpine AS base

WORKDIR /src

COPY .. .
RUN go mod download

FROM base AS builder

RUN go build ./cmd/bookinfo

FROM alpine:3

WORKDIR /app
COPY --from=builder /src/bookinfo /usr/bin/bookinfo
CMD ["bookinfo"]