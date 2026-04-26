package services

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"

	"go-crypto-wallet/internal/wallet"
	clientpkg "go-crypto-wallet/pkg/client"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

type WalletService interface {
	CreateWallet(ctx context.Context, req wallet.CreateWalletRequest) (wallet.WalletResponse, error)
	GetBalance(ctx context.Context, address string) (wallet.BalanceResponse, error)
	SendTransaction(ctx context.Context, req wallet.SendTransactionRequest) (map[string]any, error)
	GetTransactions(ctx context.Context, address string) (wallet.TransactionsResponse, error)
}

type walletService struct {
	ethClient *clientpkg.EthereumClient
}

func NewWalletService(ethClient *clientpkg.EthereumClient) WalletService {
	return &walletService{ethClient: ethClient}
}

func (s *walletService) CreateWallet(_ context.Context, req wallet.CreateWalletRequest) (wallet.WalletResponse, error) {
	privateKey, err := ethcrypto.GenerateKey()
	if err != nil {
		return wallet.WalletResponse{}, err
	}

	privateKeyBytes := ethcrypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)

	if !ok {
		return wallet.WalletResponse{}, errors.New("invalid public key")
	}

	address := ethcrypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return wallet.WalletResponse{
		Address:    address,
		PrivateKey: privateKeyHex,
		Label:      req.Label,
		Nonce:      0,
	}, nil
}

func (s *walletService) GetBalance(ctx context.Context, address string) (wallet.BalanceResponse, error) {
	client := s.ethClient

	account := common.HexToAddress(address)

	balence, err := client.Client().BalanceAt(ctx, account, nil)
	if err != nil {
		return wallet.BalanceResponse{}, err
	}

	// convert wei -> ETH
	fbalence := new(big.Float)
	fbalence.SetString(balence.String())

	ethValue := new(big.Float).Quo(fbalence, big.NewFloat(math.Pow10(18)))

	return wallet.BalanceResponse{
		Address:     common.HexToAddress(address).Hex(),
		BalanceWei:  balence.String(),
		BalanceETH:  ethValue.String(),
		NetworkName: "ethereum-sepolia",
	}, nil
}

func (s *walletService) SendTransaction(ctx context.Context, req wallet.SendTransactionRequest) (map[string]any, error) {

	privateKey, err := ethcrypto.HexToECDSA(req.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	// derive sender
	fromAddress := ethcrypto.PubkeyToAddress(privateKey.PublicKey)

	// nounce
	nounce, err := s.ethClient.Client().PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return nil, err
	}

	// amount (already in wei)
	value := new(big.Int)
	value.SetString(req.AmountWei, 10)

	// gas price
	gasPrice, err := s.ethClient.Client().SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}

	// transaction
	to := common.HexToAddress(req.ToAddress)
	tx := types.NewTransaction(nounce, to, value, 21000, gasPrice, nil)

	// chain id
	chainID, err := s.ethClient.Client().NetworkID(ctx)
	if err != nil {
		return nil, err
	}

	//  sign
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	// send
	err = s.ethClient.Client().SendTransaction(ctx, signedTx)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"transaction_hash": signedTx.Hash().Hex(),
		"from":             fromAddress.Hex(),
		"to":               common.HexToAddress(req.ToAddress).Hex(),
		"amount_wei":       req.AmountWei,
		"nounce":           nounce,
		"status":           "queued",
	}, nil
}

func (s *walletService) GetTransactions(ctx context.Context, address string) (wallet.TransactionsResponse, error) {

	normalizedAddress := common.HexToAddress(address).Hex()

	url := os.Getenv("ALCHEMY_RPC_URL")
	fetch := func(filter map[string]any) ([]wallet.TransactionResponse, error) {
		payload := map[string]any{
			"jsonrpc": "2.0",
			"id":      1,
			"method":  "alchemy_getAssetTransfers",
			"params":  []any{filter},
		}
		body, _ := json.Marshal(payload)
		req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var raw map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return nil, err
		}
		resultRaw, ok := raw["result"].(map[string]any)
		if !ok {
			return nil, nil
		}
		transfersRaw, ok := resultRaw["transfers"].([]any)
		if !ok {
			return nil, nil
		}
		var txs []wallet.TransactionResponse
		for _, t := range transfersRaw {
			tx, ok := t.(map[string]any)
			if !ok {
				continue
			}
			hash, _ := tx["hash"].(string)
			from, _ := tx["from"].(string)
			to, _ := tx["to"].(string)
			// use rawContract.value (wei, hex)
			rawContract, _ := tx["rawContract"].(map[string]any)
			hexVal, _ := rawContract["value"].(string)
			valueWei := new(big.Int)
			if len(hexVal) > 2 {
				valueWei.SetString(hexVal[2:], 16)
			}
			// direction
			direction := "outgoing"
			if strings.EqualFold(to, normalizedAddress) {
				direction = "incoming"
			}
			txs = append(txs, wallet.TransactionResponse{
				Hash:      hash,
				From:      from,
				To:        to,
				AmountWei: valueWei.String(),
				Status:    "confirmed",
				Direction: direction, // add this field in struct
			})
		}
		return txs, nil
	}
	// outgoing
	outgoing, _ := fetch(map[string]any{
		"fromBlock":   "0x0",
		"toBlock":     "latest",
		"fromAddress": normalizedAddress,
		"category":    []string{"external"},
	})
	// incoming
	incoming, _ := fetch(map[string]any{
		"fromBlock": "0x0",
		"toBlock":   "latest",
		"toAddress": normalizedAddress,
		"category":  []string{"external"},
	})
	allTxs := append(outgoing, incoming...)
	return wallet.TransactionsResponse{
		Address:      normalizedAddress,
		Transactions: allTxs,
	}, nil
}
