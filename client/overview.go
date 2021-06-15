package client

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type overview struct {
	total int64
}

func (o overview) View() string {
	nameStyle := lipgloss.NewStyle().Width(15)

	total := fmt.Sprintf("%d", o.total)
	str := nameStyle.Render(total)

	return str
}
