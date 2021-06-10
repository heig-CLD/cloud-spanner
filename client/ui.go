package client

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

const (
	headerHeight = 3
	footerHeight = 3
)

type errMsg error

type model struct {
	db db

	ready      bool
	err        error
	richPeople []Rich
	viewport   viewport.Model
}

func initialModel(db db) model {
	var richPeople []Rich
	viewport := viewport.Model{Width: 100, Height: 300}

	return model{
		db:         db,
		richPeople: richPeople,
		viewport:   viewport,
	}
}

func (m model) Init() tea.Cmd {
	return m.db.retrieveUsers()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "esc" {
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg
		return m, nil

	case userMsg:
		m.richPeople = msg
		m.updateViewportContent()
		return m, nil

	case tea.WindowSizeMsg:
		margins := 10

		if !m.ready {
			m.viewport = viewport.Model{Width: msg.Width - margins, Height: msg.Height - margins}
			m.viewport.YPosition = 100
			m.ready = true
		} else {
			m.viewport.Width = msg.Width - margins
			m.viewport.Height = msg.Height - margins
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.ready == false {
		return "Init\n"
	}

	header := " List of rich people "
	header += strings.Repeat("─", m.viewport.Width-runewidth.StringWidth(header))

	footer := fmt.Sprintf(" %3.f%% ", m.viewport.ScrollPercent()*100)
	footer += strings.Repeat("─", m.viewport.Width-runewidth.StringWidth(footer))

	return fmt.Sprintf("%s\n%s\n%s", header, m.viewport.View(), footer)
}

func (m model) updateViewportContent() {
	content := ""

	for _, r := range m.richPeople {
		content += fmt.Sprintf("%s", r.View())
		content += "\n"
	}

	m.viewport.SetContent(content)
}
