package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type model struct {
	db db

	richPeople []Rich
	overview   overview
}

func initialModel(db db) model {
	var richPeople []Rich

	return model{
		db:         db,
		richPeople: richPeople,
		overview:   overview{total: 1000, poorest: 42, richest: 845},
	}
}

func (m model) Init() tea.Cmd {
	return m.db.retrieveUsers()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case userMsg:
		m.richPeople = []Rich(msg)
		return m, m.db.retrieveUsers()
	}

	return m, nil
}

func (m model) View() string {
	mainTitle := lipgloss.NewStyle().
		Width(1000).
		Bold(true).
		Background(lipgloss.Color("#FF7CCB")).
		Foreground(lipgloss.Color("#000000")).
		PaddingLeft(2).
		Align(lipgloss.Left).
		Render("CLD: Cloud Spanner")

	row := lipgloss.JoinHorizontal(
		0,
		m.richPeopleView(),
		m.richPeopleView(),
	)

	topRow := lipgloss.JoinHorizontal(
		0,
		m.overview.View(),
	)

	body := lipgloss.JoinVertical(
		0,
		mainTitle,
		topRow,
		row,
	)
	return body
}

func (m model) richPeopleView() string {
	style := lipgloss.NewStyle().Padding(1)

	content := ""
	for _, r := range m.richPeople {
		content += fmt.Sprintf("%s", r.View())
		content += "\n"
	}

	return style.Render(content)
}
