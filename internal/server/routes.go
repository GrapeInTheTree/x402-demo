package server

import (
	x402 "github.com/coinbase/x402/go"
	x402http "github.com/coinbase/x402/go/http"
)

// BuildRoutes creates payment-protected route configurations.
func BuildRoutes(payTo string, network x402.Network) x402http.RoutesConfig {
	return x402http.RoutesConfig{
		"GET /weather": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					PayTo:   payTo,
					Price:   x402.Price("$0.1"),
					Network: network,
				},
			},
			Description: "Current weather data",
			MimeType:    "application/json",
		},
		"GET /joke": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:  "exact",
					PayTo:   payTo,
					Price:   x402.Price("$0.1"),
					Network: network,
				},
			},
			Description: "Random programming joke",
			MimeType:    "application/json",
		},
		"GET /premium-data": {
			Accepts: x402http.PaymentOptions{
				{
					Scheme:            "exact",
					PayTo:             payTo,
					Price:             x402.Price("$0.1"),
					Network:           network,
					MaxTimeoutSeconds: 120,
				},
			},
			Description: "Premium analytics data",
			MimeType:    "application/json",
		},
	}
}
