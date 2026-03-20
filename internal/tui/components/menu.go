package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// MenuItem represents a selectable menu entry.
type MenuItem struct {
	Title       string
	Description string
	Icon        string
}

// Menu renders a vertical list of items with a cursor.
type Menu struct {
	Items  []MenuItem
	Cursor int
	Width  int
}

// NewMenu creates a new Menu from the given items.
func NewMenu(items []MenuItem) Menu {
	return Menu{Items: items, Width: 60}
}

// Up moves the cursor up.
func (m *Menu) Up() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

// Down moves the cursor down.
func (m *Menu) Down() {
	if m.Cursor < len(m.Items)-1 {
		m.Cursor++
	}
}

// Selected returns the currently selected item index.
func (m *Menu) Selected() int {
	return m.Cursor
}

// View renders the menu with cursor highlighting.
func (m Menu) View() string {
	var b strings.Builder

	rowWidth := max(m.Width-4, 20)

	for i, item := range m.Items {
		icon := item.Icon
		if icon == "" {
			icon = " "
		}

		if i == m.Cursor {
			// Selected: highlight bar with cursor
			line := fmt.Sprintf(" ▸ %s %s", icon, item.Title)
			row := lipgloss.NewStyle().
				Background(tui.ColorHighlight).
				Foreground(lipgloss.Color("#A78BFA")).
				Bold(true).
				Width(rowWidth).
				Padding(0, 1).
				Render(line)
			b.WriteString(row + "\n")

			desc := lipgloss.NewStyle().
				Foreground(tui.ColorSecondary).
				Width(rowWidth).
				PaddingLeft(6).
				Render(item.Description)
			b.WriteString(desc + "\n")
		} else {
			// Unselected: same width, no background
			line := fmt.Sprintf("   %s %s", icon, item.Title)
			row := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9CA3AF")).
				Width(rowWidth).
				Padding(0, 1).
				Render(line)
			b.WriteString(row + "\n")

			desc := lipgloss.NewStyle().
				Foreground(tui.ColorMuted).
				Width(rowWidth).
				PaddingLeft(6).
				Render(item.Description)
			b.WriteString(desc + "\n")
		}

		b.WriteString("\n")
	}

	return b.String()
}
