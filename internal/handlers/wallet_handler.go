package handlers

import (
	"net/http"

	"go-crypto-wallet/internal/services"
	"go-crypto-wallet/internal/utils"
	"go-crypto-wallet/internal/wallet"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service services.WalletService
}

func NewWalletHandler(service services.WalletService) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) CreateWallet(ctx *gin.Context) {
	var req wallet.CreateWalletRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Failure(ctx, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	result, err := h.service.CreateWallet(ctx.Request.Context(), req)
	if err != nil {
		utils.Failure(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"wallet": result})
}

func (h *WalletHandler) GetBalance(ctx *gin.Context) {
	address := ctx.Param("address")
	result, err := h.service.GetBalance(ctx.Request.Context(), address)
	if err != nil {
		utils.Failure(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"balance": result})
}

func (h *WalletHandler) SendTransaction(ctx *gin.Context) {
	var req wallet.SendTransactionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Failure(ctx, http.StatusBadRequest, "invalid request body: "+err.Error())
		return
	}

	result, err := h.service.SendTransaction(ctx.Request.Context(), req)
	if err != nil {
		utils.Failure(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"transaction": result})
}

func (h *WalletHandler) GetTransactions(ctx *gin.Context) {
	address := ctx.Param("address")
	result, err := h.service.GetTransactions(ctx.Request.Context(), address)
	if err != nil {
		utils.Failure(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, gin.H{"transactions": result})
}
