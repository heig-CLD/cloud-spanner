package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	db db

	richPeople []Rich
}

func initialModel(db db) model {
	var richPeople []Rich

	return model{
		db:         db,
		richPeople: richPeople,
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
	content := m.richPeopleView()
	return content
}

func (m model) richPeopleView() string {
	content := ""

	for _, r := range m.richPeople {
		content += fmt.Sprintf("%s", r.View())
		content += "\n"
	}

	return content
}
