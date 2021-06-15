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

	total := formatMoney(o.total)
	poorest := formatMoney(o.poorest)
	richest := formatMoney(o.richest)

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

func formatMoney(money int64) string {
	return fmt.Sprintf("%d $", money)
}
