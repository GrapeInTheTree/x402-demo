package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
	ginmw "github.com/coinbase/x402/go/http/gin"
	evmserver "github.com/coinbase/x402/go/mechanisms/evm/exact/server"
	"github.com/gin-gonic/gin"

	"github.com/GrapeInTheTree/x402-demo/internal/config"
	"github.com/GrapeInTheTree/x402-demo/internal/server"
)

func main() {
	cfg, err := config.LoadResource()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config error: %v\n", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	network := x402.Network(cfg.Network)

	// Create facilitator client
	facilitatorClient := x402http.NewHTTPFacilitatorClient(&x402http.FacilitatorConfig{
		URL:     cfg.FacilitatorURL,
		Timeout: 60 * time.Second,
	})

	// Build payment-protected routes
	routes := server.BuildRoutes(cfg.PayToAddress, network)

	logger.Info("resource server initialized",
		"network", cfg.Network,
		"payTo", cfg.PayToAddress,
		"facilitatorURL", cfg.FacilitatorURL,
	)

	// Set up Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Health endpoint (no payment required)
	r.GET("/health", server.HealthHandler("resource", cfg.Network))

	// Register EVM exact scheme with fallback money parser for non-default networks.
	// SDK has built-in configs for Base Sepolia, Base Mainnet, etc.
	// This parser only activates for networks the SDK doesn't know about (e.g. Chiliz).
	evmScheme := evmserver.NewExactEvmScheme()
	evmScheme.RegisterMoneyParser(func(amount float64, net x402.Network) (*x402.AssetAmount, error) {
		// Return nil for SDK-supported networks — let the default parser handle them
		knownNetworks := map[string]bool{
			"eip155:8453":  true, // Base Mainnet
			"eip155:84532": true, // Base Sepolia
			"eip155:137":   true, // Polygon
			"eip155:42161": true, // Arbitrum
		}
		if knownNetworks[string(net)] {
			return nil, nil // delegate to SDK default
		}

		// Custom parser for unknown networks (Chiliz, etc.)
		if cfg.USDCAddress != "" {
			atomicAmount := int64(amount * 1_000_000) // USDC 6 decimals
			return &x402.AssetAmount{
				Asset:  cfg.USDCAddress,
				Amount: fmt.Sprintf("%d", atomicAmount),
				Extra: map[string]interface{}{
					"name":    "USDC",
					"version": "2",
				},
			}, nil
		}
		return nil, nil
	})

	// Apply x402 payment middleware
	r.Use(ginmw.X402Payment(ginmw.Config{
		Routes:      routes,
		Facilitator: facilitatorClient,
		Schemes: []ginmw.SchemeConfig{
			{Network: network, Server: evmScheme},
		},
		SyncFacilitatorOnStart: true,
		Timeout:                60 * time.Second,
		SettlementHandler: func(c *gin.Context, resp *x402.SettleResponse) {
			logger.Info("payment settled",
				"txHash", resp.Transaction,
				"network", resp.Network,
				"payer", resp.Payer,
			)
		},
	}))

	// Protected endpoints
	r.GET("/weather", server.WeatherHandler)
	r.GET("/joke", server.JokeHandler)
	r.GET("/premium-data", server.PremiumDataHandler)

	// Graceful shutdown
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("resource server starting", "port", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "error", err)
	}
}
