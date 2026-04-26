package main

import (
	"log"
	"net/http"

	"go-crypto-wallet/config"
	"go-crypto-wallet/internal/handlers"
	"go-crypto-wallet/internal/services"
	clientpkg "go-crypto-wallet/pkg/client"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	var ethClient *clientpkg.EthereumClient
	if cfg.AlchemyRPCURL != "" {
		ethClient, err = clientpkg.NewEthereumClient(cfg.AlchemyRPCURL)
		if err != nil {
			log.Fatalf("initialize ethereum client: %v", err)
		}
		defer ethClient.Close()
	} else {
		log.Println("ALCHEMY_RPC_URL is not set; starting without an active ethereum RPC client")
	}

	walletService := services.NewWalletService(ethClient)
	walletHandler := handlers.NewWalletHandler(walletService)

	router := setupRouter(walletHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Printf("server listening on %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}

func setupRouter(walletHandler *handlers.WalletHandler) *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://crypto-wallet-frontend-peach.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data": gin.H{
				"service": "crypto-wallet-backend",
				"status":  "healthy",
			},
			"error": "",
		})
	})

	walletGroup := router.Group("/wallet")
	{
		walletGroup.POST("/create", walletHandler.CreateWallet)
		walletGroup.GET("/balance/:address", walletHandler.GetBalance)
		walletGroup.POST("/send", walletHandler.SendTransaction)
		walletGroup.GET("/transactions/:address", walletHandler.GetTransactions)
	}

	return router
}
