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
	// Header
	headerStyle := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder())

	headerString := lipgloss.JoinHorizontal(
		0,
		statHeader("(strong)", transactions.strongAmount),
		statHeader("(stale)", transactions.staleAmount),
	)

	content := []string{
		headerStyle.Render(headerString),
	}

	for _, t := range transactions.strong {
		content = append(content, t.View())
	}

	bodyStyle := lipgloss.NewStyle().Padding(0, 1)
	bodyString := lipgloss.JoinVertical(
		0,
		content...,
	)

	return bodyStyle.Render(bodyString)
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
		boldStyle.Render("# transfers"),
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

func (t transaction) View() string {
	timeStyle := lipgloss.NewStyle().Width(10)
	style := lipgloss.NewStyle().Width(14)
	smallStyle := lipgloss.NewStyle().Width(3).MarginRight(5)

	time := fmt.Sprintf("%02d:%02d:%02d", t.timestamp.Hour(), t.timestamp.Minute(), t.timestamp.Second())
	amount := fmt.Sprintf("%d $", t.amount)

	content := lipgloss.JoinHorizontal(
		0,
		timeStyle.Render(time),
		style.Render(t.from),
		smallStyle.Render("->"),
		style.Render(t.to),
		style.Render(amount),
	)
	return content
}
