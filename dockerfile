FROM golang:1.23-alpine AS builder

WORKDIR /builder

ENV CGO_ENABLED=1
RUN apk add --no-cache \
    gcc \
    musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

RUN mkdir -p database

COPY --from=builder /builder/bin/main .
COPY --from=builder /builder/database ./database
COPY --from=builder /builder/.env* ./