package components

import (
	"encoding/json"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/GrapeInTheTree/x402-playground/internal/tui"
)

// JSONView renders JSON data with syntax highlighting.
func JSONView(data []byte, width int) string {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return string(data)
	}

	formatted, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return string(data)
	}

	return colorizeJSON(string(formatted))
}

func colorizeJSON(s string) string {
	keyStyle := lipgloss.NewStyle().Foreground(tui.ColorSecondary)
	strStyle := lipgloss.NewStyle().Foreground(tui.ColorSuccess)
	numStyle := lipgloss.NewStyle().Foreground(tui.ColorAccent)
	boolStyle := lipgloss.NewStyle().Foreground(tui.ColorPrimary).Bold(true)

	var result strings.Builder

	for _, line := range strings.Split(s, "\n") {
		trimmed := strings.TrimSpace(line)
		indent := line[:len(line)-len(trimmed)]

		if idx := strings.Index(trimmed, ":"); idx > 0 && strings.HasPrefix(trimmed, "\"") {
			key := trimmed[:idx]
			val := strings.TrimSpace(trimmed[idx+1:])
			result.WriteString(indent)
			result.WriteString(keyStyle.Render(key))
			result.WriteString(": ")
			result.WriteString(colorizeValue(val, strStyle, numStyle, boolStyle))
			result.WriteString("\n")
		} else {
			result.WriteString(indent)
			result.WriteString(colorizeValue(trimmed, strStyle, numStyle, boolStyle))
			result.WriteString("\n")
		}
	}

	return result.String()
}

func colorizeValue(val string, strStyle, numStyle, boolStyle lipgloss.Style) string {
	trailing := ""
	if strings.HasSuffix(val, ",") {
		trailing = ","
		val = strings.TrimSuffix(val, ",")
	}

	clean := strings.TrimSpace(val)

	switch {
	case strings.HasPrefix(clean, "\""):
		return strStyle.Render(clean) + trailing
	case clean == "true" || clean == "false":
		return boolStyle.Render(clean) + trailing
	case clean == "null":
		return lipgloss.NewStyle().Foreground(tui.ColorMuted).Render(clean) + trailing
	case clean == "{" || clean == "}" || clean == "[" || clean == "]" ||
		clean == "{}" || clean == "[]":
		return clean + trailing
	default:
		if isNumeric(clean) {
			return numStyle.Render(clean) + trailing
		}
		return clean + trailing
	}
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	start := 0
	if s[0] == '-' {
		start = 1
	}
	dotSeen := false
	for i := start; i < len(s); i++ {
		if s[i] == '.' && !dotSeen {
			dotSeen = true
		} else if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return start < len(s)
}
