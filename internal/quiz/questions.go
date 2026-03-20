package quiz

// AllQuestions returns the full set of x402 protocol quiz questions.
func AllQuestions() []Question {
	return []Question{
		decodePaymentRequired(),
		usdcAmountConversion(),
		buildVerifyRequest(),
		eip712DomainFields(),
		parseSettlementResponse(),
	}
}

func decodePaymentRequired() Question {
	return Question{
		ID:         "decode-header",
		Title:      "Decode PAYMENT-REQUIRED Header",
		Difficulty: "easy",
		Description: `When a resource server returns HTTP 402, it includes a
PAYMENT-REQUIRED header containing base64-encoded JSON.

Write a function that decodes this header and extracts the
payTo address from the first "accepts" entry.`,
		Template: `package x402quiz

import (
	"encoding/base64"
	"encoding/json"
)

// DecodePayTo takes a base64-encoded PAYMENT-REQUIRED header value
// and returns the payTo address from the first accepts entry.
//
// The decoded JSON structure looks like:
//   {
//     "accepts": [
//       { "scheme": "exact", "payTo": "0x1234...", "network": "eip155:84532", ... }
//     ]
//   }
func DecodePayTo(headerValue string) (string, error) {
	// TODO: 1. Base64 decode the headerValue
	// TODO: 2. Unmarshal the JSON
	// TODO: 3. Return the payTo field from the first accepts entry

	_ = base64.StdEncoding.DecodeString // hint: use this
	_ = json.Unmarshal                  // hint: use this

	return "", nil // replace this
}
`,
		TestCode: `package x402quiz

import (
	"encoding/base64"
	"testing"
)

func TestDecodePayTo_Basic(t *testing.T) {
	raw := ` + "`" + `{"accepts":[{"scheme":"exact","payTo":"0xABCD1234","network":"eip155:84532"}]}` + "`" + `
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	got, err := DecodePayTo(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "0xABCD1234" {
		t.Errorf("expected payTo '0xABCD1234', got %q", got)
	}
}

func TestDecodePayTo_InvalidBase64(t *testing.T) {
	_, err := DecodePayTo("not-valid-base64!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestDecodePayTo_EmptyAccepts(t *testing.T) {
	raw := ` + "`" + `{"accepts":[]}` + "`" + `
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	_, err := DecodePayTo(encoded)
	if err == nil {
		t.Error("expected error for empty accepts")
	}
}
`,
		Hints: []string{
			"base64.StdEncoding.DecodeString() returns []byte and error",
			"Define a struct with Accepts []struct{ PayTo string `json:\"payTo\"` }",
			"Check len(decoded.Accepts) > 0 before accessing [0]",
		},
	}
}

func usdcAmountConversion() Question {
	return Question{
		ID:         "usdc-amount",
		Title:      "USDC Amount Conversion",
		Difficulty: "easy",
		Description: `USDC uses 6 decimal places. In x402, amounts are expressed in
the smallest unit (like wei for ETH).

Write functions to convert between dollar amounts and USDC
smallest units. For example: $0.10 = 100000 units.`,
		Template: `package x402quiz

// DollarsToUSDC converts a dollar amount (e.g., 0.10) to USDC
// smallest units (6 decimals). For example: 0.10 → 100000
func DollarsToUSDC(dollars float64) uint64 {
	// TODO: multiply by 10^6 and convert to uint64
	return 0
}

// USDCToDollars converts USDC smallest units back to dollars.
// For example: 100000 → 0.10
func USDCToDollars(units uint64) float64 {
	// TODO: divide by 10^6
	return 0.0
}

// FormatUSDC returns a human-readable string like "$0.10" for a
// given amount in USDC smallest units.
func FormatUSDC(units uint64) string {
	// TODO: convert to dollars and format as "$X.XX"
	return ""
}
`,
		TestCode: `package x402quiz

import "testing"

func TestDollarsToUSDC(t *testing.T) {
	tests := []struct{ dollars float64; want uint64 }{
		{0.10, 100000},
		{1.00, 1000000},
		{0.01, 10000},
		{100.0, 100000000},
	}
	for _, tt := range tests {
		got := DollarsToUSDC(tt.dollars)
		if got != tt.want {
			t.Errorf("DollarsToUSDC(%v) = %d, want %d", tt.dollars, got, tt.want)
		}
	}
}

func TestUSDCToDollars(t *testing.T) {
	tests := []struct{ units uint64; want float64 }{
		{100000, 0.10},
		{1000000, 1.00},
		{10000, 0.01},
	}
	for _, tt := range tests {
		got := USDCToDollars(tt.units)
		diff := got - tt.want
		if diff < -0.001 || diff > 0.001 {
			t.Errorf("USDCToDollars(%d) = %v, want %v", tt.units, got, tt.want)
		}
	}
}

func TestFormatUSDC(t *testing.T) {
	tests := []struct{ units uint64; want string }{
		{100000, "$0.10"},
		{1000000, "$1.00"},
		{1500000, "$1.50"},
	}
	for _, tt := range tests {
		got := FormatUSDC(tt.units)
		if got != tt.want {
			t.Errorf("FormatUSDC(%d) = %q, want %q", tt.units, got, tt.want)
		}
	}
}
`,
		Hints: []string{
			"USDC has 6 decimal places: 1 USDC = 1,000,000 units",
			"Use uint64(dollars * 1_000_000) for conversion",
			`Use fmt.Sprintf("$%.2f", dollars) for formatting`,
		},
	}
}

func buildVerifyRequest() Question {
	return Question{
		ID:         "build-verify",
		Title:      "Build Facilitator /verify Request",
		Difficulty: "medium",
		Description: `The facilitator's /verify endpoint expects a JSON body with:
- x402Version (int): always 2
- paymentPayload (JSON): the raw payment payload
- paymentRequirements (JSON): the selected requirements

Write a function that constructs this JSON body from the given inputs.`,
		Template: `package x402quiz

import "encoding/json"

// BuildVerifyBody constructs the JSON body for a facilitator /verify request.
// It must return a JSON object with these exact fields:
//   { "x402Version": 2, "paymentPayload": <payload>, "paymentRequirements": <requirements> }
//
// IMPORTANT: payload and requirements are raw JSON — they must be embedded
// directly, not re-encoded as strings.
func BuildVerifyBody(payload, requirements []byte) ([]byte, error) {
	// TODO: Build the verify request body
	// Hint: use json.RawMessage to embed raw JSON without re-encoding

	_ = json.RawMessage{} // hint
	_ = json.Marshal      // hint

	return nil, nil
}
`,
		TestCode: `package x402quiz

import (
	"encoding/json"
	"testing"
)

func TestBuildVerifyBody_Structure(t *testing.T) {
	payload := []byte(` + "`" + `{"from":"0xClient","signature":"0xABC"}` + "`" + `)
	requirements := []byte(` + "`" + `{"scheme":"exact","payTo":"0x1234"}` + "`" + `)

	body, err := BuildVerifyBody(payload, requirements)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	// Check x402Version
	var version int
	json.Unmarshal(result["x402Version"], &version)
	if version != 2 {
		t.Errorf("x402Version = %d, want 2", version)
	}

	// Check paymentPayload is embedded JSON, not a string
	if result["paymentPayload"][0] == '"' {
		t.Error("paymentPayload should be embedded JSON, not a string")
	}

	// Check paymentRequirements is embedded JSON
	if result["paymentRequirements"][0] == '"' {
		t.Error("paymentRequirements should be embedded JSON, not a string")
	}
}

func TestBuildVerifyBody_RoundTrip(t *testing.T) {
	payload := []byte(` + "`" + `{"value":"100000"}` + "`" + `)
	requirements := []byte(` + "`" + `{"network":"eip155:84532"}` + "`" + `)

	body, err := BuildVerifyBody(payload, requirements)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed struct {
		Version      int             ` + "`" + `json:"x402Version"` + "`" + `
		Payload      json.RawMessage ` + "`" + `json:"paymentPayload"` + "`" + `
		Requirements json.RawMessage ` + "`" + `json:"paymentRequirements"` + "`" + `
	}
	json.Unmarshal(body, &parsed)

	if parsed.Version != 2 {
		t.Errorf("version = %d, want 2", parsed.Version)
	}
	if string(parsed.Payload) != string(payload) {
		t.Errorf("payload mismatch: got %s", parsed.Payload)
	}
	if string(parsed.Requirements) != string(requirements) {
		t.Errorf("requirements mismatch: got %s", parsed.Requirements)
	}
}
`,
		Hints: []string{
			"json.RawMessage lets you embed raw JSON bytes without re-encoding",
			"Build a struct or map with json.RawMessage fields",
			"map[string]any{\"x402Version\": 2, \"paymentPayload\": json.RawMessage(payload)}",
		},
	}
}

func eip712DomainFields() Question {
	return Question{
		ID:         "eip712-domain",
		Title:      "EIP-712 Domain Separator",
		Difficulty: "medium",
		Description: `EIP-712 requires a domain separator to prevent signature replay
across different contracts and chains.

For USDC on Base Sepolia, fill in the correct domain fields.
The domain name must match the token contract's name() return value exactly.`,
		Template: `package x402quiz

// EIP712Domain represents the domain separator for EIP-712 typed data signing.
type EIP712Domain struct {
	Name              string
	Version           string
	ChainID           uint64
	VerifyingContract string
}

// USDCDomain returns the correct EIP-712 domain for USDC on Base Sepolia.
//
// Key facts:
// - Base Sepolia chain ID: 84532
// - USDC contract: 0x036CbD53842c5426634e7929541eC2318f3dCF7e
// - The domain "name" must match the token's name() return value EXACTLY
//   (Hint: Base Sepolia USDC returns "USDC", NOT "USD Coin")
func USDCDomain() EIP712Domain {
	return EIP712Domain{
		// TODO: Fill in the correct values
		Name:              "", // What does the USDC contract's name() return?
		Version:           "", // FiatTokenV2 version string
		ChainID:           0,  // Base Sepolia chain ID
		VerifyingContract: "", // USDC contract address
	}
}

// Permit2Domain returns the correct EIP-712 domain for Permit2 on Base Sepolia.
//
// Key facts:
// - Permit2 is deployed via CREATE2 at the same address on all chains
// - Permit2 address: 0x000000000022D473030F116dDEE9F6B43aC78BA3
func Permit2Domain() EIP712Domain {
	return EIP712Domain{
		// TODO: Fill in the correct values
		Name:              "",
		Version:           "",
		ChainID:           0,
		VerifyingContract: "",
	}
}
`,
		TestCode: `package x402quiz

import "testing"

func TestUSDCDomain(t *testing.T) {
	d := USDCDomain()

	if d.Name != "USDC" {
		t.Errorf("Name = %q, want \"USDC\" (not \"USD Coin\"!)", d.Name)
	}
	if d.Version != "2" {
		t.Errorf("Version = %q, want \"2\"", d.Version)
	}
	if d.ChainID != 84532 {
		t.Errorf("ChainID = %d, want 84532 (Base Sepolia)", d.ChainID)
	}
	if d.VerifyingContract != "0x036CbD53842c5426634e7929541eC2318f3dCF7e" {
		t.Errorf("VerifyingContract = %q, want USDC address", d.VerifyingContract)
	}
}

func TestPermit2Domain(t *testing.T) {
	d := Permit2Domain()

	if d.Name != "Permit2" {
		t.Errorf("Name = %q, want \"Permit2\"", d.Name)
	}
	if d.ChainID != 84532 {
		t.Errorf("ChainID = %d, want 84532", d.ChainID)
	}
	if d.VerifyingContract != "0x000000000022D473030F116dDEE9F6B43aC78BA3" {
		t.Errorf("VerifyingContract = %q, want Permit2 address", d.VerifyingContract)
	}
}
`,
		Hints: []string{
			"USDC domain name is \"USDC\" — the contract returns this from name()",
			"FiatTokenV2 version is \"2\"",
			"Base Sepolia chain ID is 84532",
		},
	}
}

func parseSettlementResponse() Question {
	return Question{
		ID:         "parse-settlement",
		Title:      "Parse PAYMENT-RESPONSE Header",
		Difficulty: "easy",
		Description: `After successful payment, the resource server returns a
PAYMENT-RESPONSE header (base64-encoded JSON) containing
the settlement result.

Parse it and return a structured result with the tx hash
and success status.`,
		Template: `package x402quiz

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Settlement holds the parsed payment settlement result.
type Settlement struct {
	Success     bool
	Transaction string
	Network     string
	Payer       string
}

// ParsePaymentResponse decodes a base64-encoded PAYMENT-RESPONSE header
// and returns the structured settlement data.
//
// Example decoded JSON:
//   { "success": true, "transaction": "0xabc...", "network": "eip155:84532", "payer": "0x..." }
func ParsePaymentResponse(headerValue string) (*Settlement, error) {
	// TODO: 1. Base64 decode
	// TODO: 2. JSON unmarshal into Settlement
	// TODO: 3. Validate that transaction is not empty when success is true

	_ = base64.StdEncoding.DecodeString
	_ = json.Unmarshal
	_ = fmt.Errorf

	return nil, nil
}
`,
		TestCode: `package x402quiz

import (
	"encoding/base64"
	"testing"
)

func TestParsePaymentResponse_Success(t *testing.T) {
	raw := ` + "`" + `{"success":true,"transaction":"0xABC123","network":"eip155:84532","payer":"0xDEF456"}` + "`" + `
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	s, err := ParsePaymentResponse(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.Success {
		t.Error("expected Success=true")
	}
	if s.Transaction != "0xABC123" {
		t.Errorf("Transaction = %q, want \"0xABC123\"", s.Transaction)
	}
	if s.Network != "eip155:84532" {
		t.Errorf("Network = %q, want \"eip155:84532\"", s.Network)
	}
}

func TestParsePaymentResponse_Failed(t *testing.T) {
	raw := ` + "`" + `{"success":false,"transaction":"","network":"eip155:84532","payer":""}` + "`" + `
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	s, err := ParsePaymentResponse(encoded)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Success {
		t.Error("expected Success=false")
	}
}

func TestParsePaymentResponse_InvalidBase64(t *testing.T) {
	_, err := ParsePaymentResponse("not-valid!!!")
	if err == nil {
		t.Error("expected error for invalid base64")
	}
}

func TestParsePaymentResponse_SuccessWithoutTx(t *testing.T) {
	raw := ` + "`" + `{"success":true,"transaction":"","network":"eip155:84532"}` + "`" + `
	encoded := base64.StdEncoding.EncodeToString([]byte(raw))

	_, err := ParsePaymentResponse(encoded)
	if err == nil {
		t.Error("expected error: success=true but no transaction hash")
	}
}
`,
		Hints: []string{
			"Same base64 → JSON pattern as DecodePayTo",
			"Use json struct tags: `json:\"success\"`",
			"Validate: if Success && Transaction == \"\" → return error",
		},
	}
}
