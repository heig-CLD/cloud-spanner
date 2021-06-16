package client

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	db db

	richPeople   []Rich
	overview     overview
	transactions Transactions
}

func initialModel(db db) model {
	var richPeople []Rich

	return model{
		db:           db,
		richPeople:   richPeople,
		overview:     overview{},
		transactions: Transactions{},
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.db.retrieveTotalUsers(),
		m.db.retrieveTotalMoney(),
		m.db.retrieveRichest(),
		m.db.retrievePoorest(),
		m.db.retrieveUsers(),
		m.db.retrieveTransactions(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case msgUser:
		m.richPeople = msg
		return m, m.db.retrieveUsers()

	case msgTotalUsers:
		m.overview.totalUsers = int64(msg)
		return m, m.db.retrieveTotalUsers()

	case msgTotalMoney:
		m.overview.totalMoney = int64(msg)
		return m, m.db.retrieveTotalMoney()

	case msgRichest:
		m.overview.richest = int64(msg)
		return m, m.db.retrieveRichest()

	case msgPoorest:
		m.overview.poorest = int64(msg)
		return m, m.db.retrievePoorest()

	case msgTransactions:
		m.transactions.strong = msg
		return m, m.db.retrieveTransactions()
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
		m.transactions.View(),
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
		content += fmt.Sprintf("%s\n", r.View())
	}

	return style.Render(content)
}
