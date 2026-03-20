package explore

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

// TypedDataModel shows EIP-712 TypedData structure with field exploration.
type TypedDataModel struct {
	explorer components.FieldExplorer
	permit2  bool // false=EIP-3009, true=Permit2
	width    int
	height   int
}

// NewTypedDataModel creates a new EIP-712 typed data explorer.
func NewTypedDataModel(width, height int) *TypedDataModel {
	return &TypedDataModel{
		explorer: components.NewFieldExplorer(eip3009Fields()),
		width:    width,
		height:   height,
	}
}

func eip3009Fields() []components.Field {
	return []components.Field{
		{Name: "[Domain] name", Value: "USDC", Description: "EIP-712 domain name. Return value of USDC contract's name(). On Base Sepolia it is 'USDC' (not 'USD Coin')."},
		{Name: "[Domain] version", Value: "2", Description: "EIP-712 domain version. Version string of the FiatTokenV2 contract."},
		{Name: "[Domain] chainId", Value: "84532", Description: "Base Sepolia Chain ID. Used for replay protection."},
		{Name: "[Domain] contract", Value: "0x036CbD...", Description: "Contract address that verifies the signature (USDC). Key element of domain separation."},
		{Name: "[Msg] from", Value: "0xClient...", Description: "Signer and USDC holder. Signature recovery verifies this address matches."},
		{Name: "[Msg] to", Value: "0xPayTo...", Description: "USDC recipient address. Must match payTo in PAYMENT-REQUIRED."},
		{Name: "[Msg] value", Value: "100000", Description: "Transfer amount (0.1 USDC). Must be >= maxAmountRequired to pass verification."},
		{Name: "[Msg] validAfter", Value: "0", Description: "Signature validity start time (Unix). 0 = valid immediately."},
		{Name: "[Msg] validBefore", Value: "1718000000", Description: "Signature expiration time (Unix). Transaction cannot execute after this time."},
		{Name: "[Msg] nonce", Value: "0xabcd...1234", Description: "Random 32 bytes. Prevents double-spend — once used, a nonce cannot be reused."},
	}
}

func permit2Fields() []components.Field {
	return []components.Field{
		{Name: "[Domain] name", Value: "Permit2", Description: "EIP-712 domain name for the Permit2 contract."},
		{Name: "[Domain] chainId", Value: "84532", Description: "Base Sepolia Chain ID."},
		{Name: "[Domain] contract", Value: "0x000...22D4", Description: "Permit2 contract address (CREATE2, same on all chains)."},
		{Name: "[Msg] permitted.token", Value: "0x036CbD...", Description: "Token address to transfer (USDC). Unlike EIP-3009, works with any ERC-20."},
		{Name: "[Msg] permitted.amount", Value: "100000", Description: "Permitted transfer amount (0.1 USDC)."},
		{Name: "[Msg] spender", Value: "0x4020...0001", Description: "x402Permit2Proxy address. Permit2 allows this address to move tokens."},
		{Name: "[Msg] nonce", Value: "12345", Description: "Permit2 nonce. Unlike EIP-3009's random nonce, this increments sequentially."},
		{Name: "[Msg] deadline", Value: "1718000000", Description: "Signature expiration time (Unix)."},
		{Name: "[Witness] to", Value: "0xPayTo...", Description: "Final token recipient. x402Permit2Proxy transfers tokens to this address."},
		{Name: "[Witness] validAfter", Value: "0", Description: "Validity start time. Same role as in EIP-3009."},
	}
}

// Update handles key events for field navigation and EIP-3009/Permit2 switching.
func (m *TypedDataModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			m.explorer.Up()
		case "down", "j":
			m.explorer.Down()
		case "tab":
			m.permit2 = !m.permit2
			if m.permit2 {
				m.explorer = components.NewFieldExplorer(permit2Fields())
			} else {
				m.explorer = components.NewFieldExplorer(eip3009Fields())
			}
		}
	}
	return nil
}

// SetSize updates the model dimensions.
func (m *TypedDataModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.explorer.Width = width
}

// View renders the EIP-712 typed data field explorer.
func (m *TypedDataModel) View() string {
	mode := "EIP-3009 (USDC Direct)"
	if m.permit2 {
		mode = "Permit2 (Universal ERC-20)"
	}

	title := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		Render("EIP-712 TypedData — " + mode)

	hint := tui.MutedStyle.
		Render("Tab to switch EIP-3009/Permit2  ↑/↓ select field")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		hint,
		"",
		m.explorer.View(),
	)
}
