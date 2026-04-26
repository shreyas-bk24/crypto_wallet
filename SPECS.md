# 📄 Go Crypto Wallet — Technical Specification

---

## 🧠 Objective

Design and implement a lightweight crypto wallet backend in Go that enables:

- Wallet generation (keypair creation)
- Balance retrieval from Ethereum
- Transaction creation, signing, and broadcasting
- Transaction history retrieval via external APIs

The system is intended for demonstration and learning purposes, focusing on blockchain interaction and backend design.

---

## 🏗️ System Overview

The system acts as a stateless backend service that interfaces with the Ethereum blockchain through an RPC provider.

### High-Level Flow

1. Client sends request to API  
2. Backend processes request using Go  
3. Blockchain interaction via go-ethereum  
4. RPC provider relays request to Ethereum network  
5. Response returned to client  

---

## 🧱 Architecture

[ Client ]      ↓ [ Go REST API ]      ↓ [ go-ethereum (ethclient, crypto) ]      ↓ [ RPC Provider (Alchemy) ]      ↓ [ Ethereum Sepolia Testnet ]

---

## ⚙️ Functional Requirements

### 1. Wallet Generation

- Generate ECDSA private key using secp256k1
- Derive public key
- Compute Ethereum address
- Return address and private key (hex encoded)

---

### 2. Balance Retrieval

- Accept wallet address
- Query Ethereum node via RPC
- Fetch balance in wei
- Convert wei → ETH
- Return formatted response

---

### 3. Transaction Sending

- Accept:
  - Private key
  - Recipient address
  - Amount (ETH)
- Derive sender address
- Fetch account nonce
- Estimate gas price
- Construct transaction
- Sign transaction using private key
- Broadcast to network
- Return transaction hash

---

### 4. Transaction History

- Accept wallet address
- Fetch transaction data via RPC provider API
- Filter incoming/outgoing transactions
- Return structured JSON response

---

## 🔐 Cryptography Details

- Algorithm: ECDSA
- Curve: secp256k1
- Address derivation:
  - Keccak-256 hash of public key
  - Last 20 bytes used as address

---

## 📡 External Dependencies

- Ethereum RPC Provider (Alchemy)
- Ethereum test network (Sepolia)

---

## 📊 Non-Functional Requirements

### Performance
- Low-latency API responses (< 300ms excluding blockchain latency)

### Scalability
- Stateless design allows horizontal scaling

### Reliability
- Dependent on RPC provider uptime

### Security (Demo Scope)
- No persistent storage of private keys
- No authentication implemented

---

## ⚠️ Constraints

- Uses public RPC provider (rate limits may apply)
- No transaction indexing (relies on external API)
- Not suitable for production use

---

## ❗ Assumptions

- Users understand private key responsibility
- Transactions occur on testnet only
- RPC provider is available and responsive

---

## 🔄 Data Flow Example (Send Transaction)

1. Client sends /wallet/send request  
2. Backend parses private key  
3. Connects to RPC provider  
4. Fetches nonce + gas price  
5. Constructs transaction  
6. Signs using ECDSA  
7. Broadcasts transaction  
8. Returns transaction hash  

---

## 📁 Project Structure (Suggested)

/cmd   main.go /internal   wallet/   handlers/   services/   utils/ /pkg   client/ /config

---

## 🔮 Future Enhancements

- Secure key vault integration
- Multi-chain support (Polygon, BSC)
- WebSocket support for real-time updates
- Transaction indexing service
- Role-based access control

---

## 📌 Summary

This system demonstrates core blockchain wallet operations using Go, focusing on:

- Cryptographic key handling  
- Ethereum RPC interaction  
- Transaction lifecycle management  
- Clean and minimal backend architecture  

--