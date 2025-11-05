package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	confirmPromptStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(gruvboxRed).
				Padding(1, 2).
				MarginTop(2).
				MarginBottom(2)

	confirmButtonStyle = lipgloss.NewStyle().
				Padding(0, 2).
				MarginRight(2)

	confirmButtonSelectedStyle = lipgloss.NewStyle().
					Background(gruvboxRed).
					Foreground(gruvboxBg0).
					Bold(true).
					Padding(0, 2).
					MarginRight(2)
)

// renderDeleteConfirm renders the delete confirmation dialog
func (m Model) renderDeleteConfirm() string {
	if m.deleteTask == nil {
		return "No task selected"
	}

	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(gruvboxRed).
		Render("⚠  Delete Task?")
	b.WriteString(title + "\n\n")

	// Task info
	taskInfo := lipgloss.NewStyle().
		Foreground(gruvboxGray).
		Render(fmt.Sprintf("Task: %s", m.deleteTask.Description))
	b.WriteString(taskInfo + "\n\n")

	// Warning
	warning := lipgloss.NewStyle().
		Foreground(gruvboxYellow).
		Render("This action cannot be undone!")
	b.WriteString(warning + "\n\n")

	// Confirmation prompt
	var prompt strings.Builder
	prompt.WriteString("Are you sure?\n\n")

	// Options: Yes / No
	yesText := "Yes, delete"
	noText := "No, cancel"

	// m.moveSelection: 0 = No, 1 = Yes
	if m.moveSelection == 0 {
		prompt.WriteString(confirmButtonSelectedStyle.Render(noText))
		prompt.WriteString(confirmButtonStyle.Render(yesText))
	} else {
		prompt.WriteString(confirmButtonStyle.Render(noText))
		prompt.WriteString(confirmButtonSelectedStyle.Render(yesText))
	}

	b.WriteString(confirmPromptStyle.Render(prompt.String()))
	b.WriteString("\n")

	// Help view
	helpView := m.help.View(m.keys)
	b.WriteString(helpStyle.Render(helpView))

	return b.String()
}
