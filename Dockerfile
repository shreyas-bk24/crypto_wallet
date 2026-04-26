FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /out/crypto-wallet ./cmd/main.go

FROM alpine:3.20 AS runtime

RUN apk add --no-cache ca-certificates && addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /out/crypto-wallet /usr/local/bin/crypto-wallet

EXPOSE 8080

USER app

ENTRYPOINT ["/usr/local/bin/crypto-wallet"]