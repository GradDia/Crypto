FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates build-base postgresql-client

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -trimpath -o /cryptoapp ./cmd/api

FROM alpine:3.19

RUN apk add --no-cache tzdata postgresql-client

COPY --from=builder /cryptoapp /app/cryptoapp
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY pkg/migrations/postgres/ /app/migrations/

WORKDIR /app
EXPOSE 8080