# Stage 1 - Install dependencies
FROM golang:1.24-alpine as base

WORKDIR /src

COPY .. .
RUN go mod download

# Stage 2 - Linter
#FROM golangci/golangci-lint:v2.1.5-alpine as linter
#
#WORKDIR /src
#COPY .. .
#COPY ../.golangci.yaml ./
#RUN golangci-lint run

# Stage 4 - Build application
FROM base as builder

#COPY --from=linter /src/lint_report.json .
RUN go build ./cmd/bookinfo

# Stafe 5 - Build final image
FROM alpine:3

WORKDIR /app
COPY --from=builder /src/bookinfo /usr/bin/bookinfo
CMD ["bookinfo"]