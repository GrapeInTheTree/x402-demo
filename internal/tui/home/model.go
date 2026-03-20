package home

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
	"github.com/GrapeInTheTree/x402-playground/internal/tui/components"
)

// Model is the home page TUI model with a navigation menu.
type Model struct {
	menu   components.Menu
	width  int
	height int
}

var menuItems = []components.MenuItem{
	{Title: "Learn", Description: "Learn x402 protocol with coding quizzes", Icon: "◈"},
	{Title: "Explore", Description: "Inspect protocol data structures live", Icon: "◎"},
	{Title: "Practice", Description: "Execute payment flows (EIP-3009 / Permit2)", Icon: "▶"},
	{Title: "Dashboard", Description: "Wallet balances & transaction status", Icon: "◫"},
}

var pageMap = []tui.Page{
	tui.PageLearn,
	tui.PageExplore,
	tui.PagePractice,
	tui.PageDashboard,
}

// New creates a new home page model with the given dimensions.
func New(width, height int) *Model {
	m := &Model{
		menu:   components.NewMenu(menuItems),
		width:  width,
		height: height,
	}
	m.menu.Width = min(width, 60)
	return m
}

// Init implements the SubModel interface.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles key events for menu navigation.
func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.menu.Up()
		case "down", "j":
			m.menu.Down()
		case "enter":
			idx := m.menu.Selected()
			if idx >= 0 && idx < len(pageMap) {
				return m, func() tea.Msg {
					return tui.NavigateMsg{Page: pageMap[idx]}
				}
			}
		}
	}
	return m, nil
}

// SetSize updates the model dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.menu.Width = min(width, 60)
}

// View renders the home page — centered landing screen.
func (m *Model) View() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(tui.ColorPrimary).
		Render("x402 Protocol Explorer")

	subtitle := lipgloss.NewStyle().
		Foreground(tui.ColorMuted).
		Render("Interactive learning tool for the x402 payment protocol")

	// Center title block
	titleBlock := lipgloss.JoinVertical(lipgloss.Center, title, subtitle)

	body := lipgloss.JoinVertical(lipgloss.Center,
		titleBlock,
		"",
		m.menu.View(),
	)

	// Center everything in the available space
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		body)
}
