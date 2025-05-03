# Stage 1 - Install dependencies
FROM golang:1.24-alpine as base

WORKDIR /src

COPY ../go.mod go.sum ./
RUN go mod download
COPY ../.. .

# Stage 2 - Linter
FROM golangci/golangci-lint:v2.1.5-alpine as linter

WORKDIR /src
COPY ../.. .
RUN golangci-lint run

# Stage 3 - Build application
FROM base as builder

COPY --from=linter /src/lint_report.json .
RUN go build ./cmd/greeter-app

# Stafe 4 - Build final image
FROM alpine:3

WORKDIR /app
COPY --from=builder /src/book /usr/bin/book
CMD ["book"]