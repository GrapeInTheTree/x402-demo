package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	ColorPrimary   = lipgloss.Color("#7C3AED") // violet
	ColorSecondary = lipgloss.Color("#06B6D4") // cyan
	ColorAccent    = lipgloss.Color("#F59E0B") // amber
	ColorSuccess   = lipgloss.Color("#10B981") // green
	ColorError     = lipgloss.Color("#EF4444") // red
	ColorMuted     = lipgloss.Color("#6B7280") // gray
	ColorBorder    = lipgloss.Color("#374151") // dark gray
	ColorBg        = lipgloss.Color("#111827") // near-black

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(ColorSecondary).
			Bold(true)

	MutedStyle = lipgloss.NewStyle().
			Foreground(ColorMuted)

	AccentStyle = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(ColorSuccess)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(ColorError)

	// Layout
	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(1, 2)

	ActiveBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(1, 2)

	// Menu
	MenuItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	SelectedMenuItemStyle = lipgloss.NewStyle().
				PaddingLeft(1).
				Foreground(ColorPrimary).
				Bold(true).
				SetString("> ")

	// Status bar
	StatusBarStyle = lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1)
)

// LayoutPage renders body inside a bordered container, centered in the terminal,
// with a hints bar pinned at the bottom outside the border.
func LayoutPage(body, hints string, width, height int) string {
	// Hints bar at bottom (outside border)
	hintsRendered := MutedStyle.Render("  " + hints)
	hintsH := lipgloss.Height(hintsRendered) + 1 // +1 for gap

	// Border container dimensions (with margin from terminal edges)
	const hMargin = 2
	boxW := width - hMargin*2
	boxH := height - hintsH - 2 // 2 for top/bottom margin
	if boxW < 40 {
		boxW = width
	}
	if boxH < 10 {
		boxH = height - hintsH
	}

	// Inner content area (border takes 2 chars each side + padding)
	innerW := boxW - 6 // 2 border + 4 padding (2 each side)
	innerH := boxH - 4 // 2 border + 2 padding (1 top + 1 bottom)
	if innerW < 20 {
		innerW = boxW - 2
	}
	if innerH < 5 {
		innerH = boxH - 2
	}

	// Center body inside the inner area
	innerContent := lipgloss.Place(innerW, innerH,
		lipgloss.Center, lipgloss.Center,
		body)

	// Draw the bordered box
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorBorder).
		Padding(1, 2).
		Width(boxW).
		Height(boxH).
		Render(innerContent)

	// Center the box in the terminal width
	centeredBox := lipgloss.PlaceHorizontal(width, lipgloss.Center, box)

	return centeredBox + "\n" + hintsRendered
}
