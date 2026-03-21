package explore

import (
	"fmt"
	"strings"

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
	availW := m.width - 4 // RootModel padding
	gap := 1
	innerTotal := availW - gap - 4 // 2 borders
	leftW := innerTotal * 2 / 5
	rightW := innerTotal - leftW

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorBorder).
		Padding(0, 1)

	// Left: field list
	leftTitle := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorSecondary).
		Render("PAYMENT-REQUIRED Fields")
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
	leftContent := lipgloss.JoinVertical(lipgloss.Left, leftTitle, "", fieldList.String())
	leftBox := boxStyle.Width(leftW).Height(m.height - 4).Render(leftContent)

	// Right: description of selected field
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
