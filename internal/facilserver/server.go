package facilserver

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	x402 "github.com/coinbase/x402/go"
	"github.com/gin-gonic/gin"
)

// Server wraps an x402 facilitator with HTTP handlers.
type Server struct {
	facilitator *x402.X402Facilitator
	logger      *slog.Logger
}

// New creates a new facilitator HTTP server.
func New(facilitator *x402.X402Facilitator, logger *slog.Logger) *Server {
	return &Server{
		facilitator: facilitator,
		logger:      logger,
	}
}

type verifySettleRequest struct {
	X402Version         int             `json:"x402Version"`
	PaymentPayload      json.RawMessage `json:"paymentPayload"`
	PaymentRequirements json.RawMessage `json:"paymentRequirements"`
}

// HandleVerify handles POST /verify.
func (s *Server) HandleVerify(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var req verifySettleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	s.logger.Info("verify request", "x402Version", req.X402Version)

	resp, err := s.facilitator.Verify(c.Request.Context(), req.PaymentPayload, req.PaymentRequirements)
	if err != nil {
		s.logger.Error("verify failed", "error", err)
		// Return 400 with error details (matching SDK convention)
		verifyErr, ok := err.(*x402.VerifyError)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"isValid":        false,
				"invalidReason":  verifyErr.InvalidReason,
				"invalidMessage": verifyErr.InvalidMessage,
				"payer":          verifyErr.Payer,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleSettle handles POST /settle.
func (s *Server) HandleSettle(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body"})
		return
	}

	var req verifySettleRequest
	if err := json.Unmarshal(body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	s.logger.Info("settle request", "x402Version", req.X402Version)

	resp, err := s.facilitator.Settle(c.Request.Context(), req.PaymentPayload, req.PaymentRequirements)
	if err != nil {
		s.logger.Error("settle failed", "error", err)
		settleErr, ok := err.(*x402.SettleError)
		if ok {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":      false,
				"errorReason":  settleErr.ErrorReason,
				"errorMessage": settleErr.ErrorMessage,
				"transaction":  settleErr.Transaction,
				"network":      settleErr.Network,
				"payer":        settleErr.Payer,
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// HandleSupported handles GET /supported.
func (s *Server) HandleSupported(c *gin.Context) {
	resp := s.facilitator.GetSupported()
	c.JSON(http.StatusOK, resp)
}
