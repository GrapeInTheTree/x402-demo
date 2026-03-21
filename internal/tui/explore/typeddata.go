package explore

import (
	"fmt"
	"strings"

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

	availW := m.width - 4
	gap := 1
	innerTotal := availW - gap - 4
	leftW := innerTotal * 2 / 5
	rightW := innerTotal - leftW

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorBorder).
		Padding(0, 1)

	// Left: field list
	leftTitle := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
		Render("EIP-712 — " + mode)
	tabHint := tui.MutedStyle.Render("Tab to switch")
	var fieldList strings.Builder
	nameW := max(leftW/2, 12)
	for i, f := range m.explorer.Fields {
		nameStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary).Width(nameW)
		valStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB"))
		if i == m.explorer.Cursor {
			nameStyle = lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true).Width(nameW)
			valStyle = lipgloss.NewStyle().Foreground(tui.ColorAccent)
			fmt.Fprintf(&fieldList, " > %s %s\n", nameStyle.Render(f.Name+":"), valStyle.Render(f.Value))
		} else {
			fmt.Fprintf(&fieldList, "   %s %s\n", nameStyle.Render(f.Name+":"), valStyle.Render(f.Value))
		}
	}
	leftContent := lipgloss.JoinVertical(lipgloss.Left, leftTitle, tabHint, "", fieldList.String())
	leftBox := boxStyle.Width(leftW).Height(m.height - 4).Render(leftContent)

	// Right: description
	rightTitle := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render("Field Details")
	desc := ""
	if m.explorer.Cursor >= 0 && m.explorer.Cursor < len(m.explorer.Fields) {
		f := m.explorer.Fields[m.explorer.Cursor]
		name := lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).Render(f.Name)
		val := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(f.Value)
		body := lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(rightW - 2).Render(f.Description)
		desc = lipgloss.JoinVertical(lipgloss.Left, name, val, "", body)
	}
	rightContent := lipgloss.JoinVertical(lipgloss.Left, rightTitle, "", desc)
	rightBox := boxStyle.Width(rightW).Height(m.height - 4).Render(rightContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftBox, " ", rightBox)
}
