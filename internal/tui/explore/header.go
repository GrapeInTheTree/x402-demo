package explore

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

// HeaderModel shows the PAYMENT-REQUIRED header structure with field exploration.
type HeaderModel struct {
	explorer components.FieldExplorer
	width    int
	height   int
}

// NewHeaderModel creates a new header explorer with PAYMENT-REQUIRED field definitions.
func NewHeaderModel(width, height int) *HeaderModel {
	fields := []components.Field{
		{Name: "scheme", Value: "exact", Description: "Payment scheme. 'exact' means exact amount payment. Only scheme supported in x402."},
		{Name: "network", Value: "eip155:84532", Description: "CAIP-2 network identifier. 'eip155' = EVM chain, '84532' = Base Sepolia Chain ID."},
		{Name: "maxAmountRequired", Value: "100000", Description: "Maximum payment amount (smallest unit). USDC has 6 decimals, so 100000 = 0.1 USDC = $0.10."},
		{Name: "resource", Value: "https://.../weather", Description: "Target resource URL. The API endpoint the client wants to access."},
		{Name: "payTo", Value: "0x1234...abcd", Description: "Address to receive USDC. Receive-only address — no private key needed."},
		{Name: "asset", Value: "0x036CbD...3dCF7e", Description: "ERC-20 token contract address for payment (Base Sepolia USDC)."},
		{Name: "extra.name", Value: "USDC", Description: "EIP-712 domain name field. Must exactly match the token contract's name() return value."},
		{Name: "extra.version", Value: "2", Description: "EIP-712 domain version field. USDC v2 (FiatTokenV2) contract."},
		{Name: "extra.assetTransferMethod", Value: "(optional)", Description: "If 'permit2', uses Permit2 method. Defaults to EIP-3009 when unset. Client SDK uses this to determine signing method."},
	}

	return &HeaderModel{
		explorer: components.NewFieldExplorer(fields),
		width:    width,
		height:   height,
	}
}

// Update handles key events for field navigation.
func (m *HeaderModel) Update(msg tea.Msg) tea.Cmd {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up", "k":
			m.explorer.Up()
		case "down", "j":
			m.explorer.Down()
		}
	}
	return nil
}

// SetSize updates the model dimensions.
func (m *HeaderModel) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.explorer.Width = width
}

// View renders the header field explorer with descriptions.
func (m *HeaderModel) View() string {
	title := lipgloss.NewStyle().
		Foreground(tui.ColorSecondary).
		Bold(true).
		Render("PAYMENT-REQUIRED Header Fields")

	subtitle := tui.MutedStyle.
		Render("↑/↓ to select fields — see description below")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		subtitle,
		"",
		m.explorer.View(),
	)
}
