package learn

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/quiz"
	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

type page int

const (
	pageQuestionList page = iota
	pageQuestion
	pageResult
)

type editorFinishedMsg struct{ err error }

type testResultMsg struct{ result *quiz.Result }

// Model is the quiz-based learning page model.
type Model struct {
	questions []quiz.Question
	cursor    int
	current   page
	runner    *quiz.Runner
	results   map[int]*quiz.Result
	score     quiz.Score
	width     int
	height    int
}

// New creates a new quiz learning model.
func New(width, height int) *Model {
	questions := quiz.AllQuestions()
	return &Model{
		questions: questions,
		results:   make(map[int]*quiz.Result),
		score:     quiz.Score{Questions: len(questions)},
		width:     width,
		height:    height,
	}
}

// Init implements the SubModel interface.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles navigation, editor launch, and test results.
func (m *Model) Update(msg tea.Msg) (tui.SubModel, tea.Cmd) {
	switch msg := msg.(type) {
	case editorFinishedMsg:
		if msg.err != nil {
			return m, nil
		}
		return m, m.runTests()

	case testResultMsg:
		m.results[m.cursor] = msg.result
		m.score.Answered++
		if msg.result.Passed == msg.result.Total && msg.result.Total > 0 {
			m.score.Correct++
		}
		m.current = pageResult
		return m, nil

	case tea.KeyMsg:
		switch m.current {
		case pageQuestionList:
			return m.updateList(msg)
		case pageQuestion:
			return m.updateQuestion(msg)
		case pageResult:
			return m.updateResult(msg)
		}
	}
	return m, nil
}

func (m *Model) updateList(msg tea.KeyMsg) (tui.SubModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.questions)-1 {
			m.cursor++
		}
	case "enter":
		m.current = pageQuestion
	case "esc", "q":
		return m, func() tea.Msg { return tui.BackMsg{} }
	}
	return m, nil
}

func (m *Model) updateQuestion(msg tea.KeyMsg) (tui.SubModel, tea.Cmd) {
	switch msg.String() {
	case "enter", "e":
		return m, m.openEditor()
	case "esc":
		m.current = pageQuestionList
	}
	return m, nil
}

func (m *Model) updateResult(msg tea.KeyMsg) (tui.SubModel, tea.Cmd) {
	switch msg.String() {
	case "enter", "r":
		return m, m.openEditor()
	case "n":
		if m.cursor < len(m.questions)-1 {
			m.cursor++
			m.current = pageQuestion
		} else {
			m.current = pageQuestionList
		}
	case "esc":
		m.current = pageQuestionList
	}
	return m, nil
}

func (m *Model) openEditor() tea.Cmd {
	if m.runner == nil {
		r, err := quiz.NewRunner()
		if err != nil {
			return nil
		}
		m.runner = r
	}

	q := m.questions[m.cursor]
	if err := os.WriteFile(m.runner.TemplatePath(), []byte(q.Template), 0644); err != nil {
		return nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		} else {
			editor = "nano"
		}
	}

	c := exec.Command(editor, m.runner.TemplatePath())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err: err}
	})
}

func (m *Model) runTests() tea.Cmd {
	runner := m.runner
	q := m.questions[m.cursor]
	return func() tea.Msg {
		solution, err := os.ReadFile(runner.TemplatePath())
		if err != nil {
			return testResultMsg{result: &quiz.Result{Error: err.Error()}}
		}
		result := runner.Run(string(solution), q.TestCode)
		return testResultMsg{result: result}
	}
}

// SetSize updates the model dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// View renders the current quiz page.
func (m *Model) View() string {
	switch m.current {
	case pageQuestionList:
		return m.viewList()
	case pageQuestion:
		return m.viewQuestion()
	case pageResult:
		return m.viewResult()
	default:
		return m.viewList()
	}
}

func (m *Model) viewList() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render("Learn — x402 Protocol Quiz")
	scoreText := lipgloss.NewStyle().Foreground(tui.ColorMuted).
		Render(fmt.Sprintf("Score: %d/%d", m.score.Correct, m.score.Questions))
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", min(m.width-8, 60)))

	var items strings.Builder
	for i, q := range m.questions {
		icon := "○"
		iconStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted)
		if r, ok := m.results[i]; ok {
			if r.Passed == r.Total && r.Total > 0 {
				icon = "✓"
				iconStyle = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
			} else {
				icon = "✗"
				iconStyle = lipgloss.NewStyle().Foreground(tui.ColorError)
			}
		}

		diffStyle := lipgloss.NewStyle().Foreground(tui.ColorMuted)
		switch q.Difficulty {
		case "easy":
			diffStyle = lipgloss.NewStyle().Foreground(tui.ColorSuccess)
		case "medium":
			diffStyle = lipgloss.NewStyle().Foreground(tui.ColorAccent)
		case "hard":
			diffStyle = lipgloss.NewStyle().Foreground(tui.ColorError)
		}

		if i == m.cursor {
			cursor := lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true).Render(">")
			name := lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true).Render(q.Title)
			fmt.Fprintf(&items, " %s %s %s  %s\n", cursor, iconStyle.Render(icon), name, diffStyle.Render("["+q.Difficulty+"]"))
		} else {
			name := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(q.Title)
			fmt.Fprintf(&items, "   %s %s  %s\n", iconStyle.Render(icon), name, diffStyle.Render("["+q.Difficulty+"]"))
		}
	}

	body := lipgloss.JoinVertical(lipgloss.Left, title, scoreText, divider, "", items.String())
	return tui.LayoutPage(body, "↑/↓ navigate  enter select  ? help  esc back", m.width, m.height)
}

func (m *Model) viewQuestion() string {
	q := m.questions[m.cursor]

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(fmt.Sprintf("Question %d/%d: %s", m.cursor+1, len(m.questions), q.Title))
	subtitle := lipgloss.NewStyle().Foreground(tui.ColorAccent).Render("[" + q.Difficulty + "]")
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", min(m.width-8, 60)))

	desc := lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).
		Width(min(m.width-10, 70)).Render(q.Description)

	prompt := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorAccent).
		Render("Press Enter to open $EDITOR...")

	var hintsView string
	if len(q.Hints) > 0 {
		var hb strings.Builder
		hb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).Render("Hints:") + "\n")
		for i, h := range q.Hints {
			hb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).
				Render(fmt.Sprintf("  %d. %s", i+1, h)) + "\n")
		}
		hintsView = hb.String()
	}

	body := lipgloss.JoinVertical(lipgloss.Left, title, subtitle, divider, "", desc, "", prompt, "", hintsView)
	return tui.LayoutPage(body, "enter open editor  esc back", m.width, m.height)
}

func (m *Model) viewResult() string {
	q := m.questions[m.cursor]
	r := m.results[m.cursor]

	title := lipgloss.NewStyle().Bold(true).Foreground(tui.ColorPrimary).
		Render(fmt.Sprintf("Result: %s", q.Title))
	divider := lipgloss.NewStyle().Foreground(tui.ColorBorder).
		Render(strings.Repeat("─", min(m.width-8, 60)))

	var status string
	if r == nil {
		status = tui.ErrorStyle.Render("No result")
	} else if r.Error != "" {
		status = lipgloss.JoinVertical(lipgloss.Left,
			tui.ErrorStyle.Render("✗ "+r.Error),
			"",
			lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(min(m.width-10, 80)).Render(r.Output),
		)
	} else {
		var sb strings.Builder
		if r.Compiled {
			sb.WriteString(tui.SuccessStyle.Render("✓ Compilation: PASS") + "\n")
		} else {
			sb.WriteString(tui.ErrorStyle.Render("✗ Compilation: FAIL") + "\n")
		}
		if r.Passed == r.Total && r.Total > 0 {
			sb.WriteString(tui.SuccessStyle.Render(fmt.Sprintf("✓ Tests: %d/%d PASSED", r.Passed, r.Total)) + "\n")
		} else {
			sb.WriteString(tui.ErrorStyle.Render(fmt.Sprintf("✗ Tests: %d/%d passed", r.Passed, r.Total)) + "\n")
		}
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(tui.ColorMuted).Width(min(m.width-10, 80)).Render(r.Output))
		status = sb.String()
	}

	scoreText := lipgloss.NewStyle().Foreground(tui.ColorAccent).Bold(true).
		Render(fmt.Sprintf("Score: %d/%d", m.score.Correct, m.score.Questions))

	body := lipgloss.JoinVertical(lipgloss.Left, title, divider, "", status, "", scoreText)
	return tui.LayoutPage(body, "r retry  n next question  esc back to list", m.width, m.height)
}
