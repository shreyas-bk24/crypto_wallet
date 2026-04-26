package wallet

type CreateWalletRequest struct {
	Label string `json:"label"`
}

type SendTransactionRequest struct {
	PrivateKey string `json:"private_key" binding:"required"`
	ToAddress  string `json:"to_address" binding:"required"`
	AmountWei  string `json:"amount_wei" binding:"required"`
}

type WalletResponse struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey,omitempty"`
	Label      string `json:"label,omitempty"`
	Nonce      uint64 `json:"nonce,omitempty"`
}

type BalanceResponse struct {
	Address     string `json:"address"`
	BalanceWei  string `json:"balance_wei"`
	BalanceETH  string `json:"balance_eth"`
	NetworkName string `json:"network_name"`
}

type TransactionResponse struct {
	Hash      string `json:"hash"`
	From      string `json:"from"`
	To        string `json:"to"`
	AmountWei string `json:"amount_wei"`
	Status    string `json:"status"`
	Direction string `json:"direction"`
}

type TransactionsResponse struct {
	Address      string                `json:"address"`
	Transactions []TransactionResponse `json:"transactions"`
}
