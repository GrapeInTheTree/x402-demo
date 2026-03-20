package demo

import "encoding/json"

// WalletInfo identifies a named wallet address.
type WalletInfo struct {
	Name    string
	Address string
}

// WalletBalance holds ETH and USDC balances for a wallet.
type WalletBalance struct {
	Wallet  WalletInfo
	ETH     string // formatted with 6 decimal places
	USDC    string // formatted with 6 decimal places
	ETHRaw  string // raw wei value
	USDCRaw string // raw smallest-unit value
}

// DecodedPaymentRequired is the decoded PAYMENT-REQUIRED header.
type DecodedPaymentRequired struct {
	Accepts     []json.RawMessage `json:"accepts"`
	Resource    string            `json:"resource,omitempty"`
	Description string            `json:"description,omitempty"`
	MimeType    string            `json:"mimeType,omitempty"`
}

// AcceptItem is a single entry in the accepts array.
type AcceptItem struct {
	Scheme             string                 `json:"scheme"`
	Network            string                 `json:"network"`
	MaxAmountRequired  string                 `json:"maxAmountRequired"`
	Resource           string                 `json:"resource"`
	Description        string                 `json:"description"`
	MimeType           string                 `json:"mimeType"`
	PayTo              string                 `json:"payTo"`
	Asset              string                 `json:"asset"`
	Extra              map[string]interface{} `json:"extra"`
}

// FlowState tracks the state of a payment flow execution.
type FlowState struct {
	TransferMethod  string
	Wallets         []WalletInfo
	BalancesBefore  []WalletBalance
	BalancesAfter   []WalletBalance
	PaymentRequired *DecodedPaymentRequired
	PaymentPayload  json.RawMessage // raw payload JSON
	VerifyResponse  json.RawMessage // raw verify response
	SettleResponse  json.RawMessage // raw settle response
	TxHash          string
	CurrentStep     int
	TotalSteps      int
	Error           error
}

// NewFlowState creates a new flow state with default values.
func NewFlowState(transferMethod string) *FlowState {
	return &FlowState{
		TransferMethod: transferMethod,
		TotalSteps:     10,
	}
}

// StepDescription returns the description for each step.
func StepDescription(step int) string {
	descriptions := map[int]string{
		1:  "Check wallet addresses & balances",
		2:  "Facilitator /supported (service discovery)",
		3:  "Client → Resource Server: API call without payment",
		4:  "Decode PAYMENT-REQUIRED header from 402 response",
		5:  "Client: create payment signature (off-chain)",
		6:  "Client → Resource Server: retry with PAYMENT-SIGNATURE",
		7:  "Resource Server → Facilitator /verify (off-chain verification)",
		8:  "Verification passed → return data + /settle request",
		9:  "Facilitator /settle → on-chain settlement + PAYMENT-RESPONSE",
		10: "Check balances after settlement",
	}
	return descriptions[step]
}
