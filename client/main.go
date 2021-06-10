package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-runewidth"
)

const (
	useHighPerformanceRenderer = false

	headerHeight = 3
	footerHeight = 3
)

type errMsg error

type model struct {
	ready      bool
	err        error
	richPeople []Rich
	viewport   viewport.Model
}

func initialModel() model {
	richPeople := RandomRichPeople(100)
	viewport := viewport.Model{Width: 100, Height: 300}

	return model{richPeople: richPeople, viewport: viewport}
}

func (m model) Init() tea.Cmd {
	return nil
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

	case tea.WindowSizeMsg:
		verticalMargins := 10

		if !m.ready {
			content := ""
			for _, r := range m.richPeople {
				content += fmt.Sprintf("%s", r.View())
				content += "\n"
			}

			m.viewport = viewport.Model{Width: msg.Width, Height: msg.Height - verticalMargins}
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(content)
			m.ready = true
		} else {
			m.viewport.Width = msg.Width
			m.viewport.Height = msg.Height - verticalMargins
		}

		if useHighPerformanceRenderer {
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	}

	m.viewport, cmd = m.viewport.Update(msg)
	if useHighPerformanceRenderer {
		cmds = append(cmds, cmd)
	}

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

func StartClient() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if err := p.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
