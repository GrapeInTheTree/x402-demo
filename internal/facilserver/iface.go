package facilserver

import (
	"context"

	x402 "github.com/coinbase/x402/go"
)

// Facilitator defines the interface for payment verification and settlement.
// This decouples the HTTP handlers from the concrete x402 SDK implementation,
// making the server testable with mock facilitators.
type Facilitator interface {
	Verify(ctx context.Context, payloadBytes []byte, requirementsBytes []byte) (*x402.VerifyResponse, error)
	Settle(ctx context.Context, payloadBytes []byte, requirementsBytes []byte) (*x402.SettleResponse, error)
	GetSupported() x402.SupportedResponse
}
