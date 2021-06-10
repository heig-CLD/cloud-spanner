package client

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	quitting   bool
	err        error
	richPeople []Rich
}

func initialModel() model {
	rich := Rich{
		name:              "michael",
		percentOfAllMoney: 0.78,
	}

	richPeople := []Rich{rich, rich}
	return model{richPeople: richPeople}
}

func (m model) Init() tea.Cmd {
	return spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		return m, cmd
	}

}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	str := ""
	for i, r := range m.richPeople {
		str += fmt.Sprintf("%d %s", i, r.View())
		str += "\n"
	}

	return str
}

func StartClient() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
