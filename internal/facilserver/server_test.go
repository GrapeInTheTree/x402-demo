package facilserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	x402 "github.com/coinbase/x402/go"
	"github.com/gin-gonic/gin"
)

// mockFacilitator implements the Facilitator interface for testing.
type mockFacilitator struct {
	verifyFunc    func(ctx context.Context, payload, requirements []byte) (*x402.VerifyResponse, error)
	settleFunc    func(ctx context.Context, payload, requirements []byte) (*x402.SettleResponse, error)
	supportedResp x402.SupportedResponse
}

func (m *mockFacilitator) Verify(ctx context.Context, payload, requirements []byte) (*x402.VerifyResponse, error) {
	return m.verifyFunc(ctx, payload, requirements)
}

func (m *mockFacilitator) Settle(ctx context.Context, payload, requirements []byte) (*x402.SettleResponse, error) {
	return m.settleFunc(ctx, payload, requirements)
}

func (m *mockFacilitator) GetSupported() x402.SupportedResponse {
	return m.supportedResp
}

func init() {
	gin.SetMode(gin.TestMode)
}

func newTestServer(mock *mockFacilitator) (*Server, *gin.Engine) {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := New(mock, logger)
	r := gin.New()
	r.POST("/verify", srv.HandleVerify)
	r.POST("/settle", srv.HandleSettle)
	r.GET("/supported", srv.HandleSupported)
	return srv, r
}

func makeRequest(t *testing.T, r *gin.Engine, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		json.NewEncoder(&buf).Encode(body)
	}
	req := httptest.NewRequest(method, path, &buf)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ─── /verify tests ──────────────────────────────────────────

func TestHandleVerify_Success(t *testing.T) {
	mock := &mockFacilitator{
		verifyFunc: func(_ context.Context, _, _ []byte) (*x402.VerifyResponse, error) {
			return &x402.VerifyResponse{
				IsValid: true,
				Payer:   "0x1234567890abcdef1234567890abcdef12345678",
			}, nil
		},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "POST", "/verify", map[string]any{
		"x402Version":         2,
		"paymentPayload":      map[string]any{"test": true},
		"paymentRequirements": map[string]any{"scheme": "exact"},
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp x402.VerifyResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.IsValid {
		t.Error("expected isValid=true")
	}
	if resp.Payer == "" {
		t.Error("expected payer address")
	}
}

func TestHandleVerify_InvalidJSON(t *testing.T) {
	mock := &mockFacilitator{}
	_, r := newTestServer(mock)

	req := httptest.NewRequest("POST", "/verify", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleVerify_VerifyError(t *testing.T) {
	mock := &mockFacilitator{
		verifyFunc: func(_ context.Context, _, _ []byte) (*x402.VerifyResponse, error) {
			return nil, x402.NewVerifyError("insufficient_balance", "0xabc", "not enough USDC")
		},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "POST", "/verify", map[string]any{
		"x402Version":         2,
		"paymentPayload":      map[string]any{},
		"paymentRequirements": map[string]any{},
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["isValid"] != false {
		t.Error("expected isValid=false")
	}
	if resp["invalidReason"] != "insufficient_balance" {
		t.Errorf("expected insufficient_balance, got %v", resp["invalidReason"])
	}
}

func TestHandleVerify_InternalError(t *testing.T) {
	mock := &mockFacilitator{
		verifyFunc: func(_ context.Context, _, _ []byte) (*x402.VerifyResponse, error) {
			return nil, fmt.Errorf("rpc connection failed")
		},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "POST", "/verify", map[string]any{
		"x402Version":         2,
		"paymentPayload":      map[string]any{},
		"paymentRequirements": map[string]any{},
	})

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

// ─── /settle tests ──────────────────────────────────────────

func TestHandleSettle_Success(t *testing.T) {
	mock := &mockFacilitator{
		settleFunc: func(_ context.Context, _, _ []byte) (*x402.SettleResponse, error) {
			return &x402.SettleResponse{
				Success:     true,
				Transaction: "0xabcdef1234567890",
				Network:     "eip155:84532",
				Payer:       "0x1234567890abcdef1234567890abcdef12345678",
			}, nil
		},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "POST", "/settle", map[string]any{
		"x402Version":         2,
		"paymentPayload":      map[string]any{},
		"paymentRequirements": map[string]any{},
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp x402.SettleResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Transaction != "0xabcdef1234567890" {
		t.Errorf("unexpected tx: %s", resp.Transaction)
	}
}

func TestHandleSettle_SettleError(t *testing.T) {
	mock := &mockFacilitator{
		settleFunc: func(_ context.Context, _, _ []byte) (*x402.SettleResponse, error) {
			return nil, x402.NewSettleError("transaction_reverted", "0xabc", "eip155:84532", "", "out of gas")
		},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "POST", "/settle", map[string]any{
		"x402Version":         2,
		"paymentPayload":      map[string]any{},
		"paymentRequirements": map[string]any{},
	})

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["success"] != false {
		t.Error("expected success=false")
	}
	if resp["errorReason"] != "transaction_reverted" {
		t.Errorf("expected transaction_reverted, got %v", resp["errorReason"])
	}
}

func TestHandleSettle_InvalidJSON(t *testing.T) {
	mock := &mockFacilitator{}
	_, r := newTestServer(mock)

	req := httptest.NewRequest("POST", "/settle", bytes.NewBufferString("{invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

// ─── /supported tests ───────────────────────────────────────

func TestHandleSupported_WithKinds(t *testing.T) {
	mock := &mockFacilitator{
		supportedResp: x402.SupportedResponse{
			Kinds: []x402.SupportedKind{
				{X402Version: 2, Scheme: "exact", Network: "eip155:84532"},
			},
			Extensions: []string{},
			Signers:    map[string][]string{"eip155:*": {"0xfacilitator"}},
		},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "GET", "/supported", nil)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp x402.SupportedResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Kinds) != 1 {
		t.Fatalf("expected 1 kind, got %d", len(resp.Kinds))
	}
	if resp.Kinds[0].Scheme != "exact" {
		t.Errorf("expected scheme=exact, got %s", resp.Kinds[0].Scheme)
	}
}

func TestHandleSupported_Empty(t *testing.T) {
	mock := &mockFacilitator{
		supportedResp: x402.SupportedResponse{Kinds: []x402.SupportedKind{}},
	}
	_, r := newTestServer(mock)

	w := makeRequest(t, r, "GET", "/supported", nil)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for empty supported, got %d", w.Code)
	}
}
