package client

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type overview struct {
	totalUsers int64
	totalMoney int64
	poorest    int64
	richest    int64
}

func (o overview) View() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		MarginRight(2)

	contentStyle := lipgloss.NewStyle().
		MarginRight(6)

	totalUsers := fmt.Sprintf("%d", o.totalUsers)
	totalMoney := formatMoney(o.totalMoney)
	poorest := formatMoney(o.poorest)
	richest := formatMoney(o.richest)

	blockStyle := lipgloss.NewStyle().
		Margin(0, 1).
		Padding(0, 1).
		Border(lipgloss.RoundedBorder())

	blockString := lipgloss.JoinHorizontal(
		0,
		titleStyle.Render("Users"),
		contentStyle.Render(totalUsers),
		titleStyle.Render("Total Money"),
		contentStyle.Render(totalMoney),
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
