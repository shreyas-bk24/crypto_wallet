package client

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
)

type EthereumClient struct {
	client *ethclient.Client
	rpcURL string
}

func NewEthereumClient(rpcURL string) (*EthereumClient, error) {
	rpcURL = strings.TrimSpace(rpcURL)
	if rpcURL == "" {
		return nil, fmt.Errorf("rpc url is required")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("dial ethereum rpc: %w", err)
	}

	return &EthereumClient{client: client, rpcURL: rpcURL}, nil
}

func (c *EthereumClient) Client() *ethclient.Client {
	if c == nil {
		return nil
	}
	return c.client
}

func (c *EthereumClient) Close() {
	if c == nil || c.client == nil {
		return
	}
	c.client.Close()
}
