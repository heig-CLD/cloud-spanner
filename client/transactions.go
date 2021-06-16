package client

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type transaction struct {
	from      string
	to        string
	amount    int64
	timestamp time.Time
}

type Transactions struct {
	strong       []transaction
	strongAmount int64
	staleAmount  int64
}

func (transactions Transactions) View() string {
	rowStyle := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder())

	rowString := lipgloss.JoinHorizontal(
		0,
		statHeader("(strong)", 150),
		statHeader("(stale)", 125),
	)

	return rowStyle.Render(rowString)
}

func statHeader(transactionType string, amount int64) string {
	boldStyle := lipgloss.NewStyle().Bold(true)
	commentStyle := lipgloss.NewStyle().Italic(true)

	titleStyle := lipgloss.NewStyle().
		Width(15).
		MarginRight(3).
		Align(lipgloss.Center)

	titleString := lipgloss.JoinVertical(
		0,
		boldStyle.Render("# transactions"),
		commentStyle.Render(transactionType),
	)

	content := lipgloss.NewStyle().
		Render(fmt.Sprintf("%d", amount))

	blockStyle := lipgloss.NewStyle().MarginLeft(3).MarginRight(3)
	blockString := lipgloss.JoinHorizontal(
		0,
		titleStyle.Render(titleString),
		content,
	)

	return blockStyle.Render(blockString)
}
