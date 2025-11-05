package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/lucasefe/ontop/internal/models"
)

var (
	movePromptStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(gruvboxGreen).
			Padding(1, 2).
			MarginTop(2).
			MarginBottom(2)

	moveOptionStyle = lipgloss.NewStyle().
			Padding(0, 2).
			MarginRight(2)

	moveSelectedStyle = lipgloss.NewStyle().
				Background(gruvboxGreen).
				Foreground(gruvboxBg0).
				Bold(true).
				Padding(0, 2).
				MarginRight(2)
)

// renderMovePrompt renders the move task prompt
func (m Model) renderMovePrompt() string {
	if m.moveTask == nil {
		return "No task selected"
	}

	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(gruvboxGreen).
		Render("Move Task")
	b.WriteString(title + "\n\n")

	// Task info
	taskInfo := lipgloss.NewStyle().
		Foreground(gruvboxGray).
		Render(fmt.Sprintf("Task: %s", m.moveTask.Description))
	b.WriteString(taskInfo + "\n")

	currentColumn := lipgloss.NewStyle().
		Foreground(gruvboxGray).
		Render(fmt.Sprintf("Current: %s", formatColumnName(m.moveTask.Column)))
	b.WriteString(currentColumn + "\n\n")

	// Column selection prompt
	var prompt strings.Builder
	prompt.WriteString("Select destination column:\n\n")

	// Column options
	columns := []string{"INBOX", "IN PROGRESS", "DONE"}
	columnValues := []string{models.ColumnInbox, models.ColumnInProgress, models.ColumnDone}

	for i, col := range columns {
		if i == m.moveSelection {
			prompt.WriteString(moveSelectedStyle.Render(fmt.Sprintf("  %s  ", col)))
		} else {
			prompt.WriteString(moveOptionStyle.Render(col))
		}
	}
	prompt.WriteString("\n\n")

	// Show if moving to same column
	if columnValues[m.moveSelection] == m.moveTask.Column {
		sameColumnMsg := lipgloss.NewStyle().
			Foreground(gruvboxYellow).
			Render("(already in this column)")
		prompt.WriteString(sameColumnMsg + "\n\n")
	}

	b.WriteString(movePromptStyle.Render(prompt.String()))
	b.WriteString("\n")

	// Help view
	helpView := m.help.View(m.keys)
	b.WriteString(helpStyle.Render(helpView))

	return b.String()
}
