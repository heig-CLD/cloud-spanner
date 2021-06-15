package client

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type overview struct {
	total   int64
	poorest int64
	richest int64
}

func (o overview) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		MarginRight(2)

	contentStyle := lipgloss.NewStyle().
		MarginRight(4)

	total := fmt.Sprintf("%d$", o.total)
	poorest := fmt.Sprintf("%d$", o.poorest)
	richest := fmt.Sprintf("%d$", o.richest)

	blockStyle := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder())
	blockString := lipgloss.JoinHorizontal(
		0,
		titleStyle.Render("Total Money"),
		contentStyle.Render(total),
		titleStyle.Render("Richest"),
		contentStyle.Render(richest),
		titleStyle.Render("Poorest"),
		contentStyle.Render(poorest),
	)

	return blockStyle.Render(blockString)
}
