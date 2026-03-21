package explore

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// CompareModel shows EIP-3009 vs Permit2 side-by-side.
type CompareModel struct {
	width  int
	height int
}

// NewCompareModel creates a new side-by-side comparison model.
func NewCompareModel(width, height int) *CompareModel {
	return &CompareModel{width: width, height: height}
}

// Update is a no-op since the compare view is static.
func (m *CompareModel) Update(msg tea.Msg) tea.Cmd {
	return nil
}

// SetSize updates the model dimensions.
func (m *CompareModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

const minCompareColWidth = 30

// View renders the EIP-3009 vs Permit2 side-by-side comparison.
func (m *CompareModel) View() string {
	availW := m.width - 4 // RootModel padding
	gap := 2
	// (leftW + 2) + gap + (rightW + 2) <= availW
	innerTotal := availW - gap - 4
	colWidth := max(innerTotal/2, minCompareColWidth)

	leftStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorSecondary).
		Width(colWidth).
		Padding(0, 1)

	rightStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(tui.ColorAccent).
		Width(colWidth).
		Padding(0, 1)

	left := buildComparePanel("EIP-3009", "transferWithAuthorization", tui.ColorSecondary, eip3009Rows)
	right := buildComparePanel("Permit2", "permitWitnessTransferFrom", tui.ColorAccent, permit2Rows)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftStyle.Render(left),
		"  ",
		rightStyle.Render(right),
	)
}

type compareRow struct{ Key, Val string }

var eip3009Rows = []compareRow{
	{"Token Support", "EIP-3009 tokens only (USDC, EURC)"},
	{"Prerequisites", "None"},
	{"Domain", "USDC contract"},
	{"Primary Type", "TransferWithAuthorization"},
	{"Contracts", "Direct token contract call"},
	{"Nonce", "Random 32-byte (single use)"},
	{"On-chain Call", "USDC.transferWithAuthorization(...)"},
	{"Gas", "Sponsored by Facilitator"},
}

var permit2Rows = []compareRow{
	{"Token Support", "Any ERC-20 token"},
	{"Prerequisites", "One-time approve(Permit2, amount)"},
	{"Domain", "Permit2 contract"},
	{"Primary Type", "PermitWitnessTransferFrom"},
	{"Contracts", "Permit2 + x402Permit2Proxy"},
	{"Nonce", "Sequential Permit2 nonce"},
	{"On-chain Call", "Proxy.settle(...) → Permit2"},
	{"Gas", "Sponsored by Facilitator"},
}

func buildComparePanel(title, subtitle string, color lipgloss.Color, rows []compareRow) string {
	var b strings.Builder

	b.WriteString(lipgloss.NewStyle().Foreground(color).Bold(true).Render(title) + "\n")
	b.WriteString(tui.MutedStyle.Render(subtitle) + "\n\n")

	keyStyle := lipgloss.NewStyle().Foreground(color).Bold(true)
	for _, r := range rows {
		b.WriteString(keyStyle.Render(r.Key+":") + "\n")
		b.WriteString("  " + r.Val + "\n\n")
	}

	return b.String()
}
