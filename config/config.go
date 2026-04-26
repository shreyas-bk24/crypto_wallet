package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AlchemyRPCURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	rpcURL := strings.TrimSpace(os.Getenv("ALCHEMY_RPC_URL"))

	return &Config{AlchemyRPCURL: rpcURL}, nil
}
