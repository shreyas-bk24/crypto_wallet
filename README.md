# Crypto Wallet Backend

Production-style Go backend scaffold for a crypto wallet service.

## Setup

1. Copy `.env.example` to `.env` and set `ALCHEMY_RPC_URL`.
2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Run the server:

   ```bash
   go run cmd/main.go
   ```

## Docker

Build the image:

```bash
docker build -t crypto-wallet-backend .
```

Run the container:

```bash
docker run --rm -p 8080:8080 --env-file .env crypto-wallet-backend
```

## Endpoints

- `GET /health`
- `POST /wallet/create`
- `GET /wallet/balance/:address`
- `POST /wallet/send`
- `GET /wallet/transactions/:address`